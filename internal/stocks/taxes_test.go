package stocks

import (
	"testing"

	"github.com/shopspring/decimal"
)

func Test_sell_limit_price(t *testing.T) {
	expected, _ := decimal.NewFromString("20000")
	if !sell_limit_price.Equal(expected) {
		t.Errorf("incorrect value for constant 'sell_limit_price', expected %v, got %v", expected, sell_limit_price)
	}
}

func Test_tax_percentage(t *testing.T) {
	expected := decimal_from_string("0.2")
	if !tax_percentage.Equal(expected) {
		t.Errorf("incorrect value for constant 'tax_percentage', expected %v, got %v", expected, tax_percentage)
	}
}

func Test_AddOperation(t *testing.T) {
	calc := Calculator{}
	if len(calc.Operations) != 0 {
		t.Errorf("Operations field is expected to be empty on initialization")
	}

	op_price := decimal_from_string("20.00")
	op_quantity := decimal_from_string("1000")
	add_op := Operation{
		Type:     "buy",
		Price:    op_price,
		Quantity: op_quantity,
	}
	calc.AddOperation(add_op)
	if len(calc.Operations) != 1 {
		t.Errorf("AddOperation should add the element")
	}
	if calc.Operations[0].Type != "buy" {
		t.Errorf("AddOperation added wrong Operation")
	}
}

func Test_add_tax(t *testing.T) {
	calc := Calculator{}
	if len(calc.Taxes) != 0 {
		t.Errorf("Taxes field is expected to be empty on initialization")
	}

	tax_value := decimal_from_string("1000")
	add_tax := Tax{tax_value}
	calc.add_tax(add_tax)
	if len(calc.Taxes) != 1 {
		t.Errorf("add_tax should add the element")
	}
	if !calc.Taxes[0].Value.Equal(tax_value) {
		t.Errorf("add_tax added wrong Tax")
	}
}

func Test_CalculateTaxes(t *testing.T) {
	data := []struct {
		test_name       string
		operations      []Operation
		expected_result []Tax
	}{
		{
			"Case 01",
			[]Operation{
				{"buy", decimal_from_string("10.00"), decimal_from_string("100")},
				{"sell", decimal_from_string("15.00"), decimal_from_string("50")},
				{"sell", decimal_from_string("15.00"), decimal_from_string("50")},
			},
			[]Tax{{decimal.Zero}, {decimal.Zero}, {decimal.Zero}},
		},
		{
			"Case 02",
			[]Operation{
				{"buy", decimal_from_string("10.00"), decimal_from_string("10000")},
				{"sell", decimal_from_string("20.00"), decimal_from_string("5000")},
				{"sell", decimal_from_string("5.00"), decimal_from_string("5000")},
			},
			[]Tax{{decimal.Zero}, {decimal_from_string("10000")}, {decimal.Zero}},
		},
		{
			"Case 03",
			[]Operation{
				{"buy", decimal_from_string("10.00"), decimal_from_string("10000")},
				{"sell", decimal_from_string("5.00"), decimal_from_string("5000")},
				{"sell", decimal_from_string("20.00"), decimal_from_string("3000")},
			},
			[]Tax{{decimal.Zero}, {decimal.Zero}, {decimal_from_string("1000")}},
		},
		{
			"Case 04",
			[]Operation{
				{"buy", decimal_from_string("10.00"), decimal_from_string("10000")},
				{"buy", decimal_from_string("25.00"), decimal_from_string("5000")},
				{"sell", decimal_from_string("15.00"), decimal_from_string("10000")},
			},
			[]Tax{{decimal.Zero}, {decimal.Zero}, {decimal.Zero}},
		},
		{
			"Case 05",
			[]Operation{
				{"buy", decimal_from_string("10.00"), decimal_from_string("10000")},
				{"buy", decimal_from_string("25.00"), decimal_from_string("5000")},
				{"sell", decimal_from_string("15.00"), decimal_from_string("10000")},
				{"sell", decimal_from_string("25.00"), decimal_from_string("5000")},
			},
			[]Tax{{decimal.Zero}, {decimal.Zero}, {decimal.Zero}, {decimal_from_string("10000")}},
		},
		{
			"Case 06",
			[]Operation{
				{"buy", decimal_from_string("10.00"), decimal_from_string("10000")},
				{"sell", decimal_from_string("2.00"), decimal_from_string("5000")},
				{"sell", decimal_from_string("20.00"), decimal_from_string("2000")},
				{"sell", decimal_from_string("20.00"), decimal_from_string("2000")},
				{"sell", decimal_from_string("25.00"), decimal_from_string("1000")},
			},
			[]Tax{{decimal.Zero}, {decimal.Zero}, {decimal.Zero}, {decimal.Zero}, {decimal_from_string("3000")}},
		},
		{
			"Case 07",
			[]Operation{
				{"buy", decimal_from_string("10.00"), decimal_from_string("10000")},
				{"sell", decimal_from_string("2.00"), decimal_from_string("5000")},
				{"sell", decimal_from_string("20.00"), decimal_from_string("2000")},
				{"sell", decimal_from_string("20.00"), decimal_from_string("2000")},
				{"sell", decimal_from_string("25.00"), decimal_from_string("1000")},
				{"buy", decimal_from_string("20.00"), decimal_from_string("10000")},
				{"sell", decimal_from_string("15.00"), decimal_from_string("5000")},
				{"sell", decimal_from_string("30.00"), decimal_from_string("4350")},
				{"sell", decimal_from_string("30.00"), decimal_from_string("650")},
			},
			[]Tax{
				{decimal.Zero},
				{decimal.Zero},
				{decimal.Zero},
				{decimal.Zero},
				{decimal_from_string("3000")},
				{decimal.Zero},
				{decimal.Zero},
				{decimal_from_string("3700")},
				{decimal.Zero},
			},
		},
		{
			"Case 08",
			[]Operation{
				{"buy", decimal_from_string("10.00"), decimal_from_string("10000")},
				{"sell", decimal_from_string("50.00"), decimal_from_string("10000")},
				{"buy", decimal_from_string("20.00"), decimal_from_string("10000")},
				{"sell", decimal_from_string("50.00"), decimal_from_string("10000")},
			},
			[]Tax{{decimal.Zero}, {decimal_from_string("80000")}, {decimal.Zero}, {decimal_from_string("60000")}},
		},
	}

	for _, d := range data {
		t.Run(d.test_name, func(t *testing.T) {
			calc := Calculator{}
			for _, op := range d.operations {
				calc.AddOperation(op)
			}

			calc.CalculateTaxes()

			for index, tax := range d.expected_result {
				if !calc.Taxes[index].Value.Equal(tax.Value) {
					t.Errorf("Unexpected Tax calculated %+v", calc.Taxes[index])
				}
			}
		})
	}
}

func decimal_from_string(str string) decimal.Decimal {
	dec, _ := decimal.NewFromString(str)
	return dec
}
