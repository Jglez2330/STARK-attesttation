package polynomial

import (
	field "jglez2330/stark-attestation/field"
	"math/big"
	"testing"
)

func TestDistributiveProperty(t *testing.T) {
	// Test the distributive property of polynomials

	Field_252 := field.Field_252{}
	// Create the field elements

	zero := Field_252.Zero()
	one := Field_252.One()
	a := field.FieldElement{big.NewInt(2), Field_252}
	b := field.FieldElement{big.NewInt(3), Field_252}

	// Create the polynomials
	p1 := NewPolynomial([]field.FieldElement{a, b, zero})
	p2 := NewPolynomial([]field.FieldElement{a, b, one})
	p3 := NewPolynomial([]field.FieldElement{a, b, a})

	// Test the distributive property
	// (p1 + p2) * p3 = p1 * p3 + p2 * p3
	left := p1.Add(p2).Multiply(p3)
	right := p1.Multiply(p3).Add(p2.Multiply(p3))

	if !left.Equals(right) {
		t.Errorf("Distributive property failed")
	}

}

func TestDivision(t *testing.T) {
	// Test the division of polynomials

	Field_252 := field.Field_252{}
	// Create the field elements

	zero := Field_252.Zero()

	one := Field_252.One()
	a := field.FieldElement{big.NewInt(2), Field_252}
	b := field.FieldElement{big.NewInt(3), Field_252}

	// Create the polynomials
	p1 := NewPolynomial([]field.FieldElement{a, b, one})
	p2 := NewPolynomial([]field.FieldElement{a.Multiply(a), b, one})

	print(p1.Degree())
	print(p2.Degree())

	// Test the division
	// p1 / p2 = 1
	result, rem := p1.Divide(p2)
	for i := 0; i < len(result.coefficients); i++ {
		a := result.coefficients[i].Value.Int64()
		print(a)
	}
	for i := 0; i < len(rem.coefficients); i++ {
		a := rem.coefficients[i].Value.Uint64()
		print(a)
	}

	if !rem.Equals(NewPolynomial([]field.FieldElement{a.Negate(), zero, zero})) {
		t.Errorf("Division failed")
	}

}

func TestPolynomial_Interpolate_domain(t *testing.T) {
	// Test the interpolation of polynomials

	Field_252 := field.Field_252{}
	// Create the field elements

	zero := Field_252.Zero()
	one := Field_252.One()
	a := field.FieldElement{big.NewInt(2), Field_252}
	b := field.FieldElement{big.NewInt(3), Field_252}
	c := field.FieldElement{big.NewInt(4), Field_252}
	d := field.FieldElement{big.NewInt(5), Field_252}

	// Create the domain and codomain
	domain := []field.FieldElement{a, b, c, d}
	codomain := []field.FieldElement{a, b, c, d}

	// Create the polynomial
	p := NewPolynomial([]field.FieldElement{one, zero, zero, zero})

	// Test the interpolation
	// p = 0
	p = p.Interpolate_domain(domain, codomain)
	if !(p.Degree() == 1) {
		t.Errorf("Interpolation failed")
	}
}
