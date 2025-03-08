package logger

import (
	"os"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger initializes the zap logger
func InitLogger(development bool, noLogs bool, logPath string) error {
	var cfg zap.Config
	var options []zap.Option

	if development {
		cfg = zap.NewDevelopmentConfig()
		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	} else {
		cfg = zap.NewProductionConfig()
	}

	// Handle log suppression
	if noLogs {
		cfg.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	}

	// Handle log file output
	if logPath != "" {
		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return errors.Wrap(err, "failed to create log file")
		}

		fileEncoder := zapcore.NewJSONEncoder(cfg.EncoderConfig)
		fileWriteSyncer := zapcore.AddSync(logFile)
		fileCore := zapcore.NewCore(fileEncoder, fileWriteSyncer, cfg.Level)

		// Create logger with both stdout and file outputs
		consoleEncoder := zapcore.NewConsoleEncoder(cfg.EncoderConfig)
		consoleWriteSyncer := zapcore.AddSync(os.Stdout)
		consoleCore := zapcore.NewCore(consoleEncoder, consoleWriteSyncer, cfg.Level)

		core := zapcore.NewTee(fileCore, consoleCore)
		logger := zap.New(core, options...)
		zap.ReplaceGlobals(logger)
	} else {
		// Create standard logger
		logger, err := cfg.Build(options...)
		if err != nil {
			return errors.Wrap(err, "failed to initialize logger")
		}
		zap.ReplaceGlobals(logger)
	}

	return nil
}

// Sync flushes any buffered log entries
func Sync() {
	_ = zap.L().Sync()
	_ = zap.S().Sync()
}
