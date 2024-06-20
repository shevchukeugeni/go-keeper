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

func (ro *router) createNote(w http.ResponseWriter, r *http.Request) {
	var req types.CreateNoteRequest

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

	note, err := req.Validate()
	if err != nil {
		http.Error(w, "incorrect data: "+err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.NewV4().String()

	err = ro.notesRepo.Create(r.Context(), userID, id, note)
	if err != nil {
		http.Error(w, "failed to create: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}

func (ro *router) getNote(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	noteID := chi.URLParam(r, "id")

	note, err := ro.notesRepo.Get(r.Context(), userID, noteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "no such record", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get from db: "+err.Error(), http.StatusInternalServerError)
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

func (ro *router) getNotesKeys(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	noteKeys, err := ro.notesRepo.GetKeysList(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "no recorded notes found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get from db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(noteKeys)
	if err != nil {
		http.Error(w, "Can't marshal data: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ro *router) updateNote(w http.ResponseWriter, r *http.Request) {
	var req types.UpdateNoteRequest

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

	note, err := req.Validate()
	if err != nil {
		http.Error(w, "incorrect data: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = ro.notesRepo.Update(r.Context(), userID, req.ID, note)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "nothing to update", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (ro *router) deleteNote(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	noteID := chi.URLParam(r, "id")

	err = ro.notesRepo.Delete(r.Context(), userID, noteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "nothing to delete", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete from db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
