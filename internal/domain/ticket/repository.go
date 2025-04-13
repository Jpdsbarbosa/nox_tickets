package ticket

import "time"

type Repository interface {
	// Criar novo ticket
	Create(ticket *Ticket) error

	// Buscar ticket por ID
	GetByID(id string) (*Ticket, error)

	// Listar tickets com filtros
	List(filtros TicketFiltros) ([]*Ticket, error)

	// Atualizar ticket
	Update(ticket *Ticket) error

	// Deletar ticket (se necessário)
	Delete(id string) error
}

// TicketFiltros define os filtros possíveis para busca
type TicketFiltros struct {
	Status      []Status
	Categoria   []Categoria
	Responsavel string
	AbertoPor   string
	DataInicio  time.Time
	DataFim     time.Time
	// ... outros filtros necessários
}
