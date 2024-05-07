package handlers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	sumnjivoLice.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.granicnaPolicijaRepo.CreateSumnjivoLice(ctx, &sumnjivoLice)
	if err != nil {
		http.Error(w, "Error creating Sumnjivo lice", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *GranicnaPolicijaHandler) CreatePrelazHandler(w http.ResponseWriter, r *http.Request) {
	var prelaz data.Prelaz
	err := prelaz.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	prelaz.ID = primitive.NewObjectID()
	prelaz.Datum = primitive.NewDateTimeFromTime(time.Now())

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.granicnaPolicijaRepo.CreatePrelaz(ctx, &prelaz)
	if err != nil {
		http.Error(w, "Error creating Prelaz", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *GranicnaPolicijaHandler) CreateKrivicnaPrijavaHandler(w http.ResponseWriter, r *http.Request) {
	var krivicnaPrijava data.KrivicnaPrijava
	err := krivicnaPrijava.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	krivicnaPrijava.ID = primitive.NewObjectID()
	krivicnaPrijava.Datum = primitive.NewDateTimeFromTime(time.Now())

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.granicnaPolicijaRepo.CreateKrivicnaPrijava(ctx, &krivicnaPrijava)
	if err != nil {
		http.Error(w, "Error creating Krivicna prijava", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
