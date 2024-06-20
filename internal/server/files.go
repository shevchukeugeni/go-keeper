package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/go-units"
	"github.com/go-chi/chi/v5"

	"keeper-project/internal/auth"
	"keeper-project/types"
)

func (ro *router) getFile(w http.ResponseWriter, r *http.Request) {
	fileId := chi.URLParam(r, "id")
	if fileId == "" {
		http.Error(w, "empty fileID", http.StatusBadRequest)
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	f, err := ro.fileService.GetFile(r.Context(), userID, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Meta", f.Metadata)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", f.Name))
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	w.Write(f.Bytes)
}

func (ro *router) getFiles(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	list, err := ro.fileService.GetFilesList(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(list) == 0 {
		http.Error(w, "no files", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		http.Error(w, "Can't marshal data: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ro *router) createFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "form/json")

	err := r.ParseMultipartForm(32 * units.MiB)
	if err != nil {
		http.Error(w, "file size limit exceeded", http.StatusBadRequest)
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	files, ok := r.MultipartForm.File["file"]
	if !ok || len(files) == 0 {
		http.Error(w, "file required", http.StatusBadRequest)
		return
	}
	fileInfo := files[0]
	fileReader, err := fileInfo.Open()
	dto := types.CreateFileDTO{
		Name:     fileInfo.Filename,
		Size:     fileInfo.Size,
		Reader:   fileReader,
		Metadata: r.Form.Get("Metadata"),
	}

	err = ro.fileService.Create(r.Context(), userID, dto)
	if err != nil {
		http.Error(w, "unable to store file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (ro *router) deleteFile(w http.ResponseWriter, r *http.Request) {
	fileId := chi.URLParam(r, "id")
	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	err = ro.fileService.Delete(r.Context(), userID, fileId)
	if err != nil {
		http.Error(w, "unable to delete: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
