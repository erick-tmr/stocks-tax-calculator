package stocks

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Operation struct {
	Type     string          `json:"operation"`
	Price    decimal.Decimal `json:"unit-cost"`
	Quantity decimal.Decimal `json:"quantity"`
}

type Tax struct {
	Value decimal.Decimal
}

func (t Tax) MarshalJSON() ([]byte, error) {
	str_repr := fmt.Sprintf("{\"tax\": %v}", t.Value.StringFixed(2))
	return []byte(str_repr), nil
}

func (t *Tax) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}

	parsed, _ := decimal.NewFromString(string(b))
	*t = Tax{parsed}

	return nil
}

type Calculator struct {
	Operations     []Operation
	Taxes          []Tax
	holding_stocks decimal.Decimal
	actual_avg     decimal.Decimal
	acc_loss       decimal.Decimal
}

var sell_limit_price, _ = decimal.NewFromString("20000")
var tax_percentage, _ = decimal.NewFromString("0.2")

func (c *Calculator) AddOperation(op Operation) {
	c.Operations = append(c.Operations, op)
}

func (c *Calculator) add_tax(tax Tax) {
	c.Taxes = append(c.Taxes, tax)
}

func (c *Calculator) CalculateTaxes() {
	for _, op := range c.Operations {
		var calculated_tax Tax
		c.holding_stocks, c.actual_avg, c.acc_loss, calculated_tax = proccess_operation(
			op,
			c.actual_avg,
			c.holding_stocks,
			c.acc_loss,
		)
		c.add_tax(calculated_tax)
	}
}

func proccess_operation(
	op Operation,
	actual_avg decimal.Decimal,
	holding_stocks decimal.Decimal,
	acc_loss decimal.Decimal,
) (decimal.Decimal, decimal.Decimal, decimal.Decimal, Tax) {
	var new_holding_stocks decimal.Decimal
	var new_avg decimal.Decimal
	var new_acc_loss decimal.Decimal
	var calculated_tax Tax

	switch op.Type {
	case "buy":
		new_holding_stocks = holding_stocks.Add(op.Quantity)
		new_avg = BuyPrice(holding_stocks, actual_avg, op.Quantity, op.Price)
		calculated_tax = Tax{decimal.Zero}
	case "sell":
		new_holding_stocks = holding_stocks.Sub(op.Quantity)
		new_avg = actual_avg
		op_total_value := op.Quantity.Mul(op.Price)
		has_profit := op.Price.Sub(actual_avg)

		if has_profit.IsZero() {
			calculated_tax = Tax{decimal.Zero}
			new_acc_loss = acc_loss
		} else if has_profit.IsPositive() {
			if op_total_value.LessThanOrEqual(sell_limit_price) {
				calculated_tax = Tax{decimal.Zero}

				if acc_loss.IsPositive() {
					profit := has_profit.Mul(op.Quantity)
					new_acc_loss = calculate_loss(acc_loss, profit)
				}
			} else {
				if acc_loss.IsPositive() {
					profit := has_profit.Mul(op.Quantity)
					new_acc_loss = acc_loss.Sub(profit)

					if new_acc_loss.IsNegative() {
						real_profit := new_acc_loss.Abs()
						tax := calculate_tax(real_profit)
						calculated_tax = Tax{tax}
						new_acc_loss = decimal.Zero
					} else {
						calculated_tax = Tax{decimal.Zero}
					}
				} else {
					new_acc_loss = acc_loss
					profit := has_profit.Mul(op.Quantity)
					tax := calculate_tax(profit)
					calculated_tax = Tax{tax}
				}
			}
		} else {
			calculated_tax = Tax{decimal.Zero}
			profit := has_profit.Mul(op.Quantity)
			new_acc_loss = acc_loss.Add(profit.Abs())
		}
	}

	return new_holding_stocks, new_avg, new_acc_loss, calculated_tax
}

func calculate_tax(profit decimal.Decimal) decimal.Decimal {
	return profit.Mul(tax_percentage)
}

func calculate_loss(acc_loss decimal.Decimal, profit decimal.Decimal) decimal.Decimal {
	new_acc_loss := acc_loss.Sub(profit)

	if new_acc_loss.IsNegative() {
		new_acc_loss = decimal.Zero
	}

	return new_acc_loss
}
