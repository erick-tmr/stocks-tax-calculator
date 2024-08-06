package stocks

import (
	"testing"

	"github.com/shopspring/decimal"
)

func Test_BuyPrice(t *testing.T) {
	holding_stocks, _ := decimal.NewFromString("10000")
	actual_avg, _ := decimal.NewFromString("10.00")
	bought_stocks, _ := decimal.NewFromString("5000")
	bought_price, _ := decimal.NewFromString("25.00")

	result := BuyPrice(holding_stocks, actual_avg, bought_stocks, bought_price)

	expected_price, _ := decimal.NewFromString("15.00")
	if !result.Equal(expected_price) {
		t.Errorf("incorrect result, expected %v, got %v", expected_price, result)
	}
}
