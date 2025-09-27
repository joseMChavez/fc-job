package pdf_gofpdf

import (
	"bytes"
	"fmt"
	"time"

	"github.com/joseMChavez/fc-job/internal/domain"
	"github.com/joseMChavez/fc-job/internal/ports"
	"github.com/jung-kurt/gofpdf"
)

type Generator struct{}

func NewPDFGenerator() ports.PDFGenerator {
	return &Generator{}
}

func (g *Generator) Generate(inv *domain.Invoice) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Factura")

	pdf.Ln(20)
	pdf.SetFont("Arial", "", 12)

	pdf.Cell(40, 10, fmt.Sprintf("ID: %d", inv.ID))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Cliente: %s", inv.CustomerEmail))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Monto: %d ", inv.TotalAmount))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Fecha: %s", inv.DueDate.Format(time.RFC822)))

	// Generar PDF en memoria
	buf := new(bytes.Buffer)
	if err := pdf.Output(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
