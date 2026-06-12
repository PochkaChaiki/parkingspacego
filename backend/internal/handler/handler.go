package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/pochkachaiki/parkingspace/internal/model"
)

const (
	userString500 = "internal server error"
)

type Handler struct {
	log *log.Logger
	srv Service
}

func NewHandler(srv Service, logger *log.Logger) *Handler {
	return &Handler{srv: srv, log: logger}
}

func addHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
}

func getPhoneFromPath(path string) string {
	return strings.TrimPrefix(path, "/api/sessions/")
}

// StartSession обрабатывает POST /api/sessions.
func (h *Handler) StartSession(w http.ResponseWriter, r *http.Request) {
	// h.log.Printf("StartSession: POST /api/sessions от %s", r.RemoteAddr)

	var req *model.RecordDto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Printf("StartSession: Unmarshalling JSON error: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	status, err := h.srv.StartSession(r.Context(), req)
	if err != nil {
		h.log.Printf("StartSession: service error: %v", err)
		http.Error(w, userString500, http.StatusInternalServerError)
		return
	}

	switch status {
	case model.Success:
		w.WriteHeader(http.StatusCreated)
	case model.Failure:
		w.WriteHeader(http.StatusOK)
	case model.Occupied:
		w.WriteHeader(http.StatusOK)
	}

	addHeaders(&w)

	json.NewEncoder(w).Encode(&model.Response{
		Status: status,
	})
}

// GetSession обрабатывает GET /api/sessions/{phone_number}.
func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	phone := getPhoneFromPath(r.URL.Path)
	record, err := h.srv.GetSession(r.Context(), phone)
	if err != nil {
		h.log.Printf("GetSession: error for %s: %v", phone, err)
		http.Error(w, userString500, http.StatusInternalServerError)
		return
	}

	if record == nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	addHeaders(&w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(record)
}

// ProlongSession обрабатывает PATCH /api/sessions/{phone_number}.
func (h *Handler) ProlongSession(w http.ResponseWriter, r *http.Request) {
	// Извлекаем phone из пути /api/sessions/{phone}

	phone := getPhoneFromPath(r.URL.Path)

	var req model.ProlongSessionDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Printf("ProlongSession: Unmarshalling JSON error for %s: %v", phone, err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := h.srv.ProlongSession(r.Context(), phone, req.Duration)
	if err != nil {
		h.log.Printf("ProlongSession: service error for %s: %v", phone, err)
		http.Error(w, userString500, http.StatusInternalServerError)
		return
	}

	if updated == nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	addHeaders(&w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

// StopSession обрабатывает DELETE /api/sessions/{phone_number}.
func (h *Handler) StopSession(w http.ResponseWriter, r *http.Request) {

	phone := getPhoneFromPath(r.URL.Path)

	if err := h.srv.StopSession(r.Context(), phone); err != nil {
		h.log.Printf("StopSession: service error for %s: %v", phone, err)
		http.Error(w, userString500, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Health обрабатывает GET /health для healthchecks
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	addHeaders(&w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
