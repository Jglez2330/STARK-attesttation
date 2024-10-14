package field

import (
	"math/big"
	"testing"
)

func TestField_252(t *testing.T) {
	f := field_252{}

	fe := NewFieldElement(f.prime(), f)

	if !fe.Is_zero() {
		t.Errorf("Expected true, got false")
	}
}

func TestAdd(t *testing.T) {
	f := field_252{}

	fe1 := NewFieldElement(big.NewInt(1), f)
	fe2 := NewFieldElement(big.NewInt(1), f)

	add := fe1.add(fe2)

	if add.value.Cmp(big.NewInt(2)) != 0 {
		t.Errorf("Expected true, got false")
	}
}

func TestSubtract(t *testing.T) {
	f := field_252{}

	fe1 := NewFieldElement(big.NewInt(1), f)
	fe2 := NewFieldElement(big.NewInt(1), f)

	subtract := fe1.subtract(fe2)

	if subtract.value.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected true, got false")
	}
}

func TestMultiply(t *testing.T) {
	f := field_252{}

	fe1 := NewFieldElement(big.NewInt(2), f)
	fe2 := NewFieldElement(big.NewInt(2), f)

	multiply := fe1.multiply(fe2)

	if multiply.value.Cmp(big.NewInt(4)) != 0 {
		t.Errorf("Expected true, got false")
	}
}

func TestNegate(t *testing.T) {
	f := field_252{}

	fe := NewFieldElement(big.NewInt(1), f)

	negate := fe.negate()

	neg := new(big.Int).Sub(f.prime(), big.NewInt(1))

	if negate.value.Cmp(neg) != 0 {
		t.Errorf("Expected true, got false")
	}
}

func TestXGDC(t *testing.T) {
	tests := []struct {
		x, y string // x and y as strings to be converted to big.Int
		gcd  string // expected gcd
	}{
		{"240", "46", "2"},
		{"120", "23", "1"},
		{"56", "15", "1"},
		{"10", "0", "10"}, // GCD of any number with 0 is the number itself
	}

	for _, test := range tests {
		// Convert string inputs to big.Int
		x := new(big.Int)
		y := new(big.Int)
		expectedGcd := new(big.Int)
		x.SetString(test.x, 10)
		y.SetString(test.y, 10)
		expectedGcd.SetString(test.gcd, 10)

		// Call the xgcd function
		a, b, gcd := Xgcd(x, y)

		// Check if gcd matches expected
		if gcd.Cmp(expectedGcd) != 0 {
			t.Errorf("GCD(%s, %s) = %s; want %s", x, y, gcd, expectedGcd)
		}

		// Verify a * x + b * y = gcd
		result := new(big.Int).Mul(a, x)           // a * x
		result.Add(result, new(big.Int).Mul(b, y)) // a * x + b * y

		if result.Cmp(gcd) != 0 {
			t.Errorf("Failed equation: %s * %s + %s * %s = %s; want %s", a, x, b, y, result, gcd)
		}
	}
}

func TestInverse(t *testing.T) {
	f := field_252{}

	fe := NewFieldElement(big.NewInt(2), f)

	inverse := fe.inverse()

	inv := new(big.Int).ModInverse(big.NewInt(2), f.prime())

	if inverse.value.Cmp(inv) != 0 {
		t.Errorf("Expected true, got false")
	}
}

func TestDivide(t *testing.T) {
	f := field_252{}

	fe1 := NewFieldElement(big.NewInt(6), f)
	fe2 := NewFieldElement(big.NewInt(2), f)

	divide := fe1.divide(fe2)
	a := divide.value.Int64()

	if divide.value.Cmp(big.NewInt(3)) != 0 {
		t.Errorf("Expected true, got false %d", a)
	}
}
