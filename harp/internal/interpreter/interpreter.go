package interpreter

import (
	"math"
	"reflect"
	"strconv"

	"github.com/astraikis/harp/internal/models"
)

func Interpret(statements []models.Stmt) {
	for _, stmt := range statements {
		execute(stmt)
	}
}

func execute(stmt models.Stmt) {
	switch reflect.TypeOf(stmt).String() {
	case "models.ExprStmt":
		executeExprStmt(stmt.(models.ExprStmt))
	case "models.VarStmt":
		executeVarStmt(stmt.(models.VarStmt))
	}
}

func executeExprStmt(stmt models.ExprStmt) {
	evaluate(stmt.Expression)
}

func executeVarStmt(stmt models.VarStmt) {
	var value models.Expr
	if stmt.Initializer != nil {
		value = evaluate(stmt.Initializer)
	}

	defineValue(stmt.Name.Lexeme, value)
}

func evaluate(expr models.Expr) interface{} {
	switch reflect.TypeOf(expr).String() {
	case "models.BinaryExpr":
		return evaluateBinaryExpr(expr.(models.BinaryExpr))
	case "models.LiteralExpr":
		return evaluateLiteralExpr(expr.(models.LiteralExpr))
	case "models.GroupingExpr":
		return evaluateGroupingExpr(expr.(models.GroupingExpr))
	case "models.VarExpr":
		return evaluateVarExpr(expr.(models.VarExpr))
	case "models.AssignExpr":
		return evaluateAssignExpr(expr.(models.AssignExpr))
	}

	return ""
}

func evaluateBinaryExpr(expr models.BinaryExpr) interface{} {
	left := evaluate(expr.Left)
	right := evaluate(expr.Right)

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

		switch expr.Operator.Type {
		case models.PLUS:
			return leftFloat + rightFloat
		case models.MINUS:
			return leftFloat - rightFloat
		case models.STAR:
			return leftFloat * rightFloat
		case models.SLASH:
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

		switch expr.Operator.Type {
		case models.PLUS:
			return leftInt + rightInt
		case models.MINUS:
			return leftInt - rightInt
		case models.STAR:
			return leftInt * rightInt
		case models.SLASH:
			return leftInt / rightInt
		}
	}

	return nil
}

func evaluateGroupingExpr(expr models.GroupingExpr) interface{} {
	return evaluate(expr.Expression)
}

func evaluateLiteralExpr(expr models.LiteralExpr) interface{} {
	return expr.Literal
}

func evaluateVarExpr(expr models.VarExpr) interface{} {
	return getValue(expr.Name)
}

func evaluateAssignExpr(expr models.AssignExpr) interface{} {
	value := evaluate(expr.Value)
	assignValue(expr.Name, value)
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
