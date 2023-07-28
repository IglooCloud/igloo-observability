package api

import (
	"crypto/tls"
	"fmt"

	"net/http"
	"time"

	"github.com/IglooCloud/igloo-observability/internal/log"
	"github.com/IglooCloud/igloo-observability/internal/warehouse"
	"github.com/dyson/certman"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var logger = log.Default().Service("api")

type Config struct {
	Port    int
	SSLPort int
	SSLCert string
	SSLKey  string
}

// Run starts an HTTP and HTTPS server with automatic reloading
// of SSL certificates and compatible with automatic renewal
// of the certificates through Let's Encrypt.
func Start(gauge warehouse.Gauge, config Config) {
	router := gin.Default()

	// Remove all CORS blocks
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Serving the .well-known route to allow automatic
	// Let's Encrypt certificate renewal
	router.Static("/.well-known", "./.well-known")

	registerGaugeRoutes(router, gauge)

	// Setup server listening on SSL port using certman to
	// automatically reload the SSL certificate
	if config.SSLCert != "" {
		cm, err := certman.New(config.SSLCert, config.SSLKey)
		if err != nil {
			logger.Panic(err)
		}
		cm.Logger(logger)
		if err := cm.Watch(); err != nil {
			logger.Panic(err)
		}

		// Listen on the port specified in the config
		s := &http.Server{
			Addr:    fmt.Sprintf(":%d", config.SSLPort),
			Handler: router,
			TLSConfig: &tls.Config{
				GetCertificate: cm.GetCertificate,
			},
		}

		go (func() {
			logger.Infof("HTTPS server listening on %s", s.Addr)
			if err := s.ListenAndServeTLS("", ""); err != nil {
				logger.Panic(err)
			}
		})()
	}

	err := router.Run(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		logger.Panic(err)
	}
}
