package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"sud_service/data"
)

type KeyProduct struct{}

type SudHandler struct {
	logger  *log.Logger
	sudRepo *data.SudRepo
	tracer  trace.Tracer
}

func NewSudHandler(l *log.Logger, r *data.SudRepo, t trace.Tracer) *SudHandler {
	return &SudHandler{l, r, t}
}

func (h *SudHandler) DobaviPredmete(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "SudHandler.DobaviPredmete")
	defer span.End()

	predmeti, err := h.sudRepo.DobaviPredmete(ctx)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska"))
		span.SetStatus(codes.Error, "Greska")
	}

	if predmeti == nil {
		return
	}

	err = predmeti.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom konvertovanja u JSON"))
		span.SetStatus(codes.Error, "Greska prilikom konvertovanja u JSON")
	}
}

func (h *SudHandler) DobaviPredmetPoId(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "SudHandler.DobaviPredmetPoId")
	defer span.End()

	vars := mux.Vars(r)
	predmetId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Greska")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska"))
		return
	}

	predmet, err := h.sudRepo.DobaviPredmetPoID(ctx, predmetId)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska dobavljanja predmeta po ID"))
		span.SetStatus(codes.Error, "Greska dobavljanja predmeta po ID")
	}

	if predmet == nil {
		return
	}

	err = predmet.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom konvertovanja u JSON"))
		span.SetStatus(codes.Error, "Greska prilikom konvertovanja u JSON")
	}
}

func (h *SudHandler) DodajPredmet(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "SudHandler.DodajPredmet")
	defer span.End()

	predmet, ok := req.Context().Value(KeyProduct{}).(*data.Predmet)
	if !ok || predmet == nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Pogresan format zahteva"))
		span.SetStatus(codes.Error, "Pogresan format zahteva")
		return
	}

	err := h.sudRepo.DodajPredmet(ctx, predmet)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Greska prilikom dodavanja predmeta"))
		span.SetStatus(codes.Error, "Greska prilikom dodavanja predmeta")
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (s *SudHandler) MiddlewareDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		predmet := &data.Predmet{}
		err := predmet.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			s.logger.Fatal(err)
			return
		}
		ctx := context.WithValue(h.Context(), KeyProduct{}, predmet)
		h = h.WithContext(ctx)
		next.ServeHTTP(rw, h)
	})
}
