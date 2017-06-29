package mdb

//MongoDB business logic

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type StockDataSums struct {
	InvestedSum float32 `json:"isum" bson:"isum"`
	PriceSum    float32 `json:"psum" bson:"psum"`
	QuantitySum int     `json:"qsum" bson:"qsum"`
	Count       int     `json:"cnt" bson:"cnt"`
	Symbol      string  `json:"id" bson:"_id,omitempty"`
}

//StockData to hold stock details
type StockData struct {
	Symbol     string        `json:"symbol"`
	Price      float32       `json:"Price"`
	Quantity   int           `json:"Quantity"`
	InsertDate string        `json:"InsertDate"`
	Tax        float32       `json:"Tax"`
	Operation  string        `json:"Operation"`
	Id         bson.ObjectId `json:"id" bson:"_id,omitempty"`
}

type StockDataTables struct {
	Data [][]string `json:"data"`
}

//MongoDBConn connection details for the DB
type MongoDBConn struct {
	host string
	db   string
	coll string
}

func CreateMongoDBConn(host string, db string, coll string) *MongoDBConn {
	return &MongoDBConn{host, db, coll}
}

//Stock returns a slice of stocks for the given symbol from MongoDB
func (mconn *MongoDBConn) Stock(symbol string, data *[]StockData) {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)
	var sf = bson.M{}
	if symbol != "" {
		sf = bson.M{"symbol": symbol}
	}
	q := c.Find(sf)
	if q == nil {
		log.Fatal(err)
	}

	err = q.All(data)
	if err != nil {
		log.Fatal(err)
	}
}

//StockDataTables returns a slice of stocks for the given symbol from MongoDB
//The JSON formatted so it can be used by DataTables straight away.
//It adds Delete button for the html table as well.
func (mconn *MongoDBConn) StockDataTables(symbol string, data *StockDataTables) {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)
	var sf = bson.M{}
	if symbol != "" {
		sf = bson.M{"symbol": symbol}
	}
	q := c.Find(sf)
	if q == nil {
		log.Fatal(err)
	}
	d := []StockData{}
	err = q.All(&d)
	if err != nil {
		log.Fatal(err)
	}
	dlen := len(d)
	data.Data = make([][]string, 0, dlen-1)
	for i := 0; i < dlen; i++ {
		cd := []string{d[i].Symbol, fmt.Sprintf("%f", d[i].Price), fmt.Sprintf("%d", d[i].Quantity), fmt.Sprintf("%f", d[i].Tax), d[i].InsertDate, d[i].Operation, fmt.Sprintf("<button id='%s' class='delstock btn-danger btn-xs'>Delete</button>", d[i].Id.Hex())}
		data.Data = append(data.Data, cd)
	}

}

//StockSums groups ticks and returns sums for fields
func (mconn *MongoDBConn) StockSums(data *[]StockDataSums) {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)
	// var sf = []bson.M{{"$group": bson.M{"_id": "$symbol",
	// 	"isum": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$price", "$quantity"}}},
	// 	"cnt":  bson.M{"$sum": 1},
	// 	"psum": bson.M{"$sum": "$price"},
	// 	"qsum": bson.M{"$sum": "$quantity"}}}}
	var sf = []bson.M{{"$group": bson.M{"_id": "$symbol",
		"isum": bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$operation", "Buy"}},
			"then": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$price", "$quantity"}}},
			"else": bson.M{"$sum": bson.M{"$multiply": []interface{}{bson.M{"$multiply": []interface{}{"$price", "$quantity"}}, -1}}}}}},
		"cnt": bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$operation", "Buy"}},
			"then": 1,
			"else": 0}}},
		"psum": bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$operation", "Buy"}},
			"then": "$price",
			"else": 0}}},
		"qsum": bson.M{"$sum": bson.M{"$cond": bson.M{"if": bson.M{"$eq": []interface{}{"$operation", "Buy"}},
			"then": "$quantity",
			"else": bson.M{"$multiply": []interface{}{-1, "$quantity"}}}}}}}}

	p := c.Pipe(sf)
	if p == nil {
		log.Fatal(err)
	}
	err = p.All(data)
	if err != nil {
		log.Fatal(err)
	}
}

//StockUniqueSymbols returns a slice of unique stock ticker from MongoDB
func (mconn *MongoDBConn) StockUniqueSymbols(data *[]string) {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)

	q := c.Find(nil)
	if q == nil {
		log.Fatal(err)
	}
	var n int
	n, err = q.Count()
	if err != nil {
		log.Fatal("Failed to count result", err)
	}

	log.Println("Count = ", n)
	err = q.Distinct("symbol", data)
	if err != nil {
		log.Fatal(err)
	}
}

//StockByID returns a slice of stocks for the given Object ID from MongoDB
func (mconn *MongoDBConn) StockByID(id string, data *[]StockData) error {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		log.Println("Failed to connect to DB:", err)
		return err
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)

	q := c.Find(bson.M{"_id": bson.ObjectIdHex(id)})
	if q == nil {
		log.Println("Failed to find by ID:", err)
		return err
	}

	err = q.All(data)
	if err != nil {
		log.Println("Failed to parse result:", err)
		return err
	}
	return nil
}

//RemoveByID removes a document by ID from MongoDB
func (mconn *MongoDBConn) RemoveByID(id string) error {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		log.Println("Failed to connec to DB:", err)
		return err
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)

	err = c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Println("Failed to delete:", err)
		return err
	}
	return nil
}

//Place adds a stock to DB
func (mconn *MongoDBConn) Place(symbol string, price float32, quantity int, oper string, inDate time.Time) {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)
	// sqr := StockQueryResult{}
	// err = YahooStockData("GSK.L", &sqr)
	// if err != nil {
	// 	log.Fatal("test failed")
	// }
	data := StockData{
		Symbol:     symbol,
		Price:      price,
		Quantity:   quantity,
		Operation:  oper,
		InsertDate: "",
	}
	err = c.Insert(&data)
	if err != nil {
		log.Fatal(err)
	}
}

//Place2 adds a stock to DB
func (mconn *MongoDBConn) Place2(sd StockData) error {
	session, err := mgo.Dial(mconn.host)
	if err != nil {
		log.Println("Error Place2, DB connection:", err)
		return err
	}
	defer session.Close()
	c := session.DB(mconn.db).C(mconn.coll)

	//sd.InsertDate = ""
	err = c.Insert(&sd)
	if err != nil {
		log.Println("Error Place2, insert operation:", err)
		return err
	}
	return nil
}
