package configuration

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"time"

	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/pkg/msend"
	"ewintr.nl/gte/pkg/mstore"
)

var (
	ErrUnableToRead = errors.New("unable to read configuration")
)

type Configuration struct {
	IMAPURL      string
	IMAPUsername string
	IMAPPassword string

	SMTPURL      string
	SMTPUsername string
	SMTPPassword string

	FromName    string
	FromAddress string

	ToName    string
	ToAddress string

	LocalDBPath string
}

type LocalConfiguration struct {
	MinSyncInterval time.Duration
}

func New(src io.Reader) *Configuration {
	conf := &Configuration{}
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		if len(line) != 2 {

			continue
		}

		key, value := strings.TrimSpace(line[0]), strings.TrimSpace(line[1])
		switch key {
		case "imap_url":
			conf.IMAPURL = value
		case "imap_username":
			conf.IMAPUsername = value
		case "imap_password":
			conf.IMAPPassword = value
		case "smtp_url":
			conf.SMTPURL = value
		case "smtp_username":
			conf.SMTPUsername = value
		case "smtp_password":
			conf.SMTPPassword = value
		case "to_name":
			conf.ToName = value
		case "to_address":
			conf.ToAddress = value
		case "from_name":
			conf.FromName = value
		case "from_address":
			conf.FromAddress = value
		case "local_db_path":
			conf.LocalDBPath = value
		}
	}

	return conf
}

func (c *Configuration) IMAP() *mstore.IMAPConfig {
	return &mstore.IMAPConfig{
		IMAPURL:      c.IMAPURL,
		IMAPUsername: c.IMAPUsername,
		IMAPPassword: c.IMAPPassword,
	}
}

func (c *Configuration) SMTP() *msend.SSLSMTPConfig {
	return &msend.SSLSMTPConfig{
		URL:      c.SMTPURL,
		Username: c.SMTPUsername,
		Password: c.SMTPPassword,
		From:     c.FromAddress,
		To:       c.ToAddress,
	}
}

func (c *Configuration) Sqlite() *storage.SqliteConfig {
	return &storage.SqliteConfig{
		DBPath: c.LocalDBPath,
	}
}
