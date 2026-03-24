package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

//we enumerate the type of nodes
type NodeType int

const (
	NodeLiteral NodeType = iota
	NodeConcat
	NodeAlteration
	NodeQuantified
)

//the repetition semantics
type QuantifierKind int

const (
	QuantExact QuantifierKind = iota
	QuantOptional
	QuantStar
	QuantPlus
)

//upper bound
const MaxRepeat = 5

type Node struct {
	Type     NodeType
	Char     rune
	Children []*Node
	Child    *Node
	QKind    QuantifierKind
	QCount   int
}

//lexer
//taggign each token
type tokenKind int

const (
	tokChar tokenKind = iota
	tokLParen
	tokRParen
	tokPipe
	tokQMark
	tokStar
	tokPlus
	tokPower
	tokEOF
)

type token struct {
	kind tokenKind
	ch   rune
	val  int
}

//super script digit map
var superDigit = map[rune]int{
	'⁰': 0, '¹': 1, '²': 2, '³': 3, '⁴': 4,
	'⁵': 5, '⁶': 6, '⁷': 7, '⁸': 8, '⁹': 9,
}

//lexer conversion from regex string to slice of tokens
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
		case r == '|':
			tokens = append(tokens, token{kind: tokPipe})
			i++
		case r == '?':
			tokens = append(tokens, token{kind: tokQMark})
			i++
		case r == '*':
			tokens = append(tokens, token{kind: tokStar})
			i++
		case r == '+' || r == '⁺':
			tokens = append(tokens, token{kind: tokPlus})
			i++
		default:
			// Check for superscript digits
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
				// ordinary character literal
				tokens = append(tokens, token{kind: tokChar, ch: r})
				i++
			}
		}
	}
	tokens = append(tokens, token{kind: tokEOF})
	return tokens
}

// recursive-descent parser
type Parser struct {
	tokens []token
	pos    int
}

func NewParser(input string) *Parser {
	return &Parser{tokens: lex(input)}
}

func (p *Parser) peek() token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return token{kind: tokEOF}
}

func (p *Parser) advance() token {
	t := p.peek()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return t
}

func (p *Parser) expect(k tokenKind) token {
	t := p.advance()
	if t.kind != k {
		panic(fmt.Sprintf("parser: expected token kind %d, got %d", k, t.kind))
	}
	return t
}

// parsing the entry point, returns the root node
func (p *Parser) Parse() *Node {
	node := p.parseRegex()
	p.expect(tokEOF)
	return node
}

// concatenation of one or more terms
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

func (p *Parser) parseTerm() *Node {
	atom := p.parseAtom()
	// check for quantifier
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

func (p *Parser) parseAlteration() *Node {
	first := p.parseRegex()
	if p.peek().kind != tokPipe {
		return first
	}
	alts := []*Node{first}
	for p.peek().kind == tokPipe {
		p.advance()
		alts = append(alts, p.parseRegex())
	}
	return &Node{Type: NodeAlteration, Children: alts}
}

//thingie to pretty print
func (n *Node) String() string {
	switch n.Type {
	case NodeLiteral:
		return fmt.Sprintf("Literal('%c')", n.Char)
	case NodeConcat:
		parts := make([]string, len(n.Children))
		for i, c := range n.Children {
			parts[i] = c.String()
		}
		return fmt.Sprintf("Concat(%s)", strings.Join(parts, ", "))
	case NodeAlteration:
		parts := make([]string, len(n.Children))
		for i, c := range n.Children {
			parts[i] = c.String()
		}
		return fmt.Sprintf("Alt(%s)", strings.Join(parts, " | "))
	case NodeQuantified:
		return fmt.Sprintf("%s{%s}", n.Child.String(), quantStr(n.QKind, n.QCount))
	}
	return "?"
}

func quantStr(qk QuantifierKind, count int) string {
	switch qk {
	case QuantExact:
		return fmt.Sprintf("%d", count)
	case QuantOptional:
		return "?"
	case QuantStar:
		return "*"
	case QuantPlus:
		return "+"
	}
	return "?"
}

//generator to make valid strings
func Generate(n *Node) string {
	var sb strings.Builder
	generate(n, &sb)
	return sb.String()
}

func generate(n *Node, sb *strings.Builder) {
	switch n.Type {
	case NodeLiteral:
		sb.WriteRune(n.Char)

	case NodeConcat:
		for _, child := range n.Children {
			generate(child, sb)
		}

	case NodeAlteration:
		idx := rand.Intn(len(n.Children))
		generate(n.Children[idx], sb)

	case NodeQuantified:
		count := pickCount(n.QKind, n.QCount)
		for i := 0; i < count; i++ {
			generate(n.Child, sb)
		}
	}
}

func pickCount(qk QuantifierKind, exact int) int {
	switch qk {
	case QuantExact:
		return exact
	case QuantOptional:
		//random 0 or 1
		return rand.Intn(2)
	case QuantStar:
		return rand.Intn(MaxRepeat + 1)
	case QuantPlus:
		return 1 + rand.Intn(MaxRepeat)
	}
	return 1
}

// bonus point
func GenerateWithTrace(n *Node) (string, []string) {
	var sb strings.Builder
	var trace []string
	step := 0
	traceGen(n, &sb, &trace, &step)
	trace = append(trace, fmt.Sprintf("Step %d: Result: \"%s\"", step+1, sb.String()))
	return sb.String(), trace
}

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

func qDesc(qk QuantifierKind) string {
	switch qk {
	case QuantExact:
		return "exact"
	case QuantOptional:
		return "?"
	case QuantStar:
		return "*"
	case QuantPlus:
		return "+"
	}
	return "?"
}

//main
func main() {
	rand.Seed(time.Now().UnixNano())
	// v2
	// unicode superscripts as exponents
	regexes := []string{
		"M?N²(O|P)³Q*R⁺",
		"(X|Y|Z)³8⁺(9|0)",
		"(H|I)(J|K)L*N?",
	}

	separator := strings.Repeat("─", 60)

	for idx, re := range regexes {
		fmt.Printf("  Regex #%d:  %s\n", idx+1, re)
		fmt.Printf("  Length (runes): %d\n", utf8.RuneCountInString(re))

		//parse
		parser := NewParser(re)
		ast := parser.Parse()

		//show ast
		fmt.Printf("\n  ► AST: %s\n", ast.String())

		//generate samples
		fmt.Printf("\n  ► Generated samples (10 strings):\n")
		for i := 0; i < 10; i++ {
			fmt.Printf("      %2d) %s\n", i+1, Generate(ast))
		}

		//bonus
		result, trace := GenerateWithTrace(ast)
		fmt.Printf("\n  ► Trace:\n")
		for _, line := range trace {
			fmt.Printf("      %s\n", line)
		}
		fmt.Printf("  Result: %s\n", result)

		fmt.Printf("\n%s\n", separator)
	}
}
