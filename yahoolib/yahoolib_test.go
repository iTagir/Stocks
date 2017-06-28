package yahoolib

import "testing"
import "log"

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
