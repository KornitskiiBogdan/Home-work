package internalhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/hw12_13_14_15_calendar/internal/domain"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	var dto eventDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	event, err := fromDTO(dto)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	created, err := s.app.CreateEvent(r.Context(), event)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toDTO(created))
}

func (s *Server) getEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	event, err := s.app.GetEvent(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toDTO(event))
}

func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	var dto eventDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	event, err := fromDTO(dto)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	event.ID = id

	if err := s.app.UpdateEvent(r.Context(), event); err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, struct{}{})
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := s.app.DeleteEvent(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, struct{}{})
}

func (s *Server) listOnDay(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	day, err := time.Parse(time.RFC3339, r.URL.Query().Get("day"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	events, err := s.app.ListOnDay(r.Context(), userID, day)
	if err != nil {
		writeAppError(w, err)
		return
	}
	out := make([]eventDTO, 0, len(events))
	for _, e := range events {
		out = append(out, toDTO(e))
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) listOnWeek(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	week, err := time.Parse(time.RFC3339, r.URL.Query().Get("week"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	events, err := s.app.ListOnWeek(r.Context(), userID, week)
	if err != nil {
		writeAppError(w, err)
		return
	}
	out := make([]eventDTO, 0, len(events))
	for _, e := range events {
		out = append(out, toDTO(e))
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) listOnMonth(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	month, err := time.Parse(time.RFC3339, r.URL.Query().Get("month"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	events, err := s.app.ListOnMonth(r.Context(), userID, month)
	if err != nil {
		writeAppError(w, err)
		return
	}
	out := make([]eventDTO, 0, len(events))
	for _, e := range events {
		out = append(out, toDTO(e))
	}
	writeJSON(w, http.StatusOK, out)
}

func writeAppError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, domain.ErrTitleRequired), errors.Is(err, domain.ErrUserIDRequired), errors.Is(err, domain.ErrEndTimeAfterStart):
		writeError(w, http.StatusBadRequest, err)
	case errors.Is(err, domain.ErrDateBusy), errors.Is(err, domain.ErrIDExists):
		writeError(w, http.StatusConflict, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{
		Error: err.Error(),
		Code:  status,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}
}
