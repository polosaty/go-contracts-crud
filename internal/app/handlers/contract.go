package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/polosaty/go-contracts-crud/internal/app/storage"
	"log"
	"net/http"
	"strconv"
)

func (h *mainHandler) createContract() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var contract storage.Contract
		if err := json.NewDecoder(r.Body).Decode(&contract); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := h.repository.CreateContract(ctx, &contract)
		if err != nil {
			log.Println("create contract error", err)
			if errors.Is(err, storage.ErrDuplicateContract) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		contract.ID = &id
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(contract)
		if err != nil {
			log.Println("marshal response error: ", err)
		}
	}
}

func (h *mainHandler) readContract() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idString := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			log.Println("parse contract_id error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		contract, err := h.repository.ReadContract(ctx, id)
		if err != nil {
			if errors.Is(err, storage.ErrContractNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			log.Println("read contract error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(contract)
		if err != nil {
			log.Println("marshal response error: ", err)
		}
	}
}

// readContractList
//  method: get
//  path: /api/contract/
func (h *mainHandler) readContractList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		limitStr := chi.URLParam(r, "limit")
		offsetStr := chi.URLParam(r, "offset")
		var (
			limit  uint64
			offset uint64
			err    error
		)
		if limitStr != "" {
			limit, err = strconv.ParseUint(limitStr, 10, 32)
			if err != nil {
				log.Println("parse limit parameter error", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		if offsetStr != "" {
			offset, err = strconv.ParseUint(offsetStr, 10, 64)
			if err != nil {
				log.Println("parse offset parameter error", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		companies, err := h.repository.ReadContractList(ctx, uint(limit), offset)
		if err != nil {
			log.Println(err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if len(companies) == 0 {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		err = json.NewEncoder(w).Encode(companies)
		if err != nil {
			log.Println("marshal response error: ", err)
		}
	}
}

// updateContract
//  method: post
//  path: /api/contract/{id}
func (h *mainHandler) updateContract() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var contract storage.Contract
		if err := json.NewDecoder(r.Body).Decode(&contract); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		idString := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			http.Error(w, "not valid id provided:"+err.Error(), http.StatusBadRequest)
			return
		}

		err = h.repository.UpdateContract(ctx, id, &contract)
		if err != nil {
			log.Println("create contract error", err)
			if errors.Is(err, storage.ErrDuplicateContract) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

	}
}

// deleteContract
//  method: delete
//  path: /api/contract/{id}
func (h *mainHandler) deleteContract() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idString := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			http.Error(w, "not valid id provided:"+err.Error(), http.StatusBadRequest)
			return
		}

		err = h.repository.DeleteContract(ctx, id)
		if err != nil {
			log.Println("create contract error", err)
			if errors.Is(err, storage.ErrDuplicateContract) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}
