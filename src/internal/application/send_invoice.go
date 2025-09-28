package application

import (
	"fmt"

	"github.com/joseMChavez/fc-job/src/internal/ports"
)

type InvoiceSender struct {
	Repo        ports.InvoiceRepository
	PdfGen      ports.PDFGenerator
	EmailSender ports.EmailSender
}

func NewInvoiceSender(r ports.InvoiceRepository, e ports.EmailSender, p ports.PDFGenerator) *InvoiceSender {
	return &InvoiceSender{Repo: r, EmailSender: e, PdfGen: p}
}
func (usercase *InvoiceSender) SendInvoice() error {
	invoices, err := usercase.Repo.FindPending()
	if err != nil {
		return err
	}
	var months = []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}

	for _, inv := range invoices {
		subj := fmt.Sprintf("Factura #%d, %s %d", inv.ID, months[inv.Month-1], inv.Year)
		body := fmt.Sprintf("Hola, Adjunto factura correspondiente al mes de %s, saludos!", months[inv.Month-1])

		var attachments []ports.Attachment
		if usercase.PdfGen != nil {
			pdf, err := usercase.PdfGen.Generate(inv)
			if err != nil {
				return err
			}
			attachments = append(attachments, ports.Attachment{
				Filename: fmt.Sprintf("factura-%d.pdf", inv.ID),
				Data:     pdf,
			})
		}

		if err := usercase.EmailSender.Send(inv.CustomerEmail, subj, body, attachments); err != nil {
			return err
		}

		if err := usercase.Repo.MarkAsSent(inv.ID); err != nil {
			return err
		}
	}
	return nil
}
