package router

import (
	"net/http"

	"nox_tickets/internal/interfaces/http/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// newRouter cria e configura um novo router
func NewRouter(ticketHandler *handler.TicketHandler) *chi.Mux {
	r := chi.NewRouter()

	// adiciona middleware de loggind
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Rota básica para healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// rotas de tickets
	r.Route("/tickets", func(r chi.Router) {
		// POST /tickets - criar novo ticket
		r.Post("/", ticketHandler.Criar)

		// GET /tickets - listar tickets
		r.Get("/", ticketHandler.Listar)

		// Rotas que precisam do ID do ticket
		r.Route("/{id}", func(r chi.Router) {
			// GET /tickets/{id} - obter ticket por ID
			r.Get("/", ticketHandler.Buscar)

			// PUT /tickets/{id} - atualizar ticket
			r.Put("/", ticketHandler.Atualizar)

			// PATCH /tickets/{id}/status - atualizar status do ticket
			r.Patch("/status", ticketHandler.AtualizarStatus)

			// POST /tickets/{id}/observacoes - adicionar observações ao ticket
			r.Post("/observacoes", ticketHandler.AdicionarObservacao)
		})
	})

	return r
}
