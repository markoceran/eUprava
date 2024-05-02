package handlers

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"time"
	"tuzilastvo_service/data"
	"tuzilastvo_service/helper"
)

type KeyProduct struct{}

type TuzilastvoHandler struct {
	logger         *log.Logger
	tuzilastvoRepo *data.TuzilastvoRepo
	tracer         trace.Tracer
}

func NewTuzilastvoHandler(l *log.Logger, r *data.TuzilastvoRepo, t trace.Tracer) *TuzilastvoHandler {
	return &TuzilastvoHandler{l, r, t}
}

func (h *TuzilastvoHandler) KreirajZahtevZaSudskiPostupak(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "TuzilastvoHandler.KreirajZahtevZaSudskiPostupak")
	defer span.End()

	// PROVERA DA LI VEC POSTOJI ZAHTEV ZA S POSTUPAK ZA NEKU K PRIJAVU

	var zahtev data.ZahtevZaSudskiPostupak
	if err := json.NewDecoder(req.Body).Decode(&zahtev); err != nil {
		span.SetStatus(codes.Error, "Pogresan format zahteva")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Pogresan format zahteva"))
		return
	}

	zahtev.ID = primitive.NewObjectID()

	zahtev.Datum = primitive.NewDateTimeFromTime(time.Now())

	claims := helper.ExtractClaims(req)
	logovaniKorisnikId, err := primitive.ObjectIDFromHex(claims["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Id korisnika nije procitan")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Id korisnika nije procitan"))
		return
	}

	zahtev.IdTuzioca = logovaniKorisnikId
	zahtev.KrivicnaPrijava = data.KrivicnaPrijava{}

	err = h.tuzilastvoRepo.DodajZahtevZaSudskiPostupak(ctx, &zahtev)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Greska prilikom kreiranja zahteva za sudski postupak"))
		span.SetStatus(codes.Error, "Greska prilikom kreiranja zahteva za sudski postupak")
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Zahtev za sudski postupak je uspešno kreiran"))

}

func (h *TuzilastvoHandler) DobaviZahteveZaSudskiPostupak(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "TuzilastvoHandler.DobaviZahteveZaSudskiPostupak")
	defer span.End()

	zahtevi, err := h.tuzilastvoRepo.DobaviZahteveZaSudskiPostupak(ctx)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska"))
		span.SetStatus(codes.Error, "Greska")
	}

	if zahtevi == nil {
		return
	}

	err = zahtevi.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom konvertovanja u JSON"))
		span.SetStatus(codes.Error, "Greska prilikom konvertovanja u JSON")
	}
}
