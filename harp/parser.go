package main

import (
	"errors"
	"fmt"
	"reflect"
)

var p_tokens = []token{}
var p_stmts = []stmt{}
var p_current = 0

// parse parses a list of tokens and returns
// the corresponding list of stmts.
func parse(tokens []token) []stmt {
	p_tokens = tokens
	for {
		if p_isAtEnd() {
			break
		}
		next := declaration()
		p_stmts = append(p_stmts, next)
	}

	return p_stmts
}

func declaration() stmt {
	if p_match([]tokenType{INT_VAR, DOUBLE_VAR, BOOL_VAR, STRING_VAR}) {
		return varDeclaration()
	}

	return statement()
}

func varDeclaration() stmt {
	name, err := p_consume([]tokenType{IDENTIFIER}, "Expect variable name.")
	if err != nil {
		panic(err)
	}

	var initializer expr
	if p_match([]tokenType{EQUAL}) {
		initializer = expression()
	}

	p_consume([]tokenType{SEMICOLON}, "Expect ';' after variable declaration.")
	return varStmt{name: *name, initializer: initializer}
}

func statement() stmt {
	return expressionStatement()
}

func expressionStatement() stmt {
	expr := expression()
	p_consume([]tokenType{SEMICOLON}, "Expect ';' after expression.")
	return exprStmt{expression: expr}
}

func expression() expr {
	return assignment()
}

func assignment() expr {
	expr := equality()

	if p_match([]tokenType{EQUAL}) {
		equals := p_previous()
		value := assignment()

		if reflect.TypeOf(expr).String() == "main.varExpr" {
			name := expr.(varExpr).name
			return assignExpr{name: name, value: value}
		}

		harpError(equals.Lexeme, "Invalid assignment target.", 0, 0)
	}

	return expr
}

func equality() expr {
	parsed := comparison()

	for {
		if !p_match([]tokenType{BANG_EQUAL, EQUAL_EQUAL}) {
			break
		}

		operator := p_previous()
		right := comparison()
		parsed = binaryExpr{left: parsed, operator: operator, right: right}
	}

	return parsed
}

func comparison() expr {
	parsed := term()

	for {
		if !p_match([]tokenType{GREATER, GREATER_EQUAL, LESS, LESS_EQUAL}) {
			break
		}

		operator := p_previous()
		right := term()
		parsed = binaryExpr{left: parsed, operator: operator, right: right}
	}

	return parsed
}

func term() expr {
	parsed := factor()

	for {
		if !p_match([]tokenType{MINUS, PLUS}) {
			break
		}

		operator := p_previous()
		right := factor()
		parsed = binaryExpr{left: parsed, operator: operator, right: right}
	}

	return parsed
}

func factor() expr {
	parsed := unary()

	for {
		if !p_match([]tokenType{STAR, SLASH}) {
			break
		}

		operator := p_previous()
		right := unary()
		parsed = binaryExpr{left: parsed, operator: operator, right: right}
	}

	return parsed
}

func unary() expr {
	if p_match([]tokenType{BANG, MINUS}) {
		operator := p_previous()
		right := unary()
		return unaryExpr{operator: operator, right: right}
	}

	return primary()
}

func primary() expr {
	if p_match([]tokenType{TRUE}) {
		return literalExpr{literal: true}
	}
	if p_match([]tokenType{FALSE}) {
		return literalExpr{literal: false}
	}
	if p_match([]tokenType{NULL}) {
		return literalExpr{literal: nil}
	}
	if p_match([]tokenType{INT, DOUBLE, STRING}) {
		return literalExpr{literal: p_previous().Literal}
	}
	if p_match([]tokenType{IDENTIFIER}) {
		return varExpr{name: p_previous()}
	}
	if p_match([]tokenType{LEFT_PAREN}) {
		inner := expression()
		_, err := p_consume([]tokenType{RIGHT_PAREN}, "Expect ')' after expression.")
		if err != nil {
			sync()
			harpError(tokens[current].Lexeme, err.Error(), tokens[current].Line, tokens[current].Line)
		}
		return groupingExpr{expression: inner}
	}

	return literalExpr{literal: false}
}

// p_advance returns the next token.
func p_advance() token {
	if !p_isAtEnd() {
		p_current += 1
	}

	return p_previous()
}

// p_consume checks if expected is equal to the next token's type
// consumes it if it does.
func p_consume(expectedTypes []tokenType, message string) (*token, error) {
	for _, expected := range expectedTypes {
		if expected == p_peek().Type {
			token := p_advance()
			return &token, nil
		}
	}

	return nil, errors.New(message)
}

// p_match reports whether the current token
// matches tokenType and consumes it if it does.
func p_match(expectedTypes []tokenType) bool {
	for _, expected := range expectedTypes {
		if p_check(expected) {
			p_advance()
			return true
		}
	}

	return false
}

// p_check reports whether the current token
// matches tokenType, but does not consume it.
func p_check(tokenType tokenType) bool {
	if p_isAtEnd() {
		return false
	}

	return p_peek().Type == tokenType
}

// p_peek returns the current token.
func p_peek() token {
	return p_tokens[p_current]
}

// p_previous returns the previous token.
func p_previous() token {
	return p_tokens[p_current-1]
}

// p_isAtEnd() reports whether we're at the end of tokens.
func p_isAtEnd() bool {
	return p_current >= len(p_tokens) || p_peek().Type == EOF
}

func sync() {
	p_advance()

	for {
		if isAtEnd() {
			return
		}

		if p_previous().Type == SEMICOLON {
			return
		}

		switch p_peek().Type {
		case FUNC:
		case INT_VAR:
		case DOUBLE_VAR:
		case STRING_VAR:
		case BOOL_VAR:
		case FOR:
		case IF:
		case WHILE:
		case RETURN:
			return
		}

		advance()
	}
}

func printStatements() {
	for _, stmt := range p_stmts {
		printStatement(stmt, 0)
	}
}

func printStatement(stmt expr, depth int) {
	switch reflect.TypeOf(stmt).String() {
	case "main.binaryExpr":
		be := stmt.(binaryExpr)
		for i := 0; i < depth; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("Binary expression {\n")

		printStatement(be.left, depth+1)
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf(be.operator.Lexeme + "\n\n")
		printStatement(be.right, depth+1)
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		fmt.Println("}")
	case "main.literalExpr":
		le := stmt.(literalExpr)
		for i := 0; i < depth; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("Literal expression: %s\n\n", le.literal.(string))
	case "main.groupingExpr":
		ge := stmt.(groupingExpr)
		printStatement(ge.expression, depth+1)
	case "main.varStmt":
		vs := stmt.(varStmt)
		for i := 0; i < depth; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("Variable statement {\n")
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf(vs.name.Lexeme)
		fmt.Printf(" {\n")
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		printStatement(vs.initializer, depth+1)
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("}\n")
		fmt.Println("}")
	case "main.assignExpr":
		ae := stmt.(assignExpr)
		for i := 0; i < depth; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("Assignment expression {\n")
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf(ae.name.Lexeme)
		fmt.Printf(" { \n")
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		printStatement(ae.value, depth+1)
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("}\n")
		fmt.Println("}")
	case "main.exprStmt":
		es := stmt.(exprStmt)
		for i := 0; i < depth; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("Expression statement {\n")
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		printStatement(es.expression, depth+1)
		for i := 0; i < depth; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("}\n")
	}
}
