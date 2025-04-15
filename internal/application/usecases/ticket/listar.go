package ticket

import (
	"nox_tickets/internal/domain/ticket"
)

// input - estrutura que define os filtros e paginação para listar tickets
type ListarTicketsInput struct {
	// Filtros opcionais
	Status      *ticket.Status
	Categoria   *ticket.Categoria
	Responsavel *string

	// paginação
	Pagina         int // número da página (1-based)
	ItensPorPagina int // número de itens por página
}

// output resumido de cada ticket na listagem
type TicketResumoOutput struct {
	ID           string
	Titulo       string
	Status       ticket.Status
	Categoria    ticket.Categoria
	Urgencia     int
	Gravidade    int
	AbertoPor    string
	Responsavel  string
	DataAbertura string
}

// Output da listagem com paginação
type ListarTicketsOutput struct {
	Tickets      []TicketResumoOutput
	Total        int // número total de tickets
	TotalPaginas int // número total de páginas
	PaginaAtual  int // número da página atual
}

// Caso de uso de listar tickets
type ListarTicketsUseCase struct {
	ticketRepository ticket.Repository
}

// Construtor do usecase
func NewListarTicketsUseCase(repo ticket.Repository) *ListarTicketsUseCase {
	return &ListarTicketsUseCase{
		ticketRepository: repo,
	}
}

// Execute executa o caso de uso de listar tickets
func (uc *ListarTicketsUseCase) Execute(input ListarTicketsInput) (*ListarTicketsOutput, error) {
	// Validação básica de paginação
	if input.Pagina < 1 {
		input.Pagina = 1
	}
	if input.ItensPorPagina < 1 {
		input.ItensPorPagina = 10
	}

	// Cria filtro para a busca
	filtros := ticket.TicketFiltros{
		Status:      []ticket.Status{},
		Categoria:   []ticket.Categoria{},
		Responsavel: "",
	}

	// Adiciona filtros se fornecidos
	if input.Status != nil {
		filtros.Status = append(filtros.Status, *input.Status)
	}
	if input.Categoria != nil {
		filtros.Categoria = append(filtros.Categoria, *input.Categoria)
	}
	if input.Responsavel != nil {
		filtros.Responsavel = *input.Responsavel
	}

	// Busca os tickets com filtros
	tickets, err := uc.ticketRepository.List(filtros)
	if err != nil {
		return nil, err
	}

	// implementa a paginação manual
	total := len(tickets)
	inicio := (input.Pagina - 1) * input.ItensPorPagina
	fim := inicio + input.ItensPorPagina
	if fim > total {
		fim = total
	}

	// Aplica paginação
	ticketsPaginados := tickets
	if inicio < total {
		ticketsPaginados = tickets[inicio:fim]
	} else {
		ticketsPaginados = []*ticket.Ticket{}
	}

	// Converte para o formato de saída
	ticketsOutput := make([]TicketResumoOutput, len(ticketsPaginados))
	for i, t := range ticketsPaginados {
		ticketsOutput[i] = TicketResumoOutput{
			ID:           t.ID,
			Titulo:       t.Titulo,
			Status:       t.Status,
			Categoria:    t.Categoria,
			Urgencia:     t.Urgencia,
			Gravidade:    t.Gravidade,
			AbertoPor:    t.AbertoPor,
			Responsavel:  t.Responsavel,
			DataAbertura: t.DataAbertura.Format("2006-01-02 15:04:05"),
		}
	}

	// calcula total de páginas
	totalPaginas := (total + input.ItensPorPagina - 1) / input.ItensPorPagina

	return &ListarTicketsOutput{
		Tickets:      ticketsOutput,
		Total:        total,
		TotalPaginas: totalPaginas,
		PaginaAtual:  input.Pagina,
	}, nil
}
