package service

import (
	"crypto/tls"
	"errors"
	"fmt"
	"goapp/config"
	"goapp/lang"
	"net/smtp"
	"os"
	"strconv"
	"strings"

	"github.com/jordan-wright/email"
)

type MailAttachment struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type Mail struct {
	From        string
	To          string
	Bcc         []string
	Cc          []string
	Bc          []string
	Subject     string
	Text        string
	HTML        string
	ReadReceipt []string
	Attachments []MailAttachment
	locale      lang.Locale
}

// Send sendet eine Email
func (m Mail) SendMail(config *config.Smtp) (*string, error) {

	e := email.NewEmail()

	if m.To == "" {
		return &m.locale.Error.Mail.To_is_missing, errors.New("mail 'to' is missing")
	}

	if m.Subject == "" {
		return &m.locale.Error.Mail.Subject_is_missing, errors.New("mail 'Subject' is missing")
	}

	e.From = config.Sender
	e.To = strings.Split(m.To, ",")
	if m.From != "" {
		e.ReplyTo = strings.Split(m.From, ",")
	}
	e.Subject = m.Subject

	if m.Bcc != nil {
		e.Bcc = m.Bcc
	}
	if m.Cc != nil {
		e.Cc = m.Cc
	}

	if m.ReadReceipt != nil {
		e.ReadReceipt = m.ReadReceipt
	}

	if m.Attachments != nil {

		m.Text = m.Text + "\n"

		if m.HTML != "" {
			m.HTML = m.HTML + "<br>"
		}

		for _, attachment := range m.Attachments {

			osFile, err := os.Open(attachment.Path)
			if err != nil {
				return &m.locale.Error.Mail.Add_attachment_read_file, errors.New("could not add mail attachment, file not found or broke")
			}
			defer osFile.Close()

			r := (*os.File)(osFile)
			// this is doing a type conversion from *os.File to io.Reader
			_, err = e.Attach(r, attachment.Name, attachment.Type)

			if err != nil {
				return &m.locale.Error.Mail.Add_attachment_on_attach, errors.New("could not add mail attachment")
			}
		}

	}

	e.Text = []byte(m.Text)

	if m.HTML != "" {
		e.HTML = []byte(m.HTML)
	}

	var err error

	var auth smtp.Auth
	if config.User != "" {
		auth = smtp.PlainAuth("", config.User, config.Password, config.Host)
	}

	if config.SendWithTLS {
		err = e.SendWithTLS(config.Host+":"+strconv.Itoa(config.Port), auth, &tls.Config{InsecureSkipVerify: config.SkipSSLVerify})
	} else if config.SendWithStartTLS {
		err = e.SendWithStartTLS(config.Host+":"+strconv.Itoa(config.Port), auth, &tls.Config{InsecureSkipVerify: config.SkipSSLVerify})
	} else {
		err = e.Send(config.Host+":"+strconv.Itoa(config.Port), auth)
	}

	errMsg := fmt.Sprintf("%v", err)
	return &errMsg, err

}
