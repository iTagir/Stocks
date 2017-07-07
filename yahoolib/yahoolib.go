package yahoolib

/*
Queries stock data from yahoo

*/

import (
	"fmt"
	"github.com/iTagir/stocks/common"
	"log"
	"net/http"
)

type stockItem struct {
	Symbol string `json:"symbol"`
	Ask    string `json:"Ask"`
}

type result struct {
	Quote []stockItem `json:"quote"`
}

type queryResult struct {
	Count   int    `json:"count"`
	Results result `json:"results"`
}

type StockQueryResult struct {
	Query queryResult `json:"query"`
}

//Stock returns slice of stock from Yahoo FS API
func YahooStockData(symbol string, data *StockQueryResult) error {

	qurl := "http://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20in%20%28" + symbol + "%29&env=store://datatables.org/alltableswithkeys&format=json"
	fmt.Println(qurl)
	resp, err := http.Get(qurl)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("YahooStockData, Request failed: ", resp.StatusCode)
		return fmt.Errorf("Request failed: %d", resp.StatusCode)
	}

	err = common.ParseResponse(resp, data)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	log.Println("YL: ", symbol, data.Query.Results.Quote[0].Ask)
	if data.Query.Count == 0 {
		return fmt.Errorf("Not found.")
	}
	return nil
}

//StringArr2HTML - convert string slice to comma searated, quoted string
func StringArr2HTML(sarr []string) string {
	var out string
	for i, v := range sarr {
		if i == 0 {
			out = fmt.Sprintf("%%22%s%%22", v)
		} else {
			out = fmt.Sprintf("%s%%2C%%22%s%%22", out, v)
		}

	}
	return out
}
