package merkle

import (
	"crypto/sha256"
	"jglez2330/stark-attestation/field"
	"math/big"
	"testing"
)

func TestMerkleBuild(t *testing.T) {
	// Test the MerkleBuild function
	// Create a new Merkle tree

	size := 4
	leaves := make([]field.FieldElement, size)
	for i := 0; i < size; i++ {
		leaves[i] = field.NewFieldElement(big.NewInt(int64(i)), field.Field_252{})
	}

	merkle := NewMerkle(leaves)

	// Test open leaves
	for i := 0; i < size; i++ {
		path := merkle.Open_path(i)
		if !merkle.Verify(merkle.Root, i, leaves[i], path) {
			t.Errorf("Bad path")
		}
	}

}

func TestWrongLeafOpen(t *testing.T) {
	// Test the MerkleBuild function
	// Create a new Merkle tree

	size := 4
	leaves := make([]field.FieldElement, size)
	for i := 0; i < size; i++ {
		leaves[i] = field.NewFieldElement(big.NewInt(int64(i)), field.Field_252{})
	}

	merkle := NewMerkle(leaves)

	// Test open leaves
	for i := 0; i < size; i++ {
		path := merkle.Open_path(i)
		if merkle.Verify(merkle.Root, i, field.NewFieldElement(big.NewInt(10), field.Field_252{}), path) {
			t.Errorf("Bad path")
		}
	}

}

func TestOpenFalseRoot(t *testing.T) {
    // Test the MerkleBuild function
    // Create a new Merkle tree

    size := 4
    leaves := make([]field.FieldElement, size)
    for i := 0; i < size; i++ {
        leaves[i] = field.NewFieldElement(big.NewInt(int64(i)), field.Field_252{})
    }

    merkle := NewMerkle(leaves)
    invalid_root := field.NewFieldElement(big.NewInt(10), field.Field_252{})
    h := sha256.New()
    h.Write(invalid_root.Value.Bytes())
    hash_root := h.Sum(nil)

    // Test open leaves
    for i := 0; i < size; i++ {
        path := merkle.Open_path(i)
        if merkle.Verify(hash_root, i, leaves[i], path) {
            t.Errorf("Bad path")
        }
    }

}

func TestOpenFalsePath(t *testing.T) {
    // Test the MerkleBuild function
    // Create a new Merkle tree

    size := 4
    leaves := make([]field.FieldElement, size)
    for i := 0; i < size; i++ {
        leaves[i] = field.NewFieldElement(big.NewInt(int64(i)), field.Field_252{})
    }

    merkle := NewMerkle(leaves)

    // Test open leaves
    for i := 0; i < size; i++ {
        path := merkle.Open_path(i)
        invalid := field.NewFieldElement(big.NewInt(10), field.Field_252{})
        h := sha256.New()
        h.Write(invalid.Value.Bytes())
        hash_invalid := h.Sum(nil)
        path[0] = hash_invalid
        if merkle.Verify(merkle.Root, i, leaves[i], path) {
            t.Errorf("Bad path")
        }
    }

}
