package utils

import (
    "os"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
    config := zap.NewProductionConfig()
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    var err error
    Logger, err = config.Build()
    if err != nil {
        panic(err)
    }

    // Set log level based on environment variable
    logLevel := os.Getenv("LOG_LEVEL")
    if logLevel != "" {
        var level zapcore.Level
        err := level.UnmarshalText([]byte(logLevel))
        if err == nil {
            Logger = Logger.WithOptions(zap.IncreaseLevel(level))
        }
    }
}

func GetLogger() *zap.Logger {
    if Logger == nil {
        InitLogger()
    }
    return Logger
}