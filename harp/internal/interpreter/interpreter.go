package interpreter

import (
	"math"
	"reflect"
	"strconv"

	"github.com/astraikis/harp/internal/models"
)

var Globals = &Environment{values: map[string]interface{}{}, parent: nil}
var currEnvironment = Globals

type Interpreter struct {
	intGlobals *Environment
	intCurr    *Environment
}

type Function struct {
	*models.Function
	Interpreter *Interpreter
}

func (f *Function) Call(arguments []models.Expr) interface{} {
	env := &Environment{values: map[string]interface{}{}, parent: currEnvironment}

	for i, param := range f.Params {
		DefineValue(param.Name, arguments[i], env)
	}

	executeBlockStmt(f.Body, env)
	return nil
}

func Interpret(statements []models.Stmt) {
	DefineValue("clock", models.Clock{}, currEnvironment)
	DefineValue("print", models.Print{}, currEnvironment)

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
	case "models.BlockStmt":
		executeBlockStmt(stmt.(models.BlockStmt).Statements, &Environment{values: map[string]interface{}{}, parent: currEnvironment})
	case "models.IfStmt":
		executeIfStmt(stmt.(models.IfStmt))
	case "models.WhileStmt":
		executeWhileStmt(stmt.(models.WhileStmt))
	case "models.FuncStmt":
		executeFuncStmt(stmt.(models.FuncStmt))
	}
}

func executeFuncStmt(stmt models.FuncStmt) {
	//function := models.Function{Name: stmt.Name.Lexeme, Params: stmt.Params, Body: stmt.Body}
	function := &Function{
		Function: &models.Function{
			Name:   stmt.Name.Lexeme,
			Params: stmt.Params,
			Body:   stmt.Body,
		},
	}
	DefineValue(stmt.Name.Lexeme, function, currEnvironment)
}

func executeExprStmt(stmt models.ExprStmt) {
	evaluate(stmt.Expression)
}

func executeVarStmt(stmt models.VarStmt) {
	var value models.Expr
	if stmt.Initializer != nil {
		value = evaluate(stmt.Initializer)
	}

	DefineValue(stmt.Name.Lexeme, value, currEnvironment)
}

func executeBlockStmt(blockStmts []models.Stmt, blockEnvironment *Environment) {
	prevEnvironment := currEnvironment
	currEnvironment = blockEnvironment

	for _, stmt := range blockStmts {
		execute(stmt)
	}

	currEnvironment = prevEnvironment
}

func executeIfStmt(stmt models.IfStmt) {
	if isTruthy(evaluate(stmt.Condition)) {
		execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		execute(stmt.ElseBranch)
	}
}

func executeWhileStmt(stmt models.WhileStmt) {
	for {
		if !isTruthy(evaluate(stmt.Condition)) {
			break
		}
		execute(stmt.Body)
	}
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
	case "models.LogicExpr":
		return evaluateLogicExpr(expr.(models.LogicExpr))
	case "models.CallExpr":
		return evaluateCallExpr(expr.(models.CallExpr))
	}

	return ""
}

func evaluateCallExpr(expr models.CallExpr) interface{} {
	callee := evaluate(expr.Callee)

	var arguments []models.Expr
	for i := 0; i < len(expr.Arguments); i++ {
		arguments = append(arguments, evaluate(expr.Arguments[i]))
	}

	function := callee.(models.Callable)
	if len(arguments) != function.Arity() {
		// Error
	}

	return function.Call(arguments)
}

func evaluateLogicExpr(expr models.LogicExpr) interface{} {
	left := evaluate(expr.Left)

	if expr.Operator.Type == models.OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return evaluate(expr.Right)
}

func evaluateBinaryExpr(expr models.BinaryExpr) interface{} {
	left := evaluate(expr.Left)
	right := evaluate(expr.Right)

	switch expr.Operator.Type {
	case models.EQUAL_EQUAL:
		return left == right
	case models.BANG_EQUAL:
		return left != right
	}

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
		case models.LESS:
			return leftFloat < rightFloat
		case models.LESS_EQUAL:
			return leftFloat <= rightFloat
		case models.GREATER:
			return leftFloat > rightFloat
		case models.GREATER_EQUAL:
			return leftFloat >= rightFloat
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
		case models.LESS:
			return leftInt < rightInt
		case models.LESS_EQUAL:
			return leftInt <= rightInt
		case models.GREATER:
			return leftInt > rightInt
		case models.GREATER_EQUAL:
			return leftInt >= rightInt
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
	return GetValue(expr.Name.Lexeme, currEnvironment)
}

func evaluateAssignExpr(expr models.AssignExpr) interface{} {
	value := evaluate(expr.Value)
	AssignValue(expr.Name.Lexeme, value, currEnvironment)
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

func isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}

	if boolValue, ok := value.(bool); ok {
		return boolValue
	}

	return true
}
