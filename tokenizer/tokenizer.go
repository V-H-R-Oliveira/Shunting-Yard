package tokenizer

import "log"

type TokenType int

const (
	NUMBER TokenType = iota
	PLUS
	MINUS
	STAR
	SLASH
	POWER
	LEFT_PARENTHESIS
	RIGHT_PARENTHESIS
)

type Token struct {
	Lexeme    string
	TokenType TokenType
}

func NewToken(lexeme string, tokenType TokenType) *Token {
	return &Token{
		Lexeme:    lexeme,
		TokenType: tokenType,
	}
}

type Tokenizer struct {
	tokens []*Token
}

func (tokenizer *Tokenizer) Tokenize(expression []byte) []*Token {
	current := 0
	end := false

	var token *Token
	expressionLength := len(expression)

	for !end {
		char := expression[current]

		if isNumber(char) {
			token, continueAt := numberToken(expression, current)
			current = continueAt
			end = continueAt == 0
			tokenizer.tokens = append(tokenizer.tokens, token)
			continue
		}

		switch char {
		case '+':
			token = NewToken(string(char), PLUS)
		case '-':
			lookaheadChar := lookahead(expression, current)
			tokensLength := len(tokenizer.tokens)

			// example: 4-3 == 4 - 3, should result 4,-,3
			// example: -3 should result -3 and --3 should result -, -3
			isLastTokenNumber := tokensLength > 0 && tokenizer.tokens[tokensLength-1].TokenType == NUMBER

			if isNumber(lookaheadChar) && (!isLastTokenNumber || tokensLength == 0) {
				token, continueAt := numberToken(expression, current)
				current = continueAt
				end = continueAt == 0
				tokenizer.tokens = append(tokenizer.tokens, token)
				continue
			} else {
				token = NewToken(string(char), MINUS)
			}
		case '*':
			token = NewToken(string(char), STAR)
		case '/':
			token = NewToken(string(char), SLASH)
		case '^':
			token = NewToken(string(char), POWER)
		case '(':
			token = NewToken(string(char), LEFT_PARENTHESIS)
		case ')':
			token = NewToken(string(char), RIGHT_PARENTHESIS)
		default:
			token = nil
		}

		if token != nil {
			tokenizer.tokens = append(tokenizer.tokens, token)
		}

		current++

		if current >= expressionLength {
			end = true
		}
	}

	return tokenizer.tokens
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{
		tokens: []*Token{},
	}
}

func isNumber(char byte) bool {
	return char >= '0' && char <= '9'
}

// parse a stream of characters to an integer number
func numberToken(expression []byte, start int) (*Token, int) {
	continueAt := 0
	number := []byte{}

	firstChar := expression[start]

	if firstChar == '-' {
		number = append(number, firstChar)
		start++
	}

	for idx, char := range expression[start:] {
		if isNumber(char) {
			number = append(number, char)
			continue
		}

		continueAt = start + idx
		break
	}

	return NewToken(string(number), NUMBER), continueAt
}

// look one character ahead of the current pointer
func lookahead(expression []byte, current int) byte {
	expressionLength := len(expression)

	if (current + 1) > expressionLength {
		log.Fatal("Tried to lookahead after the end of the expression")
	}

	return expression[current+1]
}