package component

import (
	"go-mod.ewintr.nl/gte/internal/configuration"
	"go-mod.ewintr.nl/gte/internal/storage"
	"go-mod.ewintr.nl/gte/pkg/msend"
	"go-mod.ewintr.nl/gte/pkg/mstore"
	"fyne.io/fyne/v2"
)

type Configuration struct {
	prefs fyne.Preferences
	conf  configuration.Configuration
}

func NewConfigurationFromPreferences(prefs fyne.Preferences) *Configuration {
	return &Configuration{
		prefs: prefs,
		conf:  configuration.Configuration{},
	}
}

func (c *Configuration) Load() {
	c.conf.IMAPURL = c.prefs.String("ConfigIMAPURL")
	c.conf.IMAPUsername = c.prefs.String("ConfigIMAPUser")
	c.conf.IMAPPassword = c.prefs.String("ConfigIMAPPassword")
	c.conf.IMAPFolderPrefix = c.prefs.String("ConfigIMAPFolderPrefix")
	c.conf.SMTPURL = c.prefs.String("ConfigSMTPURL")
	c.conf.SMTPUsername = c.prefs.String("ConfigSMTPUser")
	c.conf.SMTPPassword = c.prefs.String("ConfigSMTPPassword")
	c.conf.FromName = c.prefs.String("ConfigGTEFromName")
	c.conf.FromAddress = c.prefs.String("ConfigGTEFromAddress")
	c.conf.ToName = c.prefs.String("ConfigGTEToName")
	c.conf.ToAddress = c.prefs.String("ConfigGTEToAddress")
	c.conf.LocalDBPath = c.prefs.String("ConfigGTELocalDBPath")
}

func (c *Configuration) Fields() map[string]string {
	return map[string]string{
		"ConfigIMAPURL":          c.conf.IMAPURL,
		"ConfigIMAPUser":         c.conf.IMAPUsername,
		"ConfigIMAPPassword":     c.conf.IMAPPassword,
		"ConfigIMAPFolderPrefix": c.conf.IMAPFolderPrefix,
		"ConfigSMTPURL":          c.conf.SMTPURL,
		"ConfigSMTPUser":         c.conf.SMTPUsername,
		"ConfigSMTPPassword":     c.conf.SMTPPassword,
		"ConfigGTEToName":        c.conf.ToName,
		"ConfigGTEToAddress":     c.conf.ToAddress,
		"ConfigGTEFromName":      c.conf.FromName,
		"ConfigGTEFromAddress":   c.conf.FromAddress,
		"ConfigGTELocalDBPath":   c.conf.LocalDBPath,
	}
}

func (c *Configuration) Set(key, value string) {
	switch key {
	case "ConfigIMAPURL":
		c.conf.IMAPURL = value
		c.prefs.SetString("ConfigIMAPURL", c.conf.IMAPURL)
	case "ConfigIMAPUser":
		c.conf.IMAPUsername = value
		c.prefs.SetString("ConfigIMAPUser", c.conf.IMAPUsername)
	case "ConfigIMAPPassword":
		c.conf.IMAPPassword = value
		c.prefs.SetString("ConfigIMAPPassword", c.conf.IMAPPassword)
	case "ConfigIMAPFolderPrefix":
		c.conf.IMAPFolderPrefix = value
		c.prefs.SetString("ConfigIMAPFolderPrefix", c.conf.IMAPFolderPrefix)
	case "ConfigSMTPURL":
		c.conf.SMTPURL = value
		c.prefs.SetString("ConfigSMTPURL", c.conf.SMTPURL)
	case "ConfigSMTPUser":
		c.conf.SMTPUsername = value
		c.prefs.SetString("ConfigSMTPUser", c.conf.SMTPUsername)
	case "ConfigSMTPPassword":
		c.conf.SMTPPassword = value
		c.prefs.SetString("ConfigSMTPPassword", c.conf.SMTPPassword)
	case "ConfigGTEToName":
		c.conf.ToName = value
		c.prefs.SetString("ConfigGTEToName", c.conf.ToName)
	case "ConfigGTEToAddress":
		c.conf.ToAddress = value
		c.prefs.SetString("ConfigGTEToAddress", c.conf.ToAddress)
	case "ConfigGTEFromName":
		c.conf.FromName = value
		c.prefs.SetString("ConfigGTEFromName", c.conf.FromName)
	case "ConfigGTEFromAddress":
		c.conf.FromAddress = value
		c.prefs.SetString("ConfigGTEFromAddress", c.conf.FromAddress)
	case "ConfigGTELocalDBPath":
		c.conf.LocalDBPath = value
		c.prefs.SetString("ConfigGTELocalDBPath", c.conf.LocalDBPath)
	}
}

func (c *Configuration) IMAP() *mstore.IMAPConfig {
	return c.conf.IMAP()
}

func (c *Configuration) SMTP() *msend.SSLSMTPConfig {
	return c.conf.SMTP()
}

func (c *Configuration) Sqlite() *storage.SqliteConfig {
	return c.conf.Sqlite()
}
