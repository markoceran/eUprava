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

func (h *TuzilastvoHandler) KreirajZahtevZaSklapanjeSporazuma(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "TuzilastvoHandler.KreirajZahtevZaSklapanjeSporazuma")
	defer span.End()

	// PROVERA DA LI VEC POSTOJI ZAHTEV

	var zahtev data.ZahtevZaSklapanjeSporazuma
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
	zahtev.Prihvacen = false

	err = h.tuzilastvoRepo.DodajZahtevZaSklapanjeSporazuma(ctx, &zahtev)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Greska prilikom kreiranja zahteva za sklapanje sporazuma"))
		span.SetStatus(codes.Error, "Greska prilikom kreiranja zahteva za sklapanje sporazuma")
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Zahtev za sklapanje sporazuma je uspešno kreiran"))

}

func (h *TuzilastvoHandler) DobaviZahteveZaSklapanjeSporazuma(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "TuzilastvoHandler.DobaviZahteveZaSklapanjeSporazuma")
	defer span.End()

	zahtevi, err := h.tuzilastvoRepo.DobaviZahteveZaSklapanjeSporazuma(ctx)
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

func (h *TuzilastvoHandler) KreirajSporazum(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "TuzilastvoHandler.KreirajSporazum")
	defer span.End()

	// PROVERA DA LI VEC POSTOJI

	//var sporazum data.Sporazum
	//if err := json.NewDecoder(req.Body).Decode(&sporazum); err != nil {
	//	span.SetStatus(codes.Error, "Pogresan format zahteva")
	//	writer.WriteHeader(http.StatusBadRequest)
	//	writer.Write([]byte("Pogresan format zahteva"))
	//	return
	//}

	sporazum := data.Sporazum{}

	sporazum.ID = primitive.NewObjectID()

	sporazum.Datum = primitive.NewDateTimeFromTime(time.Now())

	sporazum.Zahtev = data.ZahtevZaSklapanjeSporazuma{}

	err := h.tuzilastvoRepo.DodajSporazum(ctx, &sporazum)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Greska prilikom kreiranja sporazuma"))
		span.SetStatus(codes.Error, "Greska prilikom kreiranja sporazuma")
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Sporazum je uspešno kreiran"))

}

func (h *TuzilastvoHandler) DobaviSporazume(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "TuzilastvoHandler.DobaviSporazume")
	defer span.End()

	sporazumi, err := h.tuzilastvoRepo.DobaviSporazume(ctx)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska"))
		span.SetStatus(codes.Error, "Greska")
	}

	if sporazumi == nil {
		return
	}

	err = sporazumi.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom konvertovanja u JSON"))
		span.SetStatus(codes.Error, "Greska prilikom konvertovanja u JSON")
	}
}
