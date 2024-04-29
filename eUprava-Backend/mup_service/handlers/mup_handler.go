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
	"math/rand"
	"mup_service/data"
	"net/http"
	"os"
	"time"
)

var (
	authServiceHost = os.Getenv("AUTH_SERVICE_HOST")
	authServicePort = os.Getenv("AUTH_SERVICE_PORT")
)

type KeyProduct struct{}
type KorisnikId struct{ id primitive.ObjectID }

type MupHandler struct {
	logger  *log.Logger
	mupRepo *data.MupRepo
	tracer  trace.Tracer
}

func NewMupHandler(l *log.Logger, r *data.MupRepo, t trace.Tracer) *MupHandler {
	return &MupHandler{l, r, t}
}

func (h *MupHandler) DobaviKorisnikaOdAuthServisa(ctx context.Context, korisnikId primitive.ObjectID) (data.Korisnik, error) {
	dobaviKorisnikaEndpoint := fmt.Sprintf("http://%s:%s/korisnik/%s", authServiceHost, authServicePort, korisnikId.Hex())

	// Create an HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", dobaviKorisnikaEndpoint, nil)
	if err != nil {
		fmt.Println("Greska prilikom kreiranja zahteva:", err)
		return data.Korisnik{}, err
	}

	// Make the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Greska prilikom kreiranja zahteva:", err)
		return data.Korisnik{}, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return data.Korisnik{}, fmt.Errorf("Greska: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Greska prilikom citanja tela odgovora:", err)
		return data.Korisnik{}, err
	}

	// Unmarshal the response body into a Korisnik object
	var korisnik data.Korisnik
	if err := json.Unmarshal(body, &korisnik); err != nil {
		fmt.Println("Unmarshal greska tela odgovora:", err)
		return data.Korisnik{}, err
	}

	return korisnik, nil
}

func (h *MupHandler) KreirajLicnuKartu(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "MupHandler.KreirajLicnuKartu")
	defer span.End()

	vars := mux.Vars(req)
	korisnikId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Id korisnika nije procitan")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Id korisnika nije procitan"))
		return
	}

	korisnik, err := h.DobaviKorisnikaOdAuthServisa(ctx, korisnikId)
	if err != nil {
		span.SetStatus(codes.Error, "Greska pilikom dobavljanja korisnika iz auth servisa")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Greska pilikom dobavljanja korisnika iz auth servisa"))
		return
	}

	var licnaKarta data.LicnaKarta
	if err := json.NewDecoder(req.Body).Decode(&licnaKarta); err != nil {
		span.SetStatus(codes.Error, "Pogresan format zahteva")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Pogresan format zahteva"))
		return
	}

	licnaKarta.ID = primitive.NewObjectID()
	licnaKarta.Dokument.ID = primitive.NewObjectID()
	licnaKarta.Dokument.Izdato = primitive.NewDateTimeFromTime(time.Now().Truncate(24 * time.Hour))

	istice := time.Now().AddDate(5, 0, 0).Truncate(24 * time.Hour)
	licnaKarta.Dokument.Istice = primitive.NewDateTimeFromTime(istice)

	rand.Seed(time.Now().UnixNano())
	brojLicneKarte := generateBrojLicneKarte()
	licnaKarta.BrojLicneKarte = brojLicneKarte

	rand.Seed(time.Now().UnixNano())
	jmbg := generateJMBG()
	licnaKarta.JMBG = jmbg

	korisnik.LicnaKarta = &licnaKarta

	err = h.mupRepo.DodajKorisnika(ctx, &korisnik)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Greska prilikom dodavanja korisnika"))
		span.SetStatus(codes.Error, "Greska prilikom dodavanja korisnika")
		return
	}

}

func (h *MupHandler) DobaviKorisnike(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "AuthHandler.DobaviKorisnike")
	defer span.End()

	korisnici, err := h.mupRepo.DobaviKorisnike(ctx)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska"))
		span.SetStatus(codes.Error, "Greska")
	}

	if korisnici == nil {
		return
	}

	err = korisnici.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom konvertovanja u JSON"))
		span.SetStatus(codes.Error, "Greska prilikom konvertovanja u JSON")
	}
}

func generateJMBG() string {
	var jmbg string
	for i := 0; i < 13; i++ {
		// Generate a random number between 0 and 9
		randomDigit := rand.Intn(10)

		// Append the random digit to the JMBG string
		jmbg += fmt.Sprint(randomDigit)
	}
	return jmbg
}

func generateBrojLicneKarte() string {
	var brojLicneKarte string

	// Ensure the first digit is 0
	brojLicneKarte += "0"

	for i := 0; i < 8; i++ {
		// Generate a random number between 0 and 9
		randomDigit := rand.Intn(10)

		// Append the random digit to the Broj Licne Karte string
		brojLicneKarte += fmt.Sprint(randomDigit)
	}
	return brojLicneKarte
}
