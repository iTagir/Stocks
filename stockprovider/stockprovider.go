package main

//a webservice which queries mongo DB and returns JSON with single or multiple
//stock data.
import (
	"encoding/json"
	"fmt"
	"github.com/iTagir/stocks/common"
	"github.com/iTagir/stocks/mdb"
	"github.com/iTagir/stocks/yahoolib"
	"log"
	"net/http"
	"os"
)

//handleStocksProgress queries for stocks in MongoDB and then queries current process in Yahoo
func handleStocksProgress(dbhost string, dbname string, dbcoll string) common.HTTPResponseFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := make([]string, 0, 100)
		mdbConn := mdb.CreateMongoDBConn(dbhost, dbname, dbcoll)
		mdbConn.StockUniqueSymbols(&d)
		sqr := yahoolib.StockQueryResult{}
		for _, sym := range d {
			err := yahoolib.YahooStockData(sym, &sqr)
			if err != nil {
				log.Fatal("test failed")
			}
		}
	}
}

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

	host := os.Getenv("STOCK_HOST")
	port := os.Getenv("STOCK_PORT")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoDB := os.Getenv("MONGO_DB")
	mongoColl := os.Getenv("MONGO_COLLECTION")

	if port == "" {
		log.Fatal("Port variable STOCK_PORT was not set.")
		return
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	handle := handleStocks(mongoHost, mongoDB, mongoColl)
	http.HandleFunc("/stocksDataTable", handle)
	handle = handleStocksProgress(mongoHost, mongoDB, mongoColl)
	http.HandleFunc("/stocksProgress", handle)

	log.Println("Start listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
