package lexer

import (
	"regexp"
	"strings"
)

// compiled regexes used for token type classification
var (
	reMethod      = regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS|CONNECT|TRACE)$`)
	reHTTPVersion = regexp.MustCompile(`^HTTP/\d+(\.\d+)?$`)
	reInteger     = regexp.MustCompile(`^[+-]?\d+$`)
	reFloat       = regexp.MustCompile(`^[+-]?\d+\.\d+$`)
	reBoolean     = regexp.MustCompile(`^(true|false)$`)
	reQuoted      = regexp.MustCompile(`^".*"$`)
)

// lexer holds the DFA input buffer and cursor state
type Lexer struct {
	input  []rune
	pos    int
	line   int
	col    int
	tokens []Token
}


func New(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
		pos:   0,
		line:  1,
		col:   1,
	}
}

// peek returns the current rune without consuming it
func (l *Lexer) peek() (rune, bool) {
	if l.pos >= len(l.input) {
		return 0, false
	}
	return l.input[l.pos], true
}

// advance consumes the current rune and moves the cursor forward
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

// tokenise runs the lexer and returns all tokens
func (l *Lexer) Tokenise() []Token {
	l.lexRequestLine()

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
	if reMethod.MatchString(word) {
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

func (l *Lexer) lexHTTPVersion() {
	startLine, startCol := l.line, l.col
	word := l.readUntilSpace()
	if reHTTPVersion.MatchString(word) {
		l.emit(TOKEN_HTTP_VERSION, word, startLine, startCol)
	}
}

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

func (l *Lexer) lexHeaderValue() {
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

	switch {
	case reBoolean.MatchString(raw):
		l.emit(TOKEN_BOOLEAN, raw, startLine, startCol)
	case reInteger.MatchString(raw):
		l.emit(TOKEN_INTEGER, raw, startLine, startCol)
	case reFloat.MatchString(raw):
		l.emit(TOKEN_FLOAT, raw, startLine, startCol)
	case reQuoted.MatchString(raw):
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
		if !ok || ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			break
		}
		sb.WriteRune(l.advance())
	}
	return sb.String(), startLine, startCol
}

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

func (l *Lexer) skipSpaces() {
	for {
		ch, ok := l.peek()
		if !ok || (ch != ' ' && ch != '\t') {
			return
		}
		l.advance()
	}
}

func (l *Lexer) skipNewlines() {
	for {
		ch, ok := l.peek()
		if !ok || (ch != '\n' && ch != '\r') {
			return
		}
		l.advance()
	}
}