package mdb

import (
	"log"
	"testing"
)

func TestMongoDBConn_Stock(t *testing.T) {

	type args struct {
		symbol string
		data   []StockData
	}

	d := make([]StockData, 0, 10)
	tests := []struct {
		name  string
		mconn *MongoDBConn
		args  args
	}{
		// TODO: Add test cases.
		{"GSK.L",
			&MongoDBConn{"tagir-tosh", "test", "testcoll"},
			args{"GSK.L", d}},
	}
	//log.Println("test started")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mconn.Stock(tt.args.symbol, &tt.args.data)
			log.Println("data:", tt.args.data)
		})
	}
}

func TestMongoDBConn_StockByID(t *testing.T) {
	type args struct {
		id   string
		data *[]StockData
	}
	d := make([]StockData, 0, 10)
	tests := []struct {
		name  string
		mconn *MongoDBConn
		args  args
	}{
		// TODO: Add test cases.
		{"GSK.L",
			&MongoDBConn{"tagir-tosh", "test", "testcoll"},
			args{"5949291b95b95e084bd134b9", &d}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mconn.StockByID(tt.args.id, tt.args.data)
		})
		log.Println("data=", tt.args.data)
	}
}

func TestMongoDBConn_RemoveById(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name  string
		mconn *MongoDBConn
		args  args
	}{
		// TODO: Add test cases.
		{"GSK.L",
			&MongoDBConn{"tagir-tosh", "test", "testcoll"},
			args{id: "5949291b95b95e084bd134b9"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mconn.RemoveByID(tt.args.id)
		})
	}
}

func TestMongoDBConn_StockUniqueSymbols(t *testing.T) {
	type args struct {
		data *[]string
	}
	d := make([]string, 0, 100)
	tests := []struct {
		name  string
		mconn *MongoDBConn
		args  args
	}{
		// TODO: Add test cases.
		{"GSK.L",
			&MongoDBConn{"tagir-tosh", "test", "testcoll"},
			args{&d}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mconn.StockUniqueSymbols(tt.args.data)
		})
		a := *tt.args.data
		log.Println("data=", a, tt.args.data)
	}
}
