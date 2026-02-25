# Intro to formal languages. Regular grammars. Finite Automata.

### Course: Formal Languages & Finite Automata
### Author: Kushnirenko Ecaterina FAF-243

----

## Theory
formal language composed of alphabet, valid words, grammar that shows how the words are made up. 

regular grammar at most one nonterminal on right at the end

each regular grammar has equivalent automata wich can be reciprocaly converted 

## Objectives

Implement a Grammar class that can generate valid strings and convert itself into a Finite Automaton. The Finite Automaton should then be able to verify whether a given string belongs to the language.

in this code i an following the stepts of:

a. Implement a type/class for my grammar

b. Add one function that would generate 5 valid strings from the language expressed by my given grammar 

c. Implement some functionality that would convert and object of type Grammar to one of type Finite Automaton

d. For the Finite Automaton, adding a method that checks if an input string can be obtained via the state transition from it

## Implementation description

The Grammar class stores non-terminals, terminals, production rules, and a start symbol. The generate_string method starts from the start symbol and repeatedly picks a random production rule, appending the terminal character and following the next non-terminal until a terminal-only production is chosen.

```python
def generate_string(self):
    current = self.start_symbol
    result = []
    while current in self.productions:
        options = self.productions[current]
        chosen = random.choice(options)
        terminal = chosen[0]
        result.append(terminal)
        if len(chosen) > 1:
            current = chosen[1]
        else:
            current = None
    return ''.join(result)
```

The to_finite_automation method maps each non-terminal to a state and adds a final accept state X for terminal-only productions. It builds the transition table by iterating over all production rules and grouping next states by (non-terminal, terminal) pairs.

```python
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
            key = (non_terminal, terminal)
            if key not in transitions:
                transitions[key] = []
            transitions[key].append(next_state)
    from finite_automation import FiniteAutomation
    return FiniteAutomation(states, alphabet, transitions, start_state, accept_states)
```

The FiniteAutomation class simulates an NFA by tracking a set of current states. For each input symbol it computes all reachable next states through the transition table. If at any point no states are reachable the string is rejected, otherwise it checks whether any final state was reached.

```python
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
            return False
    return bool(current_states & self.accept_states)
```

## Grammar to finite automaton conversion process

Firstly we do the state mapping. Each non-terminal symbol becomes a state in the automaton. The accept state X does not correspond to any non-terminal. It is introduced to handle terminal-only productions like A → c, where the derivation ends without transitioning to another non-terminal. 

then we do the transition table  where each production rule is converted into a transition. For a rule like S → bA, the transition is δ(S, b) = {A}, meaning from state S on input 'b', move to state A. For terminal-only rules like A → c, the transition becomes δ(A, c) = {X}, directing to the accept state.

## Example case 

Since the automaton is non-deterministic, the string_belongs_to_language method must track all possible states simultaneously. For example, processing the string "bcc":

Start: current states = {S}
Read 'b': δ(S, b) = {A} → current states = {A}
Read 'c': δ(A, c) = {A, B, X} → current states = {A, B, X}
Read 'c': δ(A, c) = {A, B, X}, δ(B, c) = ∅, δ(X, c) = ∅ → current states = {A, B, X}
End: {A, B, X} ∩ {X} = {X} → accepted

## Conclusions

correctly generates  Variant 14 grammar and validates using the finite automaton 

All generated strings are accepted, and manually chosen invalid strings are properly rejected. The rejected ones as you can see from my wariant are the ones that end with anything else tha letter "c".


```
testing strings ^^
'bc' -> yes
'bcc' -> yes
'abc' -> yes
'bcabc' -> yes
'bac' -> yes
'b' -> no
'bbb' -> no
'' -> no
check for generated strings
'bc' -> yes
'bacac' -> yes
'bac' -> yes
'bcabc' -> yes
'abcabcaaaabc' -> yes
```

And after the strings that are generated by the grammar are then checked by the automata rules, and as we cann see each randomly created string is validated and considered correct, having the mark "yes". 

## References

Hopcroft, Motwani, Ullman, Introduction to Automata Theory, Languages, and Computation
