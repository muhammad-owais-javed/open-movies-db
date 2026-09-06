package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"movies-api/internal/models"
	"movies-api/internal/service"
)

type ActorHandler struct {
	service *service.ActorService
}

func NewActorHandler(service *service.ActorService) *ActorHandler {
	return &ActorHandler{service: service}
}

func (h *ActorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var actor models.Actor
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(&actor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(actor)
}

func (h *ActorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	
	name := r.URL.Query().Get("name")
	
	var actors []models.Actor
	var err error

	if name != "" {

		actors, err = h.service.GetByName(name)
	
	} else {
		
		actors, err = h.service.GetAll()
	
	}

	if err != nil {
		http.Error(w, "Failed to fetch actors", http.StatusInternalServerError )
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actors)

}

func (h *ActorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	actor, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(actor)
}

func (h *ActorHandler) Update(w http.ResponseWriter, r *http.Request) {
	
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var actor models.Actor
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(id, &actor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedActor, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to fetch updated actor", http.StatusInternalServerError )
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedActor)
	w.Write([]byte("\nActor updated successfully\n"))

}

func (h *ActorHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
