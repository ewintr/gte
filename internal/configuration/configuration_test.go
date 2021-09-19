package configuration_test

import (
	"strings"
	"testing"

	"ewintr.nl/go-kit/test"
	"ewintr.nl/gte/internal/configuration"
	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/pkg/msend"
	"ewintr.nl/gte/pkg/mstore"
)

func TestNew(t *testing.T) {
	for _, tc := range []struct {
		name   string
		source string
		exp    *configuration.Configuration
	}{
		{
			name: "empty",
			exp:  &configuration.Configuration{},
		},
		{
			name:   "lines without values",
			source: "test\n\n876\nkey=",
			exp:    &configuration.Configuration{},
		},
		{
			name:   "trim space",
			source: " imap_url\t= value",
			exp: &configuration.Configuration{
				IMAPURL: "value",
			},
		},
		{
			name:   "value with space",
			source: "imap_url=one two three",
			exp: &configuration.Configuration{
				IMAPURL: "one two three",
			},
		},
		{
			name:   "imap",
			source: "imap_url=url\nimap_username=username\nimap_password=password",
			exp: &configuration.Configuration{
				IMAPURL:      "url",
				IMAPUsername: "username",
				IMAPPassword: "password",
			},
		},
		{
			name:   "smtp",
			source: "smtp_url=url\nsmtp_username=username\nsmtp_password=password",
			exp: &configuration.Configuration{
				SMTPURL:      "url",
				SMTPUsername: "username",
				SMTPPassword: "password",
			},
		},
		{
			name:   "addresses",
			source: "to_name=to_name\nto_address=to_address\nfrom_name=from_name\nfrom_address=from_address",
			exp: &configuration.Configuration{
				ToName:      "to_name",
				ToAddress:   "to_address",
				FromName:    "from_name",
				FromAddress: "from_address",
			},
		},
		{
			name:   "local",
			source: "local_db_path=path",
			exp: &configuration.Configuration{
				LocalDBPath: "path",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, configuration.New(strings.NewReader(tc.source)))
		})
	}
}

func TestConfigs(t *testing.T) {
	conf := &configuration.Configuration{
		IMAPURL:      "imap_url",
		IMAPUsername: "imap_username",
		IMAPPassword: "imap_password",
		SMTPURL:      "smtp_url",
		SMTPUsername: "smtp_username",
		SMTPPassword: "smtp_password",
		ToName:       "to_name",
		ToAddress:    "to_address",
		FromName:     "from_name",
		FromAddress:  "from_address",
		LocalDBPath:  "db_path",
	}

	t.Run("imap", func(t *testing.T) {
		exp := &mstore.IMAPConfig{
			IMAPURL:      "imap_url",
			IMAPUsername: "imap_username",
			IMAPPassword: "imap_password",
		}

		test.Equals(t, exp, conf.IMAP())
	})

	t.Run("smtp", func(t *testing.T) {
		exp := &msend.SSLSMTPConfig{
			URL:      "smtp_url",
			Username: "smtp_username",
			Password: "smtp_password",
			To:       "to_address",
			From:     "from_address",
		}
		test.Equals(t, exp, conf.SMTP())
	})

	t.Run("sqlite", func(t *testing.T) {
		exp := &storage.SqliteConfig{
			DBPath: "db_path",
		}
		test.Equals(t, exp, conf.Sqlite())
	})
}
