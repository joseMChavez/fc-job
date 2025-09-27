package ports

type Attachment struct {
	Filename string
	Data     []byte
}

type EmailSender interface {
	// Send envia un mail. Implementación concreta maneja TLS, auth y attachments.
	Send(to, subject, body string, attachments []Attachment) error
}
