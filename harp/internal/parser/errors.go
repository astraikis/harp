package parser

import "fmt"

var parseErrors []error

func reportError(err error) {
	parseErrors = append(parseErrors, err)
}

type ParseError struct {
	Line    int
	Column  int
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("[Line %d:%d] Error: %s", e.Line, e.Column, e.Message)
}
