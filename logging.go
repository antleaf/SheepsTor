package sheepstor

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func init() {
	logger, _ = ConfigureZapSugarLogger(false)
}

func EnableDefaultLogger() {
	logger, _ = ConfigureZapSugarLogger(true)
}

func SetLogger(l *zap.SugaredLogger) {
	logger = l
}

func ConfigureZapSugarLogger(debugging bool) (*zap.SugaredLogger, error) {
	var zapLogger *zap.Logger
	var err error
	zapLogger, err = ConfigureZapLogger(debugging)
	return zapLogger.Sugar(), err
}

func ConfigureZapLogger(debugging bool) (*zap.Logger, error) {
	level := zapcore.FatalLevel
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		TimeKey:     "",
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
	}
	if debugging == true {
		level = zapcore.DebugLevel
		encoderConfig = zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			TimeKey:      "",
			EncodeLevel:  zapcore.CapitalColorLevelEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		}
	}
	zapConfig := zap.Config{
		Encoding:      "console",
		Level:         zap.NewAtomicLevelAt(level),
		OutputPaths:   []string{"stdout"},
		EncoderConfig: encoderConfig,
	}
	return zapConfig.Build()
}
