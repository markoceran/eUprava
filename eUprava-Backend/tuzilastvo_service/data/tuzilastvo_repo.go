package data

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	DATABASE                         = "tuzilastvo"
	COLLECTIONZAHTEVZASUDSKIPOSTUPAK = "zahtevZaSudskiPostupak"
)

type TuzilastvoRepo struct {
	cli    *mongo.Client
	logger *log.Logger
	client *http.Client
	tabela *mongo.Database
}

func New(ctx context.Context, logger *log.Logger) (*TuzilastvoRepo, error) {
	dburi := fmt.Sprintf("mongodb://%s:%s/", os.Getenv("TUZILASTVO_DB_HOST"), os.Getenv("TUZILASTVO_DB_PORT"))

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
	return &TuzilastvoRepo{
		cli:    client,
		logger: logger,
		client: httpClient,
		tabela: tabela,
	}, nil
}

// Disconnect from database
func (pr *TuzilastvoRepo) DisconnectMongo(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (rr *TuzilastvoRepo) Ping() {
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

func (rr *TuzilastvoRepo) DodajZahtevZaSudskiPostupak(ctx context.Context, zahtev *ZahtevZaSudskiPostupak) error {

	_, err := rr.tabela.Collection(COLLECTIONZAHTEVZASUDSKIPOSTUPAK).InsertOne(context.TODO(), zahtev)

	if err != nil {
		log.Println("Greska prilikom dodavanja zahteva za sudski postupak")
		return err
	}
	return nil
}

func (rr *TuzilastvoRepo) DobaviZahteveZaSudskiPostupak(ctx context.Context) (ZahteviZaSudskiPostupak, error) {
	filter := bson.D{{}}
	return rr.filterZahteviZaSudskiPostupak(ctx, filter)
}

func (rr *TuzilastvoRepo) filterZahteviZaSudskiPostupak(ctx context.Context, filter interface{}) (ZahteviZaSudskiPostupak, error) {
	cursor, err := rr.tabela.Collection(COLLECTIONZAHTEVZASUDSKIPOSTUPAK).Find(ctx, filter)
	if err != nil {
		log.Println("Ne postoje zahtevi za sudski postupak za dati filter")
		return nil, err
	}
	defer cursor.Close(ctx)

	return decodeZahteviZaSudskiPostupak(cursor)
}

func decodeZahteviZaSudskiPostupak(cursor *mongo.Cursor) (zahtevi ZahteviZaSudskiPostupak, err error) {
	for cursor.Next(context.TODO()) {
		var zahtev ZahtevZaSudskiPostupak
		err = cursor.Decode(&zahtev)
		if err != nil {
			return
		}
		zahtevi = append(zahtevi, &zahtev)
	}
	err = cursor.Err()
	return
}
