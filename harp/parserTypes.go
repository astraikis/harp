package main

type expr interface {
}

type assignExpr struct {
	name  token
	value expr
}

type binaryExpr struct {
	left     expr
	right    expr
	operator token
}

type unaryExpr struct {
	operator token
	right    expr
}

type literalExpr struct {
	literal interface{}
}

type groupingExpr struct {
	expression expr
}

type varExpr struct {
	name token
}

type stmt interface {
}

type exprStmt struct {
	expression expr
}

type varStmt struct {
	name        token
	initializer expr
}
