package main

import (
	"fmt"
	"strconv"
	"unicode"
)

var source string
var tokens = []token{}
var start int = 0
var current int = 0
var line int = 1
var column int = 1

var keywords = map[string]tokenType{
	"and":    AND,
	"else":   ELSE,
	"false":  FALSE,
	"func":   FUNC,
	"for":    FOR,
	"if":     IF,
	"null":   NULL,
	"or":     OR,
	"return": RETURN,
	"true":   TRUE,
	"while":  WHILE,
	"struct": STRUCT,
	"string": STRING_VAR,
	"int":    INT_VAR,
	"double": DOUBLE_VAR,
	"bool":   BOOL_VAR,
}

type token struct {
	Type    tokenType
	Lexeme  string
	Literal interface{}
	Column  int
	Line    int
}

// scan walks over source and returns
// a corresponding list of tokens.
func scan(sourceText string) []token {
	source = sourceText
	for {
		if isAtEnd() {
			break
		}
		start = current
		scanToken()
	}

	tokens = append(tokens, token{Type: EOF, Lexeme: "", Literal: nil, Column: 1, Line: line})
	return tokens
}

// scanToken adds the next token to tokens.
func scanToken() {
	var c rune = advance()
	switch c {
	case '(':
		addToken(LEFT_PAREN, "")
		column += 1
	case ')':
		addToken(RIGHT_PAREN, "")
		column += 1
	case '{':
		addToken(LEFT_BRACE, "")
		column += 1
	case '}':
		addToken(RIGHT_BRACE, "")
		column += 1
	case ',':
		addToken(COMMA, "")
		column += 1
	case '.':
		addToken(DOT, "")
		column += 1
	case '-':
		addToken(MINUS, "")
		column += 1
	case '+':
		addToken(PLUS, "")
		column += 1
	case ';':
		addToken(SEMICOLON, "")
		column += 1
	case '*':
		addToken(STAR, "")
		column += 1
	case '!':
		if match('=') {
			addToken(BANG_EQUAL, "")
			column += 2
		} else {
			addToken(BANG, "")
			column += 1
		}
	case '=':
		if match('=') {
			addToken(EQUAL_EQUAL, "")
			column += 2
		} else {
			addToken(EQUAL, "")
			column += 1
		}
	case '>':
		if match('=') {
			addToken(GREATER_EQUAL, "")
			column += 2
		} else {
			addToken(GREATER, "")
			column += 1
		}
	case '<':
		if match('=') {
			addToken(LESS_EQUAL, "")
			column += 2
		} else {
			addToken(LESS, "")
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
			addToken(SLASH, "")
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
			harpError(string(c), "Unexpected character.", column, line)
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

	addToken(STRING, source[start:current])
}

// addToken adds a token to tokens.
func addToken[T any](tokenType tokenType, literal T) {
	tokens = append(tokens, token{Type: tokenType, Lexeme: string(source[start:current]), Literal: literal, Column: column, Line: line})
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
		_type = IDENTIFIER
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
		addToken(DOUBLE, val)
	} else {
		addToken(INT, source[start:current])
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

func printTokens() {
	for i := 0; i < len(tokens); i++ {
		printToken(tokens[i])
	}
}

func printToken(token token) {
	fmt.Printf("%s %s\n", tokenTypesNames[token.Type], token.Lexeme)
}
