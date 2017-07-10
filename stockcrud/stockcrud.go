package main

import (
	"encoding/json"
	"fmt"
	"github.com/iTagir/stocks/common"
	"github.com/iTagir/stocks/mdb"
	"log"
	"net/http"
	"os"
)

func addStock(dbhost string, dbname string, dbcoll string) common.HTTPResponseFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type resp struct {
			Status string `json:"status"`
		}
		d := resp{Status: "OK"}
		w.Header().Set("Content-Type", "application/json")

		b := mdb.StockData{}
		err := common.ParseJSONBody(r.Body, &b)
		if err != nil {
			log.Println("Request parse error ", err)
			d.Status = "FAILED"
			//return err
		} else {
			log.Println("Inserted symbol:", b.Symbol)
			mdbConn := mdb.CreateMongoDBConn(dbhost, dbname, dbcoll)
			err = mdbConn.Place2(b)
			if err != nil {
				log.Println("DB insert error ", err)
				d.Status = "FAILED"
			}
		}
		err = json.NewEncoder(w).Encode(d)
		if err != nil {
			log.Println("addStock response encode error: ", err)
		}
	}
}

func delStock(dbhost string, dbname string, dbcoll string) common.HTTPResponseFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type resp struct {
			Status string `json:"status"`
		}
		d := resp{Status: "OK"}
		w.Header().Set("Content-Type", "application/json")

		type delReq struct {
			Id string `json:"id"`
		}
		b := delReq{}
		err := common.ParseJSONBody(r.Body, &b)
		if err != nil {
			log.Println("Request parse error ", err)
			d.Status = "FAILED"
			//return err
		} else {
			log.Println("Deleted symbol:", b.Id)
			mdbConn := mdb.CreateMongoDBConn(dbhost, dbname, dbcoll)
			err = mdbConn.RemoveByID(b.Id)
			if err != nil {
				log.Println("DB insert error ", err)
				d.Status = "FAILED"
			}
		}
		err = json.NewEncoder(w).Encode(d)
		if err != nil {
			log.Println("delStock response encode error: ", err)
		}
	}
}

func main() {

	host := os.Getenv("CRUD_HOST")
	port := os.Getenv("CRUD_PORT")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoDB := os.Getenv("MONGO_DB")
	mongoColl := os.Getenv("MONGO_COLLECTION")

	if port == "" {
		log.Fatal("Port variable CRUD_PORT was not set.")
		return
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	handle := addStock(mongoHost, mongoDB, mongoColl)
	http.HandleFunc("/add", handle)
	handle = delStock(mongoHost, mongoDB, mongoColl)
	http.HandleFunc("/del", handle)
	log.Println("Start listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
