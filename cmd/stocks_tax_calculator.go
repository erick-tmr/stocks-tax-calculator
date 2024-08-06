package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/erick-tmr/stocks-tax-calculator/internal/stocks"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		read_line := scanner.Text()

		var json_read []stocks.Operation
		json.Unmarshal([]byte(read_line), &json_read)

		calculator := stocks.Calculator{}
		for _, op := range json_read {
			calculator.AddOperation(op)
		}
		calculator.CalculateTaxes()

		var response []byte
		response, _ = json.Marshal(calculator.Taxes)

		fmt.Println(string(response))
	}
}
