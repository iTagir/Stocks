package main

import (
	"encoding/json"
	"fmt"
	"github.com/iTagir/stocks/common"
	"github.com/iTagir/stocks/mdb"
	"log"
	"net/http"
)

func handleStocks(dbhost string, dbname string, dbcoll string) common.HTTPResponseFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stockSymbol := r.FormValue("symbol")
		log.Println("Requested symbol:", stockSymbol)
		//d := make([]mdb.StockData, 0, 10)
		d := mdb.StockDataTables{}
		mdbConn := mdb.CreateMongoDBConn(dbhost, dbname, dbcoll)
		mdbConn.StockDataTables(stockSymbol, &d)
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(d)
		if err != nil {
			log.Println("StockProvider response encode error: ", err)
		}
	}
}

func main() {

	host := "localhost"       //os.Getenv("STOCK_HOST")
	port := "33001"           //os.Getenv("STOCK_PORT")
	mongoHost := "tagir-tosh" //os.Getenv("MONGO_HOST")
	mongoDB := "test"         //os.Getenv("MONGO_DB")
	mongoColl := "testcoll"   //os.Getenv("MONGO_COLLECTION")

	if port == "" {
		log.Fatal("Port variable STOCK_PORT was not set.")
		return
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	handle := handleStocks(mongoHost, mongoDB, mongoColl)
	http.HandleFunc("/stocks", handle)
	log.Println("Start listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
