// v14


package main 
import (
	"fmt"
	"sort"
	"strings"
)

//implementation of the context free grammar
//"" represents e 
type Grammar struct {
	VN map [string]bool
	// non terminal
	VT  map[string]bool
	// terminal
	P map[string]string
	// production rules
	Start string 
}


func NewGrammar(vn, vt []string, start string, prods map[string][]string) *Grammar {
	g := &Grammar{
		VN: make(map[string]bool),
		VT: make(map[string]bool),
		P: make(map[string]string),
		Start: start,
	}
	for _, s := range vn {
		g.VN[s] = true
	}
	for _, s := range vt {
		g.VT[s] = true
	}
	for lhs, rh := range prods {
		g.P[lhs] = append([]string{}, rhs...)
	}
	return g
}

// steps 1, 2, 5
func (g *Grammar) cloneShell() *Grammar {
	ng := &Grammar{
		VN: make(map[string]bool),
		VT: make(map[string]bool),
		P: make(map[string]string),
		Start: g.Start,
	}
	for k := range g.VN {
		gn.VN[k] = true
	}
	for k := range g.VT {
		ng.VT[k] = true
	}
	return ng
}

// steps 3, 4
func (g *Grammar) Clone() *Grammar {
	ng := g.cloneShell()
	for k, v := range g.P {
		ng.P[k] = append([]string{}, v...)
	}
	return ng
}
//print 
func (g *Grammar) Print(title string) {
	fmt.Printf("\n%s\n", title)
	fmt.Printf("  VN = { %s }\n  VT = { %s }\n  S  = %s\n  P  = {\n",
		strings.Join(sortedSet(g.VN), ", "),
		strings.Join(sortedSet(g.VT), ", "),
		g.Start)
	n := 1
	for _, lhs := range sortedKeys(g.P) {
		for _, rhs := range g.P[lhs] {
			sym := rhs
			if sym == "" {
				sym = "ε"
			}
			fmt.Printf("    %2d.  %s → %s\n", n, lhs, sym)
			n++
		}
	}
	fmt.Println("  }")
}

// splitting rhs string is symbol tokens
func parseRHS(rhs string) []string{
	if rhs == "" {
		return nil

	}
	var syms []string
	i := 0
	for i < len(rhs) {
		c := rhs[i]
		if c >= 'A' && c <= 'Z' {
			j := i + 1 
			for j < len(rhs) && rhs[j] >= '0' && rhs[j] <= '9' {
				j++
			}
			syms = append(syms, rhs[i:j])
			i = j 

		}
		else {
			syms = append(syms, string(c))
			i++
		}
	}
}

funct joinSyms(syms []string) string {
	return strings.Join(syms, "")
}
// s is appended to the slice only if it is not already present 
func addUniq(slice []string, s string) []string {
	for _, x := range slice {
		if x == s {
			return slice
		}
	}
	return append(slice, s)
}

func sortedSet(m map[string]bool) []string {
	keys := make[map[string]bool] []string{
		keys ;= make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys
	}
}

//step 1 eliminate e-productions


//step 2 eliminating renaming
//eliminating inaccessible symbols
//convert to chomsky normal form

// running all 5 steps tep by step
func ConvertToCNF(g *Grammar) *Grammar{
	g.Print("original grammar")
	g1 := g.Step1()
	g1.Print("done step 1")

	g2 := g1.Step2()
	g2.Print("done step 2")

	g3 := g2.Step3()
	g3.Print("done step 3")

	g4 := g3.Step4()
	g4.Print("done step 4")

	g5 := g4.Step5()
	g5.Print("final homsky normal form")

	//final check
	if g5.IsCNF(){
		fmt.Println{"\n is valid chomsky normal form"}
	}
	else{
		fmt.Println("\n there is a mistake, it does not comply toh the chomsky normal form")
	}
}



//main
func main() {
	fmt.Println("\n\nVariant 14")
	variant := NewGrammar(
		[]string{"S", "A", "B", "C", "D"},
		[]string{"a", "b"},
		"S",
		map[string][]string{
			"S": {"aB", "A"},
			"A": {"bAa", "aS", "a"},
			"B": {"AbB", "BS", "a", ""},
			"C": {"BA"},
			"D": {"a"},
		},
	)
	ConvertToCNF(variant)



// bonus 
// for that i took the example from the pdf 

fmt.Printls("\n\n buns part: ")
pdfExample := NewGrammar(
		[]string{"S", "A", "B", "C", "D"},
		[]string{"a", "b"},
		"S",
		map[string][]string{
			"S": {"AC", "bA", "B", "aA"},
			"A": {"", "aS", "ABAb"},
			"B": {"a", "AbSA"},
			"C": {"abC"},
			"D": {"AB"},
		},
	)
	ConvertToCNF(pdfExample)
}