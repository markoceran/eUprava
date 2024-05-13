package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
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

func (h *GranicnaPolicijaHandler) CreatePrelazHandler(w http.ResponseWriter, r *http.Request) {

	var prelaz data.Prelaz

	if err := json.NewDecoder(r.Body).Decode(&prelaz); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Pogresan format zahteva"))
		return
	}

	prelaz.ID = primitive.NewObjectID()
	prelaz.Datum = primitive.NewDateTimeFromTime(time.Now())

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := h.granicnaPolicijaRepo.CreatePrelaz(ctx, &prelaz)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Greska prilikom kreiranja prelaza"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

//func (h *GranicnaPolicijaHandler) CreatePrelazHandler(w http.ResponseWriter, r *http.Request) {
//	var prelaz data.Prelaz
//
//	if err := json.NewDecoder(r.Body).Decode(&prelaz); err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write([]byte("Pogresan format zahtjeva"))
//		return
//	}
//
//	// Validacija dokumenata prije kreiranja prelaza
//	resp, err := http.Post("http://localhost:8002/validirajDokumente", "application/json", bytes.NewBufferString(fmt.Sprintf(`{"JMBG": "%s", "Ime": "%s", "Prezime": "%s", "BrojLicneKarte": "%s", "BrojPasosa": "%s", "Drzavljanstvo": "%s"}`, prelaz.JMBGPutnika, prelaz.ImePutnika, prelaz.PrezimePutnika, prelaz.BrojLicneKartePutnika, prelaz.BrojPasosaPutnika, prelaz.DrzavljanstvoPutnika)))
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte("Greska prilikom slanja zahtjeva za validaciju"))
//		return
//	}
//	defer resp.Body.Close()
//
//	// Provjera odgovora validacije
//	if resp.StatusCode != http.StatusOK {
//		w.WriteHeader(resp.StatusCode)
//		responseMessage, _ := ioutil.ReadAll(resp.Body)
//		w.Write(responseMessage)
//		return
//	}
//
//	// Ako su dokumenti validni, nastavljamo sa kreiranjem prelaza
//	prelaz.ID = primitive.NewObjectID()
//	prelaz.Datum = primitive.NewDateTimeFromTime(time.Now())
//
//	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
//	defer cancel()
//
//	err = h.granicnaPolicijaRepo.CreatePrelaz(ctx, &prelaz)
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
