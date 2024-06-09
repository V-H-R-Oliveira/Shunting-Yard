package main

import (
	"fmt"
	"shunting-yard/m/parser"
	"shunting-yard/m/tokenizer"
	"strings"
)

func main() {
	expressions := []string{
		"-10",
		"(-10)",
		"-10 + 5",
		"-10 * - 5",
		"-10 - 5",
		"9 - (3 * (-6))",
		"9 - 3 * -6",
		"5-2",
		"-5-2",
		"10--5",
		"5 * -5",
		"1 / 5 + 2 * 4 - 3 ^ 10 / 4 ^ 2",
		"1 + 2 * ( 4 - 3 * (2 + 5) )",
		"1* (2 + 4) / 3",
		"1 + 2 * (4 - 3)",
		"1 + 2 * (4 - 3) ^ 2 * (1 + 5)",
		"(1 + 2) * (4 - 3)",
		"(1 + 2) * (4-3)",
		"(1 + 2) * (((4 - 3) - 2))",
		"-10 + 5 - -1",
		"(1 + 2 + (6 - 7)) * ((4 - 3) + ((2 + 5) + 1))",
		"(1 + 2) * (4 - 3)",
		"1 + 2 *-4 - 3",
		"-((1 + 2) / ((6 * -7) + (7 * -4) / 2) - 3)",
		"-((1 + 2) / ((6 * -7) + (7 * -4) / 2) -3)",
	}

	for _, expression := range expressions {
		tokenizer := tokenizer.NewTokenizer()
		tokens := tokenizer.Tokenize([]byte(expression))

		parser := parser.NewParser()
		rpnExpr := parser.Parse(tokens)

		fmt.Printf("shuntingYard(%s) = %s\n", expression, strings.Join(rpnExpr, ""))
	}
}
