package field

import (
	"math"
	"math/big"
)

type Field_252 struct {
	_generator FieldElement
}

func Xgcd(x, y *big.Int) (*big.Int, *big.Int, *big.Int) {
	a, b := new(big.Int), new(big.Int)
	gcd := new(big.Int).GCD(a, b, x, y)

	return a, b, gcd

}

func (f Field_252) Zero() FieldElement {
	return FieldElement{big.NewInt(0), f}
}

func (f Field_252) One() FieldElement {
	return FieldElement{big.NewInt(1), f}
}

func (f Field_252) Prime() *big.Int {
	//Prime = 2^251 + 17*2^192 + 1
	first_term := new(big.Int).Exp(big.NewInt(2), big.NewInt(251), nil)
	second_term := new(big.Int).Mul(big.NewInt(17), new(big.Int).Exp(big.NewInt(2), big.NewInt(192), nil))
	return new(big.Int).Add(first_term, new(big.Int).Add(second_term, big.NewInt(1)))
}

func (f Field_252) Multiply(a, b FieldElement) FieldElement {
	value := new(big.Int).Mod(new(big.Int).Mul(a.Value, b.Value), f.Prime())
	return FieldElement{value, f}
}

func (f Field_252) Add(a, b FieldElement) FieldElement {
	value := new(big.Int).Mod(new(big.Int).Add(a.Value, b.Value), f.Prime())
	return FieldElement{value, f}
}

func (f Field_252) Subtract(a, b FieldElement) FieldElement {
	first_term := new(big.Int).Add(f.Prime(), a.Value)
	second_term := new(big.Int).Sub(first_term, b.Value)
	value := new(big.Int).Mod(second_term, f.Prime())
	return FieldElement{value, f}
}

func (f Field_252) Negate(a FieldElement) FieldElement {
	first_term := new(big.Int).Sub(f.Prime(), a.Value)
	value := new(big.Int).Mod(first_term, f.Prime())
	return FieldElement{value, f}
}

func (f Field_252) Inverse(a FieldElement) FieldElement {
	a1, _, _ := Xgcd(a.Value, f.Prime())
	first_term := new(big.Int).Mod(a1, f.Prime())
	add_p := new(big.Int).Add(first_term, f.Prime())
	value := new(big.Int).Mod(add_p, f.Prime())
	return FieldElement{value, f}
}

func (f Field_252) Divide(a, b FieldElement) FieldElement {
	if b.Is_zero() {
		panic("Division by zero")
	}
	a1, _, _ := Xgcd(b.Value, f.Prime())
	first_term := new(big.Int).Mul(a.Value, a1)
	value := new(big.Int).Mod(first_term, f.Prime())
	return FieldElement{value, f}
}

func (f Field_252) Generator() FieldElement {
	return f._generator
}

func (f Field_252) Generator_subgroup_order(order int) FieldElement {
	return f._generator
}

func (f Field_252) Is_zero(a FieldElement) bool {
	mod_value_a := new(big.Int).Mod(a.Value, f.Prime())
	return mod_value_a.Cmp(big.NewInt(0)) == 0
}

func (f Field_252) Xor(a, b FieldElement) FieldElement {
	return f.Add(a, b)
}

func (f Field_252) Pow(a FieldElement, b int) FieldElement {
	result := f.One()
	for i := 0; i < b; i++ {
		result = f.Multiply(result, a)
	}
	return result
}

func (f Field_252) Sample(b []byte) FieldElement {
    acc := 0
    for i := 0; i < len(b); i++ {
        acc = int(math.Pow(float64(acc<<8), float64(int(b[i]))))
    }
    return FieldElement{big.NewInt(int64(acc)), f}
}


