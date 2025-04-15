package ticket

import (
	"nox_tickets/internal/domain/ticket"
	"time"
)

// Input para atualizar um ticket
type AtualizarTicketInput struct {
	ID         string
	Titulo     *string
	Descricao  *string
	Categoria  *ticket.Categoria
	Urgencia   *int
	Gravidade  *int
	Merchant   *string
	NoxID      *string
	CPF        *string
	Plataforma *string
	Contato    *string
	UsuarioID  string
}

// Output para caso de uso de atualizar ticket
type AtualizarTicketOutput struct {
	ID           string
	Titulo       string
	Status       ticket.Status
	Categoria    ticket.Categoria
	Urgencia     int
	Gravidade    int
	DataAbertura string
	AbertoPor    string
	Responsavel  string
}

// Usecase de atualizar ticket
type AtualizarTicketUseCase struct {
	ticketRepository ticket.Repository
}

// Contrutor do caso de uso
func NewAtualizarTicketUseCase(repo ticket.Repository) *AtualizarTicketUseCase {
	return &AtualizarTicketUseCase{
		ticketRepository: repo,
	}
}

// Executa o caso de uso de atualizar ticket
func (uc *AtualizarTicketUseCase) Execute(input AtualizarTicketInput) (*AtualizarTicketOutput, error) {
	// 1. busca o ticket existente
	ticketExistente, err := uc.ticketRepository.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	// 2. atualiza os campos fornecidos
	if input.Titulo != nil {
		ticketExistente.Titulo = *input.Titulo
	}
	if input.Descricao != nil {
		ticketExistente.Descricao = *input.Descricao
	}
	if input.Categoria != nil {
		ticketExistente.Categoria = *input.Categoria
	}
	if input.Urgencia != nil {
		if err := ticketExistente.SetUrgencia(*input.Urgencia, input.UsuarioID); err != nil {
			return nil, err
		}
	}
	if input.Gravidade != nil {
		if err := ticketExistente.SetGravidade(*input.Gravidade, input.UsuarioID); err != nil {
			return nil, err
		}
	}

	// 3. atualiza informacoes adicionais se forem fornecidas
	merchant := ""
	if input.Merchant != nil {
		merchant = *input.Merchant
	}
	noxID := ""
	if input.NoxID != nil {
		noxID = *input.NoxID
	}
	cpf := ""
	if input.CPF != nil {
		cpf = *input.CPF
	}
	plataforma := ""
	if input.Plataforma != nil {
		plataforma = *input.Plataforma
	}
	contato := ""
	if input.Contato != nil {
		contato = *input.Contato
	}

	ticketExistente.SetInformacaoAdicional(
		merchant,
		noxID,
		cpf,
		plataforma,
		contato,
	)

	// 4. Persiste as alteracoes
	if err := uc.ticketRepository.Update(ticketExistente); err != nil {
		return nil, err
	}

	// 5. retorna o ticket atualizado
	return &AtualizarTicketOutput{
		ID:           ticketExistente.ID,
		Titulo:       ticketExistente.Titulo,
		Status:       ticketExistente.Status,
		Categoria:    ticketExistente.Categoria,
		Urgencia:     ticketExistente.Urgencia,
		Gravidade:    ticketExistente.Gravidade,
		DataAbertura: ticketExistente.DataAbertura.Format(time.DateTime),
		AbertoPor:    ticketExistente.AbertoPor,
		Responsavel:  ticketExistente.Responsavel,
	}, nil
}

// 4. atualiza status se fornecido

// Executa o caso de uso
