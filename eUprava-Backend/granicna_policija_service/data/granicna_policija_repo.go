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

type GranicnaPolicijaRepo struct {
	cli    *mongo.Client
	logger *log.Logger
	client *http.Client
}

var (
	gpServiceHost = os.Getenv("GRANICNA_POLICIJA_SERVICE_HOST")
	gpServicePort = os.Getenv("GRANICNA_POLICIJA_SERVICE_PORT")
)

func New(ctx context.Context, logger *log.Logger) (*GranicnaPolicijaRepo, error) {
	dburi := fmt.Sprintf("mongodb://%s:%s/", os.Getenv("GRANICNA_POLICIJA_DB_HOST"), os.Getenv("GRANICNA_POLICIJA_DB_PORT"))

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
	return &GranicnaPolicijaRepo{
		cli:    client,
		logger: logger,
		client: httpClient,
	}, nil
}

// Disconnect from database
func (pr *GranicnaPolicijaRepo) DisconnectMongo(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (rr *GranicnaPolicijaRepo) Ping() {
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

func (pr *GranicnaPolicijaRepo) CreateSumnjivoLice(ctx context.Context, sumnjivoLice *SumnjivoLice) error {
	collection := pr.cli.Database("granicna_policija_db").Collection("sumnjiva_lica")
	_, err := collection.InsertOne(ctx, sumnjivoLice)
	if err != nil {
		return err
	}
	return nil
}
