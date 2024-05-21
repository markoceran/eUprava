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
	authServiceHost            = os.Getenv("AUTH_SERVICE_HOST")
	authServicePort            = os.Getenv("AUTH_SERVICE_PORT")
	granicaPolicijaServiceHost = os.Getenv("GRANICNA_POLICIJA_SERVICE_HOST")
	granicaPolicijaServicePort = os.Getenv("GRANICNA_POLICIJA_SERVICE_PORT")
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

	korisnikUBazi, _ := h.mupRepo.DobaviKorisnikaPoID(ctx, korisnikId)
	if korisnikUBazi != nil {
		span.SetStatus(codes.Error, "Korisnik vec ima izdatu licnu kartu")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Korisnik vec ima izdatu licnu kartu"))
		return
	}

	korisnik, err := h.DobaviKorisnikaOdAuthServisa(ctx, korisnikId)
	if err != nil {
		span.SetStatus(codes.Error, "Greska pilikom dobavljanja korisnika iz auth servisa")
		writer.WriteHeader(http.StatusNotFound)
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

	licnaIstice := time.Now().AddDate(5, 0, 0).Truncate(24 * time.Hour)
	licnaKarta.Dokument.Istice = primitive.NewDateTimeFromTime(licnaIstice)

	rand.Seed(time.Now().UnixNano())
	brojLicneKarte := generateBrojLicneKarte()
	licnaKarta.BrojLicneKarte = brojLicneKarte

	var jmbg string
	for {
		jmbg = generateJMBG()
		korisnikPoJmbg, _ := h.mupRepo.DobaviKorisnikaPoJmbg(ctx, jmbg)
		if korisnikPoJmbg == nil {
			// JMBG is unique
			break
		}
	}
	licnaKarta.JMBG = jmbg

	korisnik.LicnaKarta = &licnaKarta

	err = h.mupRepo.DodajKorisnika(ctx, &korisnik)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Greska prilikom dodavanja korisnika"))
		span.SetStatus(codes.Error, "Greska prilikom dodavanja korisnika")
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Lična karta je uspešno kreirana"))

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

//KREIRANJE VOZACKE DOZVOLE

func (h *MupHandler) KreirajVozackuDozvolu(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "MupHandler.KreirajVozackuDozvolu")
	defer span.End()

	vars := mux.Vars(req)
	korisnikId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Id korisnika nije procitan")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Id korisnika nije procitan"))
		return
	}

	// Dobavljanje korisnika iz baze podataka
	korisnik, err := h.mupRepo.DobaviKorisnikaPoID(ctx, korisnikId)
	if err != nil {
		span.SetStatus(codes.Error, "Korisnik nema izdatu licnu kartu")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Korisnik nema izdatu licnu kartu"))
		return
	}

	// Provera da li korisnik ima već izdatu vozacku
	imaVozacku, err := h.mupRepo.ProveriVozackuDozvolu(korisnikId)
	if err != nil {
		span.SetStatus(codes.Error, "Greska pri proveri vozacke")
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Greska pri proveri vozacke"))
		return
	}

	if imaVozacku {
		span.SetStatus(codes.Error, "Korisnik vec ima izdatu vozacku dozvolu")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Korisnik vec ima izdatu vozacku dozvolu"))
		return
	}

	var vozackaDozvola data.Vozacka
	if err := json.NewDecoder(req.Body).Decode(&vozackaDozvola); err != nil {
		span.SetStatus(codes.Error, "Pogresan format zahteva")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Pogresan format zahteva"))
		return
	}

	vozackaDozvola.ID = primitive.NewObjectID()
	vozackaDozvola.Dokument.ID = primitive.NewObjectID()
	vozackaDozvola.Dokument.Izdato = primitive.NewDateTimeFromTime(time.Now().Truncate(24 * time.Hour))

	vozackaIstice := time.Now().AddDate(10, 0, 0).Truncate(24 * time.Hour)
	vozackaDozvola.Dokument.Istice = primitive.NewDateTimeFromTime(vozackaIstice)

	korisnik.Vozacka = &vozackaDozvola

	err = h.mupRepo.AzurirajKorisnika(ctx, korisnik)
	if err != nil {
		span.SetStatus(codes.Error, "Greška prilikom ažuriranja korisnika")
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Greška prilikom ažuriranja korisnika"))
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Vozačka dozvola je uspešno kreirana"))

}

//KREIRANJE SAOBRACAJNE DOZVOLE

func (h *MupHandler) KreirajSaobracajnuDozvolu(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "MupHandler.KreirajSaobracajneDozvole")
	defer span.End()

	vars := mux.Vars(req)
	korisnikId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Id korisnika nije procitan")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Id korisnika nije procitan"))
		return
	}

	// Dobavljanje korisnika iz baze podataka
	korisnik, err := h.mupRepo.DobaviKorisnikaPoID(ctx, korisnikId)
	if err != nil {
		span.SetStatus(codes.Error, "Korisnik nema izdatu licnu kartu")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Korisnik nema izdatu licnu kartu"))
		return
	}

	// Provera da li korisnik ima već izdatu saobracajnu
	imaSaobracajnu, err := h.mupRepo.ProveriSaobracajnuDozvolu(korisnikId)
	if err != nil {
		span.SetStatus(codes.Error, "Greska pri proveri saobracajne")
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Greska pri proveri saobracajne"))
		return
	}

	if imaSaobracajnu {
		span.SetStatus(codes.Error, "Korisnik vec ima izdatu saobracajnu dozvolu")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Korisnik vec ima izdatu saobracajnu dozvolu"))
		return
	}

	var saobracajnaDozvola data.Saobracajna
	if err := json.NewDecoder(req.Body).Decode(&saobracajnaDozvola); err != nil {
		span.SetStatus(codes.Error, "Pogresan format zahteva")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Pogresan format zahteva"))
		return
	}

	saobracajnaDozvola.ID = primitive.NewObjectID()
	saobracajnaDozvola.Izdato = primitive.NewDateTimeFromTime(time.Now().Truncate(24 * time.Hour))

	saobracajnaIstice := time.Now().AddDate(10, 0, 0).Truncate(24 * time.Hour)
	saobracajnaDozvola.Istice = primitive.NewDateTimeFromTime(saobracajnaIstice)

	korisnik.Saobracajna = &saobracajnaDozvola

	err = h.mupRepo.AzurirajKorisnika(ctx, korisnik)
	if err != nil {
		span.SetStatus(codes.Error, "Greška prilikom ažuriranja korisnika")
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Greška prilikom ažuriranja korisnika"))
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Saobracajna dozvola je uspešno kreirana"))

}

func (h *MupHandler) KreirajPasos(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "MupHandler.KreirajPasos")
	defer span.End()

	vars := mux.Vars(req)
	korisnikId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Id korisnika nije procitan")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Id korisnika nije procitan"))
		return
	}

	korisnik, err := h.mupRepo.DobaviKorisnikaPoID(ctx, korisnikId)
	if err != nil {
		span.SetStatus(codes.Error, "Korisnik nema izdatu licnu kartu")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Korisnik nema izdatu licnu kartu"))
		return
	}

	imaPasos, err := h.mupRepo.ProveriPasos(korisnikId)
	if err != nil {
		span.SetStatus(codes.Error, "Greska pri proveri pasosa")
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Greska pri proveri pasosa"))
		return
	}

	if imaPasos {
		span.SetStatus(codes.Error, "Korisnik vec ima izdat pasos")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Korisnik vec ima izdat pasos"))
		return
	}

	var pasos data.Pasos
	if err := json.NewDecoder(req.Body).Decode(&pasos); err != nil {
		span.SetStatus(codes.Error, "Pogresan format zahteva")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Pogresan format zahteva"))
		return
	}

	pasos.ID = primitive.NewObjectID()
	pasos.Dokument.ID = primitive.NewObjectID()
	pasos.Dokument.Izdato = primitive.NewDateTimeFromTime(time.Now().Truncate(24 * time.Hour))

	istice := time.Now().AddDate(10, 0, 0).Truncate(24 * time.Hour)
	pasos.Dokument.Istice = primitive.NewDateTimeFromTime(istice)

	brojPasosa := generateBrojPasosa()
	pasos.BrojPasosa = brojPasosa

	korisnik.Pasos = &pasos

	err = h.mupRepo.AzurirajKorisnika(ctx, korisnik)
	if err != nil {
		span.SetStatus(codes.Error, "Greška prilikom ažuriranja korisnika")
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Greška prilikom ažuriranja korisnika"))
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Pasoš je uspešno kreiran"))

}

func generateBrojPasosa() string {
	var brojPasosa string

	brojPasosa += ""

	for i := 0; i < 9; i++ {
		// Generate a random number between 0 and 9
		randomDigit := rand.Intn(10)

		brojPasosa += fmt.Sprint(randomDigit)
	}
	return brojPasosa
}

func (h *MupHandler) ValidirajDokumente(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "MupHandler.ValidacijaDokumenata")
	defer span.End()

	var podaciZaValidaciju data.PodaciZaValidaciju
	if err := json.NewDecoder(req.Body).Decode(&podaciZaValidaciju); err != nil {
		span.SetStatus(codes.Error, "Pogresan format zahteva")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Pogresan format zahteva"))
		return
	}

	korisnik, err := h.mupRepo.DobaviKorisnikaPoJmbg(ctx, podaciZaValidaciju.JMBG)
	if err != nil {
		span.SetStatus(codes.Error, "Korisnik nije pronadjen - jmbg nije validan")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Korisnik nije pronadjen - jmbg nije validan"))
		return
	}

	if korisnik == nil {
		span.SetStatus(codes.Error, "Korisnik nije pronadjen - jmbg nije validan")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Korisnik nije pronadjen - jmbg nije validan"))
		return
	}

	if korisnik.LicnaKarta != nil {
		if dokumentJeIstekao(korisnik.LicnaKarta.Dokument.Istice) {
			span.SetStatus(codes.Error, "Licna karta je istekla")
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("Licna karta je istekla"))
			return
		} else if korisnik.LicnaKarta.Dokument.Ime != podaciZaValidaciju.Ime || korisnik.LicnaKarta.Dokument.Prezime != podaciZaValidaciju.Prezime {
			span.SetStatus(codes.Error, "Ime ili prezime nije validno")
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("Ime ili prezime nije validno"))
			return
		} else if korisnik.LicnaKarta.JMBG != podaciZaValidaciju.JMBG {
			span.SetStatus(codes.Error, "Jmbg nije validan")
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("Jmbg nije validan"))
			return
		} else if korisnik.LicnaKarta.BrojLicneKarte != podaciZaValidaciju.BrojLicneKarte {
			span.SetStatus(codes.Error, "Broj licne karte nije validan")
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("Broj licne karte nije validan"))
			return
		}
	} else {
		span.SetStatus(codes.Error, "Korisnik ne poseduje licnu kartu")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Korisnik ne poseduje licnu kartu"))
		return
	}

	if korisnik.Pasos != nil {
		if dokumentJeIstekao(korisnik.Pasos.Dokument.Istice) {
			span.SetStatus(codes.Error, "Pasos je istekao")
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("Pasos je istekao"))
			return
		} else if korisnik.Pasos.Dokument.Ime != podaciZaValidaciju.Ime || korisnik.Pasos.Dokument.Prezime != podaciZaValidaciju.Prezime {
			span.SetStatus(codes.Error, "Ime ili prezime nije validno")
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("Ime ili prezime nije validno"))
			return
		} else if korisnik.Pasos.BrojPasosa != podaciZaValidaciju.BrojPasosa {
			span.SetStatus(codes.Error, "Broj pasosa nije validan")
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("Broj pasosa nije validan"))
			return
		} else if korisnik.Pasos.Drzavljanstvo != podaciZaValidaciju.Drzavljanstvo {
			span.SetStatus(codes.Error, "Drzavljanstvo nije validno")
			writer.WriteHeader(http.StatusForbidden)
			writer.Write([]byte("Drzavljanstvo nije validno"))
			return
		}
	} else {
		span.SetStatus(codes.Error, "Korisnik ne poseduje pasos")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Korisnik ne poseduje pasos"))
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Dokumenti su validni!"))

}

func dokumentJeIstekao(istice primitive.DateTime) bool {
	// Extract the time.Time value from the primitive.DateTime
	isticeTime := istice.Time()

	// Truncate the time to get only the date portion
	isticeDate := time.Date(isticeTime.Year(), isticeTime.Month(), isticeTime.Day(), 0, 0, 0, 0, time.UTC)

	// Get the current date with time set to 00:00:00
	trenutno := time.Now().Truncate(24 * time.Hour)
	log.Println("Istice", isticeDate)
	log.Println("Trenutno", trenutno)
	// Compare Istice with the current date
	return trenutno.After(isticeDate)
}

//KREIRANJE FUNKCIJE ZA NOVI NALOG

func (h *MupHandler) DodajNalogZaPracenje(lice data.SumnjivoLice) error {

	var noviNalog data.NalogZaPracenje

	noviNalog.ID = primitive.NewObjectID()

	korisnik, err := h.mupRepo.DobaviKorisnikaPoJmbg(context.Background(), lice.Prelaz.JMBGPutnika)
	if err != nil {
		return err
	}

	noviNalog.Gradjanin = korisnik
	noviNalog.Opis = lice.Opis
	noviNalog.Datum = primitive.NewDateTimeFromTime(time.Now().Truncate(24 * time.Hour))

	err = h.mupRepo.DodajNalogZaPracenje(&noviNalog)
	if err != nil {
		return err
	}
	return nil
}

func (h *MupHandler) KreirajNalogZaPracenje(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "MupHandler.DobaviSumnjivoLice")
	defer span.End()

	dobaviSumnjivoLiceEndpoint := fmt.Sprintf("http://%s:%s/sumnjivo-lice/all", granicaPolicijaServiceHost, granicaPolicijaServicePort)
	log.Println(dobaviSumnjivoLiceEndpoint)

	// Create an HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", dobaviSumnjivoLiceEndpoint, nil)
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
	var sumnjivaLica []*data.SumnjivoLice
	if errUnmarshal := json.Unmarshal(body, &sumnjivaLica); errUnmarshal != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Unmarshal greska tela odgovora"))
		span.SetStatus(codes.Error, "Unmarshal greska tela odgovora")
		return
	}

	sumnjivoLiceMap := make(map[string]*data.SumnjivoLice)
	for _, sumnjivoLice := range sumnjivaLica {
		sumnjivoLiceMap[sumnjivoLice.ID.Hex()] = sumnjivoLice
	}

	for _, sumnjivolice := range sumnjivoLiceMap {

		if nalog, _ := h.mupRepo.DobaviNalogPoSumjivomLicu(ctx, sumnjivolice.Prelaz.JMBGPutnika); nalog != nil {
			continue
		} else {
			// Pozivamo funkciju za kreiranje novog naloga u bazi podataka
			err := h.DodajNalogZaPracenje(*sumnjivolice)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("Greska prilikom kreiranja naloga za pracenje"))
				span.SetStatus(codes.Error, "Greska prilikom kreiranja naloga za pracenje")
				return
			}
		}
	}

	rw.WriteHeader(http.StatusOK)
}

func (h *MupHandler) DobaviNalogeZaPracenje(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "MupHandler.DobaviNalogeZaPracenje")
	defer span.End()

	nalozi, err := h.mupRepo.DobaviNalogeZaPracenje(ctx)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska"))
		span.SetStatus(codes.Error, "Greska")
	}

	if nalozi == nil {
		return
	}

	err = nalozi.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom konvertovanja u JSON"))
		span.SetStatus(codes.Error, "Greska prilikom konvertovanja u JSON")
	}
}
