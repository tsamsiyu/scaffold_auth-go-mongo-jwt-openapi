package dependencies

import (
	"context"
	"net"
	"net/smtp"
	"time"

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
	conn, err := net.DialTimeout("tcp", cfg.SmtpAddress, 5*time.Second)
	if err != nil {
		return nil, err
	}

	clientCh := make(chan *smtp.Client)
	defer close(clientCh)

	errCh := make(chan error)
	defer close(errCh)

	go func() {
		// NewMailer may hang forever
		client, err := smtp.NewClient(conn, cfg.SmtpAddress)
		if err != nil {
			errCh <- err
		} else {
			clientCh <- client
		}
	}()

	var client *smtp.Client

	select {
	case <-time.After(time.Second * 5):
		panic("Creating new smtp client hanged")
	case client = <-clientCh:
		break
	case err := <-errCh:
		return nil, err
	}

	// TODO: set correct hostname
	if err := client.Hello("localhost"); err != nil {
		return nil, err
	}

	return client, nil
}

func NewMailer(client *smtp.Client, logger *zap.Logger, cfg *SmtpConfig) mail.Mailer {
	return mail.NewMailer(client, logger, cfg.SmtpFrom)
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
