package dependencies

import (
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLoggingLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
)

func NewLogger(lvl *zap.AtomicLevel) *zap.Logger {
	loggerCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= lvl.Level()
		}),
	)

	return zap.New(loggerCore)
}

func LoggerFromEnv() *zap.Logger {
	zapLevel := getZapLevelFromEnv()
	return NewLogger(zapLevel)
}

func getZapLevelFromEnv() *zap.AtomicLevel {
	rawLevel := os.Getenv("LOG_LEVEL")
	if rawLevel != "" {
		zapLevel, err := zap.ParseAtomicLevel(rawLevel)
		if err != nil {
			panic(errors.Wrapf(err, "Failed while parsing zap log level: %s", rawLevel))
		}

		return &zapLevel
	}

	return &defaultLoggingLevel
}
