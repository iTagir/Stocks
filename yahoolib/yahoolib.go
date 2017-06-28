package yahoolib

/*
Queries stock data from yahoo

*/

import (
	"fmt"
	"github.com/iTagir/bf/api"

	"log"
	"net/http"
)

type stockItem struct {
	Symbol string `json:"symbol"`
	Ask    string `json:"Ask"`
}

type result struct {
	Quote stockItem `json:"quote"`
}

type queryResult struct {
	Count  int    `json:"count"`
	Result result `json:"results"`
}

type StockQueryResult struct {
	Query queryResult `json:"query"`
}

//Stock returns slice of stock from Yahoo FS API
func YahooStockData(symbol string, data *StockQueryResult) error {

	qurl := "http://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in%20%28%22" + symbol + "%22%29&env=store://datatables.org/alltableswithkeys&format=json"

	resp, err := http.Get(qurl)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("Request failed: ", resp.StatusCode)
		return fmt.Errorf("Request failed: %d", resp.StatusCode)
	}

	err = api.ParseResponse(resp, data)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	log.Println("YL: ", symbol, data.Query.Result.Quote.Ask)
	if data.Query.Result.Quote.Ask == "" {
		return fmt.Errorf("Not found.")
	}
	return nil
}
