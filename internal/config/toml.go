package config

import (
	"io"
	"os"

	"github.com/IglooCloud/igloo-observability/internal/api"
	"github.com/IglooCloud/igloo-observability/internal/collector"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	DBPath    string
	Collector collector.Config
	API       api.Config
}

type tomlConfig struct {
	Database struct {
		Path string
	}
	API struct {
		Port    int
		SSLPort int
		SSLCert string
		SSLKey  string
	}
	Collector struct {
		Endpoints []struct {
			URL      string
			Secret   string
			Schedule string
			Enabled  bool
			Bucket   string
		}
	}
}

// Loads config from TOML file
func LoadTOML(path string) (Config, error) {
	// Read file
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	// Parse TOML
	var rawConfig tomlConfig
	err = toml.Unmarshal(buffer, &rawConfig)
	if err != nil {
		return Config{}, err
	}

	// Convert to config
	config := Config{
		DBPath: rawConfig.Database.Path,
		Collector: collector.Config{
			Schedules: make([]collector.EndpointSchedule, len(rawConfig.Collector.Endpoints)),
		},
		API: api.Config{
			Port:    rawConfig.API.Port,
			SSLPort: rawConfig.API.SSLPort,
			SSLCert: rawConfig.API.SSLCert,
			SSLKey:  rawConfig.API.SSLKey,
		},
	}
	for i, endpoint := range rawConfig.Collector.Endpoints {
		config.Collector.Schedules[i] = collector.EndpointSchedule{
			Endpoint: collector.Endpoint{
				URL:    endpoint.URL,
				Secret: endpoint.Secret,
				Bucket: endpoint.Bucket,
			},
			Schedule: endpoint.Schedule,
			Enabled:  endpoint.Enabled,
		}
	}

	return config, nil
}
