package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func newLogger() *zap.Logger {
	// Open or create the log file
	file, err := os.OpenFile("dev.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	// Use the development encoder config
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create a core with console-style output to a file
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(file),
		zapcore.DebugLevel, // More verbose for development
	)

	logger := zap.New(core)
	defer logger.Sync()

	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning")
	logger.Error("This is an error")

	return logger
}

var Logger = newLogger()
