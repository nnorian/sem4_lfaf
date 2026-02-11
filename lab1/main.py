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
    status = "yes" if result else "no"
    print(f"'{s} -> {status}'")

#test for the generated strings 
print("check for generated strings")
for i in range(5):
    word = grammar.generate_string()
    res = fa.string_belongs_to_language(word)
    status = "yes" if result else "no"
    print(f"'{s} -> {status}'")