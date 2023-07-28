package log

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/jwalton/go-supportscolor"
	"github.com/rs/zerolog"
)

var defaultLogger Logger

func setupWriter(out io.Writer) ConsoleWriter {
	noColor := !supportscolor.Stderr().SupportsColor

	output := ConsoleWriter{
		Out:        out,
		TimeFormat: time.RFC3339,
		NoColor:    noColor,
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			ServiceFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
	}

	// Print caller path relative to current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}
	divider := colorize(" |", colorCyan, noColor)
	output.FormatCaller = func(i interface{}) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			if rel, err := filepath.Rel(cwd, c); err == nil {
				c = rel
			}
			c = c + divider
		}
		return c
	}

	return output
}

var DefaultOutput = os.Stderr

func Default() Logger {
	// Get visible log levels from environment
	var loggerLevel zerolog.Level
	switch os.Getenv("LOGGER_LEVEL") {
	case "TRACE":
		loggerLevel = TraceLevel
	case "DEBUG":
		loggerLevel = DebugLevel
	case "INFO":
		loggerLevel = InfoLevel
	case "WARN":
		loggerLevel = WarnLevel
	case "ERROR":
		loggerLevel = ErrorLevel
	case "FATAL":
		loggerLevel = FatalLevel
	case "PANIC":
		loggerLevel = PanicLevel
	default:
		loggerLevel = DebugLevel
	}

	logger := zerolog.
		New(setupWriter(DefaultOutput)).
		With().
		Timestamp().
		CallerWithSkipFrameCount(3).
		Logger().
		Level(loggerLevel)

	return Logger{
		logger: logger,
	}
}
