package server

import (
	"github.com/polosaty/go-contracts-crud/internal/app/handlers"
	"github.com/polosaty/go-contracts-crud/internal/app/storage"
	"net/http"
)

func Serve(addr string, db storage.Repository) error {
	handler := handlers.NewMainHandler(db)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return server.ListenAndServe()
}
