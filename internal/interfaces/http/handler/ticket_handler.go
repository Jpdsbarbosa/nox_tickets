package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	ticketUseCase "nox_tickets/internal/application/usecases/ticket"
	ticketDomain "nox_tickets/internal/domain/ticket"

	"github.com/go-chi/chi/v5"
)

// TicketHandler contém os handlers relacionados a tickets
type TicketHandler struct {
	criarTicketUseCase         *ticketUseCase.CriarTicketUseCase
	buscarTicketUseCase        *ticketUseCase.BuscarTicketUseCase
	listarTicketsUseCase       *ticketUseCase.ListarTicketsUseCase
	atualizarTicketUseCase     *ticketUseCase.AtualizarTicketUseCase
	atualizarStatusUseCase     *ticketUseCase.AtualizarStatusUseCase
	adicionarObservacaoUseCase *ticketUseCase.AdicionarObservacaoUseCase
}

// NewTicketHandler cria uma nova instancia de TicketHandler
func NewTicketHandler(
	criarTicketUseCase *ticketUseCase.CriarTicketUseCase,
	buscarTicketUseCase *ticketUseCase.BuscarTicketUseCase,
	listarTicketsUseCase *ticketUseCase.ListarTicketsUseCase,
	atualizarTicketUseCase *ticketUseCase.AtualizarTicketUseCase,
	atualizarStatusUseCase *ticketUseCase.AtualizarStatusUseCase,
	adicionarObservacaoUseCase *ticketUseCase.AdicionarObservacaoUseCase,
) *TicketHandler {
	return &TicketHandler{
		criarTicketUseCase:         criarTicketUseCase,
		buscarTicketUseCase:        buscarTicketUseCase,
		listarTicketsUseCase:       listarTicketsUseCase,
		atualizarTicketUseCase:     atualizarTicketUseCase,
		atualizarStatusUseCase:     atualizarStatusUseCase,
		adicionarObservacaoUseCase: adicionarObservacaoUseCase,
	}
}

// Request e Response para criar um novo ticket
type CriarTicketRequest struct {
	Titulo       string                    `json:"titulo"`
	Descricao    string                    `json:"descricao"`
	Categoria    ticketDomain.Categoria    `json:"categoria"`
	Subcategoria ticketDomain.Subcategoria `json:"subcategoria"`
	AbertoPor    string                    `json:"aberto_por"`
	Urgencia     int                       `json:"urgencia"`
	Gravidade    int                       `json:"gravidade"`
	// campos opcionais
	Merchant    string `json:"merchant,omitempty"`
	NoxID       string `json:"nox_id,omitempty"`
	CPF         string `json:"cpf,omitempty"`
	Plataforma  string `json:"plataforma,omitempty"`
	Contato     string `json:"contato,omitempty"`
	Responsavel string `json:"responsavel,omitempty"`
}

type CriarTicketResponse struct {
	ID           string                    `json:"id"`
	Status       ticketDomain.Status       `json:"status"`
	DataAbertura string                    `json:"data_abertura"`
	AbertoPor    string                    `json:"aberto_por"`
	Titulo       string                    `json:"titulo"`
	Categoria    ticketDomain.Categoria    `json:"categoria"`
	Subcategoria ticketDomain.Subcategoria `json:"subcategoria"`
}

// Criar é o handler para criar um novo ticket
func (h *TicketHandler) Criar(w http.ResponseWriter, r *http.Request) {
	// ler o JSON da requisição
	var req CriarTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// converter request para input do use case
	input := ticketUseCase.CriarTicketInput{
		Titulo:       req.Titulo,
		Descricao:    req.Descricao,
		Categoria:    req.Categoria,
		Subcategoria: req.Subcategoria,
		AbertoPor:    req.AbertoPor,
		Urgencia:     req.Urgencia,
		Gravidade:    req.Gravidade,
		Merchant:     req.Merchant,
		NoxID:        req.NoxID,
		CPF:          req.CPF,
		Plataforma:   req.Plataforma,
		Contato:      req.Contato,
		Responsavel:  req.Responsavel,
	}

	// execute o use case
	output, err := h.criarTicketUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// converter output para response
	resp := CriarTicketResponse{
		ID:           output.ID,
		Status:       output.Status,
		DataAbertura: output.DataAbertura,
		AbertoPor:    output.AbertoPor,
		Titulo:       output.Titulo,
		Categoria:    output.Categoria,
		Subcategoria: output.Subcategoria,
	}

	// enviar a resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Response para buscar um ticket
type BuscarTicketResponse struct {
	ID           string                    `json:"id"`
	Status       ticketDomain.Status       `json:"status"`
	DataAbertura string                    `json:"data_abertura"`
	AbertoPor    string                    `json:"aberto_por"`
	Titulo       string                    `json:"titulo"`
	Descricao    string                    `json:"descricao"`
	Categoria    ticketDomain.Categoria    `json:"categoria"`
	Subcategoria ticketDomain.Subcategoria `json:"subcategoria"`
	Urgencia     int                       `json:"urgencia"`
	Gravidade    int                       `json:"gravidade"`
	Merchant     *string                   `json:"merchant,omitempty"`
	NoxID        *string                   `json:"nox_id,omitempty"`
	CPF          *string                   `json:"cpf,omitempty"`
	Plataforma   *string                   `json:"plataforma,omitempty"`
	Contato      *string                   `json:"contato,omitempty"`
	Responsavel  *string                   `json:"responsavel,omitempty"`
	Observacoes  []ObservacaoResponse      `json:"observacoes,omitempty"`
	Modificacoes []ModificacaoResponse     `json:"modificacoes,omitempty"`
}

type ObservacaoResponse struct {
	ID          string `json:"id"`
	UsuarioID   string `json:"usuario_id"`
	Descricao   string `json:"descricao"`
	DataCriacao string `json:"data_criacao"`
}

type ModificacaoResponse struct {
	ID              string `json:"id"`
	UsuarioID       string `json:"usuario_id"`
	CampoModificado string `json:"campo_modificado"`
	ValorAnterior   string `json:"valor_anterior"`
	ValorNovo       string `json:"valor_novo"`
	DataModificacao string `json:"data_modificacao"`
}

// Buscar é o handler para buscar um ticket por ID
func (h *TicketHandler) Buscar(w http.ResponseWriter, r *http.Request) {
	// pegar o ID da URL
	id := chi.URLParam(r, "id")

	// executar o use case
	output, err := h.buscarTicketUseCase.Execute(ticketUseCase.BuscarTicketInput{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// converter output para response
	resp := BuscarTicketResponse{
		ID:           output.ID,
		Status:       output.Status,
		DataAbertura: output.DataAbertura,
		AbertoPor:    output.AbertoPor,
		Titulo:       output.Titulo,
		Descricao:    output.Descricao,
		Categoria:    output.Categoria,
		Subcategoria: output.Subcategoria,
		Urgencia:     output.Urgencia,
		Gravidade:    output.Gravidade,
		Merchant:     output.Merchant,
		NoxID:        output.NoxID,
		CPF:          output.CPF,
		Plataforma:   output.Plataforma,
	}

	// Adicionar campos opcionais apenas se não estiverem vazios
	if output.Contato != "" {
		resp.Contato = &output.Contato
	}
	if output.Responsavel != "" {
		resp.Responsavel = &output.Responsavel
	}

	// converter observações
	for _, obs := range output.Observacoes {
		resp.Observacoes = append(resp.Observacoes, ObservacaoResponse{
			ID:          obs.ID,
			UsuarioID:   obs.UsuarioID,
			Descricao:   obs.Descricao,
			DataCriacao: obs.DataCriacao,
		})
	}

	// converter modificações
	for _, mod := range output.Modificacoes {
		resp.Modificacoes = append(resp.Modificacoes, ModificacaoResponse{
			ID:              mod.ID,
			UsuarioID:       mod.UsuarioID,
			CampoModificado: mod.CampoModificado,
			ValorAnterior:   mod.ValorAnterior,
			ValorNovo:       mod.ValorNovo,
			DataModificacao: mod.DataModificacao,
		})
	}

	// enviar a resposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Request para listar tickets
type ListarTicketsRequest struct {
	Status      *ticketDomain.Status    `json:"status,omitempty"`
	Categoria   *ticketDomain.Categoria `json:"categoria,omitempty"`
	Responsavel *string                 `json:"responsavel,omitempty"`
	Pagina      int                     `json:"pagina"`
	PorPagina   int                     `json:"por_pagina"`
}

// Listar é o handler para listar tickets
func (h *TicketHandler) Listar(w http.ResponseWriter, r *http.Request) {
	// extrair query params
	var req ListarTicketsRequest
	req.Pagina = 1     // valor padrão
	req.PorPagina = 10 // valor padrão

	if page := r.URL.Query().Get("pagina"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			req.Pagina = p
		}
	}

	if limit := r.URL.Query().Get("por_pagina"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			req.PorPagina = l
		}
	}

	// converter request para input do use case
	input := ticketUseCase.ListarTicketsInput{
		Status:         req.Status,
		Categoria:      req.Categoria,
		Responsavel:    req.Responsavel,
		Pagina:         req.Pagina,
		ItensPorPagina: req.PorPagina,
	}

	// executar o use case
	output, err := h.listarTicketsUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// enviar resposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// Request para atualizar ticket
type AtualizarTicketRequest struct {
	Titulo     *string                 `json:"titulo,omitempty"`
	Descricao  *string                 `json:"descricao,omitempty"`
	Categoria  *ticketDomain.Categoria `json:"categoria,omitempty"`
	Urgencia   *int                    `json:"urgencia,omitempty"`
	Gravidade  *int                    `json:"gravidade,omitempty"`
	Merchant   *string                 `json:"merchant,omitempty"`
	NoxID      *string                 `json:"nox_id,omitempty"`
	CPF        *string                 `json:"cpf,omitempty"`
	Plataforma *string                 `json:"plataforma,omitempty"`
	Contato    *string                 `json:"contato,omitempty"`
	UsuarioID  string                  `json:"usuario_id"`
}

// Atualizar é o handler para atualizar um ticket
func (h *TicketHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	// pegar o ID da URL
	id := chi.URLParam(r, "id")

	// ler o JSON da requisição
	var req AtualizarTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// converter request para input do use case
	input := ticketUseCase.AtualizarTicketInput{
		ID:         id,
		Titulo:     req.Titulo,
		Descricao:  req.Descricao,
		Categoria:  req.Categoria,
		Urgencia:   req.Urgencia,
		Gravidade:  req.Gravidade,
		Merchant:   req.Merchant,
		NoxID:      req.NoxID,
		CPF:        req.CPF,
		Plataforma: req.Plataforma,
		Contato:    req.Contato,
		UsuarioID:  req.UsuarioID,
	}

	// executar o use case
	output, err := h.atualizarTicketUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// enviar resposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// Request para atualizar status
type AtualizarStatusRequest struct {
	Status      ticketDomain.Status `json:"status"`
	UsuarioID   string              `json:"usuario_id"`
	Responsavel *string             `json:"responsavel,omitempty"`
}

// AtualizarStatus é o handler para atualizar o status de um ticket
func (h *TicketHandler) AtualizarStatus(w http.ResponseWriter, r *http.Request) {
	// pegar o ID da URL
	id := chi.URLParam(r, "id")

	// ler o JSON da requisição
	var req AtualizarStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// converter request para input do use case
	input := ticketUseCase.AtualizarStatusTicketInput{
		ID:          id,
		Status:      req.Status,
		UsuarioID:   req.UsuarioID,
		Responsavel: req.Responsavel,
	}

	// executar o use case
	output, err := h.atualizarStatusUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// enviar resposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// Request para adicionar observação
type AdicionarObservacaoRequest struct {
	Descricao string `json:"descricao"`
	UsuarioID string `json:"usuario_id"`
}

// AdicionarObservacao é o handler para adicionar uma observação a um ticket
func (h *TicketHandler) AdicionarObservacao(w http.ResponseWriter, r *http.Request) {
	// pegar o ID da URL
	id := chi.URLParam(r, "id")

	// ler o JSON da requisição
	var req AdicionarObservacaoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// converter request para input do use case
	input := ticketUseCase.AdicionarObservacaoInput{
		ID:        id,
		Descricao: req.Descricao,
		UsuarioID: req.UsuarioID,
	}

	// executar o use case
	output, err := h.adicionarObservacaoUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// enviar resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}
