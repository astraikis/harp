package models

import (
	"fmt"
	"time"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Column  int
	Line    int
}

type TokenType int

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_SQUARE
	RIGHT_SQUARE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	IDENTIFIER
	STRING
	INT
	DOUBLE

	AND
	ELSE
	FALSE
	FUNC
	FOR
	IF
	NULL
	OR
	RETURN
	TRUE
	WHILE
	STRUCT

	STRING_VAR
	INT_VAR
	DOUBLE_VAR
	BOOL_VAR

	EOF
)

var TokenTypesNames = map[TokenType]string{
	LEFT_PAREN:   "LEFT_PAREN",
	RIGHT_PAREN:  "RIGHT_PAREN",
	LEFT_BRACE:   "LEFT_BRACE",
	RIGHT_BRACE:  "RIGHT_BRACE",
	LEFT_SQUARE:  "LEFT_SQUARE",
	RIGHT_SQUARE: "RIGHT_SQUARE",
	COMMA:        "COMMA",
	DOT:          "DOT",
	MINUS:        "MINUS",
	PLUS:         "PLUS",
	SEMICOLON:    "SEMICOLON",
	SLASH:        "SLASH",
	STAR:         "STAR",

	BANG:          "BANG",
	BANG_EQUAL:    "BANG_EQUAL",
	EQUAL:         "EQUAL",
	EQUAL_EQUAL:   "EQUAL_EQUAL",
	GREATER:       "GREATER",
	GREATER_EQUAL: "GREATER_EQUAL",
	LESS:          "LESS",
	LESS_EQUAL:    "LESS_EQUAL",

	IDENTIFIER: "IDENTIFIER",
	STRING:     "STRING",
	INT:        "INT",
	DOUBLE:     "DOUBLE",

	AND:    "AND",
	ELSE:   "ELSE",
	FALSE:  "FALSE",
	FUNC:   "FUNC",
	FOR:    "FOR",
	IF:     "IF",
	NULL:   "NULL",
	OR:     "OR",
	RETURN: "RETURN",
	TRUE:   "TRUE",
	WHILE:  "WHILE",
	STRUCT: "STRUCT",

	STRING_VAR: "STRING_VAR",
	INT_VAR:    "INT_VAR",
	DOUBLE_VAR: "DOUBLE_VAR",
	BOOL_VAR:   "BOOL_VAR",

	EOF: "",
}

type Expr interface {
}

type AssignExpr struct {
	Name  Token
	Value Expr
}

type BinaryExpr struct {
	Left     Expr
	Right    Expr
	Operator Token
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

type LiteralExpr struct {
	Literal interface{}
}

type GroupingExpr struct {
	Expression Expr
}

type VarExpr struct {
	Name Token
}

type LogicExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type CallExpr struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

type Stmt interface {
}

type ExprStmt struct {
	Expression Expr
}

type VarStmt struct {
	Name        Token
	Initializer Expr
}

type BlockStmt struct {
	Statements []Stmt
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

type Callable interface {
	Call([]Expr) interface{}
	Arity() int
}

type Clock struct{}

func (c Clock) Call(arguments []Expr) interface{} {
	return int(time.Now().UnixMilli())
}

func (c Clock) Arity() int {
	return 0
}

type Print struct{}

func (p Print) Call(arguments []Expr) interface{} {
	if len(arguments) == 0 {
		fmt.Println()
	} else {
		fmt.Println(arguments[0])
	}
	return nil
}

func (p Print) Arity() int {
	return 1
}
