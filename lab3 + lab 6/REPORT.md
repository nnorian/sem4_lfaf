# Lexer & Scanner

### Course: Formal Languages & Finite Automata
### Author: Kushnirenko Ecaterina FAF-243


---

## Theory

Lexical analysis is the first stage of a compiler or interpreter pipeline. Its job is to read a raw string of characters and produce a flat sequence of tokens structured units that carry both a type or otherwise category and a value aka the matched text.


The lexer implemented here is a Deterministic Finite Automaton . Every character consumed causes a state transition; the current state determines how the next character is interpreted. This maps directly to the formal definition of a DFA: `δ(state, input) → next_state`.

## Objectives

1. Understand what lexical analysis is.
2. Get familiar with the inner workings of a lexer / scanner / tokenizer.
3. Implement a sample lexer and show how it works.

## Implementation description

The lexer is split across two files: `tokens.go` defines all token types, and `lexer.go` contains the DFA logic.

### Token Types

Tokens are defined as a Go `iota` enum grouped into four logical sections:

```go
type TokenType int

const (
    TOKEN_METHOD TokenType = iota 
    TOKEN_PATH
    TOKEN_HTTP_VERSION
    TOKEN_QUERY_SEP
    TOKEN_QUERY_KEY
    TOKEN_QUERY_VALUE
    TOKEN_AMPERSAND
    TOKEN_EQUALS
    TOKEN_HEADER_NAME
    TOKEN_COLON
    TOKEN_HEADER_VALUE
    TOKEN_INTEGER
    TOKEN_FLOAT
    TOKEN_BOOLEAN
    TOKEN_QUOTED_STRING
    TOKEN_NEWLINE
    TOKEN_WHITESPACE
    TOKEN_UNKNOWN
    TOKEN_EOF
)
```

Each `Token` struct stores the type, matched value, and source position:

```go
type Token struct {
    Type  TokenType
    Value string
    Line  int
    Col   int
}
```

### DFA States

The automaton declares its states explicitly as a named type:

```go
type state int

const (
    STATE_START state = iota
    STATE_REQUEST_LINE
    STATE_PATH
    STATE_QUERY_STRING
    STATE_QUERY_KEY
    STATE_QUERY_VALUE
    STATE_HTTP_VERSION
    STATE_HEADER
    STATE_HEADER_NAME
    STATE_HEADER_VALUE
    STATE_NUMBER
    STATE_QUOTED_STRING
    STATE_DONE
)
```

### Core Lexer Loop

The `Tokenise()` method drives the DFA. The first line of any HTTP request is always the request line, so it is handled separately. Every subsequent non-empty line is treated as a header:

```go
func (l *Lexer) Tokenise() []Token {
    l.lexRequestLine()

    for {
        ch, ok := l.peek()
        if !ok { break }
        if ch == '\n' || ch == '\r' {
            l.skipNewlines()
            continue
        }
        l.lexHeaderLine()
    }

    l.emit(TOKEN_EOF, "", l.line, l.col)
    return l.tokens
}
```

### Request Line Lexer

The request line is tokenised in three ordered steps, the method, path (with optional query string), then HTTP version:

```go
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
```

### Query String Lexer

Once `lexPath()` encounters a `?`, it emits `TOKEN_QUERY_SEP` and delegates to `lexQueryString()`, which loops over `key=value&...` pairs:

```go
func (l *Lexer) lexQueryString() {
    for {
        // read key
        l.emit(TOKEN_QUERY_KEY, key.String(), keyLine, keyCol)
        // read '='
        l.emit(TOKEN_EQUALS, "=", eqLine, eqCol)
        // read value
        l.emit(TOKEN_QUERY_VALUE, val.String(), valLine, valCol)
        // stop or consume '&' and loop
        if c == '&' {
            l.emit(TOKEN_AMPERSAND, "&", ampLine, ampCol)
        } else {
            return
        }
    }
}
```

### Header Value Type Inference

After reading the raw header value string, the lexer classifies it into the most specific literal type available:

```go
func (l *Lexer) lexHeaderValue() {
    //collect raw string
    switch {
    case raw == "true" || raw == "false":
        l.emit(TOKEN_BOOLEAN, raw, startLine, startCol)
    case isInteger(raw):
        l.emit(TOKEN_INTEGER, raw, startLine, startCol)
    case isFloat(raw):
        l.emit(TOKEN_FLOAT, raw, startLine, startCol)
    default:
        l.emit(TOKEN_HEADER_VALUE, raw, startLine, startCol)
    }
}
```

## DFA Transition Table

The table below formalises every state transition of the automaton. `σ` represents a character class; the rightmost column shows the token emitted on that transition (if any).

| Current State        | Input `σ`                | Next State           | Token Emitted                          |
|----------------------|--------------------------|----------------------|----------------------------------------|
| `STATE_START`        | letter                   | `STATE_REQUEST_LINE` | —                                      |
| `STATE_REQUEST_LINE` | known HTTP method word   | `STATE_PATH`         | `TOKEN_METHOD`                         |
| `STATE_REQUEST_LINE` | unknown word             | `STATE_PATH`         | `TOKEN_UNKNOWN`                        |
| `STATE_PATH`         | `/`, letter, digit, `-`  | `STATE_PATH`         | —                                      |
| `STATE_PATH`         | `?`                      | `STATE_QUERY_STRING` | `TOKEN_PATH`, `TOKEN_QUERY_SEP`        |
| `STATE_PATH`         | ` ` (space)              | `STATE_HTTP_VERSION` | `TOKEN_PATH`                           |
| `STATE_QUERY_STRING` | letter / digit           | `STATE_QUERY_KEY`    | —                                      |
| `STATE_QUERY_KEY`    | letter / digit           | `STATE_QUERY_KEY`    | —                                      |
| `STATE_QUERY_KEY`    | `=`                      | `STATE_QUERY_VALUE`  | `TOKEN_QUERY_KEY`, `TOKEN_EQUALS`      |
| `STATE_QUERY_VALUE`  | letter / digit           | `STATE_QUERY_VALUE`  | —                                      |
| `STATE_QUERY_VALUE`  | `&`                      | `STATE_QUERY_KEY`    | `TOKEN_QUERY_VALUE`, `TOKEN_AMPERSAND` |
| `STATE_QUERY_VALUE`  | ` ` (space)              | `STATE_HTTP_VERSION` | `TOKEN_QUERY_VALUE`                    |
| `STATE_HTTP_VERSION` | `HTTP/` prefix + digit   | `STATE_HTTP_VERSION` | —                                      |
| `STATE_HTTP_VERSION` | `\n`                     | `STATE_HEADER`       | `TOKEN_HTTP_VERSION`                   |
| `STATE_HEADER`       | letter / `-`             | `STATE_HEADER_NAME`  | —                                      |
| `STATE_HEADER_NAME`  | letter / `-`             | `STATE_HEADER_NAME`  | —                                      |
| `STATE_HEADER_NAME`  | `:`                      | `STATE_HEADER_VALUE` | `TOKEN_HEADER_NAME`, `TOKEN_COLON`     |
| `STATE_HEADER_VALUE` | digit                    | `STATE_NUMBER`       | —                                      |
| `STATE_HEADER_VALUE` | `"`                      | `STATE_QUOTED_STRING`| —                                      |
| `STATE_HEADER_VALUE` | `true` / `false`         | `STATE_HEADER`       | `TOKEN_BOOLEAN`                        |
| `STATE_HEADER_VALUE` | other printable char     | `STATE_HEADER_VALUE` | —                                      |
| `STATE_HEADER_VALUE` | `\n`                     | `STATE_HEADER`       | `TOKEN_HEADER_VALUE`                   |
| `STATE_NUMBER`       | digit                    | `STATE_NUMBER`       | —                                      |
| `STATE_NUMBER`       | `.`                      | `STATE_NUMBER`       | — (promotes result to FLOAT)           |
| `STATE_NUMBER`       | `\n`                     | `STATE_HEADER`       | `TOKEN_INTEGER` or `TOKEN_FLOAT`       |
| `STATE_QUOTED_STRING`| any char except `"`      | `STATE_QUOTED_STRING`| —                                      |
| `STATE_QUOTED_STRING`| `"`                      | `STATE_HEADER_VALUE` | `TOKEN_QUOTED_STRING`                  |
| `STATE_HEADER`       | EOF                      | `STATE_DONE`         | `TOKEN_EOF`                            |

## Example case

Running `go run main.go` produces the following output for three sample HTTP requests.

### Sample 1 — GET with query string

**Input:**
```
GET /search?q=golang&page=2&debug=true HTTP/1.1
Host: example.com
Accept: application/json
```

**Tokens:**
```
METHOD               | "GET"               | line 1 col 1
PATH                 | "/search"           | line 1 col 5
QUERY_SEP            | "?"                 | line 1 col 12
QUERY_KEY            | "q"                 | line 1 col 13
EQUALS               | "="                 | line 1 col 14
QUERY_VALUE          | "golang"            | line 1 col 15
AMPERSAND            | "&"                 | line 1 col 21
QUERY_KEY            | "page"              | line 1 col 22
EQUALS               | "="                 | line 1 col 26
QUERY_VALUE          | "2"                 | line 1 col 27
AMPERSAND            | "&"                 | line 1 col 28
QUERY_KEY            | "debug"             | line 1 col 29
EQUALS               | "="                 | line 1 col 34
QUERY_VALUE          | "true"              | line 1 col 35
HTTP_VERSION         | "HTTP/1.1"          | line 1 col 40
HEADER_NAME          | "Host"              | line 2 col 1
COLON                | ":"                 | line 2 col 5
HEADER_VALUE         | "example.com"       | line 2 col 7
HEADER_NAME          | "Accept"            | line 3 col 1
COLON                | ":"                 | line 3 col 7
HEADER_VALUE         | "application/json"  | line 3 col 9
```

### Sample 2 — POST with typed header values

**Input:**
```
POST /api/users HTTP/1.1
Content-Length: 42
Authorization: Bearer eyJhbGciOiJIUzI1NiJ9
```

`Content-Length: 42` produces `TOKEN_INTEGER` instead of a generic `TOKEN_HEADER_VALUE`, demonstrating that the lexer performs value-type inference at the lexical level.

### Sample 3 — DELETE with FLOAT and BOOLEAN headers

**Input:**
```
DELETE /api/items/99 HTTP/2
X-Rate-Limit: 3.14
X-Enabled: true
```

`X-Rate-Limit: 3.14` into `TOKEN_FLOAT`, `X-Enabled: true` into `TOKEN_BOOLEAN`. The DFA transitions through `STATE_NUMBER` and detects the `.` to promote the result from integer to float

## Conclusions
 A DFA-based lexer is the cleanest implementation strategy for a structured text format like HTTP because every valid input character has an unambiguous next state, and the state machine can be written directly from the formal transition table. Splitting the lexer into per-section sub-routines (`lexRequestLine`, `lexPath`, `lexQueryString`, `lexHeaderLine`) keeps each piece small and independently testable without losing the DFA property.Value-type inference at the lexical level (INTEGER, FLOAT, BOOLEAN) adds useful semantic information early in the pipeline without requiring a separate parsing pass. Tracking line and column numbers in every token costs almost nothing but is essential for meaningful error messages in any real compiler or tool.

## References

1. Formal Languages and Finite Automata — course materials, TUM
2. Crafting Interpreters, R. Nystrom — https://craftinginterpreters.com (Chapter 4: Scanning)
3. HTTP/1.1 Specification — RFC 9110