package postgres

import (
	"testing"

	"nox_tickets/internal/domain/ticket"
	"nox_tickets/internal/infrastructure/database/postgres"
)

// Função auxiliar para criar uma conexão de teste
func setupTestDB(t *testing.T) *TicketRepository {
	config := postgres.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "nox_user",
		Password: "nox_password",
		DBName:   "nox_tickets",
		SSLMode:  "disable",
	}

	db, err := postgres.NewConnection(config)
	if err != nil {
		t.Fatalf("Erro ao conectar ao banco de teste: %v", err)
	}

	return NewTicketRepository(db)
}

// Função auxiliar para criar um ticket de teste
func createTestTicket() *ticket.Ticket {
	ticket, _ := ticket.NovoTicket(
		"Ticket de Teste",
		"Descrição do ticket de teste",
		ticket.CategoriaFinanceiro,
		ticket.SubcategoriaBug,
		"usuario_teste",
	)
	return ticket
}

// Teste do método Create
func TestTicketRepository_Create(t *testing.T) {
	repo := setupTestDB(t)

	// Cria um ticket de teste
	testTicket := createTestTicket()

	// Tenta criar o ticket
	err := repo.Create(testTicket)
	if err != nil {
		t.Errorf("Erro ao criar ticket: %v", err)
	}

	// Verifica se o ticket foi criado buscando ele
	saved, err := repo.GetByID(testTicket.ID)
	if err != nil {
		t.Errorf("Erro ao buscar ticket criado: %v", err)
	}

	// Verifica se os dados estão corretos
	if saved.Titulo != testTicket.Titulo {
		t.Errorf("Título diferente: esperado %s, recebido %s", testTicket.Titulo, saved.Titulo)
	}
	if saved.Descricao != testTicket.Descricao {
		t.Errorf("Descrição diferente: esperado %s, recebido %s", testTicket.Descricao, saved.Descricao)
	}
}

// Teste do método GetByID
func TestTicketRepository_GetByID(t *testing.T) {
	repo := setupTestDB(t)

	// Cria um ticket para testar
	testTicket := createTestTicket()
	err := repo.Create(testTicket)
	if err != nil {
		t.Fatalf("Erro ao criar ticket para teste: %v", err)
	}

	// Testa buscar o ticket
	found, err := repo.GetByID(testTicket.ID)
	if err != nil {
		t.Errorf("Erro ao buscar ticket: %v", err)
	}
	if found.ID != testTicket.ID {
		t.Errorf("ID diferente: esperado %s, recebido %s", testTicket.ID, found.ID)
	}

	// Testa buscar ticket que não existe
	_, err = repo.GetByID("id_que_nao_existe")
	if err == nil {
		t.Error("Esperava erro ao buscar ticket inexistente")
	}
}

// Teste do método List
func TestTicketRepository_List(t *testing.T) {
	repo := setupTestDB(t)

	// Cria alguns tickets para testar
	ticket1 := createTestTicket()
	ticket2 := createTestTicket()
	ticket2.Status = ticket.StatusEmCurso

	err := repo.Create(ticket1)
	if err != nil {
		t.Fatalf("Erro ao criar ticket1: %v", err)
	}
	err = repo.Create(ticket2)
	if err != nil {
		t.Fatalf("Erro ao criar ticket2: %v", err)
	}

	// Testa listar sem filtros
	tickets, err := repo.List(ticket.TicketFiltros{})
	if err != nil {
		t.Errorf("Erro ao listar tickets: %v", err)
	}
	if len(tickets) < 2 {
		t.Error("Esperava encontrar pelo menos 2 tickets")
	}

	// Testa listar com filtro de status
	filtroStatus := ticket.TicketFiltros{
		Status: []ticket.Status{ticket.StatusEmCurso},
	}
	ticketsFiltrados, err := repo.List(filtroStatus)
	if err != nil {
		t.Errorf("Erro ao listar tickets com filtro: %v", err)
	}
	for _, tick := range ticketsFiltrados {
		if tick.Status != ticket.StatusEmCurso {
			t.Error("Encontrou ticket com status diferente do filtrado")
		}
	}
}

// Teste do método Update
func TestTicketRepository_Update(t *testing.T) {
	repo := setupTestDB(t)

	// Cria um ticket para testar
	testTicket := createTestTicket()
	err := repo.Create(testTicket)
	if err != nil {
		t.Fatalf("Erro ao criar ticket para teste: %v", err)
	}

	// Modifica o ticket
	novoTitulo := "Título Atualizado"
	testTicket.Titulo = novoTitulo

	// Testa atualizar
	err = repo.Update(testTicket)
	if err != nil {
		t.Errorf("Erro ao atualizar ticket: %v", err)
	}

	// Verifica se atualizou
	updated, err := repo.GetByID(testTicket.ID)
	if err != nil {
		t.Errorf("Erro ao buscar ticket atualizado: %v", err)
	}
	if updated.Titulo != novoTitulo {
		t.Errorf("Título não foi atualizado: esperado %s, recebido %s", novoTitulo, updated.Titulo)
	}
}

// Teste do método Delete
func TestTicketRepository_Delete(t *testing.T) {
	repo := setupTestDB(t)

	// Cria um ticket para testar
	testTicket := createTestTicket()
	err := repo.Create(testTicket)
	if err != nil {
		t.Fatalf("Erro ao criar ticket para teste: %v", err)
	}

	// Testa deletar
	err = repo.Delete(testTicket.ID)
	if err != nil {
		t.Errorf("Erro ao deletar ticket: %v", err)
	}

	// Verifica se foi deletado
	_, err = repo.GetByID(testTicket.ID)
	if err == nil {
		t.Error("Esperava erro ao buscar ticket deletado")
	}

	// Testa deletar ticket que não existe
	err = repo.Delete("id_que_nao_existe")
	if err == nil {
		t.Error("Esperava erro ao deletar ticket inexistente")
	}
}
