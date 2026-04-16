// v14

package main

import (
	"fmt"
	"sort"
	"strings"
)

// implementation of the context free grammar
// "" represents ε
type Grammar struct {
	VN    map[string]bool
	VT    map[string]bool
	P     map[string][]string
	Start string
}

func NewGrammar(vn, vt []string, start string, prods map[string][]string) *Grammar {
	g := &Grammar{
		VN:    make(map[string]bool),
		VT:    make(map[string]bool),
		P:     make(map[string][]string),
		Start: start,
	}
	for _, s := range vn {
		g.VN[s] = true
	}
	for _, s := range vt {
		g.VT[s] = true
	}
	for lhs, rhs := range prods {
		g.P[lhs] = append([]string{}, rhs...)
	}
	return g
}

func (g *Grammar) cloneShell() *Grammar {
	ng := &Grammar{
		VN:    make(map[string]bool),
		VT:    make(map[string]bool),
		P:     make(map[string][]string),
		Start: g.Start,
	}
	for k := range g.VN {
		ng.VN[k] = true
	}
	for k := range g.VT {
		ng.VT[k] = true
	}
	return ng
}

func (g *Grammar) Clone() *Grammar {
	ng := g.cloneShell()
	for k, v := range g.P {
		ng.P[k] = append([]string{}, v...)
	}
	return ng
}

// print
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

// splitting rhs string into symbol tokens
func parseRHS(rhs string) []string {
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
		} else {
			syms = append(syms, string(c))
			i++
		}
	}
	return syms
}

func joinSyms(syms []string) string {
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
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// step 1: eliminate ε-productions
func (g *Grammar) Step1() *Grammar {
	// computing Nε
	nullable := make(map[string]bool)
	for lhs, prods := range g.P {
		for _, rhs := range prods {
			if rhs == "" {
				nullable[lhs] = true
			}
		}
	}
	for changed := true; changed; {
		changed = false
		for lhs, prods := range g.P {
			if nullable[lhs] {
				continue
			}
			for _, rhs := range prods {
				syms := parseRHS(rhs)
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

			// collect positions of nullable symbols
			var nullPos []int
			for i, s := range syms {
				if nullable[s] {
					nullPos = append(nullPos, i)
				}
			}

			n := len(nullPos)
			skip := make([]bool, len(syms))
			for mask := 0; mask < (1 << n); mask++ {
				for j := range skip {
					skip[j] = false
				}
				for bit := 0; bit < n; bit++ {
					if mask&(1<<bit) != 0 {
						skip[nullPos[bit]] = true
					}
				}
				var newSyms []string
				for i, s := range syms {
					if !skip[i] {
						newSyms = append(newSyms, s)
					}
				}
				combo := joinSyms(newSyms)
				if combo != "" {
					newP[lhs] = addUniq(newP[lhs], combo)
				}
			}
		}
	}
	ng.P = newP
	return ng
}

// step 2: eliminate renaming (unit productions)
func (g *Grammar) Step2() *Grammar {
	ng := g.cloneShell()

	reachable := make(map[string]map[string]bool)
	for nt := range g.VN {
		reachable[nt] = map[string]bool{nt: true}
	}

	for changed := true; changed; {
		changed = false
		for nt := range g.VN {
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

	for _, nt := range sortedSet(g.VN) {
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

// step 3: eliminate inaccessible symbols
func (g *Grammar) Step3() *Grammar {
	accessible := map[string]bool{g.Start: true}
	for changed := true; changed; {
		changed = false
		for lhs := range g.P {
			if !accessible[lhs] {
				continue
			}
			for _, rhs := range g.P[lhs] {
				for _, s := range parseRHS(rhs) {
					if !accessible[s] {
						accessible[s] = true
						changed = true
					}
				}
			}
		}
	}

	var inacc []string
	for nt := range g.VN {
		if !accessible[nt] {
			inacc = append(inacc, nt)
		}
	}

	sort.Strings(inacc)
	if len(inacc) == 0 {
		fmt.Println("Inaccessible: {} (none)")
	} else {
		fmt.Printf("Inaccessible: { %s }\n", strings.Join(inacc, ", "))
	}

	ng := g.Clone()
	for _, s := range inacc {
		delete(ng.VN, s)
		delete(ng.VT, s)
		delete(ng.P, s)
	}

	for lhs, prods := range ng.P {
		var kept []string
		for _, rhs := range prods {
			ok := true
			for _, s := range parseRHS(rhs) {
				if !ng.VN[s] && !ng.VT[s] {
					ok = false
					break
				}
			}
			if ok {
				kept = append(kept, rhs)
			}
		}
		ng.P[lhs] = kept
	}
	return ng
}

// step 4: eliminate non-productive symbols
func (g *Grammar) Step4() *Grammar {
	productive := make(map[string]bool)
	for t := range g.VT {
		productive[t] = true
	}
	for changed := true; changed; {
		changed = false
		for lhs, prods := range g.P {
			if productive[lhs] {
				continue
			}
			for _, rhs := range prods {
				allProd := true
				for _, s := range parseRHS(rhs) {
					if !productive[s] {
						allProd = false
						break
					}
				}
				if allProd {
					productive[lhs] = true
					changed = true
					break
				}
			}
		}
	}

	var nonProd []string
	for nt := range g.VN {
		if !productive[nt] {
			nonProd = append(nonProd, nt)
		}
	}
	sort.Strings(nonProd)
	if len(nonProd) == 0 {
		fmt.Println("non-productive: {}")
	} else {
		fmt.Printf("non-productive: { %s }\n", strings.Join(nonProd, ", "))
	}

	ng := g.Clone()
	for _, nt := range nonProd {
		delete(ng.VN, nt)
		delete(ng.P, nt)
	}
	for lhs, prods := range ng.P {
		var kept []string
		for _, rhs := range prods {
			ok := true
			for _, s := range parseRHS(rhs) {
				if !productive[s] {
					ok = false
					break
				}
			}
			if ok {
				kept = append(kept, rhs)
			}
		}
		ng.P[lhs] = kept
	}
	return ng
}

// step 5: convert to Chomsky Normal Form
func (g *Grammar) Step5() *Grammar {
	ng := g.cloneShell()
	counter := 1
	termNT := make(map[string]string)

	// part a: replace terminals in rules of length >= 2
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
				if ng.VT[s] {
					if _, ok := termNT[s]; !ok {
						nt := fmt.Sprintf("X%d", counter)
						counter++
						termNT[s] = nt
						ng.VN[nt] = true
						newP[nt] = addUniq(newP[nt], s)
						fmt.Printf("introduce %s → %s\n", nt, s)
					}
					newSyms[i] = termNT[s]
				} else {
					newSyms[i] = s
				}
			}
			newP[lhs] = addUniq(newP[lhs], joinSyms(newSyms))
		}
	}
	ng.P = newP

	// part b: break rules with more than 2 symbols
	for {
		changed := false
		newP2 := make(map[string][]string)
		for lhs, prods := range ng.P {
			for _, rhs := range prods {
				syms := parseRHS(rhs)
				if len(syms) <= 2 {
					newP2[lhs] = addUniq(newP2[lhs], rhs)
				} else {
					nt := fmt.Sprintf("Y%d", counter)
					counter++
					ng.VN[nt] = true
					newP2[lhs] = addUniq(newP2[lhs], syms[0]+nt)
					newP2[nt] = addUniq(newP2[nt], joinSyms(syms[1:]))
					fmt.Printf("introduce %s → %s\n", nt, joinSyms(syms[1:]))
					changed = true
				}
			}
		}
		ng.P = newP2
		if !changed {
			break
		}
	}
	return ng
}

// verification: is CNF
func (g *Grammar) IsCNF() bool {
	for _, prods := range g.P {
		for _, rhs := range prods {
			syms := parseRHS(rhs)
			switch len(syms) {
			case 1:
				if !g.VT[syms[0]] {
					return false
				}
			case 2:
				if !g.VN[syms[0]] || !g.VN[syms[1]] {
					return false
				}
			default:
				return false
			}
		}
	}
	return true
}

// running all 5 steps step by step
func ConvertToCNF(g *Grammar) {
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
	g5.Print("final Chomsky Normal Form")

	// final check
	if g5.IsCNF() {
		fmt.Println("\n is valid Chomsky Normal Form")
	} else {
		fmt.Println("\n there is a mistake, it does not comply with the Chomsky Normal Form")
	}
}

// main
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
	fmt.Println("\n\n bonus part: ")
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
