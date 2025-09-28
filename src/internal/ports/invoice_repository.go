package ports

import "github.com/joseMChavez/fc-job/src/internal/domain"

type InvoiceRepository interface {
	GetByID(id int64) (*domain.Invoice, error)
	MarkAsSent(id int64) error
	FindPending() ([]*domain.Invoice, error)
	GetDetails(id int64) ([]*domain.InvoiceItem, error)
}
