package main

import (
	"fmt"
	"http-lexer/src/lexer"
	"http-lexer/src/parser"
	"strings"
)

func main() {
	samples := []struct {
		label string
		input string
	}{
		{
			label: "Simple GET with query params",
			input: "GET /search?q=golang&page=2&debug=true HTTP/1.1\nHost: example.com\nAccept: application/json",
		},
		{
			label: "POST with typed header values",
			input: "POST /api/users HTTP/1.1\nHost: api.example.com\nContent-Type: application/json\nContent-Length: 42\nAuthorization: Bearer eyJhbGciOiJIUzI1NiJ9",
		},
		{
			label: "DELETE with FLOAT and BOOLEAN headers",
			input: "DELETE /api/items/99 HTTP/2\nHost: api.example.com\nX-Rate-Limit: 3.14\nX-Enabled: true",
		},
	}

	for _, sample := range samples {
		banner := strings.Repeat("─", 60)
		fmt.Printf("\n%s\n%s\n%s\n", banner, sample.label, banner)
		fmt.Printf("Input:\n%s\n", sample.input)

		//lexer
		fmt.Println("\nTokens:")
		l := lexer.New(sample.input)
		tokens := l.Tokenise()
		for _, tok := range tokens {
			if tok.Type == lexer.TOKEN_EOF {
				break
			}
			fmt.Printf("  %s\n", tok)
		}

		//parser
		fmt.Println("\nAST:")
		p := parser.New(tokens)
		ast, err := p.Parse()
		if err != nil {
			fmt.Printf("  parse error: %v\n", err)
			continue
		}
		fmt.Print(parser.Print(ast))
	}
}