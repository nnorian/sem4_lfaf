# Regular Expression Processor

### Course: Formal Languages & Finite Automata
### Author: Kushnirenko Ecaterina FAF-243


---

## Theory

A regular expression is a compact notation for describing a regular language. Every regex defines a set of strings through three core operations: concatenation (sequencing), alternation (choice via `|`), and quantification (repetition via `?`, `*`, `+`, or an exact count). These operations map directly to the algebraic properties of regular languages studied in formal language theory.

To process a regex programmatically we need a lexer that converts the raw string into tokens, a parser that builds an Abstract Syntax Tree (AST) from those tokens, and a generator that walks the AST to produce strings belonging to the language. The parser used here is a recursive-descent parser, a top-down approach where each grammar rule maps to a function.

## Objectives

1. Write a regex processor that can parse regular expressions and generate valid strings from them.
2. Implement a lexer that tokenises regex strings, including support for Unicode superscript exponents (`²`, `³`, `⁺`, etc.).
3. Build a recursive-descent parser that produces an AST from the token stream.
4. Implement a string generator that walks the AST and produces random valid strings.
5. (Bonus) Add step-by-step tracing of the generation process.

## Implementation description

The implementation is a single Go file split into four logical sections: node types, lexer, parser, and generator.

### Node Types

The AST uses four node kinds, defined as a Go `iota` enum:

```go
type NodeType int

const (
    NodeLiteral NodeType = iota
    NodeConcat
    NodeAlteration
    NodeQuantified
)
```

Each `Node` stores its type and the relevant fields: `Char` for literals, `Children` for concatenation and alternation, `Child` plus `QKind`/`QCount` for quantified nodes:

```go
type Node struct {
    Type     NodeType
    Char     rune
    Children []*Node
    Child    *Node
    QKind    QuantifierKind
    QCount   int
}
```

Quantifiers are similarly enumerated:

```go
type QuantifierKind int

const (
    QuantExact QuantifierKind = iota
    QuantOptional
    QuantStar
    QuantPlus
)
```

### Lexer

The lexer converts a regex string into a slice of tokens. It handles parentheses, pipe, `?`, `*`, `+`/`⁺`, Unicode superscript digit sequences (interpreted as exact repetition counts), and plain character literals:

```go
func lex(input string) []token {
    var tokens []token
    runes := []rune(input)
    i := 0
    for i < len(runes) {
        r := runes[i]
        switch {
        case r == '(':
            tokens = append(tokens, token{kind: tokLParen})
            i++
        case r == ')':
            tokens = append(tokens, token{kind: tokRParen})
            i++
    
        default:
            if _, ok := superDigit[r]; ok {
                val := 0
                for i < len(runes) {
                    if d, ok := superDigit[runes[i]]; ok {
                        val = val*10 + d
                        i++
                    } else {
                        break
                    }
                }
                tokens = append(tokens, token{kind: tokPower, val: val})
            } else {
                tokens = append(tokens, token{kind: tokChar, ch: r})
                i++
            }
        }
    }
    tokens = append(tokens, token{kind: tokEOF})
    return tokens
}
```

The superscript digit map allows multi-digit exponents like `¹²` to be parsed as the integer 12:

```go
var superDigit = map[rune]int{
    '⁰': 0, '¹': 1, '²': 2, '³': 3, '⁴': 4,
    '⁵': 5, '⁶': 6, '⁷': 7, '⁸': 8, '⁹': 9,
}
```

### Recursive-Descent Parser

The parser consumes the token stream and builds the AST. The grammar is:

```
regex       → term+
term        → atom quantifier?
atom        → '(' alteration ')' | CHAR
alteration  → regex ('|' regex)*
quantifier  → '?' | '*' | '+' | POWER
```

Each grammar rule maps to a method on the `Parser` struct:

```go
func (p *Parser) parseRegex() *Node {
    var parts []*Node
    for {
        t := p.peek()
        if t.kind == tokEOF || t.kind == tokRParen || t.kind == tokPipe {
            break
        }
        parts = append(parts, p.parseTerm())
    }
    if len(parts) == 1 {
        return parts[0]
    }
    return &Node{Type: NodeConcat, Children: parts}
}
```

`parseTerm` reads an atom and then checks for an optional trailing quantifier:

```go
func (p *Parser) parseTerm() *Node {
    atom := p.parseAtom()
    switch p.peek().kind {
    case tokQMark:
        p.advance()
        return &Node{Type: NodeQuantified, Child: atom, QKind: QuantOptional}
    case tokStar:
        p.advance()
        return &Node{Type: NodeQuantified, Child: atom, QKind: QuantStar}
    case tokPlus:
        p.advance()
        return &Node{Type: NodeQuantified, Child: atom, QKind: QuantPlus}
    case tokPower:
        t := p.advance()
        return &Node{Type: NodeQuantified, Child: atom, QKind: QuantExact, QCount: t.val}
    }
    return atom
}
```

`parseAtom` handles grouped sub-expressions in parentheses and single character literals:

```go
func (p *Parser) parseAtom() *Node {
    t := p.peek()
    if t.kind == tokLParen {
        p.advance()
        node := p.parseAlteration()
        p.expect(tokRParen)
        return node
    }
    if t.kind == tokChar {
        p.advance()
        return &Node{Type: NodeLiteral, Char: t.ch}
    }
    panic(fmt.Sprintf("parser: unexpected token kind %d at pos %d", t.kind, p.pos))
}
```

### String Generator

The generator walks the AST and builds a string. For alternation nodes it picks a random branch; for quantified nodes it picks a random repeat count within the allowed range (bounded by `MaxRepeat = 5`):

```go
func pickCount(qk QuantifierKind, exact int) int {
    switch qk {
    case QuantExact:
        return exact
    case QuantOptional:
        return rand.Intn(2)
    case QuantStar:
        return rand.Intn(MaxRepeat + 1)
    case QuantPlus:
        return 1 + rand.Intn(MaxRepeat)
    }
    return 1
}
```

The `Generate` function delegates to `GenerateWithTrace` and discards the trace, avoiding code duplication:

```go
func Generate(n *Node) string {
    result, _ := GenerateWithTrace(n)
    return result
}
```

### Bonus: Step-by-Step Trace

`GenerateWithTrace` records every decision made during generation, showing which alternatives were chosen, how many repetitions were picked, and which literals were appended:

```go
func traceGen(n *Node, sb *strings.Builder, trace *[]string, step *int) {
    *step++
    switch n.Type {
    case NodeLiteral:
        sb.WriteRune(n.Char)
        *trace = append(*trace, fmt.Sprintf("Step %d: Append literal '%c'", *step, n.Char))
    case NodeConcat:
        *trace = append(*trace, fmt.Sprintf("Step %d: Process concatenation (%d parts)", *step, len(n.Children)))
        for _, c := range n.Children {
            traceGen(c, sb, trace, step)
        }
    case NodeAlteration:
        idx := rand.Intn(len(n.Children))
        *trace = append(*trace, fmt.Sprintf("Step %d: Choose from alternatives → option %d", *step, idx+1))
        traceGen(n.Children[idx], sb, trace, step)
    case NodeQuantified:
        count := pickCount(n.QKind, n.QCount)
        *trace = append(*trace, fmt.Sprintf("Step %d: Repeat ×%d (%s)", *step, count, qDesc(n.QKind)))
        for i := 0; i < count; i++ {
            traceGen(n.Child, sb, trace, step)
        }
    }
}
```

## Example case

The program processes three regex patterns using Unicode superscript notation:

### Regex 1: `M?N²(O|P)³Q*R⁺`

**AST:**
```
Concat(Literal('M'){?}, Literal('N'){2}, Alt(Literal('O') | Literal('P')){3}, Literal('Q'){*}, Literal('R'){+})
```

**Generated samples:**
```
 1) NNPOOQQQQQRRRRR
 2) MNNOOORRR
 3) MNNPPPQQQQQRR
 4) MNNOPOQQQQRRR
 5) NNOOORRRR
 6) NNPPOQQQQRRRRR
 7) MNNOPOQQQQRRR
 8) NNOOPR
 9) MNNPOOQRRRR
10) MNNPOOQQQQQRR
```

**Trace (one generation):**
```
Step 1: Process concatenation (5 parts)
Step 2: Repeat ×1 (?)
Step 3: Append literal 'M'
Step 4: Repeat ×2 (exact)
Step 5: Append literal 'N'
Step 6: Append literal 'N'
Step 7: Repeat ×3 (exact)
Step 8: Choose from alternatives → option 1
Step 9: Append literal 'O'
Step 10: Choose from alternatives → option 1
Step 11: Append literal 'O'
Step 12: Choose from alternatives → option 1
Step 13: Append literal 'O'
Step 14: Repeat ×4 (*)
Step 15: Append literal 'Q'
Step 16: Append literal 'Q'
Step 17: Append literal 'Q'
Step 18: Append literal 'Q'
Step 19: Repeat ×3 (+)
Step 20: Append literal 'R'
Step 21: Append literal 'R'
Step 22: Append literal 'R'
Step 23: Result: "MNNOOOQQQQRRR"
```

### Regex 2: `(X|Y|Z)³8⁺(9|0)`

**AST:**
```
Concat(Alt(Literal('X') | Literal('Y') | Literal('Z')){3}, Literal('8'){+}, Alt(Literal('9') | Literal('0')))
```

**Generated samples:**
```
 1) ZYY889
 2) YYY888880
 3) ZZZ888889
 4) ZYY8880
 5) ZZZ89
 6) YZY88880
 7) ZZZ8889
 8) ZXX8889
 9) ZYX888889
10) YXX880
```

### Regex 3: `(H|I)(J|K)L*N?`

**AST:**
```
Concat(Alt(Literal('H') | Literal('I')), Alt(Literal('J') | Literal('K')), Literal('L'){*}, Literal('N'){?})
```

**Generated samples:**
```
 1) IJLLLLLN
 2) IKLLLLL
 3) IKLN
 4) IKLLLLL
 5) IJLLLN
 6) HKLN
 7) HKLLLLN
 8) IJLLLLN
 9) HKLLLLLN
10) HJ
```

Every generated string matches the structure defined by its regex: optional parts appear or not, exact exponents produce the right count, `+` always gives at least one repetition, and alternations pick from the correct set.

## Conclusions

The recursive-descent approach maps naturally to the regex grammar because each production rule becomes its own function, making the parser easy to read and extend. Using Unicode superscripts as exponent notation is a compact alternative to the traditional `{n}` syntax and the lexer handles multi-digit superscripts correctly by accumulating digits in a loop. Collapsing `Generate` into a wrapper around `GenerateWithTrace` eliminates duplicated tree-walking logic. The step-by-step trace makes it possible to verify that the generator is making correct decisions at each node of the AST.

## References

1. Formal Languages and Finite Automata — course materials, TUM
2. Crafting Interpreters, R. Nystrom — https://craftinginterpreters.com (Chapter 4: Scanning)
3. Hopcroft, Motwani, Ullman — Introduction to Automata Theory, Languages, and Computation
