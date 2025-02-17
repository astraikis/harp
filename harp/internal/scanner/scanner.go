package scanner

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/astraikis/harp/internal/models"
)

var source string
var tokens = []models.Token{}
var start int = 0
var current int = 0
var line int = 1
var column int = 1

var keywords = map[string]models.TokenType{
	"and":    models.AND,
	"else":   models.ELSE,
	"false":  models.FALSE,
	"func":   models.FUNC,
	"for":    models.FOR,
	"if":     models.IF,
	"null":   models.NULL,
	"or":     models.OR,
	"return": models.RETURN,
	"true":   models.TRUE,
	"while":  models.WHILE,
	"struct": models.STRUCT,
	"string": models.STRING_VAR,
	"int":    models.INT_VAR,
	"double": models.DOUBLE_VAR,
	"bool":   models.BOOL_VAR,
}

// scan walks over source and returns
// a corresponding list of tokens.
func Scan(sourceText string) []models.Token {
	source = sourceText
	for {
		if isAtEnd() {
			break
		}
		start = current
		scanToken()
	}

	tokens = append(tokens, models.Token{Type: models.EOF, Lexeme: "", Literal: nil, Column: 1, Line: line})
	return tokens
}

// scanToken adds the next token to tokens.
func scanToken() {
	var c rune = advance()
	switch c {
	case '(':
		addToken(models.LeftParen, "")
		column += 1
	case ')':
		addToken(models.RightParen, "")
		column += 1
	case '{':
		addToken(models.LeftBrace, "")
		column += 1
	case '}':
		addToken(models.RightBrace, "")
		column += 1
	case ',':
		addToken(models.COMMA, "")
		column += 1
	case '.':
		addToken(models.DOT, "")
		column += 1
	case '-':
		addToken(models.MINUS, "")
		column += 1
	case '+':
		addToken(models.PLUS, "")
		column += 1
	case ';':
		addToken(models.SEMICOLON, "")
		column += 1
	case '*':
		addToken(models.STAR, "")
		column += 1
	case '!':
		if match('=') {
			addToken(models.BANG_EQUAL, "")
			column += 2
		} else {
			addToken(models.BANG, "")
			column += 1
		}
	case '=':
		if match('=') {
			addToken(models.EQUAL_EQUAL, "")
			column += 2
		} else {
			addToken(models.EQUAL, "")
			column += 1
		}
	case '>':
		if match('=') {
			addToken(models.GREATER_EQUAL, "")
			column += 2
		} else {
			addToken(models.GREATER, "")
			column += 1
		}
	case '<':
		if match('=') {
			addToken(models.LESS_EQUAL, "")
			column += 2
		} else {
			addToken(models.LESS, "")
			column += 1
		}
	case '/':
		if match('/') {
			for {
				if peek() == '\n' || isAtEnd() {
					break
				}
				advance()
			}
			column = 0
		} else {
			addToken(models.SLASH, "")
		}
	case ' ':
		column += 1
	case '\r':
	case '\t':
	case '\n':
		line += 1
		column = 1
	case '"':
		_string()
		column += 1
	default:
		if unicode.IsDigit(c) {
			number()
		} else if unicode.IsLetter(c) || c == '_' {
			identifier()
		} else {
			if column != 1 {
				column += 1
			}
			// Error
		}
	}
}

// _string adds the next string to tokens.
func _string() {
	for {
		if peek() == '"' || isAtEnd() {
			break
		}
		if peek() == '\n' {
			line += 1
		}
		advance()
	}

	if isAtEnd() {
	}

	advance()

	addToken(models.STRING, source[start:current])
}

// addToken adds a token to tokens.
func addToken[T any](tokenType models.TokenType, literal T) {
	tokens = append(tokens, models.Token{Type: tokenType, Lexeme: string(source[start:current]), Literal: literal, Column: column, Line: line})
}

// advance consumes and returns the next rune.
func advance() rune {
	next := rune(source[current])
	current += 1
	return next
}

// identifier adds the next identifier or keyword to tokens.
func identifier() {
	for {
		next := peek()
		if !unicode.IsLetter(next) && !unicode.IsDigit(next) && next != '_' {
			break
		}

		advance()
	}

	text := source[start:current]
	_type, exists := keywords[text]
	if !exists {
		_type = models.IDENTIFIER
	}
	addToken(_type, text)
}

// isAtEnd reports whether current is
// at end of source.
func isAtEnd() bool {
	return current >= len(source)
}

// match reports if expected is equal to the
// current rune and consumes it if true.
func match(expected rune) bool {
	if isAtEnd() {
		return false
	}

	if source[current] != byte(expected) {
		return false
	}

	current += 1
	return true
}

// number adds the next number to tokens.
func number() {
	for {
		if !unicode.IsDigit(peek()) {
			break
		}
		advance()
	}

	isDouble := false
	if peek() == '.' {
		advance()
		isDouble = true

		for {
			if !unicode.IsDigit(peek()) {
				break
			}
			advance()
		}
	}

	if isDouble {
		val, _ := strconv.ParseFloat(source[start:current], 64)
		addToken(models.DOUBLE, val)
	} else {
		addToken(models.INT, source[start:current])
	}
}

// peek returns the current character, but does
// not consume it. Returns null character if at end
// of source.
func peek() rune {
	if isAtEnd() {
		return rune('\u0000')
	}
	return rune(source[current])
}

func PrintTokens(tokens []models.Token) {
	for i := 0; i < len(tokens); i++ {
		printToken(tokens[i])
	}
}

func printToken(token models.Token) {
	fmt.Printf("%s %s\n", models.TokenTypesNames[token.Type], token.Lexeme)
}
