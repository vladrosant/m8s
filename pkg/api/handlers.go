package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vladrosant/m8s/pkg/store"
	"github.com/vladrosant/m8s/pkg/types"
)

type Server struct {
	store *store.Store
}

func NewServer(store *store.Store) *Server {
	return &Server{
		store: store,
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func (s *Server) HandleCreatePod(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var pod types.Pod
	if err := json.NewDecoder(r.Body).Decode(&pod); err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("invalid JSON: %v", err))
		return
	}

	if pod.Namespace == "" {
		pod.Namespace = "default"
	}
	if pod.Status == "" {
		pod.Status = types.PodStatusPending
	}
	pod.CreatedAt = time.Now()

	if pod.Name == "" {
		respondError(w, http.StatusBadRequest, "pod image is required")
		return
	}

	if err := s.store.CreatePod(pod); err != nil {
		respondError(w, http.StatusConflict, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, pod)
}

func (s *Server) HandleGetPod(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	namespace := r.URL.Query().Get("namespace")
	name := r.URL.Query().Get("name")

	if namespace == "" {
		namespace = "default"
	}
	if name == "" {
		respondError(w, http.StatusBadRequest, "pod name is required")
		return
	}

	pod, err := s.store.GetPod(namespace, name)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, pod)
}

func (s *Server) HandleListPods(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	pods, err := s.store.ListPods()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, types.PodList{Items: pods})
}

func (s *Server) HandleDeletePod(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	namespace := r.URL.Query().Get("namespace")
	name := r.URL.Query().Get("name")

	if namespace == "" {
		namespace = "default"
	}
	if name == "" {
		respondError(w, http.StatusBadRequest, "pod name is required")
		return
	}

	if err := s.store.DeletePod(namespace, name); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "pod deleted"})
}
