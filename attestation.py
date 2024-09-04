import hashlib
from Field import *
from Polynomial import Polynomial
class Attestation():
    def attestation(self, initial_addr, final_addr, nonce):
        return hashlib.md5((nonce+final_addr).to_bytes(4, byteorder='big')).hexdigest();
    def trace(self, initial_addr, final_addr, nonce):
        trace = []

        state = [Field(int.from_bytes(hashlib.md5((nonce+initial_addr).to_bytes(4, byteorder='big')).digest(), "big"))] + [self.field.zero()]*(initial_addr- final_addr-1)
        #emulate the trace of the program
        for i in range(initial_addr, final_addr):
            state[i] = Field(int.from_bytes(hashlib.md5((nonce+i).to_bytes(4, byteorder='big')).digest(), "big"))

        trace += [[s for s in state]]
        return trace
    def transition_constraints(self, omicron, initial_addr, final_addr, nonce):
        #Generator
        domain = [omicron^r for r in range(0, final_addr-initial_addr)]
        values = [self.field.zero()] * (final_addr-initial_addr)
        variables = Polynomial.variables()
        air = []
        #emulate the transition constraints of the program
        for i in range(initial_addr, final_addr):
            
            air += [(omicron, [i], [i+1], 1)]
        
        
        
    def boundary_constraints(self, initial_addr, final_addr, nonce):
        constraints = []

        constraints += [(0, hashlib.md5((nonce+initial_addr).to_bytes(4, byteorder='big')).hexdigest())]

        constraints += [(final_addr, hashlib.md5((nonce+final_addr).to_bytes(4, byteorder='big')).hexdigest())]

        return constraints

        #emulate the boundary constraints of the program
    def __init__(self):
        self.p = 407 * (1 << 119) + 1
        self.field = Field(self.p)

    def get_attestation_data(self):
        return self.attestation_data

    def __str__(self):
        return "Attestation: " + str(self.attestation_data)
