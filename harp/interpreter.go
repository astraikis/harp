package main

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func interpret(statements []stmt) {
	for _, stmt := range statements {
		execute(stmt)
	}
}

func evaluate(expr expr) interface{} {
	switch reflect.TypeOf(expr).String() {
	case "main.binaryExpr":
		return evaluateBinaryExpr(expr.(binaryExpr))
	case "main.literalExpr":
		return evaluateLiteralExpr(expr.(literalExpr))
	case "main.groupingExpr":
		return evaluateGroupingExpr(expr.(groupingExpr))
	case "main.varExpr":
		return evaluateVarExpr(expr.(varExpr))
	case "main.assignExpr":
		return evaluateAssignExpr(expr.(assignExpr))
	}

	return ""
}

func execute(stmt stmt) {
	switch reflect.TypeOf(stmt).String() {
	case "main.exprStmt":
		executeExprStmt(stmt.(exprStmt))
	case "main.varStmt":
		executeVarStmt(stmt.(varStmt))
	}
}

func executeExprStmt(stmt exprStmt) {
	evaluate(stmt.expression)
}

func executeVarStmt(stmt varStmt) {
	var value expr
	if stmt.initializer != nil {
		value = evaluate(stmt.initializer)
	}

	define(stmt.name.Lexeme, value)
}

func evaluateBinaryExpr(expr binaryExpr) interface{} {
	left := evaluate(expr.left)
	right := evaluate(expr.right)

	if isFloat(left) || isFloat(right) {
		// Evaluate as float.
		leftFloat := 0.0
		rightFloat := 0.0

		if isString(left) {
			leftParsed, err := strconv.ParseFloat(left.(string), 64)
			if err != nil {
				panic(err)
			}
			leftFloat = leftParsed
		} else {
			leftFloat = left.(float64)
		}
		if isString(right) {
			rightParsed, err := strconv.ParseFloat(right.(string), 64)
			if err != nil {
				panic(err)
			}
			rightFloat = rightParsed
		} else {
			rightFloat = right.(float64)
		}

		switch expr.operator.Type {
		case PLUS:
			return leftFloat + rightFloat
		case MINUS:
			return leftFloat - rightFloat
		case STAR:
			return leftFloat * rightFloat
		case SLASH:
			return leftFloat / rightFloat
		}
	} else {
		// Evaluate as int.
		leftInt := 0
		rightInt := 0

		if isString(left) {
			leftParsed, err := strconv.Atoi(left.(string))
			if err != nil {
				panic(err)
			}
			leftInt = leftParsed
		} else {
			leftInt = left.(int)
		}
		if isString(right) {
			rightParsed, err := strconv.Atoi(right.(string))
			if err != nil {
				panic(err)
			}
			rightInt = rightParsed
		} else {
			rightInt = right.(int)
		}

		switch expr.operator.Type {
		case PLUS:
			fmt.Println(leftInt + rightInt)
			return leftInt + rightInt
		case MINUS:
			return leftInt - rightInt
		case STAR:
			return leftInt * rightInt
		case SLASH:
			return leftInt / rightInt
		}
	}

	return nil
}

func evaluateGroupingExpr(expr groupingExpr) interface{} {
	return evaluate(expr.expression)
}

func evaluateLiteralExpr(expr literalExpr) interface{} {
	return expr.literal
}

func evaluateVarExpr(expr varExpr) interface{} {
	return getValue(expr.name)
}

func evaluateAssignExpr(expr assignExpr) interface{} {
	value := evaluate(expr.value)
	assignValue(expr.name, value)
	return value
}

func isFloat(value interface{}) bool {
	switch v := value.(type) {
	case float64:
		return true
	case string:
		floatParsed, err := strconv.ParseFloat(v, 64)
		if err == nil && !math.IsNaN(floatParsed) && !math.IsInf(floatParsed, 0) {
			_, err := strconv.Atoi(v)
			return err != nil
		} else {
			return true
		}
	default:
		return false
	}
}

func isString(value interface{}) bool {
	return reflect.TypeOf(value) == reflect.TypeOf("")
}
