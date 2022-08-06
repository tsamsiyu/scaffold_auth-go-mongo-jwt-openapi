package logging

import (
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func FromEnv() (*zap.Logger, error) {
	zapLevel, err := getZapLevel()
	if err != nil {
		return nil, err
	}

	loggerCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= zapLevel.Level()
		}),
	)

	return zap.New(loggerCore), nil
}

func getZapLevel() (*zap.AtomicLevel, error) {
	rawLevel := os.Getenv("LOG_LEVEL")
	if rawLevel != "" {
		zapLevel, err := zap.ParseAtomicLevel(rawLevel)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed while parsing zap log level: %s", rawLevel)
		}

		return &zapLevel, nil
	}

	return nil, nil
}
