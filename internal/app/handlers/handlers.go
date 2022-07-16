package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/polosaty/go-contracts-crud/internal/app/storage"
)

type mainHandler struct {
	chiMux     *chi.Mux
	repository storage.Repository
}

func NewMainHandler(repository storage.Repository) *chi.Mux {

	h := &mainHandler{chiMux: chi.NewMux(), repository: repository}
	h.chiMux.Use(middleware.RequestID)
	h.chiMux.Use(middleware.RealIP)
	h.chiMux.Use(middleware.Logger)
	h.chiMux.Use(middleware.Recoverer)

	h.chiMux.Route("/api", func(r chi.Router) {
		r.Route("/company", func(r chi.Router) {
			r.Post("/", h.createCompany())
			r.Get("/{id}", h.readCompany())
			r.Get("/", h.readCompanyList())
			r.Post("/{id}", h.updateCompany())
			r.Delete("/", h.deleteCompany())
		})

	})

	return h.chiMux
}
