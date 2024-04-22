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

type TuzilastvoRepo struct {
	cli    *mongo.Client
	logger *log.Logger
	client *http.Client
}

var (
	tuzilastvoServiceHost = os.Getenv("TUZILASTVO_SERVICE_HOST")
	tuzilastvoServicePort = os.Getenv("TUZILASTVO_SERVICE_PORT")
)

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

	// Return repository with logger and DB client
	return &TuzilastvoRepo{
		cli:    client,
		logger: logger,
		client: httpClient,
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
