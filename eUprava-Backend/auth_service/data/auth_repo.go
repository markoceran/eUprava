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
	DATABASE   = "korisnici"
	COLLECTION = "korisnici"
)

type AuthRepo struct {
	cli    *mongo.Client
	logger *log.Logger
	client *http.Client
	tabela *mongo.Collection
}

func New(ctx context.Context, logger *log.Logger) (*AuthRepo, error) {
	dburi := fmt.Sprintf("mongodb://%s:%s/", os.Getenv("AUTH_DB_HOST"), os.Getenv("AUTH_DB_PORT"))

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
	tabela := client.Database(DATABASE).Collection(COLLECTION)
	// Return repository with logger and DB client
	return &AuthRepo{
		cli:    client,
		logger: logger,
		client: httpClient,
		tabela: tabela,
	}, nil
}

// Disconnect from database
func (pr *AuthRepo) DisconnectMongo(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (rr *AuthRepo) Ping() {
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

func (rr *AuthRepo) DobaviKorisnike(ctx context.Context) (Korisnici, error) {
	filter := bson.D{{}}
	return rr.filter(ctx, filter)
}

func (rr *AuthRepo) DodajKorisnika(ctx context.Context, koisnik *Korisnik) error {

	koisnik.ID = primitive.NewObjectID()
	rezultat, err := rr.tabela.InsertOne(context.TODO(), koisnik)

	if err != nil {
		log.Println("Greska prilikom dodavanja korisnika")
		return err
	}
	koisnik.ID = rezultat.InsertedID.(primitive.ObjectID)
	return nil
}

func (rr *AuthRepo) DobaviKorisnika(ctx context.Context, korisnickoIme string) (*Korisnik, error) {
	filter := bson.M{"korisnickoIme": korisnickoIme}

	korisnik, err := rr.filterOne(ctx, filter)
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

func (rr *AuthRepo) DobaviKorisnikaPoId(ctx context.Context, id primitive.ObjectID) (*Korisnik, error) {
	filter := bson.M{"_id": id}

	korisnik, err := rr.filterOne(ctx, filter)
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

func (rr *AuthRepo) filter(ctx context.Context, filter interface{}) (Korisnici, error) {
	cursor, err := rr.tabela.Find(ctx, filter)
	if err != nil {
		log.Println("Ne postoji korisnik za dati filter")
		return nil, err
	}
	defer cursor.Close(ctx)

	return decode1(cursor)
}

func (rr *AuthRepo) filterOne(ctx context.Context, filter interface{}) (korisnik *Korisnik, err error) {
	result := rr.tabela.FindOne(ctx, filter)
	err = result.Decode(&korisnik)
	return
}

func decode1(cursor *mongo.Cursor) (korisnici Korisnici, err error) {
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
