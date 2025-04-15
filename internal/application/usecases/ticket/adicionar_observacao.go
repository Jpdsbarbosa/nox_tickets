package ticket

import (
	"errors"
	"nox_tickets/internal/domain/ticket"
	"time"
)

var (
	ErrDescricaoVazia = errors.New("descrição da observação é obrigatória")
)

// input de usecase de adicionar observação
type AdicionarObservacaoInput struct {
	ID        string
	Descricao string
	UsuarioID string
}

// output de usecase de adicionar observação
type AdicionarObservacaoOutput struct {
	ID          string
	TicketID    string
	UsuarioID   string
	Descricao   string
	DataCriacao string
}

// usecase de adicionar observação
type AdicionarObservacaoUseCase struct {
	ticketRepository ticket.Repository
}

// construtor de usecase de adicionar observação
func NewAdicionarObservacaoUseCase(repo ticket.Repository) *AdicionarObservacaoUseCase {
	return &AdicionarObservacaoUseCase{
		ticketRepository: repo,
	}
}

// executa a usecase de adicionar observação
func (uc *AdicionarObservacaoUseCase) Execute(input AdicionarObservacaoInput) (*AdicionarObservacaoOutput, error) {
	// 1. validacao basica da observacao
	if input.Descricao == "" {
		return nil, ErrDescricaoVazia
	}

	// 2. Busca o ticket para garantir que existe
	ticketExistente, err := uc.ticketRepository.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	// 3. cria e adiciona a observação
	err = ticketExistente.AdicionarObservacao(input.Descricao, input.UsuarioID)
	if err != nil {
		return nil, err
	}

	// 4. pega a observacao que acabou de ser adicionada (ultima da lista)
	observacao := ticketExistente.Observacoes
	novaObservacao := observacao[len(observacao)-1]

	// 5. persiste alterações
	err = uc.ticketRepository.Update(ticketExistente)
	if err != nil {
		return nil, err
	}

	// 6. prepara o output
	return &AdicionarObservacaoOutput{
		ID:          input.ID,
		TicketID:    input.ID,
		UsuarioID:   input.UsuarioID,
		Descricao:   novaObservacao.Descricao,
		DataCriacao: novaObservacao.DataCriacao.Format(time.DateTime),
	}, nil
}
