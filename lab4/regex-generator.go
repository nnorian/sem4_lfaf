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
const(
	NodeLiteral NodeType=iota
	NodeConcat 
	NodeAlteration
	NodeQuantified
)

//the repetition semantics
type QuantifierKind int
const(
	QuantExact QuantifierKind = iota
	QuantOptional 
	QuantStar
	QuantPlus

)

//upper bound
consr MaxRepeat = 5 

type Node struct {
	Type NodeTypeChar rune
	Children []*Node
	Child *Node
	QKind QuantifierKind
	QCount int
}

//lexer
//taggign each token
type tokenKind int
const(
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
	kind tokenKindch rune
	val int
}

//super script digit map 
var superDigit = map[rune]int{
	'⁰': 0, '¹': 1, '²': 2, '³': 3, '⁴': 4,
	'⁵': 5, '⁶': 6, '⁷': 7, '⁸': 8, '⁹': 9,
}

//lexer conversion fro-m regex string to slike of tokens
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
