package dependencies

import (
	"apart-deal-api/pkg/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewAppConfig(logger *zap.Logger) *config.Config {
	return &config.Config{
		IsDebug: logger.Core().Enabled(zapcore.DebugLevel),
	}
}

var ConfigModule = fx.Provide(
	NewAppConfig,
)
