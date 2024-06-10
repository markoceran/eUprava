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
	DATABASE                             = "tuzilastvo"
	COLLECTIONZAHTEVZASUDSKIPOSTUPAK     = "zahtevZaSudskiPostupak"
	COLLECTIONZAHTEVZASKLAPANJESPORAZUMA = "zahtevZaSklapanjeSporazuma"
	COLLECTIONSPORAZUM                   = "sporazum"
	COLLECTIONPORUKA                     = "poruka"
	COLLECTIONKANAL                      = "kanal"
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

func (rr *TuzilastvoRepo) DodajZahtevZaSklapanjeSporazuma(ctx context.Context, zahtev *ZahtevZaSklapanjeSporazuma) error {

	_, err := rr.tabela.Collection(COLLECTIONZAHTEVZASKLAPANJESPORAZUMA).InsertOne(context.TODO(), zahtev)

	if err != nil {
		log.Println("Greska prilikom dodavanja zahteva za sklapanje sporazuma")
		return err
	}
	return nil
}

func (rr *TuzilastvoRepo) DobaviZahteveZaSudskiPostupak(ctx context.Context) (ZahteviZaSudskiPostupak, error) {
	filter := bson.D{{}}
	return rr.filterZahteviZaSudskiPostupak(ctx, filter)
}

func (rr *TuzilastvoRepo) DobaviZahteveZaSklapanjeSporazuma(ctx context.Context) (ZahteviZaSklapanjeSporazuma, error) {
	filter := bson.D{{}}
	cursor, err := rr.tabela.Collection(COLLECTIONZAHTEVZASKLAPANJESPORAZUMA).Find(ctx, filter)
	if err != nil {
		log.Println("Ne postoje zahtevi za sklapanje sporazuma za dati filter")
		return nil, err
	}
	defer cursor.Close(ctx)

	var zahtevi ZahteviZaSklapanjeSporazuma
	for cursor.Next(context.TODO()) {
		var zahtev ZahtevZaSklapanjeSporazuma
		err = cursor.Decode(&zahtev)
		if err != nil {
			return nil, nil
		}
		zahtevi = append(zahtevi, &zahtev)
	}
	err = cursor.Err()
	return zahtevi, nil
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

func (rr *TuzilastvoRepo) DodajSporazum(ctx context.Context, sporazum *Sporazum) error {

	_, err := rr.tabela.Collection(COLLECTIONSPORAZUM).InsertOne(ctx, sporazum)

	if err != nil {
		log.Println("Greska prilikom dodavanja sporazuma")
		return err
	}
	return nil
}

func (rr *TuzilastvoRepo) DobaviSporazume(ctx context.Context) (Sporazumi, error) {
	filter := bson.D{{}}
	cursor, err := rr.tabela.Collection(COLLECTIONSPORAZUM).Find(ctx, filter)
	if err != nil {
		log.Println("Ne postoje sporazumi za dati filter")
		return nil, err
	}
	defer cursor.Close(ctx)

	var sporazumi Sporazumi
	for cursor.Next(context.TODO()) {
		var sporazum Sporazum
		err = cursor.Decode(&sporazum)
		if err != nil {
			return nil, nil
		}
		sporazumi = append(sporazumi, &sporazum)
	}
	err = cursor.Err()
	return sporazumi, nil
}

func (rr *TuzilastvoRepo) DobaviZahtevZaSudskiPostupakPoPrijavi(ctx context.Context, id primitive.ObjectID) (*ZahtevZaSudskiPostupak, error) {
	filter := bson.D{{"krivicnaPrijava._id", id}}
	var zahtev ZahtevZaSudskiPostupak

	err := rr.tabela.Collection(COLLECTIONZAHTEVZASUDSKIPOSTUPAK).FindOne(ctx, filter).Decode(&zahtev)
	if err != nil {
		return nil, err
	}

	return &zahtev, nil
}

func (rr *TuzilastvoRepo) DobaviZahtevZaSklapanjeSporazumaPoPrijavi(ctx context.Context, id primitive.ObjectID) (*ZahtevZaSklapanjeSporazuma, error) {
	filter := bson.D{{"krivicnaPrijava._id", id}}
	var zahtev ZahtevZaSklapanjeSporazuma

	err := rr.tabela.Collection(COLLECTIONZAHTEVZASKLAPANJESPORAZUMA).FindOne(ctx, filter).Decode(&zahtev)
	if err != nil {
		return nil, err
	}

	return &zahtev, nil
}

func (rr *TuzilastvoRepo) DobaviZahteveZaSklapanjeSporazumaPoGradjaninu(ctx context.Context, jmbg string) (ZahteviZaSklapanjeSporazuma, error) {
	filter := bson.D{{"krivicnaPrijava.prelaz.JMBGPutnika", jmbg}}
	var zahtevi ZahteviZaSklapanjeSporazuma

	cursor, err := rr.tabela.Collection(COLLECTIONZAHTEVZASKLAPANJESPORAZUMA).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var zahtev ZahtevZaSklapanjeSporazuma
		if err := cursor.Decode(&zahtev); err != nil {
			return nil, err
		}
		zahtevi = append(zahtevi, &zahtev)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return zahtevi, nil
}

func (rr *TuzilastvoRepo) PrihvatiZahtevZaSklapanjeSporazuma(ctx context.Context, id primitive.ObjectID) (*ZahtevZaSklapanjeSporazuma, error) {
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"prihvacen", true}}}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var prihvaceniZahtev ZahtevZaSklapanjeSporazuma
	err := rr.tabela.Collection(COLLECTIONZAHTEVZASKLAPANJESPORAZUMA).FindOneAndUpdate(ctx, filter, update, opts).Decode(&prihvaceniZahtev)
	if err != nil {
		return nil, err
	}

	return &prihvaceniZahtev, nil
}

func (rr *TuzilastvoRepo) DobaviSporazumPoZahtevu(ctx context.Context, id primitive.ObjectID) (*Sporazum, error) {
	filter := bson.D{{"zahtev._id", id}}
	var sporzum Sporazum

	err := rr.tabela.Collection(COLLECTIONSPORAZUM).FindOne(ctx, filter).Decode(&sporzum)
	if err != nil {
		return nil, err
	}

	return &sporzum, nil
}

func (rr *TuzilastvoRepo) OdbijZahtevZaSklapanjeSporazuma(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.D{{"_id", id}}

	_, err := rr.tabela.Collection(COLLECTIONZAHTEVZASKLAPANJESPORAZUMA).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (rr *TuzilastvoRepo) DobaviZahtevZaSklapanjeSporazuma(ctx context.Context, id primitive.ObjectID) (*ZahtevZaSklapanjeSporazuma, error) {
	filter := bson.D{{"_id", id}}
	var zahtev ZahtevZaSklapanjeSporazuma

	err := rr.tabela.Collection(COLLECTIONZAHTEVZASKLAPANJESPORAZUMA).FindOne(ctx, filter).Decode(&zahtev)
	if err != nil {
		return nil, err
	}

	return &zahtev, nil
}

func (rr *TuzilastvoRepo) KreirajKanal(ctx context.Context, kanal *Kanal) error {

	_, err := rr.tabela.Collection(COLLECTIONKANAL).InsertOne(context.TODO(), kanal)

	if err != nil {
		log.Println("Greska prilikom kreiranja kanala za poruke")
		return err
	}
	return nil
}

func (rr *TuzilastvoRepo) DobaviKanale(ctx context.Context) (Kanali, error) {
	filter := bson.D{{}}
	cursor, err := rr.tabela.Collection(COLLECTIONKANAL).Find(ctx, filter)
	if err != nil {
		log.Println("Ne postoje kanali za dati filter")
		return nil, err
	}
	defer cursor.Close(ctx)

	var kanali Kanali
	for cursor.Next(context.TODO()) {
		var kanal Kanal
		err = cursor.Decode(&kanal)
		if err != nil {
			return nil, nil
		}
		kanali = append(kanali, &kanal)
	}
	err = cursor.Err()
	return kanali, nil
}

func (rr *TuzilastvoRepo) DobaviKanal(ctx context.Context, id primitive.ObjectID) (*Kanal, error) {
	filter := bson.D{{"_id", id}}
	var kanal Kanal

	err := rr.tabela.Collection(COLLECTIONKANAL).FindOne(ctx, filter).Decode(&kanal)
	if err != nil {
		return nil, err
	}

	return &kanal, nil
}

func (rr *TuzilastvoRepo) KreirajPoruku(ctx context.Context, poruka *Poruka) error {

	_, err := rr.tabela.Collection(COLLECTIONPORUKA).InsertOne(context.TODO(), poruka)

	if err != nil {
		log.Println("Greska prilikom kreiranja poruke")
		return err
	}
	return nil
}

func (rr *TuzilastvoRepo) DobaviPorukePoKanalu(ctx context.Context, id primitive.ObjectID) (Poruke, error) {
	filter := bson.D{{"kanalId", id}}
	cursor, err := rr.tabela.Collection(COLLECTIONPORUKA).Find(ctx, filter)
	if err != nil {
		log.Println("Ne postoje poruke za dati kanal")
		return nil, err
	}
	defer cursor.Close(ctx)

	var poruke Poruke
	for cursor.Next(context.TODO()) {
		var poruka Poruka
		err = cursor.Decode(&poruka)
		if err != nil {
			return nil, nil
		}
		poruke = append(poruke, &poruka)
	}
	err = cursor.Err()
	return poruke, nil
}
