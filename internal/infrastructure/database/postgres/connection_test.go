// internal/infrastructure/database/postgres/connection_test.go

package postgres

import (
	"testing"
)

func TestNewConnection(t *testing.T) {
	// configuração do banco de dados
	config := Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "nox_user",
		Password: "nox_password",
		DBName:   "nox_tickets",
		SSLMode:  "disable",
	}

	// NewConnection retorna um erro
	db, err := NewConnection(config)
	if err != nil {
		t.Fatalf("Erro ao criar conexão com o banco: %v", err)
	}
	defer db.Close()

	// testa se aconexão esta ok
	if err := db.Ping(); err != nil {
		t.Fatalf("Erro ao testar conexão com o banco: %v", err)
	}

	t.Log("Conexão com o banco de dados estabelecida com sucesso")
}
