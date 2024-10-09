from Field import *
import numpy as np

class Polynomial:
    def __init__(self, coefficients, field=None):
        if field is None:
            field = Field.main()
        self.coefficients = np.array([c for c in coefficients])
        self.field = field

    def degree(self):
        if len(self.coefficients) == 0:
            return -1
        if self.is_zero_coefficients():
            return -2
        non_zero_indices = np.nonzero(self.coefficients)
        return np.max(non_zero_indices)

    def is_zero_coefficients(self):
        return np.all(self.coefficients == 0)

    def __neg__(self):
        return Polynomial(-self.coefficients)
    def __add__(self, right):
        if self.degree() < 0:
            return right
        if right.degree() < 0:
            return self
        max_len = max(len(self.coefficients), len(right.coefficients))
        coeffs = np.zeros(max_len, dtype=object)

        coeffs[:len(self.coefficients)] += self.coefficients
        coeffs[:len(right.coefficients)] += right.coefficients

        return Polynomial(coeffs, self.field)

    def __sub__(self, right):
        return self + (-right)

    def __mul__(self, other):
        if self.degree() < 0 or other.degree() < 0:
            return Polynomial([], self.field)
        coeffs = np.zeros(self.degree() + other.degree() + 1, dtype=object)
        for i in range(len(self.coefficients)):
            for j in range(len(other.coefficients)):
                coeffs[i+j] += self.coefficients[i] * other.coefficients[j]
        return Polynomial(coeffs, self.field)

    def __eq__(self, other):
        if self.degree() != other.degree():
            return False
        if self.degree() < 0:
            return True
        return np.all(self.coefficients == other.coefficients)
    def __neq__(self, other):
        return not self == other
    def is_zero(self):
        return self.degree() < 0
    def leading_coefficient(self):
        return self.coefficients[self.degree()]
    def divide(self, numerator, denominator):
        if denominator.degree() < 0:
            return None
        if numerator.degree() < denominator.degree():
            return Polynomial([], self.field), numerator

        field = denominator.coefficients[0].field
        remainder = Polynomial(numerator.coefficients.copy(), self.field)
        quotient_coefficients = np.array([field.zero()] * (numerator.degree() - denominator.degree() + 1), dtype=object)

        for i in range(numerator.degree() - denominator.degree() + 1):
            if remainder.degree() < denominator.degree():
                break
            coefficient = remainder.leading_coefficient() / denominator.leading_coefficient()
            shift = remainder.degree() - denominator.degree()
            subtractee = Polynomial(np.concatenate((np.array([field.zero()] * shift, dtype=object), [coefficient])), self.field) * denominator
            quotient_coefficients[shift] = coefficient
            remainder = remainder - subtractee

        quotient = Polynomial(quotient_coefficients, self.field)
        return quotient, remainder

    def __truediv__(self, other):
        quo, rem = Polynomial.divide(self, other)
        assert rem.is_zero(), "cannot perform polynomial division because remainder is not zero"
        return quo

    def __mod__(self, other):
        quo, rem = Polynomial.divide(self, other)
        return rem
    def __xor__(self, exponent):
        if self.is_zero():
            return Polynomial([], self.field)
        if exponent == 0:
            return Polynomial([self.coefficients[0].field.one()], self.field)
        acc = Polynomial([self.coefficients[0].field.one()], self.field)
        for i in reversed(range(len(bin(exponent)[2:]))):
            acc = acc * acc
            if (1 << i) & exponent != 0:
                acc = acc * self
        return acc

    def evaluate(self, point):
        xi = point.field.one()
        value = point.field.zero()
        for c in self.coefficients:
            value = value + c * xi
            xi = xi * point
        return value

    def evaluate_domain(self, domain):
        return [self.evaluate(d) for d in domain]

    @staticmethod
    def interpolate_domain(domain, values):
        assert len(domain) == len(values), "number of elements in domain does not match number of values -- cannot interpolate"
        assert len(domain) > 0, "cannot interpolate between zero points"
        field = domain[0].field
        x = Polynomial([field.zero(), field.one()], field)
        acc = Polynomial([], field)
        for i in range(len(domain)):
            prod = Polynomial([values[i]], field)
            for j in range(len(domain)):
                if j == i:
                    continue
                prod = prod * (x - Polynomial([domain[j]], field)) * Polynomial([(domain[i] - domain[j]).inverse()], field)
            acc = acc + prod
        return acc

    def zerofier_domain(domain):
        field = domain[0].field
        x = Polynomial([field.zero(), field.one()], field)
        acc = Polynomial([field.one()], field)
        for d in domain:
            acc = acc * (x - Polynomial(np.array([d], dtype=object), field))
        return acc

    def scale(self, factor):
        exponents = np.arange(len(self.coefficients))
        scaled_coefficients = np.array([factor**i for i in exponents], dtype=object) * self.coefficients
        return Polynomial(scaled_coefficients, self.field)

    def test_colinearity(points):
        domain = np.array([p[0] for p in points], dtype=object)
        values = np.array([p[1] for p in points], dtype=object)
        polynomial = Polynomial.interpolate_domain(domain, values)
        return polynomial.degree() <= 1
#%%
