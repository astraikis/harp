package main

var p_tokens = []token{}
var p_current = 0

func parse(tokens []token) {
	p_tokens = tokens
	for {
		if p_isAtEnd() {
			break
		}
		printToken(p_tokens[p_current])
		p_current += 1
	}
}

func p_peek() token {
	return p_tokens[p_current]
}

func p_isAtEnd() bool {
	return p_peek().Type == EOF
}
