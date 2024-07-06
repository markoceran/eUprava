package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"sud_service/client"
	"sud_service/data"
	"time"
)

type KeyProduct struct{}

type SudHandler struct {
	logger           *log.Logger
	sudRepo          *data.SudRepo
	tracer           trace.Tracer
	tuzilastvoClient client.TuzilastvoClient
}

func NewSudHandler(l *log.Logger, r *data.SudRepo, t trace.Tracer, tc client.TuzilastvoClient) *SudHandler {
	return &SudHandler{l, r, t, tc}
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

func (h *SudHandler) DodajPredmetePoZahtjevima(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "SudHandler.DodajPredmet")
	defer span.End()

	zahtjevi, err := h.tuzilastvoClient.DobaviAktivneZahtjeve(req.Context())
	if err != nil {
		log.Printf("Greska prilikom dodavanja zahtjeva: %v", err)
		http.Error(writer, "Greska prilikom dodavanja zahtjeva", http.StatusServiceUnavailable)
	}

	for _, zahtjev := range zahtjevi {
		currentTime := time.Now()
		currentDateTime := primitive.NewDateTimeFromTime(currentTime)

		predmet := &data.Predmet{
			Opis:   zahtjev.Opis,
			Datum:  currentDateTime,
			Zahtev: *zahtjev,
		}

		err = h.sudRepo.DodajPredmet(ctx, predmet)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Greska prilikom dodavanja predmeta"))
			span.SetStatus(codes.Error, "Greska prilikom dodavanja predmeta")
			return
		}
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

// TERMINI
func (h *SudHandler) DobaviTermine(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "SudHandler.DobaviTermine")
	defer span.End()

	termini, err := h.sudRepo.DobaviTermine(ctx)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska"))
		span.SetStatus(codes.Error, "Greska")
	}

	if termini == nil {
		return
	}

	err = termini.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom konvertovanja u JSON"))
		span.SetStatus(codes.Error, "Greska prilikom konvertovanja u JSON")
	}
}

func (h *SudHandler) DobaviTerminPoId(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "SudHandler.DobaviTerminPoId")
	defer span.End()

	vars := mux.Vars(r)
	terminId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Greska")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska"))
		return
	}

	termin, err := h.sudRepo.DobaviTerminPoID(ctx, terminId)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska dobavljanja termina po ID"))
		span.SetStatus(codes.Error, "Greska dobavljanja termina po ID")
	}

	if termin == nil {
		return
	}

	err = termin.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom konvertovanja u JSON"))
		span.SetStatus(codes.Error, "Greska prilikom konvertovanja u JSON")
	}
}

func (h *SudHandler) DodajTermin(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "SudHandler.DodajTermin")
	defer span.End()

	termin, ok := req.Context().Value(KeyProduct{}).(*data.TerminSudjenja)
	if !ok || termin == nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Pogresan format zahteva"))
		span.SetStatus(codes.Error, "Pogresan format zahteva")
		return
	}

	err := h.sudRepo.DodajTermin(ctx, termin)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Greska prilikom dodavanja termina"))
		span.SetStatus(codes.Error, "Greska prilikom dodavanja termina")
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (s *SudHandler) TerminMiddlewareDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		termin := &data.TerminSudjenja{}
		err := termin.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			s.logger.Fatal(err)
			return
		}
		ctx := context.WithValue(h.Context(), KeyProduct{}, termin)
		h = h.WithContext(ctx)
		next.ServeHTTP(rw, h)
	})
}

// PRESUDE
func (h *SudHandler) DobaviPresude(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "SudHandler.DobaviPresude")
	defer span.End()

	presude, err := h.sudRepo.DobaviPresude(ctx)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska"))
		span.SetStatus(codes.Error, "Greska")
	}

	if presude == nil {
		return
	}

	err = presude.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom konvertovanja u JSON"))
		span.SetStatus(codes.Error, "Greska prilikom konvertovanja u JSON")
	}
}

func (h *SudHandler) DobaviPresuduPoId(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "SudHandler.DobaviPresuduPoId")
	defer span.End()

	vars := mux.Vars(r)
	presudaId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Greska")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska"))
		return
	}

	presuda, err := h.sudRepo.DobaviPresuduPoID(ctx, presudaId)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska dobavljanja presude po ID"))
		span.SetStatus(codes.Error, "Greska dobavljanja presude po ID")
	}

	if presuda == nil {
		return
	}

	err = presuda.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom konvertovanja u JSON"))
		span.SetStatus(codes.Error, "Greska prilikom konvertovanja u JSON")
	}
}

func (h *SudHandler) DodajPresudu(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "SudHandler.DodajPresudu")
	defer span.End()

	presuda, ok := req.Context().Value(KeyProduct{}).(*data.Presuda)
	if !ok || presuda == nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Pogresan format zahteva"))
		span.SetStatus(codes.Error, "Pogresan format zahteva")
		return
	}

	err := h.sudRepo.DodajPresudu(ctx, presuda)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Greska prilikom dodavanja presude"))
		span.SetStatus(codes.Error, "Greska prilikom dodavanja presude")
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (s *SudHandler) PresudaMiddlewareDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		presuda := &data.Presuda{}
		err := presuda.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			s.logger.Fatal(err)
			return
		}
		ctx := context.WithValue(h.Context(), KeyProduct{}, presuda)
		h = h.WithContext(ctx)
		next.ServeHTTP(rw, h)
	})
}
