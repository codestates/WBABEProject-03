package model

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	client     *mongo.Client
	colPersons *mongo.Collection
}

type Person struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
	Pnum string `bson:"pnum"`
}

func NewModel() (*Model, error) {
	r := &Model{}

	var err error
	mgUrl := "mongodb://127.0.0.1:27017"
	if r.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mgUrl)); err != nil {
		return nil, err
	} else if err := r.client.Ping(context.Background(), nil); err != nil {
		return nil, err
	} else {
		db := r.client.Database("go-ready")
		r.colPersons = db.Collection("tPerson")
	}

	return r, nil
}

func (p *Model) GetPerson() []Person {

	filter := bson.D{}
	cur, err := p.colPersons.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	var pers []Person

	for _, result := range pers {
		cur.Decode(&result)
		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Println(output)
	}

	return nil
}
