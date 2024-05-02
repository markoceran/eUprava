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
	DATABASE           = "sud"
	COLLECTIONPREDMETI = "predmeti"
)

type SudRepo struct {
	cli    *mongo.Client
	logger *log.Logger
	client *http.Client
	table  *mongo.Database
}

func New(ctx context.Context, logger *log.Logger) (*SudRepo, error) {
	dburi := fmt.Sprintf("mongodb://%s:%s/", os.Getenv("SUD_DB_HOST"), os.Getenv("SUD_DB_PORT"))

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

	table := client.Database(DATABASE)

	// Return repository with logger and DB client
	return &SudRepo{
		cli:    client,
		logger: logger,
		client: httpClient,
		table:  table,
	}, nil
}

// Disconnect from database
func (sr *SudRepo) DisconnectMongo(ctx context.Context) error {
	err := sr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (sr *SudRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := sr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		sr.logger.Println(err)
	}

	// Print available databases
	databases, err := sr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		sr.logger.Println(err)
	}
	fmt.Println(databases)
}

func (sr *SudRepo) DodajPredmet(ctx context.Context, predmet *Predmet) error {
	rezultat, err := sr.table.Collection(COLLECTIONPREDMETI).InsertOne(context.TODO(), predmet)

	if err != nil {
		log.Println("Greska prilikom dodavanja predmeta")
		return err
	}
	predmet.ID = rezultat.InsertedID.(primitive.ObjectID)
	return nil
}

func (sr *SudRepo) DobaviPredmete(ctx context.Context) (Predmeti, error) {
	filter := bson.D{{}}
	return sr.filterPredmeti(ctx, filter)
}

func (sr *SudRepo) filterPredmeti(ctx context.Context, filter interface{}) (Predmeti, error) {
	cursor, err := sr.table.Collection(COLLECTIONPREDMETI).Find(ctx, filter)
	if err != nil {
		log.Println("Ne postoji predmet za dati filter")
		return nil, err
	}
	defer cursor.Close(ctx)

	return decodePredmeti(cursor)
}

func (sr *SudRepo) DobaviPredmetPoID(ctx context.Context, id primitive.ObjectID) (*Predmet, error) {
	filter := bson.D{{"_id", id}}
	var predmet Predmet

	err := sr.table.Collection(COLLECTIONPREDMETI).FindOne(ctx, filter).Decode(&predmet)
	if err != nil {
		return nil, err
	}

	return &predmet, nil
}

func decodePredmeti(cursor *mongo.Cursor) (predmeti Predmeti, err error) {
	for cursor.Next(context.TODO()) {
		var predmet Predmet
		err = cursor.Decode(&predmet)
		if err != nil {
			return
		}
		predmeti = append(predmeti, &predmet)
	}
	err = cursor.Err()
	return
}
