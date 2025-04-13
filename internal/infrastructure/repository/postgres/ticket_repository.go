// internal/infrastructure/repository/postgres/ticket_repository.go

package postgres

import (
	"database/sql"
)

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}
