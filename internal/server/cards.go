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

func (ro *router) createCard(w http.ResponseWriter, r *http.Request) {
	var req types.CreateCardRequest

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

	cardInfo, err := req.Validate()
	if err != nil {
		http.Error(w, "incorrect data: "+err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.NewV4().String()

	err = ro.cardsRepo.Create(r.Context(), userID, id, cardInfo)
	if err != nil {
		http.Error(w, "failed to create: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}

func (ro *router) getCardInfo(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	cardID := chi.URLParam(r, "id")

	card, err := ro.cardsRepo.Get(r.Context(), userID, cardID)
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
	err = json.NewEncoder(w).Encode(card)
	if err != nil {
		http.Error(w, "Can't marshal data: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ro *router) getCardsList(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	cardsIDs, err := ro.cardsRepo.GetKeysList(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "no cards info", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get from db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(cardsIDs)
	if err != nil {
		http.Error(w, "Can't marshal data: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ro *router) updateCard(w http.ResponseWriter, r *http.Request) {
	var req types.CreateCardRequest

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

	cardInfo, err := req.Validate()
	if err != nil {
		http.Error(w, "incorrect data: "+err.Error(), http.StatusBadRequest)
		return
	}

	if cardInfo.ID == "" {
		http.Error(w, "incorrect data: empty id", http.StatusBadRequest)
		return
	}

	err = ro.cardsRepo.Update(r.Context(), userID, cardInfo.ID, cardInfo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "nothing to update", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ro *router) deleteCard(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	cardID := chi.URLParam(r, "id")

	err = ro.cardsRepo.Delete(r.Context(), userID, cardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "nothing to delete", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete from db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
