package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"movies-api/internal/models"
	"movies-api/internal/service"
)

type MovieHandler struct {
	service *service.MovieService
}

func NewMovieHandler(service *service.MovieService) *MovieHandler {
	return &MovieHandler{service: service}
}

func (h *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	var movies []models.Movie
	var err error
	
	genreIDStr := r.URL.Query().Get("genre")
	yearStr := r.URL.Query().Get("year")
	actorIDStr := r.URL.Query().Get("actor")

	if genreIDStr != "" {
		genreID, _ := strconv.ParseInt(genreIDStr, 10, 64)
		movies, err = h.service.GetByGenreID(genreID)
	
	} else if yearStr != "" {
		year, _ := strconv.Atoi(yearStr)
		movies, err = h.service.GetByReleaseYear(year)
	
	} else if actorIDStr != "" {
		actorID, _ := strconv.ParseInt(actorIDStr, 10, 64)
		movies, err = h.service.GetByActorID(actorID)
	
	} else {
		// Default: fetch all movies
		movies, err = h.service.GetAll()
	}

	if err != nil {
		http.Error(w, err.Error( ), http.StatusInternalServerError )
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)

}

func (h *MovieHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	movie, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) Update(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(id, &movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedMovie, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to fetch updated movie", http.StatusInternalServerError )
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedMovie)
	w.Write([]byte("\nMovie updated successfull\n"))

}

func (h *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	force := r.URL.Query().Get("force") == "true"

	if err := h.service.Delete(id, force); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MovieHandler) Search(w http.ResponseWriter, r *http.Request) {

	title := r.URL.Query().Get("title")

	movies, err := h.service.SearchByTitle(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}


func (h *MovieHandler) GetActors(w http.ResponseWriter, r *http.Request ) {

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest )
		return
	}

	movie, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error( ), http.StatusNotFound )
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK )
	json.NewEncoder(w).Encode(movie.Actors)

}