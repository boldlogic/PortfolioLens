package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func writeSyncers(cfg Config) []zapcore.WriteSyncer {
	if cfg.OutputFile == "" {
		return []zapcore.WriteSyncer{zapcore.Lock(os.Stdout)}
	}
	f, err := os.OpenFile(cfg.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return []zapcore.WriteSyncer{zapcore.Lock(os.Stdout)}
	}
	return []zapcore.WriteSyncer{zapcore.Lock(zapcore.AddSync(f))}
}

func New(cfg Config) *zap.Logger {
	lvl, _ := zapcore.ParseLevel(cfg.Level)

	atomicLevel := zap.NewAtomicLevelAt(lvl)
	writers := writeSyncers(cfg)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	var encoder zapcore.Encoder
	if strings.ToUpper(cfg.Format) == "JSON" {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), atomicLevel)
	return zap.New(core, zap.AddCaller())
}
