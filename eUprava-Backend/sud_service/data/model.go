package data

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
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

type Predmet struct {
	ID       primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Opis     string                 `bson:"opis,omitempty" json:"opis"`
	Datum    primitive.DateTime     `bson:"datum,omitempty" json:"datum"`
	IdSudije primitive.ObjectID     `bson:"idSudije,omitempty" json:"idSudije"`
	Zahtev   ZahtevZaSudskiPostupak `bson:"zahtev,omitempty" json:"zahtev"`
}
type Predmeti []*Predmet

type TerminSudjenja struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Adresa     string             `bson:"adresa,omitempty" json:"adresa"`
	Datum      primitive.DateTime `bson:"datum,omitempty" json:"datum"`
	Prostorija string             `bson:"prostorija,omitempty" json:"prostorija"`
	Predmet    Predmet            `bson:"predmet,omitempty" json:"predmet"`
}

type Presuda struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Opis           string             `bson:"opis,omitempty" json:"opis"`
	Datum          primitive.DateTime `bson:"datum,omitempty" json:"datum"`
	TerminSudjenja TerminSudjenja     `bson:"terminSudjenja,omitempty" json:"terminSudjenja"`
	IdSudije       primitive.ObjectID `bson:"idSudije,omitempty" json:"idSudije"`
}

func (o *Predmeti) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Predmeti) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *Prelaz) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Prelaz) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *KrivicnaPrijava) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *KrivicnaPrijava) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *ZahtevZaSudskiPostupak) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *ZahtevZaSudskiPostupak) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *Predmet) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Predmet) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *TerminSudjenja) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *TerminSudjenja) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *Presuda) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Presuda) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}
