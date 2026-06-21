package application

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/bboykiv/topsigner/internal/config"
)

func NewLogger(config *config.Config) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderCfg)

	stdoutFilter := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		if config.Debug {
			return level < zapcore.ErrorLevel
		}

		return level >= zapcore.InfoLevel && level < zapcore.ErrorLevel
	})

	stderrFilter := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), stdoutFilter),
		zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), stderrFilter),
	)

	return zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
}
