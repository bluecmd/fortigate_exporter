package logging

import (
	"os"
	"time"

	zaplogfmt "github.com/sykesm/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var initialized = false

func SetupStructureLogging() *zap.SugaredLogger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.UTC().Format(time.RFC3339))
	}
	logger := zap.New(zapcore.NewCore(
		zaplogfmt.NewEncoder(config),
		os.Stdout,
		zapcore.DebugLevel,
	))
	//logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
	initialized = true
	return logger.Sugar()
}

func GetSugar() *zap.SugaredLogger {
	if initialized {
		return zap.S()
	}
	return SetupStructureLogging()
}
