package ports

import "github.com/joseMChavez/fc-job/internal/domain"

type PDFGenerator interface {
	Generate(invoice *domain.Invoice) ([]byte, error)
}
