# Determinism in Finite Automata. Conversion from NDFA to DFA. Chomsky Hierarchy.

### Course: Formal Languages & Finite Automata
### Author: Kushniorenko Ecaterina FAF-243

----

## Theory

A finite automaton is a mechanism that represents processes.
It has a finite set of states, an alphabet, a transition function, a start state, and a set of final states.

When for a given state and input symbol the transition function points to more than one state, the automaton is called non-deterministic (NDFA). When every (state, symbol) pair leads to at most one state, it is called deterministic (DFA).

Every NDFA has an equivalent DFA that accepts exactly the same language. The conversion is done through the subset construction algorithm, where each DFA state represents a set of NDFA states that can be active simultaneously.

A regular grammar is one where every production has the form A → aB or A → a, with at most one non-terminal on the right side at the end.

Every regular grammar has an equivalent finite automaton, and they can be converted.

The Chomsky hierarchy classifies grammars into four types: Type 3 (Regular), Type 2 (Context-Free), Type 1 (Context-Sensitive), and Type 0 (Unrestricted), each being a subset of the one above it.

## Objectives

Implement conversion of a finite automaton to a regular grammar, determine whether the FA is deterministic or non-deterministic, implement NDFA to DFA conversion, and classify a grammar based on the Chomsky hierarchy.

The steps followed are:

a. Implement a Grammar class with a Chomsky hierarchy classifier

b. Implement conversion of a Finite Automaton to a Regular Grammar

c. Determine whether the FA is deterministic or non-deterministic

d. Implement NDFA to DFA conversion using subset construction

## Implementation description

The `Grammar` class stores non-terminals, terminals, production rules, and a start symbol. The `classify` method determines the Chomsky type by pre-computing each level and checking from most specific to least specific, since Type 3 ⊂ Type 2 ⊂ Type 1 ⊂ Type 0:

```python
def classify(self):
    is_type3 = self._is_type3()
    is_type2 = is_type3 or self._is_type2()
    is_type1 = is_type2 or self._is_type1()

    if is_type3:
        return "it is type 3, the regular grammar type"
    if is_type2:
        return "it is type 2, the context free grammar type"
    if is_type1:
        return "it is type 1, the context sensitive grammar type"
    return "it is type 0, the unrestricted grammar type"
```

Type 3 is checked by verifying every production has the form A → aB or A → a. Type 2 requires only that the left-hand side is a single non-terminal. Type 1 requires the right-hand side to never be shorter than the left-hand side.

The `FiniteAutomaton` class stores transitions as a dictionary mapping `(state, symbol)` to a list of target states. Using a list instead of a single value is what allows the same class to represent both NDFAs and DFAs naturally.

The `is_deterministic` method simply checks whether any transition list has more than one target:

```python
def is_deterministic(self):
    for (state, sym), targets in self.delta.items():
        if len(targets) > 1:
            return False
    return True
```

The `to_regular_grammar` method maps each state to a non-terminal and each transition to a right-linear production. Every (state, symbol, target) triple produces exactly one rule of the form A → aB:

```python
for (state, sym), targets in self.delta.items():
    lhs = str(state).upper()
    for t in targets:
        rhs_nt = str(t).upper()
        P[lhs].append(f"{sym}{rhs_nt}")   # A → aB
```

## NDFA to DFA conversion — subset construction

The `to_dfa` method implements the subset construction algorithm. The core idea is that one DFA state equals a set of NDFA states that can all be active at the same time.

```python
def to_dfa(self):
    start_fs  = frozenset([self.q0])
    unvisited = [start_fs]
    visited   = set()
    dfa_delta = {}
    dfa_F     = set()

    while unvisited:
        current = unvisited.pop()
        if current in visited:
            continue
        visited.add(current)

        if current & self.F:
            dfa_F.add(current)

        for sym in self.sigma:
            next_states = frozenset(
                t
                for s in current
                for t in self.delta.get((s, sym), [])
            )
            dfa_delta[(current, sym)] = [next_states]
            if next_states and next_states not in visited:
                unvisited.append(next_states)
```

A `frozenset` is used because sets are not hashable in Python and cannot be used as dictionary keys. Each new subset of NDFA states encountered becomes a new DFA state, added to the worklist if not yet processed.

## Variant 14 — the non-determinism explained

The variant defines:

```
δ(q1, c) = q1
δ(q1, c) = q2
```

Two rules for the same (state, symbol) pair make this an NDFA. In the code both targets are merged into one list:

```python
('q1', 'c'): ['q1', 'q2']
```

The problematic `δ(q1,c) = {q1, q2}` is resolved by making `{q1,q2}` a single DFA state. Now each DFA state has exactly one target per symbol, making the result deterministic. The final state in the DFA is `{q1,q2}` because it contains `q2`, which was the original final state.

The full subset construction trace for Variant 14:

```
DFA state    | on a       | on b    | on c
-------------|------------|---------|----------
{q0}         | {q0}       | {q1}    | ∅
{q1}         | {q1}       | ∅       | {q1,q2}  ← final
{q1,q2}      | {q0,q1}    | ∅       | {q1,q2}
{q0,q1}      | {q0,q1}    | {q1}    | {q1,q2}
```

## Results

```
=======================================================
ORIGINAL AUTOMATON (Variant 14)
=======================================================
  States: ['q0', 'q1', 'q2']
  Alphabet: ['a', 'b', 'c']
  Start: q0
  Final: ['q2']
  Transitions:
    δ(q0, a) = ['q0']
    δ(q0, b) = ['q1']
    δ(q1, a) = ['q1']
    δ(q1, c) = ['q1', 'q2']   ← two targets → NDFA
    δ(q2, a) = ['q0']

─── Is Deterministic? ──────────────────────────────────
  → False  (δ(q1,c) goes to both q1 and q2)

─── Regular Grammar (from NDFA) ─────────────────────────
  Q0 -> aQ0 | bQ1
  Q1 -> cQ1 | cQ2 | aQ1
  Q2 -> aQ0

  Chomsky classification: it is type 3, the regular grammar type

─── Converting NDFA → DFA (subset construction) ─────────
  States: ['{q0}', '{q1}', '{q1,q2}', '{q0,q1}']
  Start: {q0}
  Final: ['{q1,q2}']
  Transitions:
    δ({q0},    a) = ['{q0}']
    δ({q0},    b) = ['{q1}']
    δ({q1},    c) = ['{q1,q2}']
    δ({q1,q2}, a) = ['{q0,q1}']
    δ({q1,q2}, c) = ['{q1,q2}']
    δ({q0,q1}, a) = ['{q0,q1}']
    δ({q0,q1}, b) = ['{q1}']
    δ({q0,q1}, c) = ['{q1,q2}']

─── Is DFA Deterministic? ──────────────────────────────
  → True
```

## Conclusions

The implementation correctly identifies Variant 14 as an NDFA due to `δ(q1,c)` having two targets. The subset construction produces a valid DFA with four states where every (state, symbol) pair has exactly one successor. The produced regular grammar is correctly classified as Type 3 in the Chomsky hierarchy, which is consistent with the fact that every finite automaton corresponds to a regular grammar. The Chomsky classifier works generically for all four types by pre-computing each level and checking conditions from most to least restrictive.

## References

Hopcroft, Motwani, Ullman — Introduction to Automata Theory, Languages, and Computation
