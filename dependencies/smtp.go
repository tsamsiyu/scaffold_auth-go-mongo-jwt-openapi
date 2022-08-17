package dependencies

import (
	"context"
	"net/smtp"

	"apart-deal-api/pkg/mail"

	"github.com/Netflix/go-env"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SmtpConfig struct {
	SmtpAddress string `env:"SMTP_ADDR,required=true"`
	SmtpFrom    string `env:"SMTP_FROM,required=true"`
}

func NewSmtpConfig() (*SmtpConfig, error) {
	var cfg SmtpConfig

	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func NewSmtpClient(cfg *SmtpConfig) (*smtp.Client, error) {
	client, err := smtp.Dial(cfg.SmtpAddress)
	if err != nil {
		return nil, err
	}

	if err := client.Hello("localhost"); err != nil {
		return nil, err
	}

	return client, nil
}

func NewMailer(client *smtp.Client, logger *zap.Logger, cfg *SmtpConfig) *mail.Client {
	return mail.NewClient(client, logger, cfg.SmtpFrom)
}

var SmtpModule = fx.Module(
	"Smtp",
	fx.Provide(
		NewSmtpConfig,
		NewSmtpClient,
		NewMailer,
	),
	fx.Invoke(func(lc fx.Lifecycle, client *smtp.Client) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				_ = client.Quit()

				return nil
			},
		})
	}),
)
