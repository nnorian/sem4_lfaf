from grammar import Grammar

def main():
    # V 14
    vn = {'S', 'A', 'B'}
    vt = {'a', 'b', 'c'}
    productions = {
        'S': ['aS', 'bA'],    
        'A': ['cA', 'cB', 'aA', 'c'],
        'B': ['aS'],
    }
    start_symbol = 'S'

    grammar = Grammar(vn, vt, productions, start_symbol)
    fa = grammar.to_finite_automation()

    print("testing strings ^^")
    t_strings = [
        "bc",
        "bcc",
        "abc",
        "bcabc",
        "bac",
        "b",
        "bbb",
        "",
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
