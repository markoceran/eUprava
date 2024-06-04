package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"tuzilastvo_service/data"
	"tuzilastvo_service/helper"
)

var (
	gpServiceHost = os.Getenv("GRANICNA_POLICIJA_SERVICE_HOST")
	gpServicePort = os.Getenv("GRANICNA_POLICIJA_SERVICE_PORT")
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

	vars := mux.Vars(req)
	prijavaId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Id krivicne prijave nije procitan")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Id krivicne prijave nije procitan"))
		return
	}

	zahtevZaSudskiPostupakPostoji, _ := h.tuzilastvoRepo.DobaviZahtevZaSudskiPostupakPoPrijavi(ctx, prijavaId)
	if zahtevZaSudskiPostupakPostoji != nil {
		span.SetStatus(codes.Error, "Zahtev za sudski postupak za prosledjenu prijavu vec postoji")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Zahtev za sudski postupak za prosledjenu prijavu vec postoji"))
		return
	}

	zahtevZaSklapanjeSporazumaPostoji, _ := h.tuzilastvoRepo.DobaviZahtevZaSklapanjeSporazumaPoPrijavi(ctx, prijavaId)
	if zahtevZaSklapanjeSporazumaPostoji != nil {
		span.SetStatus(codes.Error, "Nije moguce kreirati zahtev za sudski postupak. Za prosledjenu prijavu vec postoji zahtev za sklapanje sporazuma")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Nije moguce kreirati zahtev za sudski postupak. Za prosledjenu prijavu vec postoji zahtev za sklapanje sporazuma"))
		return
	}

	prijava, err := h.DobaviKrivicnuPrijavuByID(prijavaId.Hex())
	if err != nil {
		span.SetStatus(codes.Error, "Greska prilikom dobavljanja krivicne prijave po id")
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("Greska prilikom dobavljanja krivicne prijave po id"))
		return
	}

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
	zahtev.KrivicnaPrijava = *prijava

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

	vars := mux.Vars(req)
	prijavaId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Id krivicne prijave nije procitan")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Id krivicne prijave nije procitan"))
		return
	}

	zahtevZaSudskiPostupakPostoji, _ := h.tuzilastvoRepo.DobaviZahtevZaSudskiPostupakPoPrijavi(ctx, prijavaId)
	if zahtevZaSudskiPostupakPostoji != nil {
		span.SetStatus(codes.Error, "Nije moguce kreirati zahtev za sklapanje sporazuma. Za prosledjenu prijavu vec postoji zahtev za sudski postupak")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Nije moguce kreirati zahtev za sklapanje sporazuma. Za prosledjenu prijavu vec postoji zahtev za sudski postupak"))
		return
	}

	zahtevZaSklapanjeSporazumaPostoji, _ := h.tuzilastvoRepo.DobaviZahtevZaSklapanjeSporazumaPoPrijavi(ctx, prijavaId)
	if zahtevZaSklapanjeSporazumaPostoji != nil {
		span.SetStatus(codes.Error, "Zahtev za sklapanje sporazuma za prosledjenu prijavu vec postoji")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Zahtev za sklapanje sporazuma za prosledjenu prijavu vec postoji"))
		return
	}

	prijava, err := h.DobaviKrivicnuPrijavuByID(prijavaId.Hex())
	if err != nil {
		span.SetStatus(codes.Error, "Greska prilikom dobavljanja krivicne prijave po id")
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("Greska prilikom dobavljanja krivicne prijave po id"))
		return
	}

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
	zahtev.KrivicnaPrijava = *prijava
	zahtev.Prihvacen = false

	err = h.tuzilastvoRepo.DodajZahtevZaSklapanjeSporazuma(ctx, &zahtev)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Greska prilikom kreiranja zahteva za sklapanje sporazuma"))
		span.SetStatus(codes.Error, "Greska prilikom kreiranja zahteva za sklapanje sporazuma")
		return
	}

	message := "Zahtev za sklapanje sporazuma je uspešno kreiran"
	// Encode and send JSON response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(map[string]string{"message": message})
	if err != nil {
		// handle error
		return
	}

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

func (h *TuzilastvoHandler) DobaviKrivicnePrijaveOdGranicnePolicjie(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "TuzilastvoHandler.DobaviKrivicnePrijaveOdGranicnePolicjie")
	defer span.End()

	dobaviPrijaveEndpoint := fmt.Sprintf("http://%s:%s/krivicna-prijava/all", gpServiceHost, gpServicePort)
	log.Println(dobaviPrijaveEndpoint)

	// Create an HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", dobaviPrijaveEndpoint, nil)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom kreiranja zahteva"))
		span.SetStatus(codes.Error, "Greska prilikom kreiranja zahteva")
		return
	}

	// Make the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom slanja zahteva"))
		span.SetStatus(codes.Error, "Greska prilikom slanja zahteva")
		return
	}
	defer resp.Body.Close()

	// Check if resp is nil before further processing
	if resp == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Nil response from server"))
		span.SetStatus(codes.Error, "Nil response from server")
		return
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom izvrsavanja zahteva u granicna policija servisu"))
		span.SetStatus(codes.Error, "Greska prilikom izvrsavanja zahteva u granicna policija servisu")
		return
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Greska prilikom citanja odgovora"))
		span.SetStatus(codes.Error, "Greska prilikom citanja odgovora")
		return
	}

	// Unmarshal the response body into a Korisnik object
	var prijave []*data.KrivicnaPrijava
	if errUnmarshal := json.Unmarshal(body, &prijave); errUnmarshal != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Unmarshal greska tela odgovora"))
		span.SetStatus(codes.Error, "Unmarshal greska tela odgovora")
		return
	}

	prijaveMap := make(map[string]*data.KrivicnaPrijava)
	for _, prijava := range prijave {
		prijaveMap[prijava.ID.Hex()] = prijava
	}

	var svePrijave []*data.KrivicnaPrijava

	for prijavaIdString, prijava := range prijaveMap {
		prijavaId, errObjId := primitive.ObjectIDFromHex(prijavaIdString)
		if errObjId != nil {
			log.Println(codes.Error, "Greska prilikom konverzije ID-ja")
			continue
		}

		if _, err := h.tuzilastvoRepo.DobaviZahtevZaSudskiPostupakPoPrijavi(ctx, prijavaId); err == nil {
			delete(prijaveMap, prijavaIdString)
		} else {
			if _, err := h.tuzilastvoRepo.DobaviZahtevZaSklapanjeSporazumaPoPrijavi(ctx, prijavaId); err == nil {
				delete(prijaveMap, prijavaIdString)
			} else {
				svePrijave = append(svePrijave, prijava)
			}
		}
	}

	jsonResponse, err := json.Marshal(svePrijave)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Greska prilikom enkodiranja u JSON format"))
		span.SetStatus(codes.Error, "Greska prilikom enkodiranja u JSON format")
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonResponse)
}

func (h *TuzilastvoHandler) DobaviKrivicnePrijave() ([]*data.KrivicnaPrijava, error) {

	dobaviPrijaveEndpoint := fmt.Sprintf("http://%s:%s/krivicna-prijava/all", gpServiceHost, gpServicePort)
	log.Println(dobaviPrijaveEndpoint)

	// Create an HTTP request
	req, err := http.NewRequestWithContext(context.Background(), "GET", dobaviPrijaveEndpoint, nil)
	if err != nil {
		log.Println(codes.Error, "Greska prilikom kreiranja zahteva")
		return nil, err
	}

	// Make the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(codes.Error, "Greska prilikom slanja zahteva")
		return nil, err
	}
	defer resp.Body.Close()

	// Check if resp is nil before further processing
	if resp == nil {
		log.Println(codes.Error, "Nil response from server")
		return nil, err
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Println(codes.Error, "Greska prilikom izvrsavanja zahteva u granicna policija servisu")
		return nil, err
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(codes.Error, "Greska prilikom citanja odgovora")
		return nil, err
	}

	// Unmarshal the response body into a Korisnik object
	var prijave []*data.KrivicnaPrijava
	if err := json.Unmarshal(body, &prijave); err != nil {
		log.Println(codes.Error, "Unmarshal greska tela odgovora")
		return nil, err
	}

	return prijave, nil
}

func (h *TuzilastvoHandler) DobaviKrivicnuPrijavuByID(id string) (*data.KrivicnaPrijava, error) {

	prijave, err := h.DobaviKrivicnePrijave()
	if err != nil {
		log.Println(codes.Error, "Greska prilikom dobavljanja krivicnih prijava")
		return nil, err
	}

	// Create a map to store prijave with their IDs as keys
	prijaveMap := make(map[string]*data.KrivicnaPrijava)
	for _, prijava := range prijave {
		prijaveMap[prijava.ID.Hex()] = prijava
	}

	// Retrieve the prijava from the map by its ID
	prijava, ok := prijaveMap[id]
	if !ok {
		return nil, fmt.Errorf("Prijava sa id %s nije pronadjena", id)
	}

	return prijava, nil
}
