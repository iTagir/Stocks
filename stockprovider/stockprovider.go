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
	"strconv"
)

//StockProgressData the data type to return
type StockProgressData struct {
	Symbol string  `json:"symbol"`
	Price  float32 `json:"price"`
}

type StockDataProgressTables struct {
	Data [][]string `json:"data"`
}

//handleStocksProgress queries for stocks in MongoDB and then queries current process in Yahoo
func handleStocksProgressTable(dbhost string, dbname string, dbcoll string) common.HTTPResponseFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//
		log.Println("p1")
		sp := StockDataProgressTables{}
		mdbConn := mdb.CreateMongoDBConn(dbhost, dbname, dbcoll)
		log.Println("p2")
		//get current sums
		stkSums := []mdb.StockDataSums{}
		err1 := mdbConn.StockSums(&stkSums)
		log.Println("p3")
		if err1 != nil {
			w.Header().Set("Content-Type", "application/json")
			err1 := json.NewEncoder(w).Encode(sp)
			if err1 != nil {
				log.Println("StockProvider response encode error: ", err1)

			}
			return
		}
		log.Println("p3.1")
		splen := len(stkSums)
		log.Println("p3.2", splen)
		if splen == 0 {
			sp.Data = make([][]string, 0, 0)
		} else {
			sp.Data = make([][]string, 0, splen-1)
		}

		log.Println("p3.3")
		//get unique stock symbols
		d := make([]string, 0, 1)
		err1 = mdbConn.StockUniqueSymbols(&d)
		if err1 != nil {
			w.Header().Set("Content-Type", "application/json")
			err1 = json.NewEncoder(w).Encode(sp)
			if err1 != nil {
				log.Println("StockUniqueSymbols response encode error: ", err1)

			}
			return
		}
		log.Println("p4")
		symStr := yahoolib.StringArr2HTML(d)
		//get yahoo for symbols
		sqr := yahoolib.StockQueryResult{}
		//get current prices from yahoo
		err := yahoolib.YahooStockData(symStr, &sqr)
		if err != nil {
			log.Println("Error: ", err)
			//continue
			w.Header().Set("Content-Type", "application/json")
			err2 := json.NewEncoder(w).Encode(sp)
			if err2 != nil {
				log.Println("StockProvider response encode error: ", err2)
			}
			return
		}

		for _, stk := range stkSums {
			//looking for symbol in yahoo array
			var ask string
			for _, v := range sqr.Query.Results.Quote {
				if v.Symbol == stk.Symbol {
					ask = v.Ask
					break
				}
			}
			if ask == "" {
				log.Println("Symbol not found", stk.Symbol)
				continue
			}
			//convert string price to float
			ap, err := strconv.ParseFloat(ask, 32)
			if err != nil {
				log.Fatal("Could not parse asking price.")
			}
			//sp = append(sp, StockProgressData{stk.Symbol, float32(ap)})
			avgPrice := stk.PriceSum / float32(stk.Count)
			invTotal := stk.InvestedSum / 100
			cd := []string{stk.Symbol, fmt.Sprintf("%d", stk.QuantitySum), fmt.Sprintf("%d", stk.Count), fmt.Sprintf("%.2f", invTotal), fmt.Sprintf("%f", avgPrice), fmt.Sprintf("%f", ap)}
			sp.Data = append(sp.Data, cd)
		}
		w.Header().Set("Content-Type", "application/json")
		err3 := json.NewEncoder(w).Encode(sp)
		if err3 != nil {
			log.Println("StockProvider response encode error: ", err3)
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
		err := mdbConn.StockDataTables(stockSymbol, &d)
		if err != nil {
			log.Println("Failed to get data from DB: ", err)
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(d)
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
	handle = handleStocksProgressTable(mongoHost, mongoDB, mongoColl)
	http.HandleFunc("/stocksProgressTable", handle)

	log.Println("Start listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
