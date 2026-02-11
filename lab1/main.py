from grammar import Grammar

def main():
    # Variant 14: q0->S, q1->A, q2->B
    vn = {'S', 'A', 'B'}
    vt = {'a', 'b', 'c'}
    productions = {
        'S': ['aS', 'bA'],       # δ(q0,a)=q0, δ(q0,b)=q1
        'A': ['cA', 'cB', 'aA', 'c'],  # δ(q1,c)=q1, δ(q1,c)=q2, δ(q1,a)=q1, + terminal 'c' to accept at q2
        'B': ['aS'],             # δ(q2,a)=q0
    }
    start_symbol = 'S'

    grammar = Grammar(vn, vt, productions, start_symbol)
    fa = grammar.to_finite_automation()

    print("testing strings ^^")
    t_strings = [
        "bc",        # valid: S->bA->bc
        "bcc",       # valid: S->bA->bcA->bcc
        "abc",       # valid: S->aS->abA->abc
        "bcabc",     # valid: S->bA->bcB->bcaS->bcabA->bcabc
        "bac",       # valid: S->bA->baA->bac
        "b",         # invalid: ends at q1
        "bbb",       # invalid: no δ(q1,b)
        "",          # invalid: empty string
    ]

    for s in t_strings:
        res = fa.string_belongs_to_language(s)
        status = "yes" if res else "no"
        print(f"'{s}' -> {status}")

    #test for the generated strings
    print("check for generated strings")
    for i in range(5):
        word = grammar.generate_string()
        res = fa.string_belongs_to_language(word)
        status = "yes" if res else "no"
        print(f"'{word}' -> {status}")

if __name__=='__main__':
    main()
