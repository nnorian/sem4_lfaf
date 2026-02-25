# variant 14

from finite_automaton import FiniteAutomaton
from grammar import Grammar

if __name__ == "__main__":

    ndfa = FiniteAutomaton(
        Q = ['q0', 'q1', 'q2'],
        sigma = ['a', 'b', 'c'],
        delta = {
            ('q0', 'a'): ['q0'],
            ('q0', 'b'): ['q1'],
            ('q1', 'c'): ['q1', 'q2'],
            ('q2', 'a'): ['q0'],
            ('q1', 'a'): ['q1'],
        },
        q0 = 'q0',
        F = ['q2']
    )

    print(ndfa)
    det = ndfa.is_deterministic()

    gram = ndfa.to_regular_grammar()
    print(gram)

    dfa = ndfa.to_dfa()
    print(dfa)
    print(f"{dfa.is_deterministic()}")
