package parser

import "fmt"



type Method  string
type Path    string
type Version string

type HTTPRequestNode struct {
	RequestLine *RequestLineNode
	Headers     []*HeaderNode
}

type RequestLineNode struct {
	Method      Method
	Path        Path
	QueryString *QueryStringNode 
	Version     Version
}


type QueryStringNode struct {
	Params []*QueryParamNode
}


type QueryParamNode struct {
	Key   string
	Value string
}


type HeaderNode struct {
	Name  string
	Value HeaderValue
}

// header value variants 

type HeaderValue interface {
	headerValue()
}


type StringHeaderValue struct{ Value string }

func (StringHeaderValue) headerValue() {}

type IntegerHeaderValue struct{ Value string }

func (IntegerHeaderValue) headerValue() {}


type FloatHeaderValue struct{ Value string }

func (FloatHeaderValue) headerValue() {}


type BooleanHeaderValue struct{ Value string }

func (BooleanHeaderValue) headerValue() {}


type QuotedStringHeaderValue struct{ Value string }

func (QuotedStringHeaderValue) headerValue() {}


// print returns a depth-indented tree view of the AST
func Print(root *HTTPRequestNode) string {
	out := "HTTPRequest\n"
	out += printRequestLine(root.RequestLine, 1)
	for _, h := range root.Headers {
		out += printHeader(h, 1)
	}
	return out
}

func pad(depth int) string {
	s := ""
	for i := 0; i < depth; i++ {
		s += "  "
	}
	return s
}

func printRequestLine(n *RequestLineNode, depth int) string {
	out := pad(depth) + "RequestLine\n"
	out += fmt.Sprintf("%sMethod(%q)\n", pad(depth+1), n.Method)
	out += fmt.Sprintf("%sPath(%q)\n", pad(depth+1), n.Path)
	if n.QueryString != nil {
		out += pad(depth+1) + "QueryString\n"
		for _, p := range n.QueryString.Params {
			out += fmt.Sprintf("%sQueryParam(key=%q, value=%q)\n", pad(depth+2), p.Key, p.Value)
		}
	}
	out += fmt.Sprintf("%sVersion(%q)\n", pad(depth+1), n.Version)
	return out
}

func printHeader(n *HeaderNode, depth int) string {
	out := fmt.Sprintf("%sHeader(%q)\n", pad(depth), n.Name)
	switch v := n.Value.(type) {
	case StringHeaderValue:
		out += fmt.Sprintf("%sValue[string](%q)\n", pad(depth+1), v.Value)
	case IntegerHeaderValue:
		out += fmt.Sprintf("%sValue[integer](%q)\n", pad(depth+1), v.Value)
	case FloatHeaderValue:
		out += fmt.Sprintf("%sValue[float](%q)\n", pad(depth+1), v.Value)
	case BooleanHeaderValue:
		out += fmt.Sprintf("%sValue[boolean](%q)\n", pad(depth+1), v.Value)
	case QuotedStringHeaderValue:
		out += fmt.Sprintf("%sValue[quoted_string](%q)\n", pad(depth+1), v.Value)
	}
	return out
}