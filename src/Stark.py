from Field import Field, FieldElement
from Merkle import MerkleTree
from Polynomial import Polynomial

class Stark():
    def __init__(self, field:Field, expansion_factor:int):
        self.field = field
        self.expansion_factor = expansion_factor

    def prove(self, trace, constraints):
        #First we need to create the generators for the trace and the LDE polynomial
        G, H = self.create_generators(len(trace))
        # Now create the trace polynomial
        trace_poly = self.create_trace_polynomial(trace, G)
        # Test the trace polynomial
        to_commit_TP = trace_poly.evaluate_domain(H)
        # constraints are a list of tuples (const_pos, const_value)
        const_poly_list = self.get_constraints_polynomial(constraints, trace_poly, G)
        # combine the const_poly in a linear combination
        const_poly = self.combine_polynomials(const_poly_list, [FieldElement(1, self.field), FieldElement(1, self.field)])
        # Test the constraints polynomial to commit into the merkle tree
        to_commit_CP = const_poly.evaluate_domain(H)
        fri_layers, fri_domain, fri_polys, merkles = self.fri(const_poly, H, to_commit_CP)
        return to_commit_TP, to_commit_CP, fri_layers, fri_domain, fri_polys, merkles

    def next_fri_domain(domain):
        return [x^2 for x in domain[:len(domain) // 2]]


    def next_fri_polynomial(poly:Polynomial, alpha):
        odd_coefficients = poly.coefficients[1::2]
        even_coefficients = poly.coefficients[::2]
        odd = Polynomial(odd_coefficients).scale(FieldElement(alpha, Field.main()))
        even = Polynomial(even_coefficients)
        return odd + even


    def next_fri_layer(poly, dom, alpha):
        next_poly = Stark.next_fri_polynomial(poly, alpha)
        next_dom = Stark.next_fri_domain(dom)
        next_layer = [next_poly.evaluate(x) for x in next_dom]
        return next_poly, next_dom, next_layer

    def fri(self, poly, domain, to_commit):
        fri_polys = [poly]
        fri_doms = [domain]
        fri_layers = [to_commit]
        merkles = [MerkleTree(to_commit)]
        while fri_polys[-1].degree() > 0:
            next_poly, next_dom, next_layer = Stark.next_fri_layer(fri_polys[-1], fri_doms[-1], 1)
            fri_polys.append(next_poly)
            fri_doms.append(next_dom)
            fri_layers.append(next_layer)
            merkles.append(MerkleTree(next_layer))
        return fri_layers, fri_doms, fri_polys, merkles

    def combine_polynomials(self, polynomials, scales):
        acc = Polynomial([])
        i = 0
        for poly in polynomials:
            acc = acc + poly.scale(scales[i])
            i += 1
        return acc


    def get_constraints_polynomial(self, constraints, trace_poly, G):
        const_poly = []
        for (index, value) in constraints:
            nom = trace_poly - Polynomial([value])
            denom = Polynomial.term(FieldElement(G[index], self.field))
            const_poly.append(nom/denom)

        return const_poly



    def create_generators(self, num_generators: int):
        next_2_power = self.next_2_power(num_generators)
        g = self.field.primitive_nth_root(next_2_power)
        h = self.field.primitive_nth_root(next_2_power*self.expansion_factor)
        G = [ g^i for i in range(num_generators)]
        H = [ h^i for i in range(num_generators*self.expansion_factor)]
        return G, H

    def create_trace_polynomial(self, trace, g):
        return Polynomial.interpolate_domain(g, trace)
    
    def next_2_power(self, num:int):
        if num % 2 == 0:
            return 1 << num.bit_length()
        else:
            return 1 << (num.bit_length() + 1)


