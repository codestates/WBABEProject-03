package model

import (
	"context"
	"fmt"
	"log"
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
	colMenus   *mongo.Collection
	colOrders  *mongo.Collection
}

type Person struct {
	Name string `json:"name" bson:"name"`
	Age  int    `json:"age" bson:"age"`
	Pnum string `json:"pnum" bson:"pnum"`
}

// 이름, 가격
type Menu struct {
	Name     string `json:"name" bson:"name"`
	Price    int    `json:"price" bson:"price"`
	IsOpened bool   `json:"isOpened" bson:"isOpened"`
	IsToday  bool   `json:"isToday" bson:"isToday"`
}
type Order struct {
	ClientName string   `json:"name" bson:"name"`
	MenuList   []string `json:"menuList" bson:"menuList"`
	TotalPrice int      `json:"totalPrice" bson:"totalPrice"`
	Status     string   `json:"status" bson:"status"`
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
		r.colMenus = db.Collection("tMenu")
		r.colOrders = db.Collection("tOrder")
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

	// var pers []bson.M
	// if err = cursor.All(ctx, &pers); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(pers)

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var pers bson.M
		if err = cursor.Decode(&pers); err != nil {
			log.Fatal(err)
		}
		fmt.Println(pers)
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

func (p *Model) FindAllMenu() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{}
	cursor, err := p.colMenus.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	var menus []bson.M
	if err = cursor.All(ctx, &menus); err != nil {
		log.Fatal(err)
	}
	fmt.Println(menus)
	return menus, err
}
func (p *Model) GetOneMenu(elem string) (Menu, error) {
	opts := []*options.FindOneOptions{}
	fmt.Println(elem)
	filter := bson.M{"name": elem}

	var menu Menu
	if err := p.colMenus.FindOne(context.TODO(), filter, opts...).Decode(&menu); err != nil {
		return menu, err
	} else {
		fmt.Println(menu)
		return menu, nil
	}
}
func (p *Model) CreateMenu(menu Menu) error {
	if _, err := p.colMenus.InsertOne(context.TODO(), menu); err != nil {
		fmt.Println("fail insert new menu")
		return fmt.Errorf("fail, insert new menu")
	}
	return nil
}
func (p *Model) UpdateMenu(name string, price int) error {
	filter := bson.M{"name": name}
	update := bson.M{
		"$set": bson.M{
			"price": price,
		},
	}

	if _, err := p.colMenus.UpdateOne(context.Background(), filter, update); err != nil {
		return err
	}
	return nil
}
func (p *Model) UpdateMenuOnTodayMenu(name string, isToday bool) error {
	filter := bson.M{"name": name}
	update := bson.M{
		"$set": bson.M{
			"isToday": isToday,
		},
	}

	if _, err := p.colMenus.UpdateMany(context.Background(), filter, update); err != nil {
		return err
	}
	return nil
}
func (p *Model) DeleteMenu(name string) error {
	filter := bson.M{"name": name}

	if res, err := p.colMenus.DeleteOne(context.TODO(), filter); res.DeletedCount <= 0 {
		return fmt.Errorf("Could not Delete, Not found name %s", name)
	} else if err != nil {
		return err
	}
	return nil
}

// FindAllOrderWithoutDone
func (p *Model) FindAllOrderWithoutDone() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"status": bson.M{
			"$ne": "배달완료",
		},
	}
	cursor, err := p.colOrders.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	var orders []bson.M
	if err = cursor.All(ctx, &orders); err != nil {
		log.Fatal(err)
	}
	fmt.Println(orders)
	return orders, err
}

// CreateOrder
func (p *Model) CreateOrder(order Order) error {
	if _, err := p.colOrders.InsertOne(context.TODO(), order); err != nil {
		fmt.Println("fail insert new order")
		return fmt.Errorf("fail, insert new order")
	}
	return nil
}

func (p *Model) UpdateOrderMenu(pnum int, menuList []string, totalPrice int) error {
	filter := bson.M{"pnum": pnum}
	update := bson.M{
		"$set": bson.M{
			"menuList":   menuList,
			"totalPrice": totalPrice,
		},
	}

	if _, err := p.colOrders.UpdateMany(context.Background(), filter, update); err != nil {
		return err
	}
	return nil
}

func (p *Model) UpdateOrderStatus(pnum int, status string) error {
	filter := bson.M{"pnum": pnum}
	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	if _, err := p.colOrders.UpdateMany(context.Background(), filter, update); err != nil {
		return err
	}
	return nil
}
