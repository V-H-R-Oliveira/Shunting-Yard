package parser

import (
	"fmt"
	"reflect"
	"shunting-yard/m/tokenizer"
	"strings"
	"testing"
)

type testContext struct {
	expressionTokens []*tokenizer.Token
	parser           *Parser
	expected         string
}

func (testCtx *testContext) test(t *testing.T) {
	result := testCtx.parser.Parse(testCtx.expressionTokens)
	resultStr := strings.Join(result, "")

	if !reflect.DeepEqual(resultStr, testCtx.expected) {
		t.Fatalf("Expected token %s. Got %s\n", testCtx.expected, resultStr)
	}
}

func newTestContext(expressionTokens []*tokenizer.Token, expected string) *testContext {
	return &testContext{
		expressionTokens: expressionTokens,
		parser:           NewParser(),
		expected:         expected,
	}
}

func testExpression(expression, expected string, t *testing.T) {
	tokenizer := tokenizer.NewTokenizer()
	expressionTokens := tokenizer.Tokenize([]byte(expression))
	testCtx := newTestContext(expressionTokens, expected)
	testCtx.test(t)
}

func TestParser_Parse(t *testing.T) {
	expressions := map[string]string{
		"5-2":                            "52-",
		"1 / 5 + 2 * 4 - 3 ^ 10 / 4 ^ 2": "15/24*+310^42^/-",
		"1 + 2 * ( 4 - 3 * (2 + 5) )":    "124325+*-*+",
		"(1+	2)":                         "12+",
		"1* (2 + 4) / 3":                 "124+3/*",
		"1 + 2 * (4 - 3)":                "1243-*+",
		"(1 + 2) * (4 - 3)":              "12+43-*",
		"(1 + 2) * (4-3)":                "12+43-*",
		"1 + 2 * (4 - 3) ^ 2 * (1 + 5)":  "1243-2^15+**+",
		"(1 + 2) * (((4 - 3) - 2))":      "12+43-2-*",
		"(1 + 2 + (6 - 7)) * ((4 - 3) + ((2 + 5) + 1))": "1267-++43-25+1++*",
	}

	for expression, expected := range expressions {
		testLabel := fmt.Sprintf("Parsing expression \"%s\"", expression)

		t.Run(testLabel, func(t *testing.T) {
			testExpression(expression, expected, t)
		})
	}
}
