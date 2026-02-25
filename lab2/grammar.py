class Grammar: 
    def __init__(self, VN, VT, P, S):
        self.VN = set(VN)
        self.VT = set(VT)
        self.P = P
        self.S = S

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
        
        return "it is type 1, the unrestricted grammar type"

    def _is_type3(self):

        for lhs, rhss in self.P.items():
            if lhs not in self.VN or len(lhs) != 1:
                return False
            for rhs in rhss:
                if rhs == "ε":
                    continue 
                if len(rhs) == 1 and rhs in self.VT:
                    continue
                if len(rhs) == 2 and self.VT and rhs[1] in self.VN:
                    continue
                return False
        return True

    def _is_type2(self):

        for lhs in self.P:
            if len(lhs) != 1 or lhs not in self.VN:
                return False
        return True

    def _is_type1(self):

        for lhs, rhss in self.P.items():
            for rhs in rhss:
                if rhs == "ε":
                    if lhs == self.S:
                        continue
                    return False
                if len(lhs) > len(rhs):
                    return False
        return True

    def __str__(self):
        lines = []
        for lhs, rhss in self.P.items():
            lines.append(f"{lhs} -> {' | '.join(rhss)}")
        return "\n". join(lines)