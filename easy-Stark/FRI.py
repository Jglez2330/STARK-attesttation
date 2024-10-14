import ProofStream
from Field import *
from Merkle import *
from ntt import *
from hashlib import blake2b

class Fri:
    def __init__( self, domain):
        self.domain = domain

    def next_fri_domain_test(self, fri_domain):
        return [x^2 for x in fri_domain[:len(fri_domain) // 2]]
    def next_fri_domain(self, domain):
        return [x^2 for x in domain[:len(domain) // 2]]

    def next_fri_polynomial(self, poly:Polynomial, beta:FieldElement) -> Polynomial:
        odd_coef = poly.coefficients[1::2]
        even_coef = poly.coefficients[::2]
        odd = beta * Polynomial(odd_coef)
        even = Polynomial(even_coef)
        return odd + even
    def next_fri_layer(self, cp, domain, beta):
        next_poly = self.next_fri_polynomial(cp, beta)
        next_dom = self.next_fri_domain(domain)
        next_layer = next_poly.evaluate_domain(next_dom)
        return next_poly, next_dom, next_layer
    def commit(self, cp:Polynomial, domain, cp_eval, cp_merkle:MerkleTree, proof_stream:ProofStream):
        fri_polys = [cp]
        fri_domains = [domain]
        fri_layers = [cp_eval]
        fri_merkles = [cp_merkle]

        while fri_polys[-1].degree() > 0:
            beta = Polynomial([FieldElement(1, Field.main())]) # random value
            next_poly, next_domain, next_layer = self.next_fri_layer(fri_polys[-1], fri_domains[-1], beta)
            print(f'next_len fri layer: {len(next_layer)}')
            fri_polys.append(next_poly)
            fri_domains.append(next_domain)
            fri_layers.append(next_layer)
            fri_merkles.append(MerkleTree(next_layer))
        return fri_polys, fri_domains, fri_layers, fri_merkles
    def decommit_on_fri(self, idx, fri_polys:list[Polynomial], fri_domains:list[list[FieldElement]], fri_layers, fri_merkles:list[MerkleTree]):
        assert len(fri_layers) == len(fri_merkles), f'layers size should be same as merkles size'
        res = []
        i = 0
        for layer, merkle in zip(fri_layers[:-1], fri_merkles[:-1]):
            length = len(layer)
            idx = idx % length
            sib_idx = (idx + length // 2) % length
            assert len(layer) == len(fri_domains[i])
            assert layer[idx] == fri_polys[i].evaluate(fri_domains[i][idx])
            assert layer[sib_idx] == fri_polys[i].evaluate(-fri_domains[i][idx])
            res.append(layer[idx])
            res.append(merkle.get_authentication_path(idx))
            res.append(layer[sib_idx])
            res.append(merkle.get_authentication_path(sib_idx))
            i = i+1
        res.append(fri_layers[-1][0])
        return res
    def decommit_on_query(self, idx, f_eval, f_merkle:MerkleTree):
        assert idx + 8 < len(f_eval), f'idx should be less than len(f_eval) - 8'
        res = []
        res.append(f_eval[idx])
        res.append(f_merkle.get_authentication_path(idx))
        res.append(f_eval[idx+8])
        res.append(f_merkle.get_authentication_path(idx+8))
        return res

