package proof

import (
	"testing"
)

func TestFSProof(t *testing.T) {
	// Test the Fiat-Shamir proof
	objects := []any{[]int{1, 2, 3}}
	fsProof := NewFiat_Shamir_Proof(objects)
	fsProof.Push(4)
	fsProof.Push(5)
	fsProof.Push(6)

	data := fsProof.Serialize()
	fsProof2 := Deserialize(data)

	pop1, ok := fsProof.Pop().([]int)
	if !ok {
		t.Errorf("First index is not slice int")
	}
	pop2, ok2 := fsProof2.Pop().([]any)
	if !ok2 {
		t.Errorf("First index is not slice int")
	}
	for i := 0; i < len(pop1); i++ {
		if int(pop2[i].(float64)) != pop1[i] {
			t.Errorf("Index are not equal")
		}

	}

}
