package tokenizer

import (
	"reflect"
	"testing"
)

type testContext struct {
	expression []byte
	tokenizer  *Tokenizer
	expected   []*Token
}

func (testCtx *testContext) test(t *testing.T) {
	tokens := testCtx.tokenizer.Tokenize(testCtx.expression)

	if len(tokens) != len(testCtx.expected) {
		t.Fatal("Result has different length from the expected\n")
	}

	for i, token := range tokens {
		expectedToken := testCtx.expected[i]

		if !reflect.DeepEqual(expectedToken.TokenType, token.TokenType) || !reflect.DeepEqual(expectedToken.Lexeme, token.Lexeme) {
			t.Fatalf("Expected token %v. Got %v\n", expectedToken, token)
		}
	}
}

func newTestContext(expression string, expected []*Token) *testContext {
	return &testContext{
		expression: []byte(expression),
		tokenizer:  NewTokenizer(),
		expected:   expected,
	}
}

func testLiteral(expression string, tokenType TokenType, t *testing.T) {
	testCtx := newTestContext(expression, []*Token{NewToken(expression, tokenType)})
	testCtx.test(t)
}

func testExpression(expression string, expectedTokens []*Token, t *testing.T) {
	testCtx := newTestContext(expression, expectedTokens)
	testCtx.test(t)
}

func TestTokenizer_Tokenize(t *testing.T) {
	t.Run("Tokenize numbers", func(t *testing.T) {
		testLiteral("10", NUMBER, t)
		testLiteral("2000", NUMBER, t)
		testLiteral("100", NUMBER, t)
	})

	t.Run("Tokenize operators", func(t *testing.T) {
		operators := map[string]TokenType{
			"+": PLUS,
			"-": MINUS,
			"*": STAR,
			"/": SLASH,
			"^": POWER,
		}

		for lexeme, tokenType := range operators {
			testLiteral(lexeme, tokenType, t)
		}
	})

	t.Run("Tokenize parenthesis", func(t *testing.T) {
		parenthesis := map[string]TokenType{
			"(": LEFT_PARENTHESIS,
			")": RIGHT_PARENTHESIS,
		}

		for lexeme, tokenType := range parenthesis {
			testLiteral(lexeme, tokenType, t)
		}
	})

	t.Run("Tokenize expressions", func(t *testing.T) {
		expressions := map[string][]*Token{
			"1 + 2": {NewToken("1", NUMBER), NewToken("+", PLUS), NewToken("2", NUMBER)},
			"1+2":   {NewToken("1", NUMBER), NewToken("+", PLUS), NewToken("2", NUMBER)},
			"1+  2": {NewToken("1", NUMBER), NewToken("+", PLUS), NewToken("2", NUMBER)},
			"4 * 5 + (5 ^ (5 - 2) / 122131)": {
				NewToken("4", NUMBER), NewToken("*", STAR), NewToken("5", NUMBER),
				NewToken("+", PLUS), NewToken("(", LEFT_PARENTHESIS), NewToken("5", NUMBER),
				NewToken("^", POWER), NewToken("(", LEFT_PARENTHESIS), NewToken("5", NUMBER),
				NewToken("-", MINUS), NewToken("2", NUMBER), NewToken(")", RIGHT_PARENTHESIS),
				NewToken("/", SLASH), NewToken("122131", NUMBER), NewToken(")", RIGHT_PARENTHESIS),
			},
		}

		for expression, expected := range expressions {
			testExpression(expression, expected, t)
		}
	})
}
