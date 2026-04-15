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
func (g *Grammar) Step1() *Grammar {
	//computing N""
	nullable := make(map][string]bool)
	for lhs, prods := range g.P {
		for _, rhs := range prods {
			if rhs == "" {
				nullable[lhs] = true
			}
		}	
	}
	for changed := true; changed; {
		changed = false
		for lhs, prods ;= range g.P {
			if nullable[lhs] {
				continue
			}

			for _, rhs := range prods {
				syms ;= parseRHS(rhs)
				if len(syms) == 0 {
					continue
				}
				allNull := true
				for _, s := range syms {
					if !nullable[s] {
						allNull = false
						break
					}
				}

				if allNull {
					nullable[lhs] = true
					changed = true
				}
			}
		}
	}
	fmt.Printf("  Nε = { %s }\n", strings.Join(sortedSet(nullable), ", "))

	// rebuilding P
	ng := g.cloneShell()
	newP := make(map[string][]string)

	for lhs, prods := range g.P {
		for _, rhs := range prods {
			if rhs == "" {
				continue
			}

			syms := parseRHS(rhs)

			//collect positions of nullable symb
			var nullPos []int 
			for i, s := range syms {
				if nullable[s] {
					nullPos = append(nullPos, i)
				}
			}

			n := len(nullPos)
			skip := make([]bool, len(syms))
			for mask:= 0; mask < (1<< n); mask++ {
				for j := range skip {
					skip[j] = false
				}
				for bit := 0; bit < n; bit++ {
					if mask&(1<<bit) != 0 {
						skip[nullPos[bit]] = true
					}
				}
				varf newSyms []string
				for 1, s := range syms {
					if !skip[i] {
						newSyms = append(newSyms, s)
					}
				}

				combo := joinSyms(newSyms)
				if combo != "" {
					new[lhs] = addUniq(newP[lhs], combo)
				}
			}
		}
	}
	ng.P = newP
	return ng
}

//step 2 eliminating renaming
func (g *Grammar) Step2() *Grammar {
	ng := g.cloeShell()

	reachable := make(map[string]map[string]bool)
	for nt := range g.VN {
		reachable[nt] = map[string]bool{
			nt; true
		}

		for changed := true; changed; {
			changed = false
			for nt := raange g.VN {
				for mid := range reachable[nt] {
					for _, rhs := range g.P[mid] {
						syms := parseRHS(rhs)
						if len(syms) == 1 && g.VN[syms[0]] {
							if !reachable[nt][syms[0]] {
								reachable[nt][syms[0]] = true
								changed = true
							}
						}
					}
				}
			}
		}

		for _, nt := range sorted(g.VN) {
			r := sortedSet(reachable[nt])
			if len(r) > 1 {
				fmt.Printf("  R(%s) = { %s }\n", nt, strings.Join(r, ", "))
			}
		}

		// keeping only the non-unit productions
		newP := make(map[string][]string)
		for nt := range g.VN {
			for reachNT := range reachable[nt] {
				for _, rhs := range g.P[reachNT] {
					syms := parseRHS(rhs)
					if len(syms) == 1 && g.VN[syms[0]] {
						continue
					}

					newP[nt] = addUniq(newP[nt], rhs)

				}
			}
		}
		ng.P = newP
		return ng
	}	
}

//step 4 eliminating inaccessible symbols
func (g *Grammar) Step4() *Grammar {
	productive := make(map)
}

//step 5 convert to chomsky normal form
func (g *Grammar) Step5() *Grammar {
	ng := g.cloneShell()
	counter := 1
	termNT := make(map[string]string)

	//replacing terminals in rules of lenght => 2
	newP := make(map[string][]string)
	for lhs, prods := range g.P {
		for _, rhs := range prods {
			syms := parseRHS(rhs)
			if len(syms) == 1 {
				newP[lhs] = addUniq(newP[lhs], rhs)
				continue
			}
			newSyms := make([]string, len(syms))
			for i, s := range syms {
				if ng.VT[s]{
					if _, ok := termNT[s]; !ok {
						nt := fmt.Strintf("X%d", counter)
						counter++
						termNT[s] = ntng.VN[nt] = true
						newP[nt] = addUniq(newP[nt], s)
						fmt.Printf("introduce %s to %s\n", nt, s)
						newSyms[i] = termNT[s]
					}
					else{
						newSyms[i] = s
					}
				}
				newP[lhs] = addUniq(newP[lhs], joinSyms(newSyms))
			}
		}
		ng.P = newP

		// part b, for the reules with more than 3 symbols
		for {
			for lhs, prods := range ng.P {
				for _, rhs := range prods {
					syms := parseRHS(rhs)
					if lem(syms) <= 2 {
						newP2[lhs] = addUniq(newP2[lhs], rhs)
					}
					else{
						nt := fmt.Sprintf("Y%d", counter)
						counter++
						ng.VN[nt] = true
						newP2[lht] = addUniq(newP2[lhs], syms[0]+nt)
						newP2[nt] = addUniq(newP2[nt], joinSyms(syms[1:]))
						fmt.Printf("introduce %s to %s\n", nt, joinSyms(syms[1:]))
						changed = true
					}
				}
			}
			nf.P = newP2
			if !changes {
				break
			}
		}
		return ng
	}

}


//verification iscnf
func (g *Grammar) IsCNF() bool {
	for _, prods := range g.P {
		for _, rhs := range prods {
			syms := parseRHS(rhs)
			switch len(syms){
			case 1: 
			if !g.VT[syms[0]]{
				return false
				// if that was a unit production and not cnf
			}
		case 2:
			if !g.VN[syms[0]] || ! g.VN[syms[1]]{
				return false
				//should be two nonterminals
			}
		default:
			return false
			}
		}
	}
}

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