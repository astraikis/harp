package parser

import (
	"fmt"
	"reflect"

	"github.com/astraikis/harp/internal/models"
)

var tokens []models.Token
var stmts []models.Stmt
var current = 0

// Parse parses a list of tokens and returns
// the corresponding list of stmts.
func Parse(scannedTokens []models.Token) ([]models.Stmt, []error) {
	tokens = scannedTokens
	for {
		if isAtEnd() {
			break
		}
		next := declaration()
		stmts = append(stmts, next)
	}

	return stmts, parseErrors
}

func declaration() models.Stmt {
	if match([]models.TokenType{models.INT_VAR, models.DOUBLE_VAR, models.BOOL_VAR, models.STRING_VAR}) {
		return varDeclaration()
	}
	if match([]models.TokenType{models.FUNC}) {
		return function()
	}

	return statement()
}

func function() models.Stmt {
	name, err := consume([]models.TokenType{models.IDENTIFIER}, "Expect function name.")
	if err != nil {
		return models.ErrorStmt{}
	}

	_, err = consume([]models.TokenType{models.LeftParen}, "Expect '(' after function name.")
	if err != nil {
		return models.ErrorStmt{}
	}

	var parameters []models.FuncParam
	if !check(models.RightParen) {
		paramType, err := consume([]models.TokenType{models.STRING_VAR, models.INT_VAR, models.DOUBLE_VAR, models.BOOL_VAR}, "Expect parameter name.")
		if err != nil {
			return models.ErrorStmt{}
		}

		paramName, err := consume([]models.TokenType{models.IDENTIFIER}, "Expect parameter name.")
		if err != nil {
			return models.ErrorStmt{}
		}

		parameters = append(parameters, models.FuncParam{Type: paramType.Type, Name: paramName.Lexeme})

		for {
			if !match([]models.TokenType{models.COMMA}) {
				break
			}

			paramType, err := consume([]models.TokenType{models.STRING_VAR, models.INT_VAR, models.DOUBLE_VAR, models.BOOL_VAR}, "Expect parameter name.")
			if err != nil {
				return models.ErrorStmt{}
			}

			paramName, err := consume([]models.TokenType{models.IDENTIFIER}, "Expect parameter name.")
			if err != nil {
				return models.ErrorStmt{}
			}

			parameters = append(parameters, models.FuncParam{Type: paramType.Type, Name: paramName.Lexeme})
		}
	}

	fmt.Println(parameters)
	_, _ = consume([]models.TokenType{models.RightParen}, "Expect ')' after function parameters.")
	_, _ = consume([]models.TokenType{models.LeftBrace}, "Expect '{' before function body.")

	body := block()
	return models.FuncStmt{Name: *name, Params: parameters, Body: body}
}

func varDeclaration() models.Stmt {
	name, err := consume([]models.TokenType{models.IDENTIFIER}, "Expect variable name.")
	if err != nil {
		return models.ErrorStmt{}
	}

	var initializer models.Expr
	if match([]models.TokenType{models.EQUAL}) {
		initializer = expression()
	}

	_, _ = consume([]models.TokenType{models.SEMICOLON}, "Expect ';' after variable declaration.")
	return models.VarStmt{Name: *name, Initializer: initializer}
}

func statement() models.Stmt {
	if match([]models.TokenType{models.IF}) {
		return ifStatement()
	}
	if match([]models.TokenType{models.LeftBrace}) {
		return models.BlockStmt{Statements: block()}
	}
	if match([]models.TokenType{models.WHILE}) {
		return whileStatement()
	}
	if match([]models.TokenType{models.FOR}) {
		return forStatement()
	}
	return expressionStatement()
}

func forStatement() models.Stmt {
	_, err := consume([]models.TokenType{models.LeftParen}, "Expect '(' after for.")
	if err != nil {
		return models.ErrorStmt{}
	}

	var initializer models.Stmt
	if match([]models.TokenType{models.INT_VAR, models.DOUBLE_VAR, models.BOOL_VAR, models.STRING_VAR}) {
		initializer = varDeclaration()
	} else {
		initializer = expressionStatement()
	}

	var condition models.Expr
	if !check(models.SEMICOLON) {
		condition = expression()
	} else {
		condition = nil
	}

	_, _ = consume([]models.TokenType{models.SEMICOLON}, "Expect ';' after loop condition.")

	var increment models.Expr
	if !check(models.RightParen) {
		increment = expression()
	} else {
		increment = nil
	}

	_, _ = consume([]models.TokenType{models.RightParen}, "Expect ')' after for clauses.")

	body := statement()

	if increment != nil {
		body = models.BlockStmt{Statements: []models.Stmt{body, models.ExprStmt{Expression: increment}}}
	}

	if condition == nil {
		condition = models.LiteralExpr{Literal: true}
	}
	body = models.WhileStmt{Condition: condition, Body: body}

	if initializer != nil {
		body = models.BlockStmt{Statements: []models.Stmt{initializer, body}}
	}

	return body
}

func whileStatement() models.Stmt {
	_, err := consume([]models.TokenType{models.LeftParen}, "Expect '(' after while.")
	if err != nil {
		return models.ErrorStmt{}
	}

	condition := expression()
	_, err = consume([]models.TokenType{models.RightParen}, "Expect ')' after while condition.")
	if err != nil {
		return models.ErrorStmt{}
	}

	body := statement()

	return models.WhileStmt{Condition: condition, Body: body}
}

func ifStatement() models.Stmt {
	_, err := consume([]models.TokenType{models.LeftParen}, "Expect '(' after 'if'.")
	if err != nil {
		return models.ErrorStmt{}
	}

	condition := expression()
	_, err = consume([]models.TokenType{models.RightParen}, "Expect ')' after 'if'.")
	if err != nil {
		return models.ErrorStmt{}
	}

	thenBranch := statement()
	var elseBranch models.Stmt
	if match([]models.TokenType{models.ELSE}) {
		elseBranch = statement()
	}

	return models.IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func block() []models.Stmt {
	var blockStmts = []models.Stmt{}

	for {
		if check(models.RightBrace) || isAtEnd() {
			break
		}

		blockStmts = append(blockStmts, declaration())
	}

	_, err := consume([]models.TokenType{models.RightBrace}, "Expect '}' after block.")
	if err != nil {
		return []models.Stmt{models.ErrorStmt{}}
	}

	return blockStmts
}

func expressionStatement() models.Stmt {
	expr := expression()

	_, err := consume([]models.TokenType{models.SEMICOLON}, "Expect ';' after expression.")
	if err != nil {
		return models.ErrorStmt{}
	}

	return models.ExprStmt{Expression: expr}
}

func expression() models.Expr {
	return assignment()
}

func assignment() models.Expr {
	expr := or()

	if match([]models.TokenType{models.EQUAL}) {
		equals := previous()
		value := assignment()

		if reflect.TypeOf(expr).String() == "models.VarExpr" {
			name := expr.(models.VarExpr).Name
			return models.AssignExpr{Name: name, Value: value}
		}

		// Error
		if equals == value {
		}
	}

	return expr
}

func or() models.Expr {
	expr := and()

	for {
		if !match([]models.TokenType{models.OR}) {
			break
		}

		operator := previous()
		right := and()
		expr = models.LogicExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func and() models.Expr {
	expr := equality()

	for {
		if !match([]models.TokenType{models.AND}) {
			break
		}

		operator := previous()
		right := equality()
		expr = models.LogicExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func equality() models.Expr {
	parsed := comparison()

	for {
		if !match([]models.TokenType{models.BANG_EQUAL, models.EQUAL_EQUAL}) {
			break
		}

		operator := previous()
		right := comparison()
		parsed = models.BinaryExpr{Left: parsed, Operator: operator, Right: right}
	}

	return parsed
}

func comparison() models.Expr {
	parsed := term()

	for {
		if !match([]models.TokenType{models.GREATER, models.GREATER_EQUAL, models.LESS, models.LESS_EQUAL}) {
			break
		}

		operator := previous()
		right := term()
		parsed = models.BinaryExpr{Left: parsed, Operator: operator, Right: right}
	}

	return parsed
}

func term() models.Expr {
	parsed := factor()

	for {
		if !match([]models.TokenType{models.MINUS, models.PLUS}) {
			break
		}

		operator := previous()
		right := factor()
		parsed = models.BinaryExpr{Left: parsed, Operator: operator, Right: right}
	}

	return parsed
}

func factor() models.Expr {
	parsed := unary()

	for {
		if !match([]models.TokenType{models.STAR, models.SLASH}) {
			break
		}

		operator := previous()
		right := unary()
		parsed = models.BinaryExpr{Left: parsed, Operator: operator, Right: right}
	}

	return parsed
}

func unary() models.Expr {
	if match([]models.TokenType{models.BANG, models.MINUS}) {
		operator := previous()
		right := unary()
		return models.UnaryExpr{Operator: operator, Right: right}
	}

	return call()
}

func call() models.Expr {
	expr := primary()

	for {
		if !match([]models.TokenType{models.LeftParen}) {
			break
		}

		expr = finishCall(expr)
	}

	return expr
}

func finishCall(callee models.Expr) models.Expr {
	var arguments []models.Expr

	if !check(models.RightParen) {
		arguments = append(arguments, expression())

		for {
			if !match([]models.TokenType{models.COMMA}) {
				break
			}
			if len(arguments) >= 255 {
				// Error
			}
			arguments = append(arguments, expression())
		}
	}

	paren, err := consume([]models.TokenType{models.RightParen}, "Expect ')' after arguments.")
	if err != nil {
		return models.ErrorExpr{}
	}

	return models.CallExpr{Callee: callee, Paren: *paren, Arguments: arguments}
}

func primary() models.Expr {
	if match([]models.TokenType{models.TRUE}) {
		return models.LiteralExpr{Literal: true}
	}
	if match([]models.TokenType{models.FALSE}) {
		return models.LiteralExpr{Literal: false}
	}
	if match([]models.TokenType{models.NULL}) {
		return models.LiteralExpr{Literal: nil}
	}
	if match([]models.TokenType{models.INT, models.DOUBLE, models.STRING}) {
		return models.LiteralExpr{Literal: previous().Literal}
	}
	if match([]models.TokenType{models.IDENTIFIER}) {
		return models.VarExpr{Name: previous()}
	}
	if match([]models.TokenType{models.LeftParen}) {
		inner := expression()
		_, err := consume([]models.TokenType{models.RightParen}, "Expect ')' after expression.")
		if err != nil {
			sync()
			// Error
		}
		return models.GroupingExpr{Expression: inner}
	}

	return models.LiteralExpr{Literal: false}
}

// advance returns the next token.
func advance() models.Token {
	if !isAtEnd() {
		current += 1
	}

	return previous()
}

// consume checks if expected is equal to the next token's type
// consumes it if it does.
func consume(expectedTypes []models.TokenType, message string) (*models.Token, error) {
	for _, expected := range expectedTypes {
		if expected == peek().Type {
			token := advance()
			return &token, nil
		}
	}

	err := &ParseError{
		Line:    peek().Line,
		Column:  peek().Column,
		Message: message,
	}
	reportError(err)
	sync()

	return nil, err
}

// match reports whether the current token
// matches tokenType and consumes it if it does.
func match(expectedTypes []models.TokenType) bool {
	for _, expected := range expectedTypes {
		if check(expected) {
			advance()
			return true
		}
	}

	return false
}

// check reports whether the current token
// matches tokenType, but does not consume it.
func check(tokenType models.TokenType) bool {
	if isAtEnd() {
		return false
	}

	return peek().Type == tokenType
}

// peek returns the current token.
func peek() models.Token {
	return tokens[current]
}

// previous returns the previous token.
func previous() models.Token {
	return tokens[current-1]
}

// isAtEnd() reports whether we're at the end of tokens.
func isAtEnd() bool {
	return current >= len(tokens) || peek().Type == models.EOF
}

func sync() {
	advance()

	for {
		if isAtEnd() {
			return
		}

		if previous().Type == models.SEMICOLON {
			return
		}

		switch peek().Type {
		case models.FUNC:
		case models.INT_VAR:
		case models.DOUBLE_VAR:
		case models.STRING_VAR:
		case models.BOOL_VAR:
		case models.FOR:
		case models.IF:
		case models.WHILE:
		case models.RETURN:
			return
		}

		advance()
	}
}
