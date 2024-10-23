package polynomial

import (
	field "jglez2330/stark-attestation/field"
	"math/big"
)

type Polynomial struct {
	coefficients []field.FieldElement
}

type MultivarPolynomial struct {
	coefficients []Polynomial
}

func NewPolynomial(coefficients []field.FieldElement) Polynomial {
	return Polynomial{coefficients}
}

// The degree of a polynomial is the highest non-zero power of the variable in the polynomial
// For example, the degree of 3x^2 + 2x + 1 is 2
func (p Polynomial) Degree() int {
	if len(p.coefficients) == 0 {
		return -1
	}
	max_degree := -1
	for i := 0; i < len(p.coefficients); i++ {
		if !p.coefficients[i].Is_zero() {
			max_degree = i
		}
	}
	return max_degree
}

func (p Polynomial) Coefficients() []field.FieldElement {
    return p.coefficients
}

func (p Polynomial) Subtract(right Polynomial) Polynomial {
	return p.Add(right.Negate())
}

func (p Polynomial) Negate() Polynomial {
	result := make([]field.FieldElement, len(p.coefficients))
	for i := 0; i < len(p.coefficients); i++ {
		result[i] = p.coefficients[i].Negate()
	}
	return NewPolynomial(result)
}

func (p Polynomial) Add(right Polynomial) Polynomial {
	if p.Degree() == -1 {
		return right
	}
	if right.Degree() == -1 {
		return p
	}
	len_result := max(len(p.coefficients), len(right.coefficients))
	result := make([]field.FieldElement, len_result)
	for i := 0; i < len_result; i++ {
		result[i] = p.coefficients[0].Field.Zero()
	}
	for i := 0; i < len(p.coefficients); i++ {
		result[i] = result[i].Add(p.coefficients[i])
	}
	for i := 0; i < len(right.coefficients); i++ {
		result[i] = result[i].Add(right.coefficients[i])
	}
	return NewPolynomial(result)
}

func (p Polynomial) Multiply(right Polynomial) Polynomial {
	if p.Degree() == -1 || right.Degree() == -1 {
		return Polynomial{make([]field.FieldElement, 0)}
	}
	result := make([]field.FieldElement, len(p.coefficients)+len(right.coefficients)-1)
	for i := 0; i < len(result); i++ {
		result[i] = p.coefficients[0].Field.Zero()
	}
	for i := 0; i < len(p.coefficients); i++ {
		if p.coefficients[i].Is_zero() {
			continue
		}
		for j := 0; j < len(right.coefficients); j++ {
			result[i+j] = result[i+j].Add(p.coefficients[i].Multiply(right.coefficients[j]))
		}
	}
	return NewPolynomial(result)
}

func (p Polynomial) Divide(right Polynomial) (Polynomial, Polynomial) {
	if right.Degree() == -1 {
		panic("Cannot divide by an empty polynomial")
	}
	if p.Degree() < right.Degree() {
		return NewPolynomial(make([]field.FieldElement, 0)), p
	}

	remainder := NewPolynomial(p.coefficients) // Copy coefficients
	quotient_coeff := make([]field.FieldElement, 0)
	for i := 0; i <= p.Degree()-right.Degree()+1; i++ {
		quotient_coeff = append(quotient_coeff, p.coefficients[0].Field.Zero())
	}

	for i := 0; i <= p.Degree()-right.Degree()+1; i++ {
		if remainder.Degree() < right.Degree() {
			break
		}

		lead := remainder.Leading_coefficient().Divide(right.Leading_coefficient())
		degree_diff := remainder.Degree() - right.Degree()
		new_poly_coeff := make([]field.FieldElement, degree_diff+1)
		for j, _ := range new_poly_coeff {
			new_poly_coeff[j] = p.coefficients[0].Field.Zero()
		}
		new_poly_coeff[degree_diff] = lead
		sub_poly := NewPolynomial(new_poly_coeff).Multiply(right)
		quotient_coeff[degree_diff] = lead
		remainder = remainder.Subtract(sub_poly)
	}
	quotient := NewPolynomial(quotient_coeff)
	return quotient, remainder
}

func (p Polynomial) Mod(right Polynomial) Polynomial {
	_, remainder := p.Divide(right)
	return remainder
}

func (p Polynomial) Equals(right Polynomial) bool {
	if p.Degree() != right.Degree() {
		return false
	}
	for i := 0; i < len(p.coefficients); i++ {
		if !p.coefficients[i].Equal(right.coefficients[i]) {
			return false
		}
	}
	return true
}

func (p Polynomial) Interpolate_domain(domain []field.FieldElement, codomain []field.FieldElement) Polynomial {
	fiels := domain[0].Field
	x := NewPolynomial([]field.FieldElement{fiels.Zero(), fiels.One()})
	accumulator := NewPolynomial(make([]field.FieldElement, 0))
	for i := 0; i < len(domain); i++ {
		product := NewPolynomial([]field.FieldElement{codomain[i]})
		for j := 0; j < len(domain); j++ {
			if i == j {
				continue
			}
			inverse_poly := NewPolynomial([]field.FieldElement{domain[i].Subtract(domain[j]).Inverse()})
			subs_x := x.Subtract(NewPolynomial([]field.FieldElement{domain[j]}))
			mult_poly := subs_x.Multiply(inverse_poly)
			product = product.Multiply(mult_poly)
		}
		accumulator = accumulator.Add(product)
	}
	return accumulator

}

func (p Polynomial) Evaluate(point field.FieldElement) field.FieldElement {
	accumulator := point.Field.Zero()
	x := point.Field.One()
	for i := 0; i < len(p.coefficients); i++ {
		accumulator = accumulator.Add(p.coefficients[i].Multiply(x))
		x = x.Multiply(point)
	}
	return accumulator
}

func (p Polynomial) Is_zero() bool {
	for i := 0; i < len(p.coefficients); i++ {
		if !p.coefficients[i].Is_zero() {
			return false
		}
	}
	return true
}

func (p Polynomial) Is_one() bool {
	if p.Degree() != 0 {
		return false
	}
	return p.coefficients[0].Is_one()
}

func (p Polynomial) Is_constant() bool {
	return p.Degree() == 0
}

func (p Polynomial) Zerofier(point field.FieldElement) Polynomial {
	x := NewPolynomial([]field.FieldElement{point.Field.Zero(), point.Field.One()})
	return x.Subtract(NewPolynomial([]field.FieldElement{point}))
}

func (p Polynomial) Zerofier_Domain(domain []field.FieldElement) Polynomial {
	accumulator := NewPolynomial(make([]field.FieldElement, 0))
	for i := 0; i < len(domain); i++ {
		accumulator = accumulator.Multiply(p.Zerofier(domain[i]))
	}
	return accumulator
}

func (p Polynomial) Evaluate_Domain(domain []field.FieldElement) []field.FieldElement {
	var result []field.FieldElement
	for i := 0; i < len(domain); i++ {
		result = append(result, p.Evaluate(domain[i]))
	}
	return result
}

func (p Polynomial) Scale(scalar field.FieldElement) Polynomial {
	var result []field.FieldElement
	for i := 0; i < len(p.coefficients); i++ {
		scaled := scalar.Exp(field.FieldElement{big.NewInt(int64(i)), scalar.Field})
		result = append(result, p.coefficients[i].Multiply(scaled))
	}
	return NewPolynomial(result)
}

func (p Polynomial) Leading_coefficient() field.FieldElement {
	return p.coefficients[p.Degree()]
}

func (p Polynomial) Coefficients() []field.FieldElement {
	return p.coefficients
}
