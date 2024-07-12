package data

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	DATABASE                  = "mup"
	COLLECTIONKORISNICI       = "korisnici"
	COLLECTIONNALOGZAPRACENJE = "nalogZaPracenje"
)

type MupRepo struct {
	cli    *mongo.Client
	logger *log.Logger
	client *http.Client
	tabela *mongo.Database
}

func New(ctx context.Context, logger *log.Logger) (*MupRepo, error) {
	dburi := fmt.Sprintf("mongodb://%s:%s/", os.Getenv("MUP_DB_HOST"), os.Getenv("MUP_DB_PORT"))

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}
	tabela := client.Database(DATABASE)
	// Return repository with logger and DB client
	return &MupRepo{
		cli:    client,
		logger: logger,
		client: httpClient,
		tabela: tabela,
	}, nil
}

// Disconnect from database
func (pr *MupRepo) DisconnectMongo(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (rr *MupRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := rr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		rr.logger.Println(err)
	}

	// Print available databases
	databases, err := rr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		rr.logger.Println(err)
	}
	fmt.Println(databases)
}

func (rr *MupRepo) DodajKorisnika(ctx context.Context, koisnik *Korisnik) error {

	rezultat, err := rr.tabela.Collection(COLLECTIONKORISNICI).InsertOne(context.TODO(), koisnik)

	if err != nil {
		log.Println("Greska prilikom dodavanja korisnika")
		return err
	}
	koisnik.ID = rezultat.InsertedID.(primitive.ObjectID)
	return nil
}

func (rr *MupRepo) DobaviKorisnike(ctx context.Context) (Korisnici, error) {
	filter := bson.D{{}}
	return rr.filterKorisnici(ctx, filter)
}

func (rr *MupRepo) DobaviKorisnikaPoJmbg(ctx context.Context, jmbg string) (*Korisnik, error) {
	filter := bson.M{"licnaKarta.jmbg": jmbg}

	korisnik, err := rr.filterJmbg(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Println("Greska dobavljanja korisnika:", err)
		return nil, err
	}

	log.Println("Korisnik:", korisnik)

	return korisnik, nil
}

func (rr *MupRepo) DobaviKorisnikaPoID(ctx context.Context, id primitive.ObjectID) (*Korisnik, error) {
	filter := bson.D{{"_id", id}}
	var korisnik Korisnik

	err := rr.tabela.Collection(COLLECTIONKORISNICI).FindOne(ctx, filter).Decode(&korisnik)
	if err != nil {
		return nil, err
	}

	return &korisnik, nil
}

func (rr *MupRepo) AzurirajKorisnika(ctx context.Context, korisnik *Korisnik) error {
	filter := bson.D{{"_id", korisnik.ID}}
	update := bson.D{{"$set", korisnik}}

	_, err := rr.tabela.Collection(COLLECTIONKORISNICI).UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Greška prilikom ažuriranja korisnika")
		return err
	}
	return nil
}

func (h *MupRepo) ProveriLicnuKartu(korisnikId primitive.ObjectID) (bool, error) {
	// Dobavljanje korisnika iz baze podataka
	korisnik, err := h.DobaviKorisnikaPoID(context.Background(), korisnikId)
	if err != nil {
		return false, err
	}

	// Provera da li korisnik ima već izdatu ličnu kartu
	return korisnik.LicnaKarta != nil, nil
}

func (h *MupRepo) ProveriVozackuDozvolu(korisnikId primitive.ObjectID) (bool, error) {
	// Dobavljanje korisnika iz baze podataka
	korisnik, err := h.DobaviKorisnikaPoID(context.Background(), korisnikId)
	if err != nil {
		return false, err
	}

	// Provera da li korisnik ima već izdatu ličnu kartu
	return korisnik.Vozacka != nil, nil
}

func (rr *MupRepo) DobaviNalogPoSumjivomLicu(ctx context.Context, jmbg string) (*NalogZaPracenje, error) {
	filter := bson.D{{"gradjanin.licnaKarta.jmbg", jmbg}}
	var nalogzaPracenje NalogZaPracenje

	err := rr.tabela.Collection(COLLECTIONNALOGZAPRACENJE).FindOne(ctx, filter).Decode(&nalogzaPracenje)
	if err != nil {
		return nil, err
	}

	return &nalogzaPracenje, nil
}

func (rr *MupRepo) DodajNalogZaPracenje(nalog *NalogZaPracenje) error {

	_, err := rr.tabela.Collection(COLLECTIONNALOGZAPRACENJE).InsertOne(context.TODO(), nalog)

	if err != nil {
		log.Println("Greska prilikom dodavanja korisnika")
		return err
	}
	return nil
}

func (rr *MupRepo) DobaviNalogeZaPracenje(ctx context.Context) (NaloziZaPracenje, error) {
	filter := bson.D{{}}
	cursor, err := rr.tabela.Collection(COLLECTIONNALOGZAPRACENJE).Find(ctx, filter)
	if err != nil {
		log.Println("Ne postoje nalozi za pracenje za dati filter")
		return nil, err
	}
	defer cursor.Close(ctx)

	var nalozi NaloziZaPracenje
	for cursor.Next(context.TODO()) {
		var nalog NalogZaPracenje
		err = cursor.Decode(&nalog)
		if err != nil {
			return nil, nil
		}
		nalozi = append(nalozi, &nalog)
	}
	err = cursor.Err()
	return nalozi, nil
}

func (h *MupRepo) ProveriSaobracajnuDozvolu(korisnikId primitive.ObjectID) (bool, error) {
	// Dobavljanje korisnika iz baze podataka
	korisnik, err := h.DobaviKorisnikaPoID(context.Background(), korisnikId)
	if err != nil {
		return false, err
	}

	// Provera da li korisnik ima već izdatu saobracajnu
	return korisnik.Saobracajna != nil, nil
}

func (h *MupRepo) ProveriPasos(korisnikId primitive.ObjectID) (bool, error) {
	korisnik, err := h.DobaviKorisnikaPoID(context.Background(), korisnikId)
	if err != nil {
		return false, err
	}

	return korisnik.Pasos != nil, nil
}

func (rr *MupRepo) filterKorisnici(ctx context.Context, filter interface{}) (Korisnici, error) {
	cursor, err := rr.tabela.Collection(COLLECTIONKORISNICI).Find(ctx, filter)
	if err != nil {
		log.Println("Ne postoji korisnik za dati filter")
		return nil, err
	}
	defer cursor.Close(ctx)

	return decodeKorisnici(cursor)
}

func (rr *MupRepo) filterJmbg(ctx context.Context, filter interface{}) (korisnik *Korisnik, err error) {
	result := rr.tabela.Collection(COLLECTIONKORISNICI).FindOne(ctx, filter)
	err = result.Decode(&korisnik)
	return
}

func decodeKorisnici(cursor *mongo.Cursor) (korisnici Korisnici, err error) {
	for cursor.Next(context.TODO()) {
		var korisnik Korisnik
		err = cursor.Decode(&korisnik)
		if err != nil {
			return
		}
		korisnici = append(korisnici, &korisnik)
	}
	err = cursor.Err()
	return
}

func (rr *MupRepo) DobaviJmbgKorisnika(ctx context.Context, id primitive.ObjectID) (string, error) {
	filter := bson.D{{"_id", id}}
	var korisnik Korisnik

	err := rr.tabela.Collection(COLLECTIONKORISNICI).FindOne(ctx, filter).Decode(&korisnik)
	if err != nil {
		return "", err
	}

	if korisnik.LicnaKarta == nil {
		return "", fmt.Errorf("Korisnik nema licnu kartu")
	}

	return korisnik.LicnaKarta.JMBG, nil
}
