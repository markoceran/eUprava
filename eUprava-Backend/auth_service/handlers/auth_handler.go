package handlers

import (
	"auth_service/data"
	"context"
	"encoding/json"
	"github.com/cristalhq/jwt/v4"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

type KeyProduct struct{}

type AuthHandler struct {
	logger   *log.Logger
	authRepo *data.AuthRepo
	tracer   trace.Tracer
}

func NewAuthHandler(l *log.Logger, r *data.AuthRepo, t trace.Tracer) *AuthHandler {
	return &AuthHandler{l, r, t}
}

func (h *AuthHandler) DobaviKorisnike(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "AuthHandler.DobaviKorisnike")
	defer span.End()

	korisnici, err := h.authRepo.DobaviKorisnike(ctx)
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

func (h *AuthHandler) DobaviKorisnikaPoId(rw http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "AuthHandler.DobaviKorisnikaPoId")
	defer span.End()

	vars := mux.Vars(r)
	korisnikId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, "Greska")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska"))
		return
	}

	korisnik, err := h.authRepo.DobaviKorisnikaPoId(ctx, korisnikId)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska dobavljanja korisnika po ID"))
		span.SetStatus(codes.Error, "Greska dobavljanja korisnika po ID")
	}

	if korisnik == nil {
		return
	}

	err = korisnik.ToJSON(rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Greska prilikom konvertovanja u JSON"))
		span.SetStatus(codes.Error, "Greska prilikom konvertovanja u JSON")
	}
}

func (h *AuthHandler) DodajKorisnika(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "AuthHandler.DodavanjeKorisnika")
	defer span.End()

	korisnik, ok := req.Context().Value(KeyProduct{}).(*data.Korisnik)
	if !ok || korisnik == nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Pogresan format zahteva"))
		span.SetStatus(codes.Error, "Pogresan format zahteva")
		return
	}

	k, err := h.authRepo.DobaviKorisnika(ctx, korisnik.KorisnickoIme)
	if k != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Korisnik vec postoji"))
		span.SetStatus(codes.Error, "Korisnik vec postoji")
		return
	}

	lozinka := []byte(korisnik.Lozinka)
	hash, err := bcrypt.GenerateFromPassword(lozinka, bcrypt.DefaultCost)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Greska prilikom hesiranja lozinke"))
		span.SetStatus(codes.Error, "Greska prilikom hesiranja lozinke")
		return
	}
	korisnik.Lozinka = string(hash)

	err = h.authRepo.DodajKorisnika(ctx, korisnik)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Greska prilikom dodavanja korisnika"))
		span.SetStatus(codes.Error, "Greska prilikom dodavanja korisnika")
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Login(writer http.ResponseWriter, req *http.Request) {
	ctx, span := h.tracer.Start(req.Context(), "AuthHandler.Login")
	defer span.End()

	var kredencijali data.Kredencijali
	err := json.NewDecoder(req.Body).Decode(&kredencijali)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Pogresan format zahteva"))
		span.SetStatus(codes.Error, "Pogresan format zahteva")
		return
	}

	if kredencijali.KorisnickoIme == "" || kredencijali.Lozinka == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Sva polja su obavezna"))
		span.SetStatus(codes.Error, "Sva polja su obavezna")
		return
	}

	korisnik, err := h.authRepo.DobaviKorisnika(ctx, kredencijali.KorisnickoIme)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Korisnik ne postoji"))
		span.SetStatus(codes.Error, "Korisnik ne postoji")
		return
	}

	if korisnik == nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Korisnik ne postoji"))
		span.SetStatus(codes.Error, "Korisnik ne postoji")
		return
	}

	lozinkaError := bcrypt.CompareHashAndPassword([]byte(korisnik.Lozinka), []byte(kredencijali.Lozinka))
	if lozinkaError != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Pogresna lozinka"))
		span.SetStatus(codes.Error, "Pogresna lozinka")
		return
	}

	token, err := GenerateJWT(korisnik)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Greska prilikom kreiranja tokena"))
		span.SetStatus(codes.Error, "Greska prilikom kreiranja tokena")
		return
	}

	writer.Write([]byte(token))
}

func GenerateJWT(user *data.Korisnik) (string, error) {

	key := []byte(os.Getenv("SECRET_KEY"))
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		log.Println(err)
	}

	builder := jwt.NewBuilder(signer)

	claims := &data.Claims{
		ID:            user.ID,
		KorisnickoIme: user.KorisnickoIme,
		Rola:          user.Rola,
		ExpiresAt:     time.Now().Add(time.Minute * 60),
	}

	log.Println("id", claims.ID)
	log.Println("korisnicko ime", claims.KorisnickoIme)
	log.Println("rola", claims.Rola)
	log.Println("expires", claims.ExpiresAt)

	token, err := builder.Build(claims)
	if err != nil {
		log.Println(err)
	}

	return token.String(), nil
}

func (s *AuthHandler) MiddlewareDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		korisnik := &data.Korisnik{}
		err := korisnik.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			s.logger.Fatal(err)
			return
		}
		ctx := context.WithValue(h.Context(), KeyProduct{}, korisnik)
		h = h.WithContext(ctx)
		next.ServeHTTP(rw, h)
	})
}
