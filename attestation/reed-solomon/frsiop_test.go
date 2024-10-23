package reedsolomon

import (
	"jglez2330/stark-attestation/field"
	"jglez2330/stark-attestation/polynomial"
	"math/big"
	"testing"
)


func TestFRI (t *testing.T) {
    // Test the FRI function
    field_252 := field.Field_252{}
    degree := 63
    exp_fac := 4

    h := field_252.Generator_subgroup_order((degree + 1) * exp_fac)
    H := make([]field.FieldElement, 0)
    for i := 0; i < (degree+1)*exp_fac ; i++ {
        H = append(H, h.Exp(field.FieldElement{big.NewInt(int64(i)), field_252}))
    }

    put := NewFastReedSolomonIOP((degree+1)*exp_fac, field_252)

    poly_coeff := make([]field.FieldElement, 0)
    for i := 0; i < degree+1; i++ {
        poly_coeff = append(poly_coeff, field.FieldElement{big.NewInt(int64(i)), field_252})
    }

    poly := polynomial.NewPolynomial(poly_coeff)

    codeword := poly.Evaluate_Domain(H)

    

    

}
