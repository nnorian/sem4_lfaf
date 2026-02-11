from grammar import Grammar

def main():
    vn = {'S', 'A', 'B'}
    vt = {'a', 'b', 'c', 'd'}
    productions = {
        'S': ['aS', 'bB'],
        'B': ['cA', 'd'],
        'A': ['bB', 'aS'],
    }
    start_symbol = 'S'

    grammar = Grammar(vn, vt, productions, start_symbol)
    fa = grammar.to_finite_automation()

    print("testing strings ^^")
    t_strings = [
        "bd",
        "abd",
        "abcd",
        "bab",
        "baab", # edgecase
        "abc", # this and rest invalid
        "bbb",
        ""
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
