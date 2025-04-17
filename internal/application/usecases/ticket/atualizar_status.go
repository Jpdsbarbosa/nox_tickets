package ticket

import (
	"errors"
	"nox_tickets/internal/domain/ticket"
)

var (
	ErrStatusInvalido = errors.New("status invalido")
)

// input do usecase de atualizar status
type AtualizarStatusTicketInput struct {
	ID          string
	Status      ticket.Status
	UsuarioID   string
	Responsavel *string
}

// output do usecase de atualizar status
type AtualizarStatusTicketOutput struct {
	ID            string
	Status        ticket.Status
	Responsavel   string
	DataInicio    string
	DataConclusao string
}

// Usecase de atualizar status
type AtualizarStatusUseCase struct {
	ticketRepository ticket.Repository
}

// construtor do usecase de atualizar status
func NewAtualizarStatusUseCase(repo ticket.Repository) *AtualizarStatusUseCase {
	return &AtualizarStatusUseCase{
		ticketRepository: repo,
	}
}

// Executa o usecase de atualizar status
func (uc *AtualizarStatusUseCase) Execute(input AtualizarStatusTicketInput) (*AtualizarStatusTicketOutput, error) {
	// 1. busca o ticket existente
	ticketExistente, err := uc.ticketRepository.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	// 2. aplica a mudança de status de acordo com o novo status
	switch input.Status {
	case ticket.StatusEmCurso:
		if input.Responsavel == nil {
			return nil, errors.New("responsavel é obrigatório para iniciar o atendimento")
		}
		err = ticketExistente.IniciarAtendimento(*input.Responsavel)

	case ticket.StatusFinalizado:
		err = ticketExistente.Concluir(input.UsuarioID)

	case ticket.StatusCancelado:
		err = ticketExistente.Cancelar(input.UsuarioID)

	default:
		return nil, ErrStatusInvalido
	}

	if err != nil {
		return nil, err
	}

	// 3. persiste as alteracoes
	if err := uc.ticketRepository.Update(ticketExistente); err != nil {
		return nil, err
	}

	// 4. prepara as data para o output
	dataInicio := ""
	if ticketExistente.DataInicio != nil {
		dataInicio = ticketExistente.DataInicio.Format("2006-01-02 15:04:05")
	}

	dataFim := ""
	if ticketExistente.DataConclusao != nil {
		dataFim = ticketExistente.DataConclusao.Format("2006-01-02 15:04:05")
	}

	// 5. retorna o output
	return &AtualizarStatusTicketOutput{
		ID:            input.ID,
		Status:        input.Status,
		Responsavel:   ticketExistente.Responsavel,
		DataInicio:    dataInicio,
		DataConclusao: dataFim,
	}, nil
}
