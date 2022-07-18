package handlers

import (
	"encoding/json"
	"errors"
	"github.com/polosaty/go-contracts-crud/internal/app/storage"
	"log"
	"net/http"
)

func (h *mainHandler) createBuy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var buy storage.Buy
		if err := json.NewDecoder(r.Body).Decode(&buy); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := h.repository.CreateBuy(ctx, &buy)
		if err != nil {
			log.Println("create buy error", err)
			if errors.Is(err, storage.ErrDuplicateBuy) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buy.ID = &id
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(buy)
		if err != nil {
			log.Println("marshal response error: ", err)
		}
	}
}
