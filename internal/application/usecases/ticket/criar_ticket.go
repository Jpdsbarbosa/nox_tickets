package ticket

import (
	"nox_tickets/internal/domain/ticket"
)

// input do use case de criar ticket
type CriarTicketInput struct {
	Titulo       string
	Descricao    string
	Categoria    ticket.Categoria
	Subcategoria ticket.Subcategoria
	AbertoPor    string
	Urgencia     int
	Gravidade    int

	//campos opcionais
	Merchant    string
	NoxID       string
	CPF         string
	Plataforma  string
	Contato     string
	Responsavel string
}

// output do use case de criar ticket
type CriarTicketOutput struct {
	ID           string
	Status       ticket.Status
	DataAbertura string
	AbertoPor    string
	Titulo       string
	Categoria    ticket.Categoria
	Subcategoria ticket.Subcategoria
}

// use case de criar ticket
type CriarTicketUseCase struct {
	ticketRepository ticket.Repository
}

// Construtor do use case de criar ticket
func NewCriarTicketUseCase(repo ticket.Repository) *CriarTicketUseCase {
	return &CriarTicketUseCase{
		ticketRepository: repo,
	}
}

// Executa o use case de criar ticket
func (uc *CriarTicketUseCase) Execute(input CriarTicketInput) (*CriarTicketOutput, error) {
	// validações adicionais
	if input.Urgencia < 1 || input.Urgencia > 5 {
		return nil, ticket.ErrUrgenciaInvalida
	}
	if input.Gravidade < 1 || input.Gravidade > 5 {
		return nil, ticket.ErrGravidadeInvalida
	}

	// Criar o ticket usando a entidade de dominio
	novoTicket, err := ticket.NovoTicket(
		input.Titulo,
		input.Descricao,
		input.Categoria,
		input.Subcategoria,
		input.AbertoPor,
	)
	if err != nil {
		return nil, err
	}

	// Define a urgencia e a gravidade
	novoTicket.Urgencia = input.Urgencia
	novoTicket.Gravidade = input.Gravidade

	// Adiciona informações adicionais se fornecidas
	novoTicket.SetInformacaoAdicional(
		input.Merchant,
		input.NoxID,
		input.CPF,
		input.Plataforma,
		input.Contato,
	)

	// Define o responsável se fornecido
	if input.Responsavel != "" {
		if err := novoTicket.IniciarAtendimento(input.Responsavel); err != nil {
			return nil, err
		}
	}

	// Persiste o ticket usando o repositório
	err = uc.ticketRepository.Create(novoTicket)
	if err != nil {
		return nil, err
	}

	// Retorna o output com as informações do ticket criado
	return &CriarTicketOutput{
		ID:           novoTicket.ID,
		Status:       novoTicket.Status,
		DataAbertura: novoTicket.DataAbertura.Format("2006-01-02 15:04:05"),
		AbertoPor:    novoTicket.AbertoPor,
		Titulo:       novoTicket.Titulo,
		Categoria:    novoTicket.Categoria,
		Subcategoria: novoTicket.Subcategoria,
	}, nil
}
