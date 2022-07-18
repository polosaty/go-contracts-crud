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
//  method: post
//  path: /api/company/
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
		company.ID = &id
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(company)
		if err != nil {
			log.Println("marshal response error: ", err)
		}
	}
}

// readCompany
//  method: get
//  path: /api/company/{id}
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
//  method: get
//  path: /api/company/
func (h *mainHandler) readCompanyList() http.HandlerFunc {
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

		companies, err := h.repository.ReadCompanyList(ctx, uint(limit), offset)
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

// updateCompany
//  method: post
//  path: /api/company/{id}
func (h *mainHandler) updateCompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var company storage.Company
		if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		idString := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			http.Error(w, "not valid id provided:"+err.Error(), http.StatusBadRequest)
			return
		}

		err = h.repository.UpdateCompany(ctx, id, &company)
		if err != nil {
			log.Println("create company error", err)
			if errors.Is(err, storage.ErrDuplicateCompany) {
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

// deleteCompany
//  method: delete
//  path: /api/company/{id}
func (h *mainHandler) deleteCompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idString := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			http.Error(w, "not valid id provided:"+err.Error(), http.StatusBadRequest)
			return
		}

		err = h.repository.DeleteCompany(ctx, id)
		if err != nil {
			log.Println("create company error", err)
			if errors.Is(err, storage.ErrDuplicateCompany) {
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
