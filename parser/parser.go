package parser

import (
	"log"
	"shunting-yard/m/tokenizer"
)

type Parser struct {
	output          []string
	stack           []*tokenizer.Token
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
		},
		stack:  []*tokenizer.Token{},
		output: []string{},
	}
}

// Push operator to the stack using the mathematical operator precedence.
func (parser *Parser) pushToStack(operator *tokenizer.Token) {
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

		if stackedOperatorPrecedence > operatorPrecedence {
			parser.output = append(parser.output, string(stackedOperator.Lexeme))
			parser.stack = parser.stack[:stackTop]
		}
	}
}

// convert an infix expression to a RPN expression using the Shunting Yard algorithm
func (parser *Parser) Parse(tokens []*tokenizer.Token) []string {
	for _, token := range tokens {

		if token.TokenType == tokenizer.NUMBER {
			parser.output = append(parser.output, string(token.Lexeme))
			continue
		} else {
			parser.pushToStack(token)
		}
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
