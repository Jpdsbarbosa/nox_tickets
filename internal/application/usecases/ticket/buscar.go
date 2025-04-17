package ticket

import (
	"nox_tickets/internal/domain/ticket"
)

// input do caso de uso de buscar ticket
type BuscarTicketInput struct {
	ID string
}

type ObservacaoOutput struct {
	ID          string
	UsuarioID   string
	Descricao   string
	DataCriacao string
}

type ModificacaoOutput struct {
	ID              string
	UsuarioID       string
	CampoModificado string
	ValorAnterior   string
	ValorNovo       string
	DataModificacao string
}

// Output do caso de uso de buscar ticket
type BuscarTicketOutput struct {
	ID              string
	Titulo          string
	Status          ticket.Status
	Categoria       ticket.Categoria
	Subcategoria    ticket.Subcategoria
	Descricao       string
	Urgencia        int
	Gravidade       int
	AbertoPor       string
	Responsavel     string
	DataAbertura    string
	DataInicio      string
	DataConclusao   string
	DuracaoTotal    string
	DuracaoExecucao string

	// Histórico de observações e modificações
	Observacoes  []ObservacaoOutput
	Modificacoes []ModificacaoOutput

	// Campos opcionais (alterando para ponteiros)
	Merchant   *string
	NoxID      *string
	CPF        *string
	Plataforma *string
	Contato    string
}

// Caso de uso de buscar ticket
type BuscarTicketUseCase struct {
	ticketRepository ticket.Repository
}

// Executa o caso de uso de buscar ticket
func (uc *BuscarTicketUseCase) Execute(input BuscarTicketInput) (*BuscarTicketOutput, error) {
	// 1. Busca o ticket no banco
	ticket, err := uc.ticketRepository.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	// 2. Formata as datas (converte nil para string vazia quando necessário)
	dataInicio := ""
	if ticket.DataInicio != nil {
		dataInicio = ticket.DataInicio.Format("2006-01-02 15:04:05")
	}

	dataConclusao := ""
	if ticket.DataConclusao != nil {
		dataConclusao = ticket.DataConclusao.Format("2006-01-02 15:04:05")
	}

	// 3. Converte as observações do ticket para o formato de saída
	observacoes := make([]ObservacaoOutput, len(ticket.Observacoes))
	for i, obs := range ticket.Observacoes {
		observacoes[i] = ObservacaoOutput{
			ID:          obs.ID,
			UsuarioID:   obs.UsuarioID,
			Descricao:   obs.Descricao,
			DataCriacao: obs.DataCriacao.Format("2006-01-02 15:04:05"),
		}
	}

	// 4. Converte as modificações do ticket para o formato de saída
	modificacoes := make([]ModificacaoOutput, len(ticket.Modificacoes))
	for i, mod := range ticket.Modificacoes {
		modificacoes[i] = ModificacaoOutput{
			ID:              mod.ID,
			UsuarioID:       mod.UsuarioID,
			CampoModificado: mod.CampoModificado,
			ValorAnterior:   mod.ValorAnterior,
			ValorNovo:       mod.ValorNovo,
			DataModificacao: mod.DataModificacao.Format("2006-01-02 15:04:05"),
		}
	}

	// 5. Retorna todos os dados formatados
	return &BuscarTicketOutput{
		ID:              ticket.ID,
		Titulo:          ticket.Titulo,
		Status:          ticket.Status,
		Categoria:       ticket.Categoria,
		Subcategoria:    ticket.Subcategoria,
		Descricao:       ticket.Descricao,
		Urgencia:        ticket.Urgencia,
		Gravidade:       ticket.Gravidade,
		AbertoPor:       ticket.AbertoPor,
		Responsavel:     ticket.Responsavel,
		DataAbertura:    ticket.DataAbertura.Format("2006-01-02 15:04:05"),
		DataInicio:      dataInicio,
		DataConclusao:   dataConclusao,
		DuracaoTotal:    ticket.DuracaoTotal.String(),
		DuracaoExecucao: ticket.DuracaoExecucao.String(),

		// adiciona as observações e modificações
		Observacoes:  observacoes,
		Modificacoes: modificacoes,

		// adiciona os campos opcionais
		Merchant:   ticket.Merchant,
		NoxID:      ticket.NoxID,
		CPF:        ticket.CPF,
		Plataforma: ticket.Plataforma,
		Contato:    ticket.Contato,
	}, nil
}

// NewBuscarTicketUseCase cria uma nova instância do caso de uso de buscar ticket
func NewBuscarTicketUseCase(ticketRepository ticket.Repository) *BuscarTicketUseCase {
	return &BuscarTicketUseCase{
		ticketRepository: ticketRepository,
	}
}
