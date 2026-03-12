"""
Q - set of states
sigma - alphabet
delta - dict
q0 - start state
F - set of final states
"""

from grammar import Grammar

class FiniteAutomaton:
    def __init__(self, Q, sigma, delta, q0, F):
        self.Q = set(Q)
        self.sigma = set(sigma)
        self.delta = delta
        self.q0 = q0
        self.F = set(F)

# 3b deterministic

    def is_deterministic(self):
        for (state, sym), targets in self.delta.items():
            if len(targets) > 1:
                return False
        return True

# 3a from finite automata to regular grammar

    def to_regular_grammar(self):
        VN = {str(s).upper() for s in self.Q}
        VT = set(self.sigma)
        S = str(self.q0).upper()
        P = {nt: [] for nt in VN}

        for (state, sym), targets in self.delta.items():
            lhs = str(state).upper()
            for t in targets:
                rhs_nt = str(t).upper()
                P[lhs].append(f"{sym}{rhs_nt}")

        P = {k: list(dict.fromkeys(v)) for k, v in P.items() if v}
        return Grammar(VN, VT, P, S)

# 3c from ndfa to dfa

    def to_dfa(self):
        start_fs = frozenset([self.q0])
        unvisited = [start_fs]
        visited = set()
        dfa_delta = {}
        dfa_F = set()

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

        # printing stuff
        def label(fs):
            if not fs:  return "∅"
            return "{" + ",".join(sorted(str(s) for s in fs)) + "}"

        return FiniteAutomaton(
            Q = [label(s) for s in visited],
            sigma = self.sigma,
            delta = {(label(s), sym): [label(t)]
                        for (s, sym), [t] in dfa_delta.items()},
            q0 = label(start_fs),
            F = {label(s) for s in dfa_F}
            )

    def __str__(self):
        lines = [f"  States: {sorted(self.Q)}",
                 f"  Alphabet: {sorted(self.sigma)}",
                 f"  Start: {self.q0}",
                 f"  Final: {sorted(self.F)}",
                 "  Transitions:"]
        for (s, sym), targets in sorted(self.delta.items()):
            lines.append(f"    δ({s}, {sym}) = {targets}")
        return "\n".join(lines)                    
