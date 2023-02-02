package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	client     *mongo.Client
	colMenus   *mongo.Collection
	colReviews *mongo.Collection
	colOrders  *mongo.Collection
}

type Menu struct {
	MenuName    string   `bson:"menuName"`
	AbleToOrder bool     `bson:"ableToOrder"`
	Origin      []string `bson:"origin"`
	Price       int      `bson:"price"`
	RegistTime  int      `bson:"registTime"`
	ViewFlag    bool     `bson:"viewFlag"`
	RecoFlag    bool     `bson:"recoFlag"`
	Score       float32  `bson:"score"`
	ReorderNum  int      `bson:"reorderNum"`
}

// 어디에 쓰더라?
type Filter struct {
	MenuName string `bson:"menuName"`
}

type Review struct {
	MenuName string `bson:"menuName"`
	ScoreNum int    `bson:"score"`
	Orderer  string `bson:"orderer"`
	Review   string `bson:"review"`
}

type Order struct {
	OrderNumber int      `bson:"orderNumber"`
	OrderedMenu []string `bson:"orderedMenu"`
	TotalPrice  int      `bson:"totalPrice"`
	Orderer     string   `bson:"orderer"`
	OrderedTime int      `bson:"orderedTime"`
	State       string   `bson:"State"`
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
		r.colMenus = db.Collection("tMenu")
	}

	return r, nil
}

func (p *Model) GetMenu(query string) []Menu {
	// sort?q=latest(최신),reorder(재주문),highest(별점높은순),reco(금일추천메뉴)

	var menus []Menu
	var opts *options.FindOptions
	var filter = bson.D{}

	if query == "reco" {
		opts = options.Find()
		filter = append(filter, bson.E{Key: "recoFlag", Value: true})
	} else if query == "" || query == "latest" {
		opts = options.Find().SetSort(bson.D{{Key: "registTime", Value: -1}}) // same as {"registTime",-1}
	} else if query == "reorder" {
		opts = options.Find().SetSort(bson.D{{Key: "reorderNum", Value: 1}})
	} else if query == "highest" {
		opts = options.Find().SetSort(bson.D{{Key: "score", Value: 1}})
	} else {
		errors.New("query를 확인하십시오.")
	}

	cursor, err := p.colMenus.Find(context.TODO(), filter, opts)
	if err != nil {
		panic(err)
	}

	for _, result := range menus {
		cursor.Decode(&result)
		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", output)
	}
	return menus
}

func ReturnBool(str string) bool {
	boolVal, err := strconv.ParseBool(str)
	if err != nil {
		panic(err)
	}
	return boolVal
}

func ReturnInt(str string) int {
	intVal, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return intVal
}

func ReturnFloat(str string) float32 {
	floatVal, err := strconv.ParseFloat(str, 32)
	if err != nil {
		panic(err)
	}
	return float32(floatVal)
}

func ReturnArray(str string) []string {
	var data []string
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		fmt.Println("Error parsing JSON:", err)
	}
	return data
}

func (p *Model) PostMenu(newMenu Menu) {
	now := time.Now()
	secs := now.Unix()
	newMenu2 := Menu{
		MenuName:    newMenu.MenuName,
		AbleToOrder: newMenu.AbleToOrder,
		Origin:      newMenu.Origin,
		Price:       newMenu.Price,
		RegistTime:  int(secs),
		ViewFlag:    true,
		RecoFlag:    newMenu.RecoFlag,
		ReorderNum:  0,
		Score:       0.0,
	}

	result, err := p.colMenus.InsertOne(context.TODO(), newMenu2)
	if err != nil {
		panic(err)
	}
	fmt.Print("new Menu!")
	fmt.Print(result)
}

func (p *Model) UpdateMenu(MenuName string, Update map[string]string) {
	filter := bson.M{"menuName": MenuName}

	updates := bson.M{}

	for key, val := range Update {
		if key == "ableToOrder" || key == "viewFlag" || key == "recoFlag" {
			updates[key] = ReturnBool(val)
		} else if key == "price" {
			updates[key] = ReturnInt(val)
		} else if key == "score" {
			updates[key] = ReturnFloat(val)
		} else if key == "origin" {
			updates[key] = ReturnArray(val)
		} else {
			updates[key] = val
		}
	}

	update := bson.M{
		"$set": updates,
	}

	result, err := p.colMenus.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	fmt.Print(result)
}

func (p *Model) DeleteMenu(MenuName string) {
	fmt.Print(MenuName)
	filter := bson.M{"menuName": MenuName}
	result, err := p.colMenus.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	fmt.Print(result)
}

func (p *Model) GetScoreByMenuName(MenuName string) float32 {
	filter := bson.D{{"menuName", MenuName}}
	var result Menu
	err := p.colMenus.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		panic(err)
	}
	return result.Score
}

func (p *Model) GetReviewByMenuName(MenuName string) []interface{} {
	filter := bson.M{"menuName": MenuName}
	var results []bson.M
	cursor, err := p.colReviews.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	// Extract the values of review from the documents
	var reviews []interface{}
	for _, result := range results {
		reviews = append(reviews, result["review"])
	}
	return reviews
}

func (p *Model) PostReview(newReview Review) {
	// part for posting new review
	newReview2 := Review{
		MenuName: newReview.MenuName,
		ScoreNum: newReview.ScoreNum,
		Orderer:  newReview.Orderer,
		Review:   newReview.Review,
	}

	result, err := p.colReviews.InsertOne(context.TODO(), newReview2)
	if err != nil {
		panic(err)
	}
	fmt.Print("new Review!")
	fmt.Print(result)

	// part for get average score of the menu to update menu's average score
	cursor, err := p.colReviews.Aggregate(context.TODO(), []bson.M{
		{
			"$match": bson.M{"name": newReview.MenuName},
		},
		{
			"$group": bson.M{
				"_id":     nil,
				"average": bson.M{"$avg": "$fieldName"},
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())

	var result2 bson.M
	for cursor.Next(context.TODO()) {
		cursor.Decode(&result)
	}

	// part for update new average score
	average := result2["average"]
	// convert interface to float64 and convert to string
	averageInStr := strconv.FormatFloat(average.(float64), 'f', -1, 32)
	model := Model{}
	model.UpdateMenu(newReview.MenuName, map[string]string{"score": averageInStr})
}

func (p *Model) GetOrderByOrderer(Orderer string) []Order {

	var orders []Order
	var opts *options.FindOptions
	filter := bson.D{{"orderer", Orderer}}

	cursor, err := p.colMenus.Find(context.TODO(), filter, opts)
	if err != nil {
		panic(err)
	}

	for _, result := range orders {
		cursor.Decode(&result)
		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", output)
	}
	return orders
}

func CalcTotalPrice(colMenus *mongo.Collection, OrderedMenu []string) int {
	var totalPrice int = 0
	var menu Menu

	for _, str := range OrderedMenu {
		var _orderedMenu []string = strings.Split(str, ":")
		filter := bson.D{{"menuName", _orderedMenu[0]}}
		err := colMenus.FindOne(context.TODO(), filter).Decode(&menu)
		if err != nil {
			panic(err)
		}

		totalPrice += menu.Price * ReturnInt(_orderedMenu[1])
	}
	return totalPrice
}

// 재주문시 메뉴에 재주문수 커운트 올라가도록 설정 필요
func (p *Model) PostOrder(newOrder Order) {
	now := time.Now()
	secs := now.Unix()

	filter := bson.D{{}}

	// idx is current OrderNumber
	idx, error := p.colOrders.CountDocuments(context.TODO(), filter)
	if error != nil {
		panic(error)
	}

	newOrder2 := Order{
		OrderNumber: int(idx + 1),
		OrderedMenu: newOrder.OrderedMenu,
		TotalPrice:  CalcTotalPrice(p.colMenus, newOrder.OrderedMenu),
		Orderer:     newOrder.Orderer,
		OrderedTime: int(secs),
		State:       "접수중",
	}

	result, err := p.colOrders.InsertOne(context.TODO(), newOrder2)
	if err != nil {
		panic(err)
	}
	fmt.Print("new Order!")
	fmt.Print(result)
}

// 메뉴변경. state변경 가능
func (p *Model) UpdateOrderByClient(newOrder Order) {

}

// state만 변경
func (p *Model) UpdateOrderByOwner(newOrder Order) {

}
