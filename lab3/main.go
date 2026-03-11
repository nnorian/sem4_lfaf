package main

import (
	"fmt"
	"http-lexer/src/lexer"
	"strings"
)

func main(){
	samples := []struct{
		label string 
		input string
	}
	{
		{
			label: "Simple GET with query params",
			input: "GET /search?q=golang&page=2&debug=true HTTP/1.1\nHost: example.com\nAccept: application/json",
		},
		{
			label: "POST with headers",
			input: "POST /api/users HTTP/1.1\nHost: api.example.com\nContent-Type: application/json\nContent-Length: 42\nAuthorization: Bearer eyJhbGciOiJIUzI1NiJ9",
		},
		{
			label: "DELETE with float and boolean headers",
			input: "DELETE /api/items/99 HTTP/2\nHost: api.example.com\nX-Rate-Limit: 3.14\nX-Enabled: true",
		},
	
	}

	for _, sample := range samples {
		banner := fmt.Sprintf("%s", sample.label)
		fmt.Println(banner)
		fmt.Printf("Input:\n%s\n\nTokens:\n", sample.input)

		l := lexer.New(sample.input)
		tokens := l.Tokenise()

		for _, tok := range tokens {
			if tok.Type == lexer.TOKEN_EOF{
				break
			}
			fmt.Printf(" %s/n", tok)
		}
		fmt.Println()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}