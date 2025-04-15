package ticket

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Status string

var (
	ErrUrgenciaInvalida  = errors.New("urgência inválida")
	ErrGravidadeInvalida = errors.New("gravidade inválida")
)

const (
	StatusAberto     Status = "aberto"
	StatusEmCurso    Status = "em_curso"
	StatusFinalizado Status = "finalizado"
	StatusCancelado  Status = "cancelado"
)

type Categoria string

const (
	CategoriaFinanceiro  Categoria = "financeiro"
	CategoriaComercial   Categoria = "comercial"
	CategoriaCompliance  Categoria = "compliance"
	CategoriaContratos   Categoria = "contratos"
	CategoriaGestores    Categoria = "gestores"
	CategoriaMeds        Categoria = "meds"
	CategoriaOnboarding  Categoria = "onboarding"
	CategoriaOperacional Categoria = "operacional"
	CategoriaReclamacoes Categoria = "reclamacoes"
	CategoriaTI          Categoria = "ti"
	CategoriaTrading     Categoria = "trading"
)

type Subcategoria string

const (
	SubcategoriaBug                  Subcategoria = "bug"
	SubcategoriaFeature              Subcategoria = "feature"
	SubcategoriaMelhoria             Subcategoria = "melhoria"
	SubcategoriaOutros               Subcategoria = "outros"
	SubcategoriaSolicitacaoEnviada   Subcategoria = "solicitacao_enviada"
	SubcategoriaDuvidas              Subcategoria = "duvidas"
	SubcategoriaSolicitacoes         Subcategoria = "solicitacoes"
	SubcategoriaFraude               Subcategoria = "fraude"
	SubcategoriaKYC                  Subcategoria = "kyc"
	SubcategoriaUncompliant          Subcategoria = "uncompliant"
	SubcategoriaCadastroDocumentacao Subcategoria = "cadastro_documentacao"
	SubcategoriaVerificacaoTransacao Subcategoria = "verificacao_de_transacao"
	SubcategoriaSolicitacaoSaque     Subcategoria = "solicitacao_de_saque"
)

type Observacao struct {
	ID          string
	TicketID    string
	UsuarioID   string
	Descricao   string
	DataCriacao time.Time
}

type Modificacao struct {
	ID              string
	TicketID        string
	UsuarioID       string
	CampoModificado string
	ValorAnterior   string
	ValorNovo       string
	DataModificacao time.Time
}

type Ticket struct {
	ID              string
	Titulo          string
	Merchant        *string
	NoxID           *string
	CPF             *string
	Status          Status
	Categoria       Categoria
	Subcategoria    Subcategoria
	Descricao       string
	Urgencia        int
	Gravidade       int
	AbertoPor       string
	Responsavel     string
	Contato         string
	Plataforma      *string
	DataAbertura    time.Time
	DataInicio      *time.Time
	DataConclusao   *time.Time
	DuracaoTotal    time.Duration
	DuracaoExecucao time.Duration
	Observacoes     []Observacao
	Modificacoes    []Modificacao
}

func NovoTicket(titulo, descricao string, categoria Categoria, subcategoria Subcategoria, abertoPor string) (*Ticket, error) {
	if titulo == "" {
		return nil, errors.New("titulo é obrigatório")
	}
	if descricao == "" {
		return nil, errors.New("descrição é obrigatória")
	}
	if categoria == "" {
		return nil, errors.New("categoria é obrigatória")
	}
	if abertoPor == "" {
		return nil, errors.New("abertura por é obrigatório")
	}

	return &Ticket{
		ID:           uuid.New().String(),
		Titulo:       titulo,
		Descricao:    descricao,
		Categoria:    categoria,
		Subcategoria: subcategoria,
		Status:       StatusAberto,
		AbertoPor:    abertoPor,
		DataAbertura: time.Now(),
		Urgencia:     1, // valor padrão, pode ser alterado depois
		Gravidade:    1, // valor padrão, pode ser alterado depois
	}, nil
}

// SetInformacaoAdicional define informações opcionais do ticket
func (t *Ticket) SetInformacaoAdicional(merchant, noxID, cpf, plataforma, contato string) {
	if merchant != "" {
		t.Merchant = &merchant
	}
	if noxID != "" {
		t.NoxID = &noxID
	}
	if cpf != "" {
		t.CPF = &cpf
	}
	if plataforma != "" {
		t.Plataforma = &plataforma
	}
	if contato != "" {
		t.Contato = contato
	}
}

// IniciarAtendimento inicia o atendimento do ticket
func (t *Ticket) IniciarAtendimento(responsavel string) error {
	if t.Status != StatusAberto {
		return errors.New("ticket não pode ser iniciado pois não está aberto")
	}

	agora := time.Now()
	t.Status = StatusEmCurso
	t.Responsavel = responsavel
	t.DataInicio = &agora

	return t.registrarModificacao("status", string(StatusAberto), string(StatusEmCurso), responsavel)
}

// Concluir finaliza o ticket
func (t *Ticket) Concluir(usuarioID string) error {
	if t.Status != StatusEmCurso {
		return errors.New("ticket não pode ser concluído pois não está em atendimento")
	}

	agora := time.Now()
	statusAnterior := t.Status
	t.Status = StatusFinalizado
	t.DataConclusao = &agora

	// Calcula a duracao do Ticket
	t.DuracaoTotal = agora.Sub(t.DataAbertura)
	if t.DataInicio != nil {
		t.DuracaoExecucao = agora.Sub(*t.DataInicio)
	}

	return t.registrarModificacao("status", string(statusAnterior), string(StatusFinalizado), usuarioID)
}

// Cancelar cancela o ticket
func (t *Ticket) Cancelar(usuarioID string) error {
	if t.Status == StatusFinalizado || t.Status == StatusCancelado {
		return errors.New("ticket não pode ser cancelado pois já foi finalizado ou cancelado")
	}

	statusAnterior := t.Status
	t.Status = StatusCancelado

	return t.registrarModificacao("status", string(statusAnterior), string(StatusCancelado), usuarioID)
}

// AdicionarObservacao - adiciona uma nova observacao no ticket
func (t *Ticket) AdicionarObservacao(descricao, usuarioID string) error {
	observacao := Observacao{
		ID:          uuid.New().String(),
		TicketID:    t.ID,
		UsuarioID:   usuarioID,
		Descricao:   descricao,
		DataCriacao: time.Now(),
	}

	t.Observacoes = append(t.Observacoes, observacao)

	return nil
}

func (t *Ticket) SetUrgencia(urgencia int, usuarioID string) error {
	if urgencia < 1 || urgencia > 5 {
		return errors.New("urgência inválida")
	}

	valorAnterior := t.Urgencia
	t.Urgencia = urgencia
	return t.registrarModificacao("urgencia", fmt.Sprintf("%d", valorAnterior), fmt.Sprintf("%d", urgencia), usuarioID)
}

func (t *Ticket) SetGravidade(gravidade int, usuarioID string) error {
	if gravidade < 1 || gravidade > 5 {
		return errors.New("gravidade inválida")
	}

	valorAnterior := t.Gravidade
	t.Gravidade = gravidade
	return t.registrarModificacao("gravidade", fmt.Sprintf("%d", valorAnterior), fmt.Sprintf("%d", gravidade), usuarioID)
}

// registrarModificacao - registra uma nova modificacao no ticket
func (t *Ticket) registrarModificacao(campo, valorAnterior, valorNovo, usuarioID string) error {
	modificacao := Modificacao{
		ID:              uuid.New().String(),
		TicketID:        t.ID,
		UsuarioID:       usuarioID,
		CampoModificado: campo,
		ValorAnterior:   valorAnterior,
		ValorNovo:       valorNovo,
		DataModificacao: time.Now(),
	}

	t.Modificacoes = append(t.Modificacoes, modificacao)
	return nil
}
