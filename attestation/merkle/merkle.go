package merkle

import (
	"bytes"
	"crypto/sha256"
	"jglez2330/stark-attestation/field"
	"math"
)

type Merkle struct {
	_leaves    []field.FieldElement
	Leaves     [][]byte
	num_leaves int
	height     int

	Root []byte
}

func NewMerkle(leaves []field.FieldElement) Merkle {
	merkle := Merkle{}
	merkle.Root = merkle.Commit(leaves)
	return merkle
}

func (m *Merkle) Commit(leaves []field.FieldElement) []byte {
	hash := sha256.New()
	// The field elements must be converted to hashes to build the merkele tree
	leaves_hashes := make([][]byte, len(leaves))
    m._leaves = leaves
	for i, leaf := range leaves {
		element_bytes := leaf.Value.Bytes()
		hash.Write(element_bytes)
		leaves_hashes[i] = hash.Sum(nil)
		hash.Reset()
	}
	exponent := math.Ceil(math.Log2(float64(len(leaves))))
	m.num_leaves = int(math.Pow(2, exponent))

	m.height = int(exponent)
	m.Leaves = leaves_hashes
	m.build_tree(leaves_hashes)
	m.Root = m.build_tree(leaves_hashes)
	return m.Root
}

func (m *Merkle) build_tree(leaves [][]byte) []byte {
	if len(leaves) == 1 {
		return leaves[0]
	} else {
		hash_left := m.build_tree(leaves[:len(leaves)/2])
		hash_right := m.build_tree(leaves[len(leaves)/2:])
		return m.hash(hash_left, hash_right)
	}
}

func (m *Merkle) hash(left []byte, right []byte) []byte {
	hash := sha256.New()
	hash.Write(left)
	hash.Write(right)
	return hash.Sum(nil)
}

func (m *Merkle) Open_path(leaf_index int, path_auth ...[][]byte) [][]byte {
	path := make([][]byte, 0)
	if len(path_auth) == 0 {
		if m.Leaves == nil {
			panic("No leaves")
		}
		path = m.Leaves
	}

	return m.auth_path(leaf_index, path)
}

func (m *Merkle) auth_path(leaf_index int, path [][]byte) [][]byte {
	if !(0 <= leaf_index && leaf_index < len(path)) {
		panic("Index out of range %i")
	}

	if len(path) == 2 {
		return [][]byte{path[1-leaf_index]}
	} else if leaf_index < len(path)/2 {
		return append(m.auth_path(leaf_index, path[:len(path)/2]), m.build_tree(path[len(path)/2:]))
	} else {
		return append(m.auth_path(leaf_index-len(path)/2, path[len(path)/2:]), m.build_tree(path[:len(path)/2]))
	}

}

func (m *Merkle) Verify(root []byte, leaf_index int, element field.FieldElement, path [][]byte) bool {
	h := sha256.New()
	h.Write(element.Value.Bytes())
	hash_element := h.Sum(nil)

	return m.verify_path(root, leaf_index, hash_element, path)
}

func (m *Merkle) verify_path(root []byte, leaf_index int, element_hash []byte, path [][]byte) bool {
	max_size := int(math.Pow(2, float64(len(path))))
	if !(0 <= leaf_index && leaf_index < max_size) {
		panic("Index out of range")
	}

	if len(path) == 1 {
		if leaf_index == 0 {
			right_hash := m.hash(element_hash, path[0])
			return bytes.Equal(root, right_hash)
		} else {
			left_hash := m.hash(path[0], element_hash)
			return bytes.Equal(root, left_hash)
		}
	}

	if leaf_index%2 == 0 {
		right_hash := m.hash(element_hash, path[0])
		return m.verify_path(root, leaf_index/2, right_hash, path[1:])
	} else {
		left_hash := m.hash(path[0], element_hash)
		return m.verify_path(root, leaf_index/2, left_hash, path[1:])
	}
}
