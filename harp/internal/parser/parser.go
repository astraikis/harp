package parser

import (
	"errors"
	"reflect"

	"github.com/astraikis/harp/internal/models"
)

var tokens = []models.Token{}
var stmts = []models.Stmt{}
var current = 0

// Parse parses a list of tokens and returns
// the corresponding list of stmts.
func Parse(scannedTokens []models.Token) []models.Stmt {
	tokens = scannedTokens
	for {
		if isAtEnd() {
			break
		}
		next := declaration()
		stmts = append(stmts, next)
	}

	return stmts
}

func declaration() models.Stmt {
	if match([]models.TokenType{models.INT_VAR, models.DOUBLE_VAR, models.BOOL_VAR, models.STRING_VAR}) {
		return varDeclaration()
	}

	return statement()
}

func varDeclaration() models.Stmt {
	name, err := consume([]models.TokenType{models.IDENTIFIER}, "Expect variable name.")
	if err != nil {
		panic(err)
	}

	var initializer models.Expr
	if match([]models.TokenType{models.EQUAL}) {
		initializer = expression()
	}

	consume([]models.TokenType{models.SEMICOLON}, "Expect ';' after variable declaration.")
	return models.VarStmt{Name: *name, Initializer: initializer}
}

func statement() models.Stmt {
	if match([]models.TokenType{models.IF}) {
		return ifStatement()
	}
	if match([]models.TokenType{models.LEFT_BRACE}) {
		return models.BlockStmt{Statements: block()}
	}
	return expressionStatement()
}

func ifStatement() models.Stmt {
	consume([]models.TokenType{models.LEFT_PAREN}, "Expect '(' after 'if'.")
	condition := expression()
	consume([]models.TokenType{models.RIGHT_PAREN}, "Expect ')' after 'if'.")

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
		if check(models.RIGHT_BRACE) || isAtEnd() {
			break
		}

		blockStmts = append(blockStmts, declaration())
	}

	consume([]models.TokenType{models.RIGHT_BRACE}, "Expect '}' after block.")
	return blockStmts
}

func expressionStatement() models.Stmt {
	expr := expression()
	consume([]models.TokenType{models.SEMICOLON}, "Expect ';' after expression.")
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

	return primary()
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
	if match([]models.TokenType{models.LEFT_PAREN}) {
		inner := expression()
		_, err := consume([]models.TokenType{models.RIGHT_PAREN}, "Expect ')' after expression.")
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

	return nil, errors.New(message)
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
