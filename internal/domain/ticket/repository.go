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

	// Adicionar observação
	AdicionarObservacao(ticketID string, observacao *Observacao) error

	// Listar observações
	ListarObservacoes(ticketID string) ([]*Observacao, error)

	// Para rastrear modificacoes
	AdicionarModificacao(ticketID string, modificacao *Modificacao) error

	// Listar modificacoes
	ListarModificacoes(ticketID string) ([]*Modificacao, error)

	// Buscar por status específicos
	ListarPorStatus(status Status) ([]*Ticket, error)

	// Para atualizar status específicos
	AtualizarStatus(ticketID string, novoStatus Status, usuarioID string) error
}

// TicketFiltros define os filtros possíveis para busca
type TicketFiltros struct {
	Status      []Status
	Categoria   []Categoria
	Responsavel string
	AbertoPor   string
	DataInicio  time.Time
	DataFim     time.Time
	Urgencia    *int
	Gravidade   *int
}
