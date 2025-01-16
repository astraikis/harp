package parser

import (
	"fmt"
	"reflect"

	"github.com/astraikis/harp/internal/models"
)

func PrintStatements(parsedStmts []models.Stmt) {
	for _, stmt := range parsedStmts {
		printStatement(stmt, 0)
	}
}

func printStatement(stmt models.Expr, depth int) {
	switch reflect.TypeOf(stmt).String() {
	case "models.BinaryExpr":
		be := stmt.(models.BinaryExpr)
		for i := 0; i < depth; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("Binary expression {\n")

		printStatement(be.Left, depth+1)
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf(be.Operator.Lexeme + "\n\n")
		printStatement(be.Right, depth+1)
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		fmt.Println("}")
	case "models.LiteralExpr":
		le := stmt.(models.LiteralExpr)
		for i := 0; i < depth; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("Literal expression: %s\n\n", le.Literal.(string))
	case "models.GroupingExpr":
		ge := stmt.(models.GroupingExpr)
		printStatement(ge.Expression, depth+1)
	case "models.VarStmt":
		vs := stmt.(models.VarStmt)
		for i := 0; i < depth; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("Variable statement {\n")
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf(vs.Name.Lexeme)
		fmt.Printf(" {\n")
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		printStatement(vs.Initializer, depth+1)
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("}\n")
		fmt.Println("}")
	case "models.AssignExpr":
		ae := stmt.(models.AssignExpr)
		for i := 0; i < depth; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("Assignment expression {\n")
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf(ae.Name.Lexeme)
		fmt.Printf(" { \n")
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		printStatement(ae.Value, depth+1)
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("}\n")
		fmt.Println("}")
	case "models.ExprStmt":
		es := stmt.(models.ExprStmt)
		for i := 0; i < depth; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("Expression statement {\n")
		for i := 0; i < depth+1; i++ {
			fmt.Printf("   ")
		}
		printStatement(es.Expression, depth+1)
		for i := 0; i < depth; i++ {
			fmt.Printf("   ")
		}
		fmt.Printf("}\n")
	}
}
