package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"nox_tickets/internal/application/usecases/ticket"
	dbpostgres "nox_tickets/internal/infrastructure/database/postgres"
	repopostgres "nox_tickets/internal/infrastructure/repository/postgres"
	"nox_tickets/internal/interfaces/http/handler"
	"nox_tickets/internal/interfaces/http/router"
)

type Server struct {
	server *http.Server
}

// NewServer cria uma nova instancia do servidor HTTP
func NewServer(port string) *Server {
	// 1. criar a conexão com o banco
	config := dbpostgres.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "nox_user",
		Password: "nox_password",
		DBName:   "nox_tickets",
		SSLMode:  "disable",
	}

	db, err := dbpostgres.NewConnection(config)
	if err != nil {
		// por enquanto vamos usar panic, depois melhoramos tratamento de erro
		panic(fmt.Sprintf("Erro ao criar conexão com o banco de dados: %v", err))
	}

	// 2. criar o repositório do pacote repository/postgres
	ticketRepo := repopostgres.NewTicketRepository(db)

	// 3. criar os use cases
	criarTicketUseCase := ticket.NewCriarTicketUseCase(ticketRepo)
	buscarTicketUseCase := ticket.NewBuscarTicketUseCase(ticketRepo)
	listarTicketsUseCase := ticket.NewListarTicketsUseCase(ticketRepo)
	atualizarTicketUseCase := ticket.NewAtualizarTicketUseCase(ticketRepo)
	atualizarStatusUseCase := ticket.NewAtualizarStatusUseCase(ticketRepo)
	adicionarObservacaoUseCase := ticket.NewAdicionarObservacaoUseCase(ticketRepo)

	// 4. criar os handlers
	ticketHandler := handler.NewTicketHandler(
		criarTicketUseCase,
		buscarTicketUseCase,
		listarTicketsUseCase,
		atualizarTicketUseCase,
		atualizarStatusUseCase,
		adicionarObservacaoUseCase,
	)

	// 5. criar o router com os handlers
	r := router.NewRouter(ticketHandler)

	// 6. criar o servidor HTTP
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		server: srv,
	}
}

// Start inicia o servidor HTTP
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown desliga o servidor graciosamente
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
