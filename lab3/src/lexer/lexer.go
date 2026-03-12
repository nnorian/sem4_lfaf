package lexer

import (
	"strings"
	"unicode"
)

// each state represents a position
type state int

const (
	STATE_START state = iota
	STATE_REQUEST_LINE
	STATE_HEADER
	STATE_PATH
	STATE_QUERY_STRING
	STATE_QUERY_KEY
	STATE_QUERY_VALUE
	STATE_HTTP_VERSION
	STATE_HEADER_NAME
	STATE_HEADER_VALUE
	STATE_NUMBER
	STATE_QUOTED_STRING
	STATE_DONE
)

// http methods that the lexer will recognise
var httpMethods = map[string]bool{
	"GET": true, "POST": true, "PUT": true,
	"DELETE": true, "PATCH": true, "HEAD": true,
	"OPTIONS": true, "CONNECT": true, "TRACE": true,
}

// dfa state and input buffer
type Lexer struct {
	input  []rune
	pos    int
	line   int
	col    int
	state  state
	tokens []Token
}

// creating lexer for the given http request string
func New(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
		pos:   0,
		line:  1,
		col:   1,
		state: STATE_START,
	}
}

// returns current rune
func (l *Lexer) peek() (rune, bool) {
	if l.pos >= len(l.input) {
		return 0, false
	}
	return l.input[l.pos], true
}

// consume the current rune and moves cursor forward
func (l *Lexer) advance() rune {
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return ch
}

func (l *Lexer) emit(t TokenType, value string, line, col int) {
	l.tokens = append(l.tokens, Token{Type: t, Value: value, Line: line, Col: col})
}

// runs the lexer and returns all tokens
func (l *Lexer) Tokenise() []Token {
	// first line is the method path http and its version
	l.lexRequestLine()

	// every subsequent non empty line is a header
	for {
		ch, ok := l.peek()
		if !ok {
			break
		}

		if ch == '\n' || ch == '\r' {
			l.skipNewlines()
			continue
		}
		l.lexHeaderLine()
	}
	l.emit(TOKEN_EOF, "", l.line, l.col)
	return l.tokens
}

func (l *Lexer) lexRequestLine() {
	word, line, col := l.readWord()
	if httpMethods[word] {
		l.emit(TOKEN_METHOD, word, line, col)
	} else {
		l.emit(TOKEN_UNKNOWN, word, line, col)
	}

	l.skipSpaces()
	l.lexPath()
	l.skipSpaces()
	l.lexHTTPVersion()
	l.skipNewlines()
}

// lexPath reads characters until a space or newline, splitting on '?' for query
func (l *Lexer) lexPath() {
	startLine, startCol := l.line, l.col
	var sb strings.Builder

	for {
		ch, ok := l.peek()
		if !ok || ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			break
		}

		if ch == '?' {
			l.emit(TOKEN_PATH, sb.String(), startLine, startCol)
			sb.Reset()
			qLine, qCol := l.line, l.col
			l.advance()
			l.emit(TOKEN_QUERY_SEP, "?", qLine, qCol)
			l.lexQueryString()
			return
		}
		sb.WriteRune(l.advance())
	}
	if sb.Len() > 0 {
		l.emit(TOKEN_PATH, sb.String(), startLine, startCol)
	}
}

func (l *Lexer) lexQueryString() {
	for {
		ch, ok := l.peek()
		if !ok || ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			return
		}

		keyLine, keyCol := l.line, l.col
		var key strings.Builder
		for {
			c, ok := l.peek()
			if !ok || c == '=' || c == '&' || c == ' ' || c == '\n' {
				break
			}
			key.WriteRune(l.advance())
		}
		l.emit(TOKEN_QUERY_KEY, key.String(), keyLine, keyCol)

		if c, ok := l.peek(); ok && c == '=' {
			eqLine, eqCol := l.line, l.col
			l.advance()
			l.emit(TOKEN_EQUALS, "=", eqLine, eqCol)
		}

		valLine, valCol := l.line, l.col
		var val strings.Builder
		for {
			c, ok := l.peek()
			if !ok || c == '&' || c == ' ' || c == '\n' {
				break
			}
			val.WriteRune(l.advance())
		}
		l.emit(TOKEN_QUERY_VALUE, val.String(), valLine, valCol)

		if c, ok := l.peek(); ok && c == '&' {
			ampLine, ampCol := l.line, l.col
			l.advance()
			l.emit(TOKEN_AMPERSAND, "&", ampLine, ampCol)
		} else {
			return
		}
	}
}

// reads the HTTP/x.y token
func (l *Lexer) lexHTTPVersion() {
	startLine, startCol := l.line, l.col
	word := l.readUntilSpace()
	if strings.HasPrefix(word, "HTTP/") {
		l.emit(TOKEN_HTTP_VERSION, word, startLine, startCol)
	}
}

// header-name: header value
func (l *Lexer) lexHeaderLine() {
	nameLine, nameCol := l.line, l.col
	var name strings.Builder
	for {
		ch, ok := l.peek()
		if !ok || ch == ':' || ch == '\n' || ch == '\r' {
			break
		}
		name.WriteRune(l.advance())
	}
	if name.Len() == 0 {
		return
	}

	l.emit(TOKEN_HEADER_NAME, name.String(), nameLine, nameCol)

	if ch, ok := l.peek(); ok && ch == ':' {
		cLine, cCol := l.line, l.col
		l.advance()
		l.emit(TOKEN_COLON, ":", cLine, cCol)
	}

	l.skipSpaces()
	l.lexHeaderValue()
	l.skipNewlines()
}

// lexHeaderValue reads the rest of the line and classifies sub-tokens
func (l *Lexer) lexHeaderValue() {
	// collect everything until end of line
	startLine, startCol := l.line, l.col
	var sb strings.Builder
	for {
		ch, ok := l.peek()
		if !ok || ch == '\n' || ch == '\r' {
			break
		}
		sb.WriteRune(l.advance())
	}
	raw := strings.TrimSpace(sb.String())
	if raw == "" {
		return
	}

	// classify value as known literal type
	switch {
	case raw == "true" || raw == "false":
		l.emit(TOKEN_BOOLEAN, raw, startLine, startCol)
	case isInteger(raw):
		l.emit(TOKEN_INTEGER, raw, startLine, startCol)
	case isFloat(raw):
		l.emit(TOKEN_FLOAT, raw, startLine, startCol)
	case len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"':
		l.emit(TOKEN_QUOTED_STRING, raw, startLine, startCol)
	default:
		l.emit(TOKEN_HEADER_VALUE, raw, startLine, startCol)
	}
}

func (l *Lexer) readWord() (string, int, int) {
	startLine, startCol := l.line, l.col
	var sb strings.Builder
	for {
		ch, ok := l.peek()
		if !ok || unicode.IsSpace(ch) {
			break
		}
		sb.WriteRune(l.advance())
	}
	return sb.String(), startLine, startCol
}

// read until whitespace or eof
func (l *Lexer) readUntilSpace() string {
	var sb strings.Builder
	for {
		ch, ok := l.peek()
		if !ok || ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			break
		}
		sb.WriteRune(l.advance())
	}
	return sb.String()
}

// skip horizontal spaces
func (l *Lexer) skipSpaces() {
	for {
		ch, ok := l.peek()
		if !ok || (ch != ' ' && ch != '\t') {
			return
		}
		l.advance()
	}
}

// skip new lines \n and \r
func (l *Lexer) skipNewlines() {
	for {
		ch, ok := l.peek()
		if !ok || (ch != '\n' && ch != '\r') {
			return
		}
		l.advance()
	}
}

// type classification helpers
func isInteger(s string) bool {
	if len(s) == 0 {
		return false
	}
	start := 0
	if s[0] == '-' || s[0] == '+' {
		start = 1
	}
	if start == len(s) {
		return false
	}
	for _, ch := range s[start:] {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func isFloat(s string) bool {
	if len(s) == 0 {
		return false
	}

	dots := 0
	start := 0
	if s[0] == '-' || s[0] == '+' {
		start = 1
	}

	if start == len(s) {
		return false
	}
	for _, ch := range s[start:] {
		if ch == '.' {
			dots++
			if dots > 1 {
				return false
			}
		} else if ch < '0' || ch > '9' {
			return false
		}
	}
	return dots == 1
}
