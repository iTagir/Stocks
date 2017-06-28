package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func handleStocks(w http.ResponseWriter, r *http.Request) {
	stockSymbol := r.FormValue("symbol")
	log.Println("Requested symbol:", stockSymbol)
	resp, err := http.Get("http://localhost:33001/stocksDataTable?symbol=" + stockSymbol)
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

func handleStocksAdd(w http.ResponseWriter, r *http.Request) {
	//reqBody := r.Body
	var req *http.Request
	var err error
	var bodyJSON []byte
	bodyJSON, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Failed reading request")
	}
	req, err = http.NewRequest("POST", "http://localhost:33002/add", bytes.NewBuffer(bodyJSON))
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

func handleStocksDel(w http.ResponseWriter, r *http.Request) {
	//reqBody := r.Body
	var req *http.Request
	var err error
	var bodyJSON []byte
	bodyJSON, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Failed reading request")
	}
	req, err = http.NewRequest("POST", "http://localhost:33002/del", bytes.NewBuffer(bodyJSON))
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

func main() {
	docDir := os.Getenv("WA_DOCUMENT_DIR")
	webHost := os.Getenv("WA_WEB_HOST") //hostname/ip where the web ui listens on
	webPort := os.Getenv("WA_WEB_PORT") //port where the web ui listens on

	//webreport address
	webReportAddr := fmt.Sprintf("%s:%s", webHost, webPort)
	//configure paths
	docDirPath := filepath.ToSlash(docDir)

	fs := http.FileServer(http.Dir(docDirPath))
	http.Handle("/", fs)

	http.HandleFunc("/stocks", handleStocks)
	http.HandleFunc("/stocks/add", handleStocksAdd)
	http.HandleFunc("/stocks/del", handleStocksDel)

	log.Println("Start listening on ", webReportAddr)
	log.Fatal(http.ListenAndServe(webReportAddr, nil))

}
