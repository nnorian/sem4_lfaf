package lexer
import "fmt"


//categories
type TockenType int

const(
	//request line
	TOKEN_METHOD TokenType=iota
	TOKEN_PATH
	TOKEN_HTTP_VERSION
	TOKEN_QUERY_SEP
	TOKEN_QUERY_KEY
	TOKEN_QUERY_VALUE
	TOKEN_AMPERSAND
	TOKEN_EQUALS

	//headers
	TOKEN_HEADER_NAME
	TOKEN_COLON
	TOKEN_HEADER_VALUE

	//body
	TOKEN_INTEGER
	TOKEN_FLOAT
	TOKEN_BOOLEAN
	TOKEN_QUOTED_STRING 
	
	//structure
	TOKEN_NEWLINE
	TOKEN_WHITESPACE
	TOKEN_UNKNOWN
	TOKEN_EOF

)

// token names maps token type to human readable label
var tokenNames = map[TokenType]string{
	TOKEN_METHOD: "METHOD",
	TOKEN_PATH: "PATH",
	TOKEN_HTTP_VERSION: "HTTP_VERSION",
	tOKEN_QUERY_KEY: "QUERY_KEY",
	TOKEN_QUERY_VALUE: "QUERY_VALUE",
	TOKEN_AMPERSAND: "AMPERSAND",
	TOKEN_EQUALS: "EQUALS",
	TOKEN_HEADER_NAME: "HEADER_NAME",
	TOKEN_COLON: "COLON",
	TOKEN_HEADER_VALUE: "HEADER_VALUE",
	TOKEN_INTEGER: "INTEGER",
	TOKEN_FLOAT: "FLOAT", 
	TOKEN_BOOLEAN: "BOOLEAN",
	TOKEN_QUOTED_STRING: "QUOTED_STRING",
	TOKEN_NEWLINE: "NEWLINE",
	TOKEN_WHITESPACE: "WHITESPACE",
	TOKEN_UNKNOWN: "UNKNOWN",
	TOKEN_EOF: "EOF",
}

func (t TockenType) String(){
	if name, ok := tokenNames[t]; ok{
		return name
	}
	return "UNKNOWN"
}

// token interpreted as single lexical unit
type token struct{
	Type TokenType
	Value string 
	Line int
	Col int
}

func (t Token) String() string {
	return fmt.Sprintf("%-20s | %-35q | line %-3d col %d", t.Typem t.Value, t.Line, t.Col)
}