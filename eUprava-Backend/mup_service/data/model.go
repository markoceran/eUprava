package data

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type EPol int

const (
	Unknown1 EPol = iota
	Muski
	Zenski
)

type Kategorija int

const (
	Unknown2 Kategorija = iota
	A
	B
	C
	D
	F
)

type Rola int

const (
	Unknown3 Rola = iota
	Policajac
	Gradjanin
	GranicniSluzbenik
	Tuzioc
	Istrazitelj
	Sudija
)

type ETip int

const (
	Unknown4 ETip = iota
	LICNAKARTA
	PASOS
	SAOBRACAJNA
	VOZACKA
)

type EStatus int

const (
	Unknown5 EStatus = iota
	POSLAT
	OBRADA
	ZAVRSEN
)

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
	Dokument       Dokument           `bson:"dokument,omitempty" json:"dokument"`
	Pol            EPol               `bson:"pol,omitempty" json:"pol"`
	JMBG           string             `bson:"jmbg,omitempty" json:"jmbg"`
	BrojLicneKarte string             `bson:"brojLicneKarte,omitempty" json:"brojLicneKarte"`
}

type Pasos struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Dokument      Dokument           `bson:"dokument,omitempty" json:"dokument"`
	Pol           EPol               `bson:"pol,omitempty" json:"pol"`
	Drzavljanstvo string             `bson:"drzavljanstvo,omitempty" json:"drzavljanstvo"`
	BrojPasosa    string             `bson:"brojPasosa,omitempty" json:"brojPasosa"`
}

type Vozacka struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Dokument   Dokument           `bson:"dokument,omitempty" json:"dokument"`
	Kategorija Kategorija         `bson:"kategorija,omitempty" json:"kategorija"`
}

type Saobracajna struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MarkaVozila string             `bson:"markaVozila,omitempty" json:"markaVozila"`
	ModelVozila string             `bson:"modelVozila,omitempty" json:"modelVozila"`
	Izdato      primitive.DateTime `bson:"izdato,omitempty" json:"izdato"`
	Istice      primitive.DateTime `bson:"istice,omitempty" json:"istice"`
}

type Korisnik struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Ime           string             `bson:"ime,omitempty" json:"ime"`
	Prezime       string             `bson:"prezime,omitempty" json:"prezime"`
	KorisnickoIme string             `bson:"korisnickoIme,omitempty" json:"korisnickoIme"`
	Lozinka       string             `bson:"lozinka,omitempty" json:"lozinka"`
	LicnaKarta    LicnaKarta         `bson:"licnaKarta,omitempty" json:"licnaKarta"`
	Pasos         Pasos              `bson:"pasos,omitempty" json:"pasos"`
	Saobracajna   Saobracajna        `bson:"saobracajna,omitempty" json:"saobracajna"`
	Vozacka       Vozacka            `bson:"vozacka,omitempty" json:"vozacka"`
	Rola          Rola               `bson:"rola,omitempty" json:"rola"`
}

type Zahtev struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Gradjanin Korisnik           `bson:"gradjanin,omitempty" json:"gradjanin"`
	Datum     primitive.DateTime `bson:"datum,omitempty" json:"datum"`
	Tip       ETip               `bson:"tip,omitempty" json:"tip"`
	Status    EStatus            `bson:"status,omitempty" json:"status"`
}

type NalogZaPracenje struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Gradjanin Korisnik           `bson:"gradjanin,omitempty" json:"gradjanin"`
	Opis      string             `bson:"opis,omitempty" json:"opis"`
	Datum     primitive.DateTime `bson:"datum,omitempty" json:"datum"`
}

//TODO: uraditi za ostale entitete ToJSON i FromJSON

func (o *LicnaKarta) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *LicnaKarta) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}
