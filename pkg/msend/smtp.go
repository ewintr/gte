package msend

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
)

var (
	ErrSMTPInvalidConfig    = errors.New("invalid smtp configuration")
	ErrSMTPConnectionFailed = errors.New("connection to smtp server failed")
	ErrSendMessageFailed    = errors.New("could not send message")
)

type SSLSMTPConfig struct {
	URL      string
	Username string
	Password string
	From     string
	To       string
}

func (ssc *SSLSMTPConfig) Valid() bool {
	if _, _, err := net.SplitHostPort(ssc.URL); err != nil {
		return false
	}

	return ssc.Username != "" && ssc.Password != "" && ssc.To != "" && ssc.From != ""
}

type SSLSMTP struct {
	config    *SSLSMTPConfig
	client    *smtp.Client
	connected bool
}

func NewSSLSMTP(config *SSLSMTPConfig) *SSLSMTP {
	return &SSLSMTP{
		config: config,
	}
}

func (s *SSLSMTP) Connect() error {
	if !s.config.Valid() {
		return ErrSMTPInvalidConfig
	}

	host, _, _ := net.SplitHostPort(s.config.URL)
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, host)
	conn, err := tls.Dial("tcp", s.config.URL, &tls.Config{ServerName: host})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSMTPConnectionFailed, err)
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSMTPConnectionFailed, err)
	}
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("%w: %v", ErrSMTPConnectionFailed, err)
	}
	s.client = client
	s.connected = true

	return nil
}

func (s *SSLSMTP) Close() error {
	if !s.connected {
		return nil
	}

	if err := s.client.Quit(); err != nil {
		return fmt.Errorf("%w: %v", ErrSMTPConnectionFailed, err)
	}
	s.connected = false

	return nil
}

func (s *SSLSMTP) Send(msg *Message) error {
	if err := s.Connect(); err != nil {
		return err
	}
	defer s.Close()

	from := mail.Address{
		Name:    "gte",
		Address: s.config.From,
	}
	to := mail.Address{
		Name:    "todo",
		Address: s.config.To,
	}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = msg.Subject

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += fmt.Sprintf("\r\n%s", msg.Body)

	if err := s.client.Mail(s.config.From); err != nil {
		return fmt.Errorf("%w: %v", ErrSendMessageFailed, err)
	}
	if err := s.client.Rcpt(s.config.To); err != nil {
		return fmt.Errorf("%w: %v", ErrSendMessageFailed, err)
	}
	wc, err := s.client.Data()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSendMessageFailed, err)
	}
	if _, err := wc.Write([]byte(message)); err != nil {
		return fmt.Errorf("%w: %v", ErrSendMessageFailed, err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("%w: %v", ErrSendMessageFailed, err)
	}

	return nil
}
