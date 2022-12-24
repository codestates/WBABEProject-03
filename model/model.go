package model

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	client     *mongo.Client
	colPersons *mongo.Collection
}

type Person struct {
	Name string `json:"name" bson:"name"`
	Age  int    `json:"age" bson:"age"`
	Pnum string `json:"pnum" bson:"pnum"`
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

func (p *Model) Check(c *gin.Context) {
	p.RespOK(c, 0)
}

func (p *Model) RespOK(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, resp)
}

func (p *Model) GetPerson() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{}
	cursor, err := p.colPersons.Find(ctx, filter)
	if err != nil {
		panic(err)
	}

	var pers []Person
	for _, result := range pers {
		cursor.Decode(&result)
		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", output)
	}
}

func (p *Model) GetOnePerson(flag, elem string) (Person, error) {
	opts := []*options.FindOneOptions{}

	var filter bson.M
	if flag == "name" {
		filter = bson.M{"name": elem}
	} else {
		filter = bson.M{"pnum": elem}
	}

	var pers Person
	if err := p.colPersons.FindOne(context.TODO(), filter, opts...).Decode(&pers); err != nil {
		return pers, err
	} else {
		return pers, nil
	}
}

func (p *Model) CreatePerson(pers Person) error {
	if _, err := p.colPersons.InsertOne(context.TODO(), pers); err != nil {
		fmt.Println("fail insert new person")
		return fmt.Errorf("fail, insert new person")
	}
	return nil
}

func (p *Model) DeletePerson(spnum string) error {
	filter := bson.M{"pnum": spnum}

	if res, err := p.colPersons.DeleteOne(context.TODO(), filter); res.DeletedCount <= 0 {
		return fmt.Errorf("Could not Delete, Not found num %s", spnum)
	} else if err != nil {
		return err
	}
	return nil
}

func (p *Model) UpdatePerson(pnum string, age int) error {
	filter := bson.M{"pnum": pnum}
	update := bson.M{
		"$set": bson.M{
			"age": age,
		},
	}

	if _, err := p.colPersons.UpdateOne(context.Background(), filter, update); err != nil {
		return err
	}
	return nil
}
