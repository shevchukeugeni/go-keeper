package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	uuid "github.com/satori/go.uuid"

	"keeper-project/internal/auth"
	"keeper-project/types"
)

func (ro *router) createCredentials(w http.ResponseWriter, r *http.Request) {
	var req types.CreateCredentialsRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Unable to decode json: "+err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	creds, err := req.Validate()
	if err != nil {
		http.Error(w, "incorrect data: "+err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.NewV4().String()

	err = ro.credsRepo.Create(r.Context(), userID, id, creds)
	if err != nil {
		http.Error(w, "failed to create: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}

func (ro *router) getCredentials(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "empty site", http.StatusBadRequest)
		return
	}

	note, err := ro.credsRepo.Get(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "no such record", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get from db: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		http.Error(w, "Can't marshal data: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ro *router) getSites(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	sitesList, err := ro.credsRepo.GetKeysList(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get from db: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(sitesList) == 0 {
		http.Error(w, "no credentials", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(sitesList)
	if err != nil {
		http.Error(w, "Can't marshal data: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ro *router) updateCredentials(w http.ResponseWriter, r *http.Request) {
	var req types.UpdateCredentialsRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Unable to decode json: "+err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	creds, err := req.Validate()
	if err != nil {
		http.Error(w, "incorrect data: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = ro.credsRepo.Update(r.Context(), userID, req.ID, creds)
	if err != nil {
		http.Error(w, "failed to update: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (ro *router) deleteCredentials(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "empty key", http.StatusBadRequest)
		return
	}

	err = ro.credsRepo.Delete(r.Context(), userID, id)
	if err != nil {
		http.Error(w, "failed to delete from db: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
