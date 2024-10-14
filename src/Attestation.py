from Field import *
import hashlib
from Polynomial import *
from Merkle import *

class Attestation():
    def __init__( self, initial_addr, final_addr, nonce):
        self.field = Field.main()
        self.initial_addr = FieldElement(initial_addr, self.field)
        self.final_addr = FieldElement(final_addr, self.field)
        self.nonce = FieldElement(nonce, self.field)

    def attest(self):
        #Emulate the attestation process
        #The attestation process is a simple hash concatenation
        #of the for all the addresses between the initial and final address
        hash = hashlib.sha256(str(self.nonce.value).encode())
        for i in range(self.initial_addr.value, self.final_addr.value):
            hash  = hashlib.sha256(hashlib.sha256(str(FieldElement(i, self.field).value).encode()).digest() + hash.digest())
        return hash.digest()

    def trace(self):
        #Emulate the trace process
        #The trace process is a simple hash concatenation
        trace = []
        hash = hashlib.sha256(str(self.nonce.value).encode())
        trace.append(hash.digest())
        for i in range(self.initial_addr.value, self.final_addr.value):
            hash  = hashlib.sha256(hashlib.sha256(str(FieldElement(i, self.field).value).encode()).digest() + hash.digest())
            trace.append(hash.digest())

        num_trace = []
        for t in trace:
            num_trace.append(FieldElement(int.from_bytes(t, byteorder='big'), self.field))
        return num_trace

    def polynomial_constrains(self, trace):
        #Emulate the polynomial constraints
        #The polynomial constraints are a simple hash concatenation
        constraints = []
        hash0 = FieldElement(int.from_bytes(hashlib.sha256(str(self.nonce.value).encode()).digest(), byteorder='big'), self.field)
        hash1 = trace[-1]
        constraints.append((0, hash0)
        #p(x0) = f(x)-hash0/(x-g0)
        constraints.append((len(trace)-1, hash1))
                           # p(x1) = f(x)-hash1/(x-gFinal)

        return constraints
    


