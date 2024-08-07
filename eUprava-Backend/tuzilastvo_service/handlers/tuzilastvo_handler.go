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
	gpServiceHost  = os.Getenv("GRANICNA_POLICIJA_SERVICE_HOST")
	gpServicePort  = os.Getenv("GRANICNA_POLICIJA_SERVICE_PORT")
	mupServiceHost = os.Getenv("MUP_SERVICE_HOST")
	mupServicePort = os.Getenv("MUP_SERVICE_PORT")
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

	message := "Zahtev za sudski postupak je uspešno kreiran"
	// Encode and send JSON response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(map[string]string{"message": message})
	if err != nil {
		// handle error
		return
	}

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

func (h *TuzilastvoHandler) KreirajSporazum(prihvaceniZahtev data.ZahtevZaSklapanjeSporazuma) bool {

	novSporazum := data.Sporazum{}
	novSporazum.ID = primitive.NewObjectID()
	novSporazum.Datum = primitive.NewDateTimeFromTime(time.Now())
	novSporazum.Zahtev = prihvaceniZahtev
	err := h.tuzilastvoRepo.DodajSporazum(context.Background(), &novSporazum)
	if err != nil {
		return false
	}

	return true
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

func (h *TuzilastvoHandler) DobaviZahteveZaSklapanjeSporazumaByGradjanin(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "TuzilastvoHandler.DobaviZahteveZaSklapanjeSporazumaByGradjanin")
	defer span.End()

	vars := mux.Vars(r)
	korisnikId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Id korisnika nije procitan")
		http.Error(rw, "Id korisnika nije procitan", http.StatusBadRequest)
		return
	}

	// Convert korisnikId to string
	korisnikIdStr := korisnikId.Hex()

	dobaviJmbgEndpoint := fmt.Sprintf("http://%s:%s/dobaviJmbgKorisnika/%s", mupServiceHost, mupServicePort, korisnikIdStr)
	log.Println(dobaviJmbgEndpoint)

	// Create an HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", dobaviJmbgEndpoint, nil)
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

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom izvrsavanja zahteva u mup servisu"))
		span.SetStatus(codes.Error, "Greska prilikom izvrsavanja zahteva u mup servisu")
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

	// Parse the JMBG from the response body
	var response struct {
		JMBG string `json:"jmbg"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Greska prilikom parsiranja odgovora"))
		span.SetStatus(codes.Error, "Greska prilikom parsiranja odgovora")
		return
	}

	// Fetch requests for agreement based on the JMBG
	zahtevi, err := h.tuzilastvoRepo.DobaviZahteveZaSklapanjeSporazumaPoGradjaninu(ctx, response.JMBG)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Greska prilikom dobijanja zahteva"))
		span.SetStatus(codes.Error, "Greska prilikom dobijanja zahteva")
		return
	}

	if zahtevi == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Nema zahteva za sklapanje sporazuma"))
		span.SetStatus(codes.Error, "Nema zahteva za sklapanje sporazuma")
		return
	}

	// Convert the requests to JSON and send them in the response
	err = json.NewEncoder(rw).Encode(zahtevi)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Greska prilikom konvertovanja u JSON"))
		span.SetStatus(codes.Error, "Greska prilikom konvertovanja u JSON")
		return
	}
}

func (h *TuzilastvoHandler) PrihvatiZahtevZaSklapanjeSporazuma(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "TuzilastvoHandler.PrihvatiZahtevZaSklapanjeSporazuma")
	defer span.End()

	vars := mux.Vars(req)
	zahtevId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Id zahteva nije procitan")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Id zahteva nije procitan"))
		return
	}

	sporazum, _ := h.tuzilastvoRepo.DobaviSporazumPoZahtevu(context.Background(), zahtevId)
	if sporazum != nil {
		span.SetStatus(codes.Error, "Sporazum za prosledjeni zahtev vec postoji. Nije moguce prihvatiti zahtev")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Sporazum za prosledjeni zahtev vec postoji. Nije moguce prihvatiti zahtev"))
		return
	}

	prihvaceniZahtev, err := h.tuzilastvoRepo.PrihvatiZahtevZaSklapanjeSporazuma(ctx, zahtevId)
	if err != nil {
		message := "Greska prilikom prihvatanja zahteva za sklapanje sporazuma"
		// Encode and send JSON response
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(writer).Encode(map[string]string{"message": message})
		if err != nil {
			// handle error
			return
		}
		return
	}

	if prihvaceniZahtev != nil {
		resultat := h.KreirajSporazum(*prihvaceniZahtev)

		if resultat == false {
			message := "Greska prilikom kreiranja sporazuma"
			// Encode and send JSON response
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(writer).Encode(map[string]string{"message": message})
			if err != nil {
				// handle error
				return
			}
			return
		}
	}

	message := "Zahtev za sklapanje sporazuma je uspešno prihvaćen"
	// Encode and send JSON response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(map[string]string{"message": message})
	if err != nil {
		// handle error
		return
	}

}

func (h *TuzilastvoHandler) OdbijZahtevZaSklapanjeSporazuma(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "TuzilastvoHandler.OdbijZahtevZaSklapanjeSporazuma")
	defer span.End()

	vars := mux.Vars(req)
	zahtevId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Id zahteva nije procitan")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Id zahteva nije procitan"))
		return
	}

	zahtev, err := h.tuzilastvoRepo.DobaviZahtevZaSklapanjeSporazuma(ctx, zahtevId)
	if err != nil {
		message := "Greska prilikom dobavljanja zahteva za sklapanje sporazuma"
		// Encode and send JSON response
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(writer).Encode(map[string]string{"message": message})
		if err != nil {
			// handle error
			return
		}
		return
	}

	novZahtevZaSudskiPostupak := data.ZahtevZaSudskiPostupak{}
	novZahtevZaSudskiPostupak.ID = primitive.NewObjectID()
	novZahtevZaSudskiPostupak.Datum = primitive.NewDateTimeFromTime(time.Now())
	novZahtevZaSudskiPostupak.Opis = "Odbijen zahtev za sklapanje sporazuma"
	novZahtevZaSudskiPostupak.KrivicnaPrijava = zahtev.KrivicnaPrijava

	err = h.tuzilastvoRepo.OdbijZahtevZaSklapanjeSporazuma(ctx, zahtevId)
	if err != nil {
		message := "Greska prilikom odbijanja zahteva za sklapanje sporazuma"
		// Encode and send JSON response
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(writer).Encode(map[string]string{"message": message})
		if err != nil {
			// handle error
			return
		}
		return
	}

	err = h.tuzilastvoRepo.DodajZahtevZaSudskiPostupak(ctx, &novZahtevZaSudskiPostupak)
	if err != nil {
		message := "Greska prilikom kreiranja zahteva za sudski postupak"
		// Encode and send JSON response
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(writer).Encode(map[string]string{"message": message})
		if err != nil {
			// handle error
			return
		}
		return
	}

	message := "Zahtev za sklapanje sporazuma je uspešno odbijen"
	// Encode and send JSON response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(map[string]string{"message": message})
	if err != nil {
		// handle error
		return
	}

}

func (h *TuzilastvoHandler) KreirajKanal(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "TuzilastvoHandler.KreirajKanal")
	defer span.End()

	var kanal data.Kanal
	if err := json.NewDecoder(req.Body).Decode(&kanal); err != nil {
		span.SetStatus(codes.Error, "Pogresan format zahteva")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Pogresan format zahteva"))
		return
	}

	//var loc *time.Location
	//
	//loc, err := time.LoadLocation("Europe/Belgrade")
	//if err != nil {
	//	log.Fatalf("Unable to load location: %v", err)
	//}

	kanal.ID = primitive.NewObjectID()
	kanal.Kreiran = time.Now()

	err := h.tuzilastvoRepo.KreirajKanal(ctx, &kanal)
	if err != nil {
		span.SetStatus(codes.Error, "Greska prilikom kreiranja kanala za poruke")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Greska prilikom kreiranja kanala za poruke"))
		return
	}

	message := "Kanal za poruke je uspešno kreiran"
	// Encode and send JSON response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(map[string]string{"message": message})
	if err != nil {
		// handle error
		return
	}

}

func (h *TuzilastvoHandler) DobaviKanale(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "TuzilastvoHandler.DobaviKanale")
	defer span.End()

	kanali, err := h.tuzilastvoRepo.DobaviKanale(ctx)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska"))
		span.SetStatus(codes.Error, "Greska")
	}

	if kanali == nil {
		return
	}

	err = kanali.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom konvertovanja u JSON"))
		span.SetStatus(codes.Error, "Greska prilikom konvertovanja u JSON")
	}
}

func (h *TuzilastvoHandler) KreirajPoruku(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "TuzilastvoHandler.KreirajKanal")
	defer span.End()

	vars := mux.Vars(req)
	kanalId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Id kanala nije procitan")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Id kanala nije procitan"))
		return
	}

	kanal, _ := h.tuzilastvoRepo.DobaviKanal(ctx, kanalId)
	if kanal == nil {
		span.SetStatus(codes.Error, "Kanal sa prosledjenim id ne postoji")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Kanal sa prosledjenim id ne postoji"))
		return
	}

	var poruka data.Poruka
	if err := json.NewDecoder(req.Body).Decode(&poruka); err != nil {
		span.SetStatus(codes.Error, "Pogresan format zahteva")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Pogresan format zahteva"))
		return
	}

	poruka.ID = primitive.NewObjectID()
	poruka.Datum = time.Now()
	poruka.KanalId = kanalId

	rola, err := helper.ExtractUserType(req)
	if err != nil {
		span.SetStatus(codes.Error, "Greska prilikom uzimanja role iz tokena")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Greska prilikom uzimanja role iz tokena"))
		return
	}

	poruka.Posiljalac = rola

	err = h.tuzilastvoRepo.KreirajPoruku(ctx, &poruka)
	if err != nil {
		span.SetStatus(codes.Error, "Greska prilikom kreiranja poruke")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Greska prilikom kreiranja poruke"))
		return
	}

	// Encode and send JSON response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(&poruka)
	if err != nil {
		// handle error
		return
	}

}

func (h *TuzilastvoHandler) DobaviPorukePoKanalu(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "TuzilastvoHandler.DobaviPorukePoKanalu")
	defer span.End()

	vars := mux.Vars(r)
	kanalId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Id kanala nije procitan")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Id kanala nije procitan"))
		return
	}

	poruke, err := h.tuzilastvoRepo.DobaviPorukePoKanalu(ctx, kanalId)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska"))
		span.SetStatus(codes.Error, "Greska")
	}

	if poruke == nil {
		return
	}

	err = poruke.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom konvertovanja u JSON"))
		span.SetStatus(codes.Error, "Greska prilikom konvertovanja u JSON")
	}
}
