package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
	"granicna_policija_service/data"
	"io/ioutil"
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
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	prelazId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Id prelaza nije procitan"))
		return
	}

	prelaz, err := h.granicnaPolicijaRepo.GetPrelazByID(ctx, prelazId)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Prelaz ne postoji"))
		return
	}

	var sumnjivoLice data.SumnjivoLice
	if err := json.NewDecoder(r.Body).Decode(&sumnjivoLice); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Pogresan format zahteva"))
		return
	}

	sumnjivoLice.ID = primitive.NewObjectID()
	sumnjivoLice.Prelaz = *prelaz

	err = h.granicnaPolicijaRepo.CreateSumnjivoLice(ctx, &sumnjivoLice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Greska prilikom kreiranja sumnjivog lica"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func validateDocuments(prelaz *data.Prelaz) error {

	podaciZaValidaciju := data.PodaciZaValidaciju{
		Ime:            prelaz.ImePutnika,
		Prezime:        prelaz.PrezimePutnika,
		JMBG:           prelaz.JMBGPutnika,
		BrojLicneKarte: prelaz.BrojLicneKartePutnika,
		BrojPasosa:     prelaz.BrojPasosaPutnika,
		Drzavljanstvo:  prelaz.DrzavljanstvoPutnika,
	}

	jsonBody, err := json.Marshal(podaciZaValidaciju)
	if err != nil {
		return fmt.Errorf("Greška prilikom marshalling-a PodaciZaValidaciju: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8002/validirajDokumente", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("Greška prilikom kreiranja HTTP zahtjeva: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Greška prilikom izvršavanja HTTP zahtjeva: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Greška prilikom čitanja odgovora: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP zahtjev nije uspio, status kod: %d, odgovor: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (h *GranicnaPolicijaHandler) CreatePrelazHandler(w http.ResponseWriter, r *http.Request) {
	var prelaz data.Prelaz

	if err := json.NewDecoder(r.Body).Decode(&prelaz); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Pogrešan format zahtjeva"))
		return
	}

	// Validiraj dokumente prije kreiranja Prelaza
	if err := validateDocuments(&prelaz); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Dokumenti nisu validni"))
		return
	}

	prelaz.ID = primitive.NewObjectID()
	prelaz.Datum = primitive.NewDateTimeFromTime(time.Now())

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := h.granicnaPolicijaRepo.CreatePrelaz(ctx, &prelaz)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Greška prilikom kreiranja prelaza"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

//func (h *GranicnaPolicijaHandler) CreatePrelazHandler(w http.ResponseWriter, r *http.Request) {
//
//	var prelaz data.Prelaz
//
//	if err := json.NewDecoder(r.Body).Decode(&prelaz); err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write([]byte("Pogresan format zahteva"))
//		return
//	}
//
//	prelaz.ID = primitive.NewObjectID()
//	prelaz.Datum = primitive.NewDateTimeFromTime(time.Now())
//
//	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
//	defer cancel()
//
//	err := h.granicnaPolicijaRepo.CreatePrelaz(ctx, &prelaz)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte("Greska prilikom kreiranja prelaza"))
//		return
//	}
//
//	w.WriteHeader(http.StatusCreated)
//}

func (h *GranicnaPolicijaHandler) CreateKrivicnaPrijavaHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	prelazId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Id prelaza nije procitan"))
		return
	}

	prelaz, err := h.granicnaPolicijaRepo.GetPrelazByID(ctx, prelazId)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Prelaz ne postoji"))
		return
	}

	var krivicnaPrijava data.KrivicnaPrijava
	if err := json.NewDecoder(r.Body).Decode(&krivicnaPrijava); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Pogresan format zahteva"))
		return
	}

	krivicnaPrijava.ID = primitive.NewObjectID()
	krivicnaPrijava.Datum = primitive.NewDateTimeFromTime(time.Now())
	krivicnaPrijava.Prelaz = *prelaz

	err = h.granicnaPolicijaRepo.CreateKrivicnaPrijava(ctx, &krivicnaPrijava)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Greska prilikom kreiranja krivicne prijave"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *GranicnaPolicijaHandler) GetSumnjivaLicaHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	sumnjivaLica, err := h.granicnaPolicijaRepo.GetSumnjivaLica(ctx)
	if err != nil {
		http.Error(w, "Error getting Sumnjiva lica", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sumnjivaLica)
}

func (h *GranicnaPolicijaHandler) GetPrelaziHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	prelazi, err := h.granicnaPolicijaRepo.GetPrelazi(ctx)
	if err != nil {
		http.Error(w, "Error getting Prelazi", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prelazi)
}

func (h *GranicnaPolicijaHandler) GetKrivicnePrijaveHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	krivicnePrijave, err := h.granicnaPolicijaRepo.GetKrivicnePrijave(ctx)
	if err != nil {
		http.Error(w, "Error getting Krivicne prijave", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(krivicnePrijave)
}
