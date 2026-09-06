package handler

import (
	"encoding/json"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"net/http"
	"strconv"
)

type GenreHandler struct {
	service *service.GenreService
}

func NewGenreHandler(service *service.GenreService) *GenreHandler {
	return &GenreHandler{
		service: service,
	}
}

func (h *GenreHandler) Create(w http.ResponseWriter, r *http.Request) {

	var genre models.Genre

	err := json.NewDecoder(r.Body).Decode(&genre)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err = h.service.Create(&genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(genre)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *GenreHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	
	var genres []models.Genre
	var err error

	if name != "" {
		genres, err = h.service.GetByName(name)
	} else {
		genres, err = h.service.GetAll()
	}

	if err != nil {
		http.Error(w, "Failed to fetch genres", http.StatusInternalServerError )
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(genres)

	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (h *GenreHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	// PathValue extract the string
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	genre, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // 404 Not Found
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(genre)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *GenreHandler) Update(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var genre models.Genre
	err = json.NewDecoder(r.Body).Decode(&genre)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	genre.ID = id
	err = h.service.Update(&genre)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedGenre, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to fetch updated genre", http.StatusInternalServerError )
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedGenre)
	w.Write([]byte("\nGenre updated successfully\n"))

}

func (h *GenreHandler) Delete(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	forceStr := r.URL.Query().Get("force")
	var force bool
	if forceStr == "true" {
		force = true
	}

	err = h.service.Delete(id, force)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
