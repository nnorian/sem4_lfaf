# Parser & Abstract Syntax Tree

### Course: Formal Languages & Finite Automata
### Author: Kushnirenko Ecaterina FAF-243

---

## Theory

Parsing is the second stage of a compiler or interpreter pipeline. Where the lexer turns raw characters into a flat token stream, the parser imposes grammatical structure on that stream and produces a tree that reflects the nesting and hierarchy of the source language.

The tree produced is called an **Abstract Syntax Tree (AST)**. It is "abstract" because it omits details that are irrelevant to semantics, punctuation tokens such as `:`, `?`, `&`, and `=` are consumed during parsing and do not appear as nodes. What remains is a hierarchy of meaningful constructs: a request that contains a request line and a list of headers, a request line that contains a method, a path, an optional query string, and a version.

A **recursive-descent parser** is the most direct implementation of a context-free grammar. Each non-terminal in the grammar becomes a method; each method peeks at the next token, decides which production to apply, and calls sub-methods for sub-expressions. The token stream is consumed left-to-right with one token of look-ahead (LL(1)).

## Objectives

1. Get familiar with parsing and how it can be programmed.
2. Get familiar with the concept of an AST.
3. Extend the lexer with a `TokenType` enum and regular-expression-based classification.
4. Implement the AST data structures for the HTTP request format processed in lab 3.
5. Implement a recursive-descent parser that extracts syntactic information from the token stream.

## Implementation description

The implementation is split across four files. `tokens.go` and `lexer.go` live under `src/lexer/`; `ast.go` and `parser.go` live under `src/parser/`.

### TokenType enum and regex classification

`TokenType` is declared as a named `int` so it behaves like an enum. Every token category is a named constant produced with `iota`:

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
    TOKEN_UNKNOWN
    TOKEN_EOF
)
```

Classification of header values is done entirely with compiled regular expressions, not hand-written character checks:

```go
var (
    reMethod      = regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS|CONNECT|TRACE)$`)
    reHTTPVersion = regexp.MustCompile(`^HTTP/\d+(\.\d+)?$`)
    reInteger     = regexp.MustCompile(`^[+-]?\d+$`)
    reFloat       = regexp.MustCompile(`^[+-]?\d+\.\d+$`)
    reBoolean     = regexp.MustCompile(`^(true|false)$`)
    reQuoted      = regexp.MustCompile(`^".*"$`)
)
```

They are applied inside `lexHeaderValue()` after the raw value string is collected:

```go
func (l *Lexer) lexHeaderValue() {

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
```

### AST node types

The AST is defined in `ast.go`. Each node type maps directly to one grammatical construct of an HTTP request.

```
HTTPRequestNode
├── RequestLineNode
│   ├── Method       (string alias)
│   ├── Path         (string alias)
│   ├── QueryStringNode?
│   │   └── []QueryParamNode  { Key, Value string }
│   └── Version      (string alias)
└── []HeaderNode
    ├── Name   string
    └── Value  HeaderValue  (interface)
```

`HeaderValue` is a sealed interface — only the five concrete types below implement it, enforcing exhaustive handling in `switch` statements:

```go
type HeaderValue interface { headerValue() }

type StringHeaderValue      struct{ Value string }
type IntegerHeaderValue     struct{ Value string }
type FloatHeaderValue       struct{ Value string }
type BooleanHeaderValue     struct{ Value string }
type QuotedStringHeaderValue struct{ Value string }
```

Typed aliases for `Method`, `Path`, and `Version` give each field a distinct Go type even though all three are backed by `string`, which prevents accidental assignment between them.

### Parser

The parser is a hand-written recursive-descent parser in `parser.go`. It holds a slice of tokens and a cursor:

```go
type Parser struct {
    tokens []lexer.Token
    pos    int
}
```

Three helpers drive every production rule:

```go
// peek returns the current token without consuming it
func (p *Parser) peek() lexer.Token { ... }

// advance consumes the current token and returns it
func (p *Parser) advance() lexer.Token { ... }

// expect consumes the token if its type matches t,otherwise returns an error
func (p *Parser) expect(t lexer.TokenType) (lexer.Token, error) { ... }
```

#### Top-level: `Parse()`

`Parse()` is the entry point. It delegates to `parseRequestLine()` and then loops consuming headers until EOF:

```go
func (p *Parser) Parse() (*HTTPRequestNode, error) {
    rl, err := p.parseRequestLine()
    if err != nil {
        return nil, err
    }
    req := &HTTPRequestNode{RequestLine: rl}
    for p.peek().Type == lexer.TOKEN_HEADER_NAME {
        h, err := p.parseHeader()
        if err != nil {
            return nil, err
        }
        req.Headers = append(req.Headers, h)
    }
    return req, nil
}
```

#### Request line: `parseRequestLine()`

The method, path, optional query string, and HTTP version are consumed in order. The query string branch is taken only when a `TOKEN_QUERY_SEP` (`?`) is present:

```go
func (p *Parser) parseRequestLine() (*RequestLineNode, error) {
    methodTok, err := p.expect(lexer.TOKEN_METHOD)
    // ...
    pathTok, err := p.expect(lexer.TOKEN_PATH)
    // ...
    var qs *QueryStringNode
    if p.peek().Type == lexer.TOKEN_QUERY_SEP {
        p.advance()
        qs, err = p.parseQueryString()
    }
    versionTok, err := p.expect(lexer.TOKEN_HTTP_VERSION)
    // ...
    return &RequestLineNode{
        Method:      Method(methodTok.Value),
        Path:        Path(pathTok.Value),
        QueryString: qs,
        Version:     Version(versionTok.Value),
    }, nil
}
```

#### Query string: `parseQueryString()` and `parseQueryParam()`

Query params are consumed in a loop. After each param the parser checks for `TOKEN_AMPERSAND` to decide whether another param follows:

```go
func (p *Parser) parseQueryString() (*QueryStringNode, error) {
    qs := &QueryStringNode{}
    for p.peek().Type == lexer.TOKEN_QUERY_KEY {
        param, err := p.parseQueryParam()
        // ...
        qs.Params = append(qs.Params, param)
        if p.peek().Type == lexer.TOKEN_AMPERSAND {
            p.advance()
        } else {
            break
        }
    }
    return qs, nil
}
```

Each param is exactly `KEY = VALUE`:

```go
func (p *Parser) parseQueryParam() (*QueryParamNode, error) {
    keyTok, err := p.expect(lexer.TOKEN_QUERY_KEY)
    // ...
    if _, err = p.expect(lexer.TOKEN_EQUALS); err != nil { ... }
    valTok, err := p.expect(lexer.TOKEN_QUERY_VALUE)
    // ...
    return &QueryParamNode{Key: keyTok.Value, Value: valTok.Value}, nil
}
```

#### Headers: `parseHeader()` and `parseHeaderValue()`

Each header is `NAME : VALUE`. The value type is determined by a switch on the incoming token type — the token type was already decided by the lexer's regex classification, so no re-parsing is needed:

```go
func (p *Parser) parseHeaderValue() (HeaderValue, error) {
    tok := p.peek()
    switch tok.Type {
    case lexer.TOKEN_INTEGER:
        p.advance(); return IntegerHeaderValue{Value: tok.Value}, nil
    case lexer.TOKEN_FLOAT:
        p.advance(); return FloatHeaderValue{Value: tok.Value}, nil
    case lexer.TOKEN_BOOLEAN:
        p.advance(); return BooleanHeaderValue{Value: tok.Value}, nil
    case lexer.TOKEN_QUOTED_STRING:
        p.advance(); return QuotedStringHeaderValue{Value: tok.Value}, nil
    case lexer.TOKEN_HEADER_VALUE:
        p.advance(); return StringHeaderValue{Value: tok.Value}, nil
    default:
        return nil, fmt.Errorf("unexpected token %s (%q) as header value", ...)
    }
}
```

### AST pretty-printer

`ast.go` also contains `Print()`, a recursive function that renders the tree as an indented text diagram, used by `main.go` to visualise the output:

```go
func Print(root *HTTPRequestNode) string {
    out := "HTTPRequest\n"
    out += printRequestLine(root.RequestLine, 1)
    for _, h := range root.Headers {
        out += printHeader(h, 1)
    }
    return out
}
```

## Grammar (informal BNF)

The grammar recognised by the parser:

```
HTTPRequest ::= RequestLine Header*
RequestLine ::= METHOD PATH QueryString? HTTP_VERSION
QueryString ::= QUERY_SEP QueryParam (AMPERSAND QueryParam)*
QueryParam ::= QUERY_KEY EQUALS QUERY_VALUE
Header ::= HEADER_NAME COLON HeaderValue
HeaderValue ::= INTEGER | FLOAT | BOOLEAN | QUOTED_STRING | HEADER_VALUE
```

## Example output

Running `go run main.go` produces the following for all three sample requests.

### Sample 1 — GET with query string

**Input:**
```
GET /search?q=golang&page=2&debug=true HTTP/1.1
Host: example.com
Accept: application/json
```

**AST:**
```
HTTPRequest
  RequestLine
    Method("GET")
    Path("/search")
    QueryString
      QueryParam(key="q", value="golang")
      QueryParam(key="page", value="2")
      QueryParam(key="debug", value="true")
    Version("HTTP/1.1")
  Header("Host")
    Value[string]("example.com")
  Header("Accept")
    Value[string]("application/json")
```

The query string is parsed as three `QueryParamNode` children of a single `QueryStringNode`, rather than appearing as raw `QUERY_KEY`/`QUERY_VALUE` tokens. This is the structural information that parsing adds over lexing.

### Sample 2 — POST with typed header values

**Input:**
```
POST /api/users HTTP/1.1
Host: api.example.com
Content-Type: application/json
Content-Length: 42
Authorization: Bearer eyJhbGciOiJIUzI1NiJ9
```

**AST:**
```
HTTPRequest
  RequestLine
    Method("POST")
    Path("/api/users")
    Version("HTTP/1.1")
  Header("Host")
    Value[string]("api.example.com")
  Header("Content-Type")
    Value[string]("application/json")
  Header("Content-Length")
    Value[integer]("42")
  Header("Authorization")
    Value[string]("Bearer eyJhbGciOiJIUzI1NiJ9")
```

`Content-Length: 42` is stored as `IntegerHeaderValue`, not a plain string. The type was assigned by the lexer's regex and preserved through the AST node, so any downstream consumer can call type-specific logic without reparsing.

### Sample 3 — DELETE with FLOAT and BOOLEAN headers

**Input:**
```
DELETE /api/items/99 HTTP/2
Host: api.example.com
X-Rate-Limit: 3.14
X-Enabled: true
```

**AST:**
```
HTTPRequest
  RequestLine
    Method("DELETE")
    Path("/api/items/99")
    Version("HTTP/2")
  Header("Host")
    Value[string]("api.example.com")
  Header("X-Rate-Limit")
    Value[float]("3.14")
  Header("X-Enabled")
    Value[boolean]("true")
```

`X-Rate-Limit` maps to `FloatHeaderValue` and `X-Enabled` maps to `BooleanHeaderValue`, both determined solely by the regular expressions in the lexer without any additional logic in the parser.

## Conclusions

A recursive-descent parser is a natural fit for a grammar as regular as HTTP request syntax: each production rule translates to one method, and the one-token look-ahead is always sufficient to decide which branch to take. Combining it with the DFA-based lexer from lab 3 produces a clean two-stage pipeline where the lexer handles character-level classification and the parser handles structural relationships between tokens.

Using a sealed `HeaderValue` interface for the five value types means the Go compiler enforces exhaustive handling at every switch site. Typed string aliases for `Method`, `Path`, and `Version` prevent fields from being swapped accidentally. Both design choices shift error detection from runtime to compile time.

The clear separation between lexing and parsing also simplifies the regex-based token classification: the lexer applies six compiled regexes once per token, stores the result as a `TokenType`, and the parser acts only on that type — neither stage needs to re-examine the raw string.

## References

1. Formal Languages and Finite Automata — course materials, TUM
2. Crafting Interpreters, R. Nystrom — https://craftinginterpreters.com (Chapters 5–6: Representing Code, Parsing Expressions)
3. HTTP/1.1 Specification — RFC 9110
4. The Go Programming Language Specification — https://go.dev/ref/spec
