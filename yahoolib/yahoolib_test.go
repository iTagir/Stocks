package yahoolib

import (
	"log"
	"testing"
)

func TestYahooStockData(t *testing.T) {
	type args struct {
		symbol string
		data   *StockQueryResult
	}
	sqr := StockQueryResult{}
	var yahooErr bool
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.

		{
			name: "Test GSK.L",
			args: args{
				symbol: "GSK.L",
				data:   &sqr,
			},
			wantErr: yahooErr,
		},
		{
			name: "Test fsdfsfsdf",
			args: args{
				symbol: "fsdfsfsdf",
				data:   &sqr,
			},
			wantErr: yahooErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.data = &StockQueryResult{}
			if err := YahooStockData(tt.args.symbol, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("YahooStockData() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				log.Println("Ask: ", tt.args.data.Query.Result.Quote.Ask)
			}

		})
	}
}

func TestStringArr2HTML(t *testing.T) {
	type args struct {
		sarr []string
	}
	data := []string{"GSK.L", "VOD.L"}
	var out string = "%22GSK.L%22%2C%22VOD.L%22"
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "t1",
			args: args{
				sarr: data,
			},
			want: out,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringArr2HTML(tt.args.sarr); got != tt.want {
				t.Errorf("StringArr2HTML() = %v, want %v", got, tt.want)
			}
		})
	}
}
