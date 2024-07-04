package data

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type Pol string

const (
	Muski  = "Muski"
	Zenski = "Zenski"
)

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
type LicnaKarta struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Dokument       *Dokument          `bson:"dokument,omitempty" json:"dokument,omitempty"`
	Pol            Pol                `bson:"pol,omitempty" json:"pol"`
	JMBG           string             `bson:"jmbg,omitempty" json:"jmbg,omitempty"`
	BrojLicneKarte string             `bson:"brojLicneKarte,omitempty" json:"brojLicneKarte,omitempty"`
}

type Pasos struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Dokument      *Dokument          `bson:"dokument,omitempty" json:"dokument,omitempty"`
	Pol           Pol                `bson:"pol,omitempty" json:"pol"`
	Drzavljanstvo string             `bson:"drzavljanstvo,omitempty" json:"drzavljanstvo"`
	BrojPasosa    string             `bson:"brojPasosa,omitempty" json:"brojPasosa,omitempty"`
}

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

type SumnjivoLice struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Prelaz Prelaz             `bson:"prelaz,omitempty" json:"prelaz"`
	Opis   string             `bson:"opis,omitempty" json:"opis"`
}

type KrivicnaPrijava struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Datum  primitive.DateTime `bson:"datum,omitempty" json:"datum"`
	Opis   string             `bson:"opis,omitempty" json:"opis"`
	Prelaz Prelaz             `bson:"prelaz,omitempty" json:"prelaz"`
}
type PodaciZaValidaciju struct {
	JMBG           string `bson:"jmbg,omitempty" json:"jmbg,omitempty"`
	Ime            string `bson:"ime,omitempty" json:"ime"`
	Prezime        string `bson:"prezime,omitempty" json:"prezime"`
	BrojLicneKarte string `bson:"brojLicneKarte,omitempty" json:"brojLicneKarte,omitempty"`
	BrojPasosa     string `bson:"brojPasosa,omitempty" json:"brojPasosa,omitempty"`
	Drzavljanstvo  string `bson:"drzavljanstvo,omitempty" json:"drzavljanstvo"`
}
type Putnik struct {
	Ime           string             `bson:"ime,omitempty" json:"ime"`
	Prezime       string             `bson:"prezime,omitempty" json:"prezime"`
	DatumRodjenja primitive.DateTime `bson:"datumRodjenja,omitempty" json:"datumRodjenja"`
	Drzavljanstvo string             `bson:"drzavljanstvo,omitempty" json:"drzavljanstvo"`
	Saobracajna   *Saobracajna       `bson:"saobracajna,omitempty" json:"saobracajna,omitempty"`
	Vozacka       *Vozacka           `bson:"vozacka,omitempty" json:"vozacka,omitempty"`
}

type Vozacka struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Dokument   *Dokument          `bson:"dokument,omitempty" json:"dokument,omitempty"`
	Kategorija Kategorija         `bson:"kategorija,omitempty" json:"kategorija"`
}

type Saobracajna struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MarkaVozila string             `bson:"markaVozila,omitempty" json:"markaVozila"`
	ModelVozila string             `bson:"modelVozila,omitempty" json:"modelVozila"`
	Izdato      primitive.DateTime `bson:"izdato,omitempty" json:"izdato"`
	Istice      primitive.DateTime `bson:"istice,omitempty" json:"istice"`
}

type Dokument struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Ime           string             `bson:"ime,omitempty" json:"ime"`
	Prezime       string             `bson:"prezime,omitempty" json:"prezime"`
	DatumRodjenja primitive.DateTime `bson:"datumRodjenja,omitempty" json:"datumRodjenja"`
	MestoRodjenja string             `bson:"mestoRodjenja,omitempty" json:"mestoRodjenja"`
	Izdato        primitive.DateTime `bson:"izdato,omitempty" json:"izdato,omitempty"`
	Istice        primitive.DateTime `bson:"istice,omitempty" json:"istice,omitempty"`
}

//TODO: uraditi za ostale entitete ToJSON i FromJSON

func (o *Prelaz) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Prelaz) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *SumnjivoLice) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *SumnjivoLice) FromJSON(r io.Reader) error {
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

func (o *PodaciZaValidaciju) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *PodaciZaValidaciju) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}
