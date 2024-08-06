package stocks

import (
	"github.com/shopspring/decimal"
)

func BuyPrice(
	holding_stocks decimal.Decimal,
	actual_avg decimal.Decimal,
	bought_stocks decimal.Decimal,
	bought_price decimal.Decimal,
) decimal.Decimal {
	if holding_stocks.IsZero() {
		return bought_price
	}

	stocks_total := holding_stocks.Add(bought_stocks)
	return holding_stocks.Mul(actual_avg).Add(bought_stocks.Mul(bought_price)).Div(stocks_total)
}
