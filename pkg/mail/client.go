package mail

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Mailer interface {
	Send(ctx context.Context, letter Letter) error
}

type Letter struct {
	From    string
	To      []string
	Body    string
	Subject string
}

type PoorMailer struct {
	sync.Mutex

	stdFrom string
	logger  *zap.Logger
	driver  *smtp.Client
}

func NewMailer(
	driver *smtp.Client,
	logger *zap.Logger,
	stdFrom string,
) Mailer {
	return &PoorMailer{
		driver:  driver,
		stdFrom: stdFrom,
		logger:  logger,
	}
}

func (c *PoorMailer) Send(ctx context.Context, letter Letter) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if err := c.driver.Reset(); err != nil {
		return err
	}

	from := letter.From
	if from == "" {
		from = c.stdFrom
	}

	if len(letter.To) == 0 {
		return errors.New("Must specify at least one recipient")
	}

	if err := c.driver.Mail(from); err != nil {
		return nil
	}

	for _, to := range letter.To {
		if err := c.driver.Rcpt(to); err != nil {
			return err
		}
	}

	bodyWriter, err := c.driver.Data()
	if err != nil {
		return nil
	}

	defer bodyWriter.Close()

	_, err = bodyWriter.Write([]byte(buildHtmlContent(letter)))
	if err != nil {
		return errors.Wrap(err, "failed while writing smtp message body")
	}

	c.logger.
		With(zap.String("from", from)).
		With(zap.String("to", strings.Join(letter.To, ","))).
		Debug("Sent email")

	return nil
}

func (c *PoorMailer) SendTxt(ctx context.Context, letter Letter) error {
	// one transaction per connection
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if err := c.driver.Reset(); err != nil {
		return err
	}

	from := letter.From
	if from == "" {
		from = c.stdFrom
	}

	if len(letter.To) == 0 {
		return errors.New("Must specify at least one recipient")
	}

	if err := c.driver.Mail(from); err != nil {
		return nil
	}

	for _, to := range letter.To {
		if err := c.driver.Rcpt(to); err != nil {
			return err
		}
	}

	bodyWriter, err := c.driver.Data()
	if err != nil {
		return nil
	}

	defer bodyWriter.Close()

	_, err = bodyWriter.Write([]byte(letter.Body))
	if err != nil {
		return errors.Wrap(err, "failed while writing smtp message body")
	}

	c.logger.
		With(zap.String("from", from)).
		With(zap.String("to", strings.Join(letter.To, ","))).
		Debug("Sent email")

	return nil
}

func buildHtmlContent(letter Letter) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", letter.From)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(letter.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", letter.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", letter.Body)

	return msg
}
