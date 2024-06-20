package server

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/docker/go-units"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"

	"keeper-project/internal/auth"
	"keeper-project/internal/store"
	"keeper-project/types"
)

type router struct {
	logger      *zap.Logger
	userRepo    store.User
	notesRepo   store.Secrets[types.Note]
	credsRepo   store.Secrets[types.Credentials]
	cardsRepo   store.Secrets[types.CardInfo]
	fileService store.FileService
}

func SetupRouter(logger *zap.Logger,
	user store.User,
	notesRepo store.Secrets[types.Note],
	credsRepo store.Secrets[types.Credentials],
	cardsRepo store.Secrets[types.CardInfo],
	fileService store.FileService) http.Handler {
	ro := &router{
		logger:      logger,
		userRepo:    user,
		notesRepo:   notesRepo,
		cardsRepo:   cardsRepo,
		credsRepo:   credsRepo,
		fileService: fileService,
	}
	return ro.Handler()
}

func (ro *router) Handler() http.Handler {
	rtr := chi.NewRouter()
	rtr.Use(middleware.Logger)
	rtr.Get("/clients/*", func(w http.ResponseWriter, r *http.Request) {
		workDir, _ := os.Getwd()
		filesDir := http.Dir(workDir)
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
		fs.ServeHTTP(w, r)
	})
	rtr.Post("/api/user/register", ro.register)
	rtr.Post("/api/user/login", ro.auth)
	rtr.Route("/api/secret", func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Use(middleware.RequestSize(32 * units.MiB))
		r.Post("/text", ro.createNote)
		r.Get("/text/{id}", ro.getNote)
		r.Get("/texts", ro.getNotesKeys)
		r.Put("/text", ro.updateNote)
		r.Delete("/text/{id}", ro.deleteNote)
		r.Post("/card", ro.createCard)
		r.Get("/card/{id}", ro.getCardInfo)
		r.Get("/cards", ro.getCardsList)
		r.Put("/card", ro.updateCard)
		r.Delete("/card/{id}", ro.deleteCard)
		r.Post("/cred", ro.createCredentials)
		r.Get("/cred/{id}", ro.getCredentials)
		r.Get("/creds", ro.getSites)
		r.Put("/cred", ro.updateCredentials)
		r.Delete("/cred/{id}", ro.deleteCredentials)
		r.Post("/file", ro.createFile)
		r.Get("/file/{id}", ro.getFile)
		r.Get("/files", ro.getFiles)
		r.Delete("/file/{id}", ro.deleteFile)
	})
	return rtr
}

func (ro *router) register(w http.ResponseWriter, r *http.Request) {
	var req types.UserLoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Unable to decode json: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Missing login or password.", http.StatusBadRequest)
		return
	}

	usr := req.User().ToDB()

	err = ro.userRepo.CreateUser(r.Context(), usr)
	if err != nil {
		if errors.Is(err, types.ErrUserAlreadyExists) {
			http.Error(w, "Unable to create user: "+err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "Unable to create user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	tokenString, err := auth.GenerateToken(usr.ID)
	if err != nil {
		http.Error(w, "Unable to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	w.WriteHeader(http.StatusOK)
}

func (ro *router) auth(w http.ResponseWriter, r *http.Request) {
	var req types.UserLoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Unable to decode json: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Missing login or password.", http.StatusBadRequest)
		return
	}

	usr, err := ro.userRepo.GetByLogin(r.Context(), req.Login)
	if err != nil {
		http.Error(w, "Unable to find user: "+err.Error(), http.StatusBadRequest)
		return
	}

	h := sha3.New512()
	h.Write([]byte(req.Password))

	if base64.StdEncoding.EncodeToString(h.Sum(nil)) != usr.Password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenString, err := auth.GenerateToken(usr.ID)
	if err != nil {
		http.Error(w, "Unable to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	w.WriteHeader(http.StatusOK)
}
