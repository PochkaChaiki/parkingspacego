package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pochkachaiki/parkingspace/internal/model"
)

type Handler struct {
	log *log.Logger
	srv Service
}

func NewHandler(srv Service, logger *log.Logger) *Handler {
	return &Handler{srv: srv, log: logger}
}

// StartSession обрабатывает POST /api/sessions.
func (h *Handler) StartSession(w http.ResponseWriter, r *http.Request) {
	h.log.Printf("StartSession: POST /api/sessions от %s", r.RemoteAddr)

	var req *model.RecordDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Fatalf("StartSession: Unmarshalling JSON error: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	created, err := h.srv.StartSession(r.Context(), req)
	if err != nil {
		h.log.Fatalf("StartSession: service error for  %s: %v", req.PhoneNumber, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if created.Status != "success" {
	// log.Printf("StartSession: Статус %s для %s (место: %d)", created.Status, req.PhoneNumber, req.SpotNumber)
	// } else {
	// log.Printf("StartSession: Успешно создана сессия для %s (место: %d)", req.PhoneNumber, req.SpotNumber)
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetSession обрабатывает GET /api/sessions/{phone_number}.
func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	// Извлекаем phone из пути /api/sessions/{phone}
	phone := strings.TrimPrefix(r.URL.Path, "/api/sessions/")
	h.log.Printf("Get session started")

	record, err := h.srv.GetSession(r.Context(), phone)
	if err != nil {
		h.log.Fatalf("GetSession: error for %s: %v", phone, err)
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	if record == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(record)
}

// ProlongSession обрабатывает PATCH /api/records/{phone_number}.
func (h *Handler) ProlongSession(w http.ResponseWriter, r *http.Request) {
	// Извлекаем phone из пути /api/records/{phone}
	phone := strings.TrimPrefix(r.URL.Path, "/api/records/")

	var req model.ProlongSessionDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Fatalf("ProlongSession: Unmarshalling JSON error for %s: %v", phone, err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Парсим duration (например "1h")
	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		h.log.Fatalf("ProlongSession: invalid duration format '%s' for %s: %v", req.Duration, phone, err)
		http.Error(w, "invalid duration format", http.StatusBadRequest)
		return
	}

	updated, err := h.srv.ProlongSession(r.Context(), phone, duration)
	if err != nil {
		h.log.Fatalf("ProlongSession: service error for %s: %v", phone, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if updated == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// newEndTime := updated.EndTime.Format("2006-01-02 15:04:05")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

// StopSession обрабатывает DELETE /api/records/{phone_number}.
func (h *Handler) StopSession(w http.ResponseWriter, r *http.Request) {
	phone := strings.TrimPrefix(r.URL.Path, "/api/records/")

	if err := h.srv.StopSession(r.Context(), phone); err != nil {
		h.log.Fatalf("StopSession: service error for %s: %v", phone, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Health обрабатывает GET /health для healthchecks
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
