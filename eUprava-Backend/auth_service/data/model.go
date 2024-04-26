package data

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"time"
)

type Korisnik struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Ime           string             `bson:"ime,omitempty" json:"ime"`
	Prezime       string             `bson:"prezime,omitempty" json:"prezime"`
	KorisnickoIme string             `bson:"korisnickoIme,omitempty" json:"korisnickoIme"`
	Lozinka       string             `bson:"lozinka,omitempty" json:"lozinka"`
	LicnaKarta    *LicnaKarta        `bson:"licnaKarta,omitempty" json:"licnaKarta,omitempty"`
	Pasos         *Pasos             `bson:"pasos,omitempty" json:"pasos,omitempty"`
	Saobracajna   *Saobracajna       `bson:"saobracajna,omitempty" json:"saobracajna,omitempty"`
	Vozacka       *Vozacka           `bson:"vozacka,omitempty" json:"vozacka,omitempty"`
	Rola          Rola               `bson:"rola,omitempty" json:"rola"`
}

type Dokument struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Ime           string             `bson:"ime,omitempty" json:"ime"`
	Prezime       string             `bson:"prezime,omitempty" json:"prezime"`
	DatumRodjenja primitive.DateTime `bson:"datumRodjenja,omitempty" json:"datumRodjenja"`
	MestoRodjenja string             `bson:"mestoRodjenja,omitempty" json:"mestoRodjenja"`
	Izdato        primitive.DateTime `bson:"izdato,omitempty" json:"izdato"`
	Istice        primitive.DateTime `bson:"istice,omitempty" json:"istice"`
}

type LicnaKarta struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Dokument       Dokument           `bson:"dokument,omitempty" json:"dokument,omitempty"`
	Pol            Pol                `bson:"pol,omitempty" json:"pol"`
	JMBG           string             `bson:"jmbg,omitempty" json:"jmbg"`
	BrojLicneKarte string             `bson:"brojLicneKarte,omitempty" json:"brojLicneKarte"`
}

type Pasos struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Dokument      Dokument           `bson:"dokument,omitempty" json:"dokument,omitempty"`
	Pol           Pol                `bson:"pol,omitempty" json:"pol"`
	Drzavljanstvo string             `bson:"drzavljanstvo,omitempty" json:"drzavljanstvo"`
	BrojPasosa    string             `bson:"brojPasosa,omitempty" json:"brojPasosa"`
}

type Vozacka struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Dokument   Dokument           `bson:"dokument,omitempty" json:"dokument,omitempty"`
	Kategorija Kategorija         `bson:"kategorija,omitempty" json:"kategorija"`
}

type Saobracajna struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MarkaVozila string             `bson:"markaVozila,omitempty" json:"markaVozila"`
	ModelVozila string             `bson:"modelVozila,omitempty" json:"modelVozila"`
	Izdato      primitive.DateTime `bson:"izdato,omitempty" json:"izdato"`
	Istice      primitive.DateTime `bson:"istice,omitempty" json:"istice"`
}

type Kategorija string

const (
	A = "A"
	B = "B"
	C = "C"
	D = "D"
	F = "F"
)

type Rola string

const (
	Policajac         = "Policajac"
	Gradjanin         = "Gradjanin"
	GranicniSluzbenik = "GranicniSluzbenik"
	Tuzioc            = "Tuzioc"
	Istrazitelj       = "Istrazitelj"
	Sudija            = "Sudija"
)

type Pol string

const (
	Muski  = "Muski"
	Zenski = "Zenski"
)

type Claims struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	KorisnickoIme string             `json:"korisnickoIme"`
	Rola          Rola               `json:"rola"`
	ExpiresAt     time.Time          `json:"expires_at"`
}

type Kredencijali struct {
	KorisnickoIme string `bson:"korisnickoIme" json:"korisnickoIme"`
	Lozinka       string `bson:"lozinka" json:"lozinka"`
}

type Korisnici []*Korisnik

func (o *Korisnici) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Korisnici) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *Korisnik) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Korisnik) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *Kredencijali) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Kredencijali) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}
