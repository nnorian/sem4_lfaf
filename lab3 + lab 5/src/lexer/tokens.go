package lexer

import "fmt"

// TokenType categorises each lexical unit
type TokenType int

const (
	// request line
	TOKEN_METHOD TokenType = iota
	TOKEN_PATH
	TOKEN_HTTP_VERSION
	TOKEN_QUERY_SEP
	TOKEN_QUERY_KEY
	TOKEN_QUERY_VALUE
	TOKEN_AMPERSAND
	TOKEN_EQUALS

	// headers
	TOKEN_HEADER_NAME
	TOKEN_COLON
	TOKEN_HEADER_VALUE

	// typed literal values
	TOKEN_INTEGER
	TOKEN_FLOAT
	TOKEN_BOOLEAN
	TOKEN_QUOTED_STRING

	// structure
	TOKEN_UNKNOWN
	TOKEN_EOF
)

var tokenNames = map[TokenType]string{
	TOKEN_METHOD:        "METHOD",
	TOKEN_PATH:          "PATH",
	TOKEN_HTTP_VERSION:  "HTTP_VERSION",
	TOKEN_QUERY_SEP:     "QUERY_SEP",
	TOKEN_QUERY_KEY:     "QUERY_KEY",
	TOKEN_QUERY_VALUE:   "QUERY_VALUE",
	TOKEN_AMPERSAND:     "AMPERSAND",
	TOKEN_EQUALS:        "EQUALS",
	TOKEN_HEADER_NAME:   "HEADER_NAME",
	TOKEN_COLON:         "COLON",
	TOKEN_HEADER_VALUE:  "HEADER_VALUE",
	TOKEN_INTEGER:       "INTEGER",
	TOKEN_FLOAT:         "FLOAT",
	TOKEN_BOOLEAN:       "BOOLEAN",
	TOKEN_QUOTED_STRING: "QUOTED_STRING",
	TOKEN_UNKNOWN:       "UNKNOWN",
	TOKEN_EOF:           "EOF",
}

func (t TokenType) String() string {
	if name, ok := tokenNames[t]; ok {
		return name
	}
	return "UNKNOWN"
}

// Token represents a single lexical unit with source position
type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

func (t Token) String() string {
	return fmt.Sprintf("%-20s | %-35q | line %-3d col %d", t.Type, t.Value, t.Line, t.Col)
}