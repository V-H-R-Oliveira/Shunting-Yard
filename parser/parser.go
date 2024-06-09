package parser

import (
	"log"
	"shunting-yard/m/tokenizer"
)

// TODO: Write tests
const UNARY_LEXEME = "n"

type Parser struct {
	output          []string
	stack           []*tokenizer.Token
	previousToken   *tokenizer.Token
	precedenceTable map[tokenizer.TokenType]int
}

func NewParser() *Parser {
	return &Parser{
		precedenceTable: map[tokenizer.TokenType]int{
			tokenizer.MINUS: 1,
			tokenizer.PLUS:  2,
			tokenizer.STAR:  3,
			tokenizer.SLASH: 4,
			tokenizer.POWER: 5,
			tokenizer.UNARY: 6, // unary '-'
		},
		previousToken: nil,
		stack:         []*tokenizer.Token{},
		output:        []string{},
	}
}

func (parser *Parser) pushUnary(operator *tokenizer.Token) {
	operator.TokenType = tokenizer.UNARY
	operator.Lexeme = UNARY_LEXEME
	parser.stack = append(parser.stack, operator)
}

func (parser *Parser) isOperator(token *tokenizer.Token) bool {
	if token != nil {
		isParenthesis := token.TokenType == tokenizer.RIGHT_PARENTHESIS || token.TokenType == tokenizer.LEFT_PARENTHESIS
		return token.TokenType != tokenizer.NUMBER && !isParenthesis
	}

	return false
}

func (parser *Parser) isUnary(next *tokenizer.Token) bool {
	if parser.previousToken != nil && parser.previousToken.TokenType == tokenizer.LEFT_PARENTHESIS && next.TokenType == tokenizer.NUMBER {
		return true
	}

	return parser.isOperator(parser.previousToken) && next.TokenType == tokenizer.NUMBER
}

// Push operator to the stack using the mathematical operator precedence.
func (parser *Parser) pushToStack(operator, next *tokenizer.Token) {
	if next != nil {
		// is the begin and next element is a number or a left parenthesis
		if parser.previousToken == nil && (next.TokenType == tokenizer.NUMBER || next.TokenType == tokenizer.LEFT_PARENTHESIS) {
			parser.pushUnary(operator)
			return
		}

		if parser.isUnary(next) {
			parser.pushUnary(operator)
			return
		}
	}

	if len(parser.stack) == 0 {
		parser.stack = append(parser.stack, operator)
		return
	}

	if operator.TokenType == tokenizer.LEFT_PARENTHESIS {
		parser.stack = append(parser.stack, operator)
		return
	}

	if operator.TokenType == tokenizer.RIGHT_PARENTHESIS {
		for {
			stackTop := len(parser.stack) - 1

			if stackTop < 0 {
				parser.stack = append(parser.stack, operator)
				return
			}

			stackedOperator := parser.stack[stackTop]

			if stackedOperator.TokenType == tokenizer.LEFT_PARENTHESIS {
				parser.stack = parser.stack[:stackTop]
				return
			}

			parser.output = append(parser.output, string(stackedOperator.Lexeme))
			parser.stack = parser.stack[:stackTop]
		}
	}

	operatorPrecedence := parser.precedenceTable[operator.TokenType]

	for {
		stackTop := len(parser.stack) - 1

		if stackTop < 0 {
			parser.stack = append(parser.stack, operator)
			return
		}

		stackedOperator := parser.stack[stackTop]
		stackedOperatorPrecedence := parser.precedenceTable[stackedOperator.TokenType]

		if operatorPrecedence >= stackedOperatorPrecedence {
			parser.stack = append(parser.stack, operator)
			return
		}

		parser.output = append(parser.output, string(stackedOperator.Lexeme))
		parser.stack = parser.stack[:stackTop]
	}
}

// convert an infix expression to a RPN expression using the Shunting Yard algorithm
func (parser *Parser) Parse(tokens []*tokenizer.Token) []string {
	for idx, token := range tokens {
		if token.TokenType == tokenizer.NUMBER {
			parser.output = append(parser.output, string(token.Lexeme))
		} else if token.TokenType == tokenizer.MINUS {
			nextToken := tokens[idx+1]
			parser.pushToStack(token, nextToken)
		} else {
			parser.pushToStack(token, nil)
		}

		parser.previousToken = token
	}

	if len(parser.stack) > 0 {
		for {
			stackTop := len(parser.stack) - 1

			if stackTop < 0 {
				break
			}

			element := parser.stack[stackTop]

			if element.TokenType == tokenizer.LEFT_PARENTHESIS || element.TokenType == tokenizer.RIGHT_PARENTHESIS {
				log.Fatal("Syntax error")
			}

			parser.output = append(parser.output, string(element.Lexeme))
			parser.stack = parser.stack[:stackTop]
		}
	}

	return parser.output
}
