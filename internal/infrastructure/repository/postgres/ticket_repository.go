// internal/infrastructure/repository/postgres/ticket_repository.go

package postgres

import (
	"database/sql"
	"fmt"
	"nox_tickets/internal/domain/ticket"
	"strings"
	"time"
)

// Função auxiliar para converter o formato de tempo do PostgreSQL para time.Duration
func parsePostgresInterval(interval string) (time.Duration, error) {
	if interval == "" || interval == "00:00:00" {
		return 0, nil
	}

	// Converte o formato "HH:MM:SS" para nanosegundos
	var hours, minutes, seconds int
	_, err := fmt.Sscanf(interval, "%d:%d:%d", &hours, &minutes, &seconds)
	if err != nil {
		return 0, err
	}

	duration := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second

	return duration, nil
}

// Função auxiliar para formatar duração para PostgreSQL
func formatDurationForPostgres(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
}

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

// criar um novo ticket
func (r *TicketRepository) Create(ticket *ticket.Ticket) error {
	// inicia uma transação
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // garante que a transação será revertida em caso de erro

	// insere o ticket principal
	_, err = tx.Exec(
		`INSERT INTO tickets (
		id, titulo, merchant, nox_id, cpf, status, categoria,
		subcategoria, descricao, urgencia, gravidade,
		aberto_por, responsavel, contato, plataforma,
		data_abertura, data_inicio, data_conclusao,
		duracao_total, duracao_execucao
		) VALUES (
		 $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19::interval, $20::interval
		)`,
		ticket.ID, ticket.Titulo, ticket.Merchant, ticket.NoxID, ticket.CPF, ticket.Status, ticket.Categoria,
		ticket.Subcategoria, ticket.Descricao, ticket.Urgencia, ticket.Gravidade,
		ticket.AbertoPor, ticket.Responsavel, ticket.Contato, ticket.Plataforma,
		ticket.DataAbertura, ticket.DataInicio, ticket.DataConclusao,
		formatDurationForPostgres(ticket.DuracaoTotal), formatDurationForPostgres(ticket.DuracaoExecucao),
	)
	if err != nil {
		return err
	}

	// confirma a transação
	return tx.Commit()
}

// buscar ticket por id
func (r *TicketRepository) GetByID(id string) (*ticket.Ticket, error) {
	t := &ticket.Ticket{}

	// Variáveis temporárias para armazenar as durações como strings
	var duracaoTotalStr, duracaoExecucaoStr string

	// Busca o ticket principal
	err := r.db.QueryRow(`
        SELECT 
            id, titulo, merchant, nox_id, cpf, status, categoria,
            subcategoria, descricao, urgencia, gravidade,
            aberto_por, responsavel, contato, plataforma,
            data_abertura, data_inicio, data_conclusao,
            duracao_total::text, duracao_execucao::text
        FROM tickets 
        WHERE id = $1
    `, id).Scan(
		&t.ID, &t.Titulo, &t.Merchant, &t.NoxID,
		&t.CPF, &t.Status, &t.Categoria, &t.Subcategoria,
		&t.Descricao, &t.Urgencia, &t.Gravidade,
		&t.AbertoPor, &t.Responsavel, &t.Contato,
		&t.Plataforma, &t.DataAbertura, &t.DataInicio,
		&t.DataConclusao, &duracaoTotalStr, &duracaoExecucaoStr,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ticket não encontrado")
	}
	if err != nil {
		return nil, err
	}

	// Converte as durações de string para time.Duration
	if duracaoTotalStr != "" {
		duration, err := parsePostgresInterval(duracaoTotalStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter duracao_total: %v", err)
		}
		t.DuracaoTotal = duration
	}

	if duracaoExecucaoStr != "" {
		duration, err := parsePostgresInterval(duracaoExecucaoStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter duracao_execucao: %v", err)
		}
		t.DuracaoExecucao = duration
	}

	// Busca as observações
	rows, err := r.db.Query(`
		SELECT id, usuario_id, descricao, data_criacao
		FROM observacoes
		WHERE ticket_id = $1
		ORDER BY data_criacao
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var obs ticket.Observacao
		err := rows.Scan(&obs.ID, &obs.UsuarioID, &obs.Descricao, &obs.DataCriacao)
		if err != nil {
			return nil, err
		}
		obs.TicketID = id
		t.Observacoes = append(t.Observacoes, obs)
	}

	// Busca as modificações
	rows, err = r.db.Query(`
		SELECT id, usuario_id, campo_modificado, valor_anterior, valor_novo, data_modificacao
		FROM modificacoes
		WHERE ticket_id = $1
		ORDER BY data_modificacao
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var mod ticket.Modificacao
		err := rows.Scan(
			&mod.ID, &mod.UsuarioID, &mod.CampoModificado,
			&mod.ValorAnterior, &mod.ValorNovo, &mod.DataModificacao,
		)
		if err != nil {
			return nil, err
		}
		mod.TicketID = id
		t.Modificacoes = append(t.Modificacoes, mod)
	}

	return t, nil
}

// listar tickets
func (r *TicketRepository) List(filtros ticket.TicketFiltros) ([]*ticket.Ticket, error) {
	// slice para guardar as condições do WHERE
	where := []string{}
	// Slice para guardar os argumentos
	args := []interface{}{}
	// contador para os placeholders
	argCount := 1

	// adiciona filtros se existirem
	if len(filtros.Status) > 0 {
		statusStrings := make([]string, len(filtros.Status))
		for i, s := range filtros.Status {
			statusStrings[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, s)
			argCount++
		}
		where = append(where, fmt.Sprintf("status = ANY(ARRAY[%s])", strings.Join(statusStrings, ",")))
	}

	if len(filtros.Categoria) > 0 {
		categStrings := make([]string, len(filtros.Categoria))
		for i, c := range filtros.Categoria {
			categStrings[i] = fmt.Sprintf("categoria ILIKE $%d", argCount)
			args = append(args, string(c))
			argCount++
		}
		where = append(where, fmt.Sprintf("(%s)", strings.Join(categStrings, " OR ")))
	}

	if filtros.Responsavel != "" {
		where = append(where, fmt.Sprintf("responsavel = $%d", argCount))
		args = append(args, filtros.Responsavel)
		argCount++
	}

	if filtros.AbertoPor != "" {
		where = append(where, fmt.Sprintf("aberto_por = $%d", argCount))
		args = append(args, filtros.AbertoPor)
		argCount++
	}

	// Construir a query
	query := `
	    SELECT
			id, titulo, merchant, nox_id, cpf, status, categoria,
			subcategoria, descricao, urgencia, gravidade,
			aberto_por, responsavel, contato, plataforma,
			data_abertura, data_inicio, data_conclusao,
			duracao_total, duracao_execucao
		FROM tickets`

	// adiciona as condições WHERE se existirem
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	// executa a query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// slice para guardar os resultados
	tickets := []*ticket.Ticket{}

	// itera sobre os resultados
	for rows.Next() {
		t := &ticket.Ticket{}
		var duracaoTotalStr, duracaoExecucaoStr string

		err := rows.Scan(
			&t.ID, &t.Titulo, &t.Merchant, &t.NoxID, &t.CPF, &t.Status, &t.Categoria,
			&t.Subcategoria, &t.Descricao, &t.Urgencia, &t.Gravidade,
			&t.AbertoPor, &t.Responsavel, &t.Contato, &t.Plataforma,
			&t.DataAbertura, &t.DataInicio, &t.DataConclusao,
			&duracaoTotalStr, &duracaoExecucaoStr,
		)
		if err != nil {
			return nil, err
		}

		// Converte as durações
		t.DuracaoTotal, err = parsePostgresInterval(duracaoTotalStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter duracao_total: %v", err)
		}

		t.DuracaoExecucao, err = parsePostgresInterval(duracaoExecucaoStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter duracao_execucao: %v", err)
		}

		tickets = append(tickets, t)
	}

	// verifica se houve erros durante a iteração
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// retorna os tickets encontrados
	return tickets, nil
}

func (r *TicketRepository) Update(ticket *ticket.Ticket) error {
	// inicia uma transação
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // garante que a transação será revertida em caso de erro

	// atualiza o ticket principal
	_, err = tx.Exec(
		`UPDATE tickets SET
		titulo = $1,
		merchant = $2,
		nox_id = $3,
		cpf = $4,
		status = $5,
		categoria = $6,
		subcategoria = $7,
		descricao = $8,
		urgencia = $9,
		gravidade = $10,
		aberto_por = $11,
		responsavel = $12,
		contato = $13,
		plataforma = $14,
		data_abertura = $15,
		data_inicio = $16,
		data_conclusao = $17,
		duracao_total = $18::interval,
		duracao_execucao = $19::interval
		WHERE id = $20
		`,
		ticket.Titulo, ticket.Merchant, ticket.NoxID, ticket.CPF, ticket.Status, ticket.Categoria,
		ticket.Subcategoria, ticket.Descricao, ticket.Urgencia, ticket.Gravidade,
		ticket.AbertoPor, ticket.Responsavel, ticket.Contato, ticket.Plataforma,
		ticket.DataAbertura, ticket.DataInicio, ticket.DataConclusao,
		formatDurationForPostgres(ticket.DuracaoTotal), formatDurationForPostgres(ticket.DuracaoExecucao),
		ticket.ID,
	)
	if err != nil {
		return err
	}

	// Salva as novas observações
	for _, obs := range ticket.Observacoes {
		// Verifica se a observação já existe
		var exists bool
		err := tx.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM observacoes WHERE id = $1)",
			obs.ID,
		).Scan(&exists)
		if err != nil {
			return err
		}

		// Se não existe, insere
		if !exists {
			_, err = tx.Exec(
				`INSERT INTO observacoes (id, ticket_id, usuario_id, descricao, data_criacao)
				VALUES ($1, $2, $3, $4, $5)`,
				obs.ID, ticket.ID, obs.UsuarioID, obs.Descricao, obs.DataCriacao,
			)
			if err != nil {
				return err
			}
		}
	}

	// Salva as novas modificações
	for _, mod := range ticket.Modificacoes {
		// Verifica se a modificação já existe
		var exists bool
		err := tx.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM modificacoes WHERE id = $1)",
			mod.ID,
		).Scan(&exists)
		if err != nil {
			return err
		}

		// Se não existe, insere
		if !exists {
			_, err = tx.Exec(
				`INSERT INTO modificacoes (id, ticket_id, usuario_id, campo_modificado, valor_anterior, valor_novo, data_modificacao)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				mod.ID, ticket.ID, mod.UsuarioID, mod.CampoModificado,
				mod.ValorAnterior, mod.ValorNovo, mod.DataModificacao,
			)
			if err != nil {
				return err
			}
		}
	}

	// confirma a transação
	return tx.Commit()
}

func (r *TicketRepository) Delete(id string) error {
	// inicia uma transação
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // garante que a transação será revertida em caso de erro

	// deleta primeiro as observacoes e modificacoes (por causa das chaves estrangeiras)
	_, err = tx.Exec(
		`DELETE FROM observacoes WHERE ticket_id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM modificacoes WHERE ticket_id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	// deleta o ticket principal
	result, err := tx.Exec(
		`DELETE FROM tickets WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	// verifica se o ticket existia
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("ticket não encontrado")
	}

	// confirma a transação
	return tx.Commit()
}

// Adicionar observação
func (r *TicketRepository) AdicionarObservacao(ticketID string, observacao *ticket.Observacao) error {
	_, err := r.db.Exec(
		`INSERT INTO observacoes (ticket_id, usuario_id, descricao, data_criacao)
		 VALUES ($1, $2, $3, $4)`,
		ticketID, observacao.UsuarioID, observacao.Descricao, observacao.DataCriacao,
	)
	return err
}

// Listar observações
func (r *TicketRepository) ListarObservacoes(ticketID string) ([]*ticket.Observacao, error) {
	rows, err := r.db.Query(
		`SELECT usuario_id, descricao, data_criacao 
		 FROM observacoes 
		 WHERE ticket_id = $1 
		 ORDER BY data_criacao`,
		ticketID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var observacoes []*ticket.Observacao
	for rows.Next() {
		obs := &ticket.Observacao{}
		err := rows.Scan(&obs.UsuarioID, &obs.Descricao, &obs.DataCriacao)
		if err != nil {
			return nil, err
		}
		observacoes = append(observacoes, obs)
	}
	return observacoes, nil
}

// Adicionar modificação
func (r *TicketRepository) AdicionarModificacao(ticketID string, modificacao *ticket.Modificacao) error {
	_, err := r.db.Exec(
		`INSERT INTO modificacoes (ticket_id, usuario_id, campo_modificado, valor_anterior, valor_novo, data_modificacao)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		ticketID, modificacao.UsuarioID, modificacao.CampoModificado,
		modificacao.ValorAnterior, modificacao.ValorNovo, modificacao.DataModificacao,
	)
	return err
}

// Listar modificações
func (r *TicketRepository) ListarModificacoes(ticketID string) ([]*ticket.Modificacao, error) {
	rows, err := r.db.Query(
		`SELECT usuario_id, campo_modificado, valor_anterior, valor_novo, data_modificacao 
		 FROM modificacoes 
		 WHERE ticket_id = $1 
		 ORDER BY data_modificacao`,
		ticketID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modificacoes []*ticket.Modificacao
	for rows.Next() {
		mod := &ticket.Modificacao{}
		err := rows.Scan(
			&mod.UsuarioID, &mod.CampoModificado, &mod.ValorAnterior,
			&mod.ValorNovo, &mod.DataModificacao,
		)
		if err != nil {
			return nil, err
		}
		modificacoes = append(modificacoes, mod)
	}
	return modificacoes, nil
}

// Listar por status
func (r *TicketRepository) ListarPorStatus(status ticket.Status) ([]*ticket.Ticket, error) {
	return r.List(ticket.TicketFiltros{Status: []ticket.Status{status}})
}

// Atualizar status
func (r *TicketRepository) AtualizarStatus(ticketID string, novoStatus ticket.Status, usuarioID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Busca o status atual
	var statusAtual ticket.Status
	err = tx.QueryRow("SELECT status FROM tickets WHERE id = $1", ticketID).Scan(&statusAtual)
	if err != nil {
		return err
	}

	// Atualiza o status
	_, err = tx.Exec("UPDATE tickets SET status = $1 WHERE id = $2", novoStatus, ticketID)
	if err != nil {
		return err
	}

	// Registra a modificação
	modificacao := &ticket.Modificacao{
		UsuarioID:       usuarioID,
		CampoModificado: "status",
		ValorAnterior:   string(statusAtual),
		ValorNovo:       string(novoStatus),
		DataModificacao: time.Now(),
	}

	_, err = tx.Exec(
		`INSERT INTO modificacoes (ticket_id, usuario_id, campo_modificado, valor_anterior, valor_novo, data_modificacao)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		ticketID, modificacao.UsuarioID, modificacao.CampoModificado,
		modificacao.ValorAnterior, modificacao.ValorNovo, modificacao.DataModificacao,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
