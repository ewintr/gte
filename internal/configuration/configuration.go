package configuration

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"code.ewintr.nl/gte/internal/storage"
	"code.ewintr.nl/gte/pkg/msend"
	"code.ewintr.nl/gte/pkg/mstore"
)

var (
	ErrUnableToRead = errors.New("unable to read configuration")
)

type Configuration struct {
	IMAPURL          string
	IMAPUsername     string
	IMAPPassword     string
	IMAPFolderPrefix string

	SMTPURL      string
	SMTPUsername string
	SMTPPassword string

	FromName    string
	FromAddress string

	ToName    string
	ToAddress string

	LocalDBPath string
	DaysAhead   int
}

type LocalConfiguration struct {
	MinSyncInterval time.Duration
}

func NewFromFile(src io.Reader) *Configuration {
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
		case "imap_folder_prefix":
			conf.IMAPFolderPrefix = value
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

func NewFromEnvironment() *Configuration {
	days, _ := strconv.Atoi(os.Getenv("GTE_DAYS_AHEAD"))

	return &Configuration{
		IMAPURL:          os.Getenv("IMAP_URL"),
		IMAPUsername:     os.Getenv("IMAP_USER"),
		IMAPPassword:     os.Getenv("IMAP_PASSWORD"),
		IMAPFolderPrefix: os.Getenv("IMAP_FOLDER_PREFIX"),
		SMTPURL:          os.Getenv("SMTP_URL"),
		SMTPUsername:     os.Getenv("SMTP_USER"),
		SMTPPassword:     os.Getenv("SMTP_PASSWORD"),
		ToName:           os.Getenv("GTE_TO_NAME"),
		ToAddress:        os.Getenv("GTE_TO_ADDRESS"),
		FromName:         os.Getenv("GTE_FROM_NAME"),
		FromAddress:      os.Getenv("GTE_FROM_ADDRESS"),
		DaysAhead:        days,
	}
}

func (c *Configuration) IMAP() *mstore.IMAPConfig {
	return &mstore.IMAPConfig{
		IMAPURL:          c.IMAPURL,
		IMAPUsername:     c.IMAPUsername,
		IMAPPassword:     c.IMAPPassword,
		IMAPFolderPrefix: c.IMAPFolderPrefix,
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
