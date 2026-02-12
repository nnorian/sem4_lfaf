import random

#vn non terminal
#vt terminal 

class Grammar:
    def __init__(self, vn, vt, productions, start_symbol):
        self.vn = vn
        self.vt = vt
        self.productions = productions
        self.start_symbol = start_symbol

    def generate_string(self):
        current =self.start_symbol
        result = []

        while current in self.productions:

            options=self.productions[current]
            chosen =random.choice(options)
            terminal = chosen[0]
            result.append(terminal)

            if len(chosen) >1:
                current = chosen[1] 
                #next non terminal taken 
            else:
                current = None #terminal only
        return ''.join(result)

    def to_finite_automation(self):
        states = set(self.vn) | {'X'}
        alphabet = set(self.vt)
        start_state = self.start_symbol
        accept_states = {'X'}
        transitions = {}

        for non_terminal, rules in self.productions.items():
            for rule in rules:
                terminal = rule[0]
                next_state = rule[1] if len(rule) > 1 else 'X'

                key=(non_terminal, terminal)
                if key not in transitions:
                    transitions[key] = []
                transitions[key].append(next_state)

        from finite_automation import FiniteAutomation
        return FiniteAutomation(states, alphabet, transitions, start_state, accept_states)
