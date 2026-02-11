"""
set of states 
input alpahbet
dict of state and symbol amd the list of next states
initial state 
accept states (f) final kinda
"""






class FiniteAutomation:
    def __init__(self, states, alphabet, transitions, start_state, accept_states):
        self.states = states
        self.alphabet = alphabet
        self.transitions = transitions
        self.start_state = start_state
        self.accept_states = accept_states


    #checks if the input string is accepted by automation
    def string_belongs_to_language(self, input_string):
        current_states = {self.start_state}
        for symbol in input_string:
            next_states = set()
            for state in current_states:
                key = (state, symbol)
                if key in self.transitions:
                    next_states.update(self.transitions[key])
            current_states = next_states

            if not current_states:
                return False # rejects the invalid transitions

        return bool(current_states & self.accept_states)
