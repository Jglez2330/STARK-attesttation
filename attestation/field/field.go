package field

import "math/big"

type Field interface {
	Zero() FieldElement
	One() FieldElement
	Prime() *big.Int
	Multiply(a, b FieldElement) FieldElement
	Add(a, b FieldElement) FieldElement
	Subtract(a, b FieldElement) FieldElement
	Negate(a FieldElement) FieldElement
	Inverse(a FieldElement) FieldElement
	Divide(a, b FieldElement) FieldElement
	Generator() FieldElement
	Generator_subgroup_order(order int) FieldElement
    Sample([]byte) FieldElement
}

type FieldElement struct {
	Value *big.Int // Value of the field element
	Field Field    // Field to which the element belongs
}

func NewFieldElement(value *big.Int, field Field) FieldElement {
	return FieldElement{value, field}
}

func (f FieldElement) Zero() FieldElement {
	return FieldElement{big.NewInt(0), f.Field}
}

func (f FieldElement) One() FieldElement {
	return FieldElement{big.NewInt(1), f.Field}
}

func (f FieldElement) Multiply(right FieldElement) FieldElement {
	return f.Field.Multiply(f, right)
}

func (f FieldElement) Add(right FieldElement) FieldElement {
	return f.Field.Add(f, right)
}

func (f FieldElement) Subtract(right FieldElement) FieldElement {
	return f.Field.Subtract(f, right)
}

func (f FieldElement) Negate() FieldElement {
	return f.Field.Negate(f)
}

func (f FieldElement) Inverse() FieldElement {
	return f.Field.Inverse(f)
}

func (f FieldElement) Divide(denom FieldElement) FieldElement {
	return f.Field.Divide(f, denom)
}

func (f FieldElement) Is_zero() bool {
	mod_value := new(big.Int).Mod(f.Value, f.Field.Prime())
	return mod_value.Cmp(big.NewInt(0)) == 0
}

func (f FieldElement) Exp(exponent FieldElement) FieldElement {
	acc := f.Field.One()
	val := f
	for i := 0; i < exponent.Value.BitLen(); i++ {
		acc = acc.Multiply(acc)
		if exponent.Value.Bit(i) == 1 {
			acc = acc.Multiply(val)
		}
	}
	return acc
}

func (f FieldElement) Equal(right FieldElement) bool {
	mod_value_a := new(big.Int).Mod(f.Value, f.Field.Prime())
	mod_value_b := new(big.Int).Mod(right.Value, f.Field.Prime())
	return mod_value_a.Cmp(mod_value_b) == 0
}

func (f FieldElement) Not_equal(right FieldElement) bool {
	return f.Value != right.Value
}

func (f FieldElement) Is_one() bool {
	mod_value_a := new(big.Int).Mod(f.Value, f.Field.Prime())
	return mod_value_a.Cmp(big.NewInt(1)) == 0
}
