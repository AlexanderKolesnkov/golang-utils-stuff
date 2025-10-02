package utils

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func NewZapLogger(dirName, fileName, level string, isLogFile bool) (logFile *os.File, logger *zap.Logger, err error) {
	if isLogFile {
		logFile, logger, err = NewZapFileConsoleLogger(dirName, fileName, level)
		if err != nil {
			return logFile, logger, err
		}
	} else {
		logger, err = NewZapDevelopLogger(level)
		if err != nil {
			return logFile, logger, err
		}
	}

	return logFile, logger, nil
}

func NewZapDevelopLogger(level string) (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()

	logLevel := switchLogLeve(level)

	cfg.Level.SetLevel(logLevel)
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	logger = logger.WithOptions(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return logger, err
}

func NewZapFileConsoleLogger(dirName, fileName, level string) (*os.File, *zap.Logger, error) {
	consoleEncoderConfig := zap.NewDevelopmentConfig()
	consoleEncoderConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig.EncoderConfig)
	consoleWriter := zapcore.AddSync(os.Stdout)

	logLevel := switchLogLeve(level)

	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, logLevel)

	logFile, err := CreateLogFile(dirName, fileName)
	if err != nil {
		return nil, nil, fmt.Errorf("CreateLogFile(%s, %s): %v", dirName, fileName, err)
	}

	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.TimeKey = "time"
	fileEncoderConfig.EncodeTime = utcTimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)
	fileWriter := zapcore.AddSync(logFile)

	fileCore := zapcore.NewCore(fileEncoder, fileWriter, zapcore.DebugLevel)

	core := zapcore.NewTee(consoleCore, fileCore)
	logger := zap.New(core)

	logger = logger.WithOptions(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return logFile, logger, nil
}

func CreateLogFile(dirName, fileName string) (*os.File, error) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("unable to create log directory: %v", err)
		}
	}

	filePath := fmt.Sprintf("%s/%s_%s.log", dirName, fileName, time.Now().Format(time.RFC3339))
	logFile, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf(filePath, err)
	}

	return logFile, nil
}

func utcTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format(time.RFC3339))
}

func switchLogLeve(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
