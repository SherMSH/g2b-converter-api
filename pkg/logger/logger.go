package logger

import (
	"converterapi/internal/config"
	"io"
	"log"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var (
	globalLogger Logger
	once         sync.Once
	mu           sync.RWMutex
)

// Init initializes the global logger with options
func Init() {
	log.Println("Zlogger init...")
	once.Do(func() {
		lvl := zerolog.Level(0)
		if !config.Config.App.DebugMode {
			lvl = zerolog.Level(1)
		}
		// Build zerolog logger
		zl := zerolog.New(os.Stdout).Level(lvl).With().Timestamp().Caller()

		if config.Config.App.Server.Name != "" {
			zl = zl.Str("service", config.Config.App.Server.Name)
		}

		globalLogger = &loggerWrapper{zl.Logger()}
	})
	globalLogger.Info().Msgf("Zlogger initialized!")
}

// Get returns the global logger instance
func Get() Logger {
	if globalLogger == nil {
		Init() // Initialize with defaults if not already initialized
	}
	return globalLogger
}

// Debug writes a new message with debug level
func Debug(format string, v ...interface{}) {
	Get().Debug().Msgf(format, v...)
}

// Info writes a new message with info level
func Info(format string, v ...interface{}) {
	Get().Info().Msgf(format, v...)
}

// Warn starts a new message with warn level
func Warn(format string, v ...interface{}) {
	Get().Warn().Msgf(format, v...)
}

// Error starts a new message with error level
func Error(format string, v ...interface{}) {
	Get().Error().Msgf(format, v...)
}

// Fatal starts a new message with fatal level
func Fatal(format string, v ...interface{}) {
	Get().Fatal().Msgf(format, v...)
}

// Panic starts a new message with panic level
func Panic(format string, v ...interface{}) {
	Get().Panic().Msgf(format, v...)
}

// Level sets logger level
func (l *loggerWrapper) Level(level zerolog.Level) Logger {
	return Get().Level(level)
}

// Output sets logger output
func (l *loggerWrapper) Output(w io.Writer) Logger {
	return Get().Output(w)
}
