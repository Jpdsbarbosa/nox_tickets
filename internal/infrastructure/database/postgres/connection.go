package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// config armazena as configurações de conexão com o banco
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConnection(config Config) (*sql.DB, error) {
	// monta a string de conexão com o banco
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	// abre a conexão com o banco
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com o banco: %v", err)
	}

	// testa a conexão com o banco
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %v", err)
	}

	return db, nil
}
