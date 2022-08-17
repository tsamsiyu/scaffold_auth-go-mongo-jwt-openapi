package mail

import (
	"context"
	"net/smtp"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type SendPayload struct {
	From string
	To   []string
	Body string
}

type Client struct {
	stdFrom string
	logger  *zap.Logger
	driver  *smtp.Client
}

func NewClient(
	driver *smtp.Client,
	logger *zap.Logger,
	stdFrom string,
) *Client {
	return &Client{
		driver:  driver,
		stdFrom: stdFrom,
		logger:  logger,
	}
}

func (c *Client) SendTxt(ctx context.Context, payload SendPayload) error {
	from := payload.From
	if from == "" {
		from = c.stdFrom
	}

	if len(payload.To) == 0 {
		return errors.New("Must specify at least one recipient")
	}

	if err := c.driver.Mail(from); err != nil {
		return nil
	}

	for _, to := range payload.To {
		if err := c.driver.Rcpt(to); err != nil {
			return err
		}
	}

	bodyWriter, err := c.driver.Data()
	if err != nil {
		return nil
	}

	defer bodyWriter.Close()

	_, err = bodyWriter.Write([]byte(payload.Body))
	if err != nil {
		return errors.Wrap(err, "failed while writing smtp message body")
	}

	c.logger.
		With(zap.Field{
			Key:    "from",
			String: from,
		}).
		With(zap.Field{
			Key:    "to",
			String: strings.Join(payload.To, ","),
		}).
		Debug("Sent email")

	return nil
}
