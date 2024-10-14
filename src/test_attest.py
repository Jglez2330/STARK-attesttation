from Attestation import *
from Stark import Stark

def test_attes():

    a = Attestation(1, 3, 1)
    print(a.attest())
    trace = a.trace()
    constraints = a.polynomial_constrains(trace)

    s = Stark(Field.main(), 8)
    print(s.prove(trace, constraints))
    print(constraints)
    print(trace)

test_attes()

