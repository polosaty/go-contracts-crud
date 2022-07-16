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

// createCompany
func (h *mainHandler) createCompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var company storage.Company
		if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := h.repository.CreateCompany(ctx, &company)
		if err != nil {
			log.Println("create company error", err)
			if errors.Is(err, storage.ErrDuplicateCompany) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		company.Id = id
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(company)
		if err != nil {
			log.Println("marshal response error: ", err)
		}
	}
}

// readCompany
func (h *mainHandler) readCompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idString := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			log.Println("parse company_id error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		company, err := h.repository.ReadCompany(ctx, id)
		if err != nil {
			if errors.Is(err, storage.ErrCompanyNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			log.Println("read company error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(company)
		if err != nil {
			log.Println("marshal response error: ", err)
		}
	}
}

// readCompanyList
func (h *mainHandler) readCompanyList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented yet", http.StatusNotFound)
	}
}

// updateCompany
func (h *mainHandler) updateCompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented yet", http.StatusNotFound)
	}
}

// deleteCompany
func (h *mainHandler) deleteCompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented yet", http.StatusNotFound)
	}
}
