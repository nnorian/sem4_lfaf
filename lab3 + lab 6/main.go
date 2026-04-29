package main

import (
	"fmt"
	"http-lexer/src/lexer"
	"http-lexer/src/parser"
)

func main() {
	samples := []string{
		"GET /search?q=golang&page=2&debug=true HTTP/1.1\nHost: example.com\nAccept: application/json",
		"POST /api/users HTTP/1.1\nHost: api.example.com\nContent-Type: application/json\nContent-Length: 42\nAuthorization: Bearer eyJhbGciOiJIUzI1NiJ9",
		"DELETE /api/items/99 HTTP/2\nHost: api.example.com\nX-Rate-Limit: 3.14\nX-Enabled: true",
	}

	for _, input := range samples {
		l := lexer.New(input)
		tokens := l.Tokenise()

		fmt.Println("Tokens:")
		for _, tok := range tokens {
			if tok.Type == lexer.TOKEN_EOF {
				break
			}
			fmt.Printf("  %s\n", tok)
		}

		p := parser.New(tokens)
		ast, err := p.Parse()
		if err != nil {
			fmt.Printf("parse error: %v\n", err)
			continue
		}

		fmt.Println("\nAST:")
		fmt.Print(parser.Print(ast))
		fmt.Println()
	}
}
