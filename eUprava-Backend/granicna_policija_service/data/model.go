package data

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type Prelaz struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Datum                primitive.DateTime `bson:"datum,omitempty" json:"datum"`
	ImePutnika           string             `bson:"imePutnika,omitempty" json:"imePutnika"`
	PrezimePutnika       string             `bson:"prezimePutnika,omitempty" json:"prezimePutnika"`
	JMBGPutnika          string             `bson:"JMBGPutnika,omitempty" json:"JMBGPutnika"`
	DrzavljanstvoPutnika string             `bson:"drzavljanstvoPutnika,omitempty" json:"drzavljanstvoPutnika"`
	MarkaVozila          string             `bson:"markaVozila,omitempty" json:"markaVozila"`
	ModelVozila          string             `bson:"modelVozila,omitempty" json:"modelVozila"`
	SvrhaPutovanja       string             `bson:"svrhaPutovanja,omitempty" json:"svrhaPutovanja"`
	Odobren              bool               `bson:"odobren,omitempty" json:"odobren"`
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

//TODO: uraditi za ostale entitete ToJSON i FromJSON

func (o *Prelaz) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *Prelaz) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}