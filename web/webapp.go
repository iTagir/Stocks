package main

import (
	"bytes"
	"fmt"
	"github.com/iTagir/stocks/common"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func handleStocks(url string) common.HTTPResponseFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stockSymbol := r.FormValue("symbol")
		log.Println("Requested symbol:", stockSymbol)
		resp, err := http.Get(url + stockSymbol)
		if err != nil {
			fmt.Println("Error: ", err)
			return
			//return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Println("Request failed: ", resp.StatusCode)
			return
			//return nil, fmt.Errorf("Request failed: %d", resp.StatusCode)
		}
		w.Header().Set("Content-Type", "application/json")
		var bs []byte
		bs, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			w.Write(bs)
		}
	}
}

func handleDashStocks(url string) common.HTTPResponseFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stockSymbol := r.FormValue("symbol")
		log.Println("Requested symbol:", stockSymbol)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error: ", err)
			return
			//return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Println("Request failed: ", resp.StatusCode)
			return
			//return nil, fmt.Errorf("Request failed: %d", resp.StatusCode)
		}
		w.Header().Set("Content-Type", "application/json")
		var bs []byte
		bs, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			w.Write(bs)
		}
	}
}

func handleStocksAdd(url string) common.HTTPResponseFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//reqBody := r.Body
		var req *http.Request
		var err error
		var bodyJSON []byte
		bodyJSON, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Failed reading request")
		}
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
		if err != nil {
			log.Fatal("New Market Catalogue Request", err)
			return
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error: ", err)
			return
			//return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Println("Request failed: ", resp.StatusCode)
			return
			//return nil, fmt.Errorf("Request failed: %d", resp.StatusCode)
		}
		w.Header().Set("Content-Type", "application/json")
		var bs []byte
		bs, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			w.Write(bs)
		}
	}
}

func handleStocksDel(url string) common.HTTPResponseFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//reqBody := r.Body
		var req *http.Request
		var err error
		var bodyJSON []byte
		bodyJSON, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Failed reading request")
		}
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
		if err != nil {
			log.Fatal("New Market Catalogue Request", err)
			return
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error: ", err)
			return
			//return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Println("Request failed: ", resp.StatusCode)
			return
			//return nil, fmt.Errorf("Request failed: %d", resp.StatusCode)
		}
		w.Header().Set("Content-Type", "application/json")
		var bs []byte
		bs, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			w.Write(bs)
		}
	}
}

func main() {
	docDir := os.Getenv("WA_DOCUMENT_DIR")
	webHost := os.Getenv("WA_WEB_HOST") //hostname/ip where the web ui listens on
	webPort := os.Getenv("WA_WEB_PORT") //port where the web ui listens on
	crudHost := os.Getenv("CRUD_HOST")  //CRUD service connection details
	crudPort := os.Getenv("CRUD_PORT")
	spHost := os.Getenv("SP_HOST") //Stock provider connection details
	spPort := os.Getenv("SP_PORT")

	//webreport address
	webReportAddr := fmt.Sprintf("%s:%s", webHost, webPort)
	//configure paths
	docDirPath := filepath.ToSlash(docDir)
	//stock provider URLs "http://localhost:33001/stocksDataTable?symbol="
	//"http://localhost:33001/stocksProgressTable"
	spDataTableURL := fmt.Sprintf("http://%s:%s/stocksDataTable?symbol=", spHost, spPort)
	spProgressTableURL := fmt.Sprintf("http://%s:%s/stocksProgressTable", spHost, spPort)

	//crud url "http://localhost:33002/del"
	crudAddURL := fmt.Sprintf("http://%s:%s/add", crudHost, crudPort)
	crudDelURL := fmt.Sprintf("http://%s:%s/del", crudHost, crudPort)

	fs := http.FileServer(http.Dir(docDirPath))
	http.Handle("/", fs)
	dtHandleFunc := handleStocks(spDataTableURL)
	http.HandleFunc("/stocks", dtHandleFunc)
	addHandleFunc := handleStocksAdd(crudAddURL)
	http.HandleFunc("/stocks/add", addHandleFunc)
	delHandleFunc := handleStocksDel(crudDelURL)
	http.HandleFunc("/stocks/del", delHandleFunc)
	dashHandleFunc := handleDashStocks(spProgressTableURL)
	http.HandleFunc("/dashallstocks", dashHandleFunc)
	log.Println("Start listening on ", webReportAddr)
	log.Fatal(http.ListenAndServe(webReportAddr, nil))

}
