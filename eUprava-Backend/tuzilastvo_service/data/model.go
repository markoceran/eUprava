package data

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"time"
)

type Prelaz struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Datum                 primitive.DateTime `bson:"datum,omitempty" json:"datum"`
	ImePutnika            string             `bson:"imePutnika,omitempty" json:"imePutnika"`
	PrezimePutnika        string             `bson:"prezimePutnika,omitempty" json:"prezimePutnika"`
	JMBGPutnika           string             `bson:"JMBGPutnika,omitempty" json:"JMBGPutnika"`
	BrojLicneKartePutnika string             `bson:"brojLicneKartePutnika,omitempty" json:"brojLicneKartePutnika,omitempty"`
	BrojPasosaPutnika     string             `bson:"brojPasosaPutnika,omitempty" json:"brojPasosaPutnika,omitempty"`
	DrzavljanstvoPutnika  string             `bson:"drzavljanstvoPutnika,omitempty" json:"drzavljanstvoPutnika"`
	MarkaVozila           string             `bson:"markaVozila,omitempty" json:"markaVozila"`
	ModelVozila           string             `bson:"modelVozila,omitempty" json:"modelVozila"`
	SvrhaPutovanja        string             `bson:"svrhaPutovanja,omitempty" json:"svrhaPutovanja"`
	Odobren               bool               `bson:"odobren,omitempty" json:"odobren"`
}

type KrivicnaPrijava struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Datum  primitive.DateTime `bson:"datum,omitempty" json:"datum"`
	Opis   string             `bson:"opis,omitempty" json:"opis"`
	Prelaz Prelaz             `bson:"prelaz,omitempty" json:"prelaz"`
}

type ZahtevZaSudskiPostupak struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Opis            string             `bson:"opis,omitempty" json:"opis"`
	Datum           primitive.DateTime `bson:"datum,omitempty" json:"datum"`
	IdTuzioca       primitive.ObjectID `bson:"idTuzioca,omitempty" json:"idTuzioca"`
	KrivicnaPrijava KrivicnaPrijava    `bson:"krivicnaPrijava,omitempty" json:"krivicnaPrijava"`
}

type ZahtevZaSklapanjeSporazuma struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Opis            string             `bson:"opis,omitempty" json:"opis"`
	Uslovi          string             `bson:"uslovi,omitempty" json:"uslovi"`
	Kazna           string             `bson:"kazna,omitempty" json:"kazna"`
	Datum           primitive.DateTime `bson:"datum,omitempty" json:"datum"`
	IdTuzioca       primitive.ObjectID `bson:"idTuzioca,omitempty" json:"idTuzioca"`
	KrivicnaPrijava KrivicnaPrijava    `bson:"krivicnaPrijava,omitempty" json:"krivicnaPrijava"`
	Prihvacen       bool               `bson:"prihvacen,omitempty" json:"prihvacen"`
}

type Sporazum struct {
	ID     primitive.ObjectID         `bson:"_id,omitempty" json:"id"`
	Zahtev ZahtevZaSklapanjeSporazuma `bson:"zahtev,omitempty" json:"zahtev"`
	Datum  primitive.DateTime         `bson:"datum,omitempty" json:"datum"`
}

type Obavestenje struct {
	ID      primitive.ObjectID         `bson:"_id,omitempty" json:"id"`
	Zahtev  ZahtevZaSklapanjeSporazuma `bson:"zahtev,omitempty" json:"zahtev"`
	Sadrzaj string                     `bson:"sadrzaj,omitempty" json:"sadrzaj"`
	Datum   time.Time                  `bson:"datum,omitempty" json:"datum"`
}

type Poruka struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	KanalId    primitive.ObjectID `bson:"kanalId" json:"kanalId"`
	Posiljalac string             `bson:"posiljalac" json:"posiljalac"`
	Sadrzaj    string             `bson:"sadrzaj" json:"sadrzaj"`
	Datum      time.Time          `bson:"datum" json:"datum"`
}

type Kanal struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Ime     string             `bson:"ime" json:"ime"`
	Opis    string             `bson:"opis" json:"opis"`
	Kreiran time.Time          `bson:"kreiran" json:"kreiran"`
}

type ZahteviZaSudskiPostupak []*ZahtevZaSudskiPostupak

type ZahteviZaSklapanjeSporazuma []*ZahtevZaSklapanjeSporazuma

type Sporazumi []*Sporazum

type Kanali []*Kanal

type Poruke []*Poruka

func (o *ZahtevZaSudskiPostupak) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *ZahtevZaSudskiPostupak) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *ZahteviZaSudskiPostupak) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *ZahteviZaSudskiPostupak) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *ZahtevZaSklapanjeSporazuma) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *ZahtevZaSklapanjeSporazuma) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *ZahteviZaSklapanjeSporazuma) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *ZahteviZaSklapanjeSporazuma) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *Sporazumi) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Sporazumi) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *Kanali) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Kanali) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *Poruke) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Poruke) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}
