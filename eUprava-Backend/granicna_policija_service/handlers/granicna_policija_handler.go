package handlers

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"granicna_policija_service/data"
	"log"
	"net/http"
	"time"
)

type GranicnaPolicijaHandler struct {
	logger               *log.Logger
	granicnaPolicijaRepo *data.GranicnaPolicijaRepo
	tracer               trace.Tracer
}

func NewGranicnaPolicijaHandler(l *log.Logger, r *data.GranicnaPolicijaRepo, t trace.Tracer) *GranicnaPolicijaHandler {
	return &GranicnaPolicijaHandler{l, r, t}
}

func (h *GranicnaPolicijaHandler) CreateSumnjivoLiceHandler(w http.ResponseWriter, r *http.Request) {
	var sumnjivoLice data.SumnjivoLice
	err := sumnjivoLice.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.granicnaPolicijaRepo.CreateSumnjivoLice(ctx, &sumnjivoLice)
	if err != nil {
		http.Error(w, "Error creating Sumnjivo lice", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
