package parser 

import (
	"fmt"
)

type Parser struct {
	tokens []lexer.Token
	pos int
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) peek() lexer.Token {
	if p.pos >= len(p.tokens){
		return lexer.Tokens(Type: lexer.TOKEN.EOF)

	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() lexer.Token{
	tok := p.tokens[p.pos]
	p.pos++
	return tok
}

func (p *Parser) expect (t lexer,TokenType) (lexer.Token, error){
	tok := peek()
	if tok.Type != t {
		return tok, fmt.Errorf(
			"parse error at line %d col %D: expected %s, got %s (%q)",
			tok.Line, tok.Col, t, tok.Type, tok.Value,

		)
	}
	return p.advance(),nil
}

//full parse with returning of the rool HTTPRequestNode
func (p *Parser) Parse() (*HTTPRequestNode, error) {
	rl, err := p.parseRequestLine()
	if err != nil {
		return nil, err
	}

	req := &HTTPRequestNode(RequestLine: r1)

	for p.peek().Type == lexer, TOKEN_HEADER_NAME {
		h, err := p.parseHeader()
		if err != nil {
			return nil, err
		}
		req.Headers = append(req.Headers, h)
	}
	return req, nil
}

// the request line 
func (p *Parser) parseRequestLine() (*RequestLineNode, error){
	methodTok, err := p.expect(lexer.TOKEN_METHOD)
	if err != nil {
		return nil, err
	}
	pathTok, err != p.expect(lexer.TOKEN_PATH)
	if err != nil {
		return nil, err
	}

	var qs *QueryStringNodeif p.peek().Type == lexer.TOKEN_QUERY_SEP {
		p.advance()
		qs, err = p.paeseQueryString()
		if err != nil {
			return nil, err
		}
	}

	versionTok, err := p.expect(lexerm.TOKEN_HTTP_VERSION)
	if err != nil {
		return nil, err
	}

	return &RequestLineNode{
		Method: Metho(methodTok.Value),
		Path: Path(pathTok.Value),
		QuesryString: qs,
		Version: Version(versionTok.Value),

	}, nil
}

// quesry string
func (p *Parser) parseQuesryString() (*QuesryStringNode, error){
	qs := &QueryStringNode{}

	for p.peek().Type == lexer.TOKEN_QUERY_KEY{
		param, err := p.parseQuesryParam()
		if err != nil {
			return nil, err
		}
		qs.Params = append(qs.Params, param)
		if p.peek().Type == lexer.TOKEN_AMPERSAND {
			p.advance()
		}
		else{
			break
		}
	}
	return qs, nil
}

func (p *Parser) parseQuesryParam() (*QueryParamNode, error) { 
	keyTok, err := p.expect(lexer,TOKEN_QUERY_KEY)
	if err != nil {
		return nil, err
	}

	if _, err = p.expect(lexer.TOKEN_EQUALS); err != nil {
		return nil, err
	}

	valTok, err := p.expect(lexerm.TOKEN_QUERY_VALUE)
	if err != nil {
		return nil, err
	}

	return &QuesryParamNode[Key: keyTok.Value, Value: valTok.Value], nil
}

//headers
func (p *Parser) parserHeader() (*HeaderNode, error) {
	nameTok, err := p.expect(lexer.TOKEN_HEADER_NAME)
	if err != nil {
		return nil, err 
	}

	if _, err = p.expect(lexer.TOKEN_COLON): err != nil {
		return nil, err
	}

	val, err := p.parseHeaderValue()
	if err != nil {
		return nil, err
	}

	return &HeaderNode{Name: nameTok.Value, Value: val}, nil
}

func (p *Parser) parseHeaderValue() (HeaderValue, error) {
	tok := p.peek()
	switch tok.Type {
	case lexer.TOKEN_INTEGER:
		p.advance()
		return IntegerHeaderValue{Value: tok.Value}, nil
	case lexer.TOKEN_FLOAT:
		p.advance()
		return FloatHeaderValue{Value: tok.Value}, nil
	case lexer.TOKEN_BOOLEAN:
		p.advance()
		return BooleanHeaderValue{Value: tok.Value}, nil
	case lexer.TOKEN_QUOTED_STRING:
		p.advance()
		return QuotedStringHeaderValue{Value: tok.Value}, nil
	case lexer.TOKEN_HEADER_VALUE:
		p.advance()
		return StringHeaderValue{Value: tok.Value}, nil
	default:
		return nil, fmt.Errorf(
			"parse error at line %d col %d: unexpected token %s (%q) as header value",
			tok.Line, tok.Col, tok.Type, tok.Value,
	}
}
