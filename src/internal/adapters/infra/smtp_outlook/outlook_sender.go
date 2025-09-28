package smtp_outlook

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/joseMChavez/fc-job/src/internal/ports"
)

// OutlookSender envía correo usando smtp.office365.com (STARTTLS)
type OutlookSender struct {
	Host     string // smtp.office365.com
	Port     string // "587"
	Username string
	Password string
}

func NewOutlookSender(host, port, user, pass string) *OutlookSender {
	return &OutlookSender{Host: host, Port: port, Username: user, Password: pass}
}

func (s *OutlookSender) Send(to, subject, body string, attachments []ports.Attachment) error {
	// Nota: este ejemplo envía texto plano y NO maneja attachments.
	from := s.Username
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)

	header := ""
	header += fmt.Sprintf("From: %s\r\n", from)
	header += fmt.Sprintf("To: %s\r\n", to)
	header += fmt.Sprintf("Subject: %s\r\n", subject)
	header += "MIME-Version: 1.0\r\n"
	header += "Content-Type: text/plain; charset=\"utf-8\"\r\n"
	header += "\r\n"

	msg := header + body

	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	// Conectar y STARTTLS
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsconfig := &tls.Config{ServerName: s.Host}
		if err := c.StartTLS(tlsconfig); err != nil {
			return err
		}
	}

	if err := c.Auth(auth); err != nil {
		return err
	}

	if err := c.Mail(from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	return c.Quit()
}
