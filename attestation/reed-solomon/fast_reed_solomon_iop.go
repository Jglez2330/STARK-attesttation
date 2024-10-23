package reedsolomon

import (
	"bytes"
	"crypto/sha256"
	proof "jglez2330/stark-attestation/Proof"
	field "jglez2330/stark-attestation/field"
	merkle "jglez2330/stark-attestation/merkle"
	"jglez2330/stark-attestation/polynomial"
	"math"
	"math/big"
)

type FastReedSolomonIOP struct {
	domain_size int
	Field       field.Field
	omega       field.FieldElement
    offset      field.FieldElement
    exp_fac     int
}

func NewFastReedSolomonIOP(domain_size, exp_fac, offset int, field_ field.Field) FastReedSolomonIOP {
	return FastReedSolomonIOP{domain_size, field_, field_.Generator_subgroup_order(domain_size*exp_fac), field.FieldElement{big.NewInt(int64(offset)), field_}, exp_fac}
}

func (r *FastReedSolomonIOP) Rounds() int {
	count := 0
	domain := r.domain_size
	for domain > r.exp_fac {
		domain = domain / 2
		count++
	}
	return count
}

// func (r FastReedSolomonIOP) Commit(cp polynomial.Polynomial, domain []field.FieldElement, codeword []field.FieldElement, cp_merkle merkle.Merkle, fs_proof proof.Fiat_Shamir_Proof) ([]polynomial.Polynomial, [][]field.FieldElement, [][]field.FieldElement,[]merkle.Merkle){
//     fri_polys  := []polynomial.Polynomial{cp}
//     fri_domains := [][]field.FieldElement{domain}
//     fri_layers := [][]field.FieldElement{codeword}
//     fri_mk := []merkle.Merkle{cp_merkle}
//     fs_proof.Push(cp_merkle.Root)
//     for i := 0; i < r.Rounds(); i++ {
//         alpha_raw := binary.BigEndian.Uint64(fs_proof.Prover_fiat_shamir(32))
//         alpha := field.FieldElement{big.NewInt(int64(alpha_raw)), r.Field}
//
//         next_poly, next_domain, next_layer := r.next_fri_layer(fri_polys[i], fri_domains[i], alpha)
//         fri_polys = append(fri_polys, next_poly)
//         fri_domains = append(fri_domains, next_domain)
//         fri_layers = append(fri_layers, next_layer)
//         cp_merkle = merkle.NewMerkle(next_layer)
//         fs_proof.Push(cp_merkle.Root)
//         fs_proof.Push(fri_mk[i].Root)
//
//     }
//     last_poly := fri_polys[len(fri_polys)-1]
//     const_last_poly := last_poly.Coefficients()[0]
//     fs_proof.Push(const_last_poly.Value.Bytes())
//
//
//
// 	//Commit to the codeword
// 	return fri_polys, fri_domains, fri_layers, fri_mk
// }

// func (r FastReedSolomonIOP) next_polynomial(poly polynomial.Polynomial, alpha field.FieldElement) polynomial.Polynomial {
//
// 	// To create the next polynomial we need to cut the degree by 2 and generate a new polynomial
// 	odd_coefficients := make([]field.FieldElement, len(poly.Coefficients())/2)
// 	even_coefficients := make([]field.FieldElement, len(poly.Coefficients())/2)
// 	for i := 0; i < len(poly.Coefficients()); i++ {
// 		if i%2 == 0 {
// 			even_coefficients = append(even_coefficients, poly.Coefficients()[i])
// 		} else {
// 			odd_coefficients = append(odd_coefficients, poly.Coefficients()[i])
// 		}
// 	}
// 	alpha_poly := polynomial.NewPolynomial([]field.FieldElement{alpha})
// 	polynomial_odd := polynomial.NewPolynomial(odd_coefficients).Multiply(alpha_poly)
// 	polynomial_even := polynomial.NewPolynomial(even_coefficients)
//
// 	return polynomial_odd.Add(polynomial_even)
// }
//
// func (r FastReedSolomonIOP) next_domain(domain []field.FieldElement) []field.FieldElement {
// 	next_domain := make([]field.FieldElement, 0)
// 	two := field.FieldElement{big.NewInt(2), r.Field}
// 	for i := 0; i < len(domain)/2; i++ {
// 		next_domain = append(next_domain, domain[i].Exp(two))
// 	}
// 	return next_domain
// }
//
// func (r FastReedSolomonIOP) next_layer(poly polynomial.Polynomial, domain []field.FieldElement) []field.FieldElement {
// 	return poly.Evaluate_Domain(domain)
// }
//
// func (r FastReedSolomonIOP) next_fri_layer(poly polynomial.Polynomial, domain []field.FieldElement, alpha field.FieldElement) (next_poly polynomial.Polynomial, next_domain []field.FieldElement, next_layer []field.FieldElement) {
// 	next_poly = r.next_polynomial(poly, alpha)
// 	next_domain = r.next_domain(domain)
// 	next_layer = r.next_layer(next_poly, next_domain)
// 	return next_poly, next_domain, next_layer
// }

// func (r FastReedSolomonIOP) Decommit(idx int, current_word, next_word []field.FieldElement, fs_proof proof.Fiat_Shamir_Proof) int {
//     sib_idx := idx + (len(current_word)/2)
//     idx_bytes := make([]byte, 4)
//     binary.BigEndian.PutUint32(idx_bytes, uint32(idx))
//     fs_proof.Push(idx_bytes)
//     sib_idx_bytes := make([]byte, 4)
//     binary.BigEndian.PutUint32(sib_idx_bytes, uint32(sib_idx))
//     fs_proof.Push(sib_idx_bytes)
//
//     merkle_idx := merkle.NewMerkle(current_word)
//     path := merkle_idx.Open_path(idx)
//     for i := 0; i < len(path); i++ {
//         fs_proof.Push(path[i])
//     }
//     path_sib := merkle_idx.Open_path(sib_idx)
//     for i := 0; i < len(path_sib); i++ {
//         fs_proof.Push(path_sib[i])
//     }
//     merkle_next := merkle.NewMerkle(next_word)
//     path_next := merkle_next.Open_path(idx)
//     for i := 0; i < len(path_next); i++ {
//         fs_proof.Push(path_next[i])
//     }
//
//     return idx + sib_idx
// }

func (r *FastReedSolomonIOP) random_index(byte_array []byte, size int) uint {
    var acc uint64
    for i := 0; i < len(byte_array); i++ {
        acc = uint64(math.Pow(float64(acc<<8), float64(int(byte_array[i]))))
    }

    return uint(acc % uint64(size))
}
func Contains[T comparable](s []T, e T) bool {
    for _, v := range s {
        if v == e {
            return true
        }
    }
    return false
}


func (r *FastReedSolomonIOP) random_indeces(seed []byte, size, reduced_size, number int) []uint {
    indeces := make([]uint, 0)
    reduced_indeces := make([]uint, 0)
    count := 0
    for i := 0; i < number; i++ {
        s := sha256.New()
        s.Write(seed)
        s.Write([]byte{byte(count)})
        index := r.random_index(s.Sum(nil), size)

        reduced_index := index % uint(reduced_size)
        count++
        if !Contains(reduced_indeces, reduced_index) {
            indeces = append(indeces, index)
            reduced_indeces = append(reduced_indeces, reduced_index)
        }
    }
    return indeces
}

func (r FastReedSolomonIOP) Commit(codeword []field.FieldElement, fs_proof proof.Fiat_Shamir_Proof) [][]field.FieldElement {
    omega := r.omega
    offset := r.offset
    codewords := make([][]field.FieldElement, 0)

    for i := 0; i < r.Rounds(); i++ {
        mk := merkle.NewMerkle(codeword)
        root := mk.Root
        fs_proof.Push(root)

        if i == r.Rounds() - 1 {
            break
        }

        alpha := r.Field.Sample(fs_proof.Prover_FS())

        codewords = append(codewords, codeword)

        codeword = r.split_fold(codeword, alpha, omega, offset)

        omega = omega.Multiply(omega)
        offset = offset.Multiply(offset)
    }

    fs_proof.Push(codeword)
    codewords = append(codewords, codeword)

    return codewords
}

func (r FastReedSolomonIOP) split_fold(codeword []field.FieldElement, alpha, omega, offset field.FieldElement) []field.FieldElement {
    one := field.FieldElement{big.NewInt(1), r.Field}
    two := field.FieldElement{big.NewInt(2), r.Field}
    two_inv := two.Inverse()
    codewords := make([]field.FieldElement, 0)
    for i := 0; i < len(codeword)/2; i++ {
        even_num := one.Add(alpha)
        omega_i := omega.Exp(field.FieldElement{big.NewInt(int64(i)), r.Field})
        even_den  := omega_i.Multiply(offset)
        even := even_num.Divide(even_den).Multiply(codeword[i])

        odd_num := one.Subtract(alpha)
        odd_den := omega_i.Multiply(offset)
        odd := odd_num.Divide(odd_den).Multiply(codeword[i + len(codeword)/2])

        res := even.Add(odd)

        codewords = append(codewords, res.Multiply(two_inv))

    }
    return codewords
}

func (r FastReedSolomonIOP) Query(current_word, next_word []field.FieldElement, index int, fs_proof proof.Fiat_Shamir_Proof) int {
    sib_index := index + len(current_word)/2

    leaf_current := current_word[index]
    leaf_sib := current_word[sib_index]
    leaf_next := next_word[index]

    fs_proof.Push(leaf_current)
    fs_proof.Push(leaf_sib)
    fs_proof.Push(leaf_next)

    mk_current := merkle.NewMerkle(current_word)
    path_current := mk_current.Open_path(index)
    path_sib := mk_current.Open_path(sib_index)
    mk_next := merkle.NewMerkle(next_word)
    path_next := mk_next.Open_path(index)

    fs_proof.Push(path_current)
    fs_proof.Push(path_sib)
    fs_proof.Push(path_next)

    return sib_index

}

func (r FastReedSolomonIOP) Prove(codeword []field.FieldElement, fs_proof proof.Fiat_Shamir_Proof) []uint {
    codewords := r.Commit(codeword, fs_proof)

    top_indeces := r.random_indeces(fs_proof.Prover_FS(), len(codewords[0])/2, len(codewords[len(codewords)-1]), 10)
    indeces := make([]uint, 0)
    for i := 0; i < len(top_indeces); i++ {
        indeces = append(indeces, top_indeces[i])
    }

    for i := 0; i < len(codewords)-1; i++ {
        //fold
        for j := 0; j < len(indeces); j++ {
            indeces[j] = uint(indeces[j] % uint(len(codewords[i])/2))
        }
        r.Query(codewords[i], codewords[i+1], int(indeces[0]), fs_proof)
    }

    return top_indeces


}

func get_root(fs_proof proof.Fiat_Shamir_Proof) []byte {
    root_any := fs_proof.Pop().([]any)

    root := make([]byte, 0)
    for i := 0; i < len(root_any); i++ {
        root = append(root, root_any[i].(byte))
    
    }

    return root
    
}

func get_codeword(fs_proof proof.Fiat_Shamir_Proof) []field.FieldElement {
    codeword_any := fs_proof.Pop().([]any)
    code := make([]field.FieldElement, 0)
    for i := 0; i < len(codeword_any); i++ {
        code = append(code, codeword_any[i].(field.FieldElement))
    }

    return code
}

func get_leaf(fs_proof proof.Fiat_Shamir_Proof) field.FieldElement {
    leaf_any := fs_proof.Pop().(any)
    leaf := leaf_any.(field.FieldElement)
    return leaf
}

func get_path(fs_proof proof.Fiat_Shamir_Proof) [][]byte {
    path_any := fs_proof.Pop().([]any)
    path := make([][]byte, 0)
    for i := 0; i < len(path_any); i++ {
        path = append(path, path_any[i].([]byte))
    }
    return path
}


func (r *FastReedSolomonIOP) Verify(fs_proof proof.Fiat_Shamir_Proof) bool {
    omega := r.omega
    offset := r.offset

    roots := make([][]byte, 0)
    alphas := make([]field.FieldElement, 0)
    for i := 0; i < r.Rounds(); i++ {
        roots = append(roots, get_root(fs_proof))
        alphas = append(alphas, r.Field.Sample(fs_proof.Verifier_FS()))
    }

    last_codeword := get_codeword(fs_proof)

    mk := merkle.NewMerkle(last_codeword)
    if bytes.Compare(mk.Root, get_root(fs_proof)) != 0 {
        return false
    }

    degree := len(last_codeword)/r.exp_fac - 1

    last_omega := omega
    last_offset := offset

    for i := 0; i < r.Rounds()-1; i++ {
        last_omega = last_omega.Multiply(last_omega)
        last_offset = last_offset.Multiply(last_offset)
    }

    last_domain := make([]field.FieldElement, 0)
    for i := 0; i < len(last_codeword); i++ {
        omega_i := last_omega.Exp(field.FieldElement{big.NewInt(int64(i)), r.Field})
        last_domain = append(last_domain, omega_i.Multiply(last_offset))
    }

    poly_new := polynomial.Polynomial{}
    poly := poly_new.Interpolate_domain(last_domain, last_codeword)

    if degree > poly.Degree() {
        return false
    }

    top_indeces := r.random_indeces(fs_proof.Verifier_FS(),r.domain_size/2, r.domain_size>>(r.Rounds()-1), 10)

    for i := 0; i < r.Rounds()-1;i++{
        c_indeces := make([]uint, 0)
        for j := 0; j < len(top_indeces); j++ {
            c_indeces = append(c_indeces, top_indeces[j] % uint(r.domain_size>>(i+1)))
        }

        a_idx := c_indeces
        b_idx := make([]uint, 0)
        for j := 0; j < len(a_idx); j++ {
            b_idx = append(b_idx, a_idx[j] + uint(r.domain_size>>(i+1)))
        }

        a_leaf := get_leaf(fs_proof)
        b_leaf := get_leaf(fs_proof)
        c_leaf := get_leaf(fs_proof)

        a_path := get_path(fs_proof)
        b_path := get_path(fs_proof)
        c_path := get_path(fs_proof)

        mk := merkle.Merkle{}

        if mk.Verify(roots[i], int(a_idx[i]), a_leaf, a_path)  {
            return false
        }

        if mk.Verify(roots[i], int(b_idx[i]), b_leaf, b_path)  {
            return false
        }

        if mk.Verify(roots[i+1], int(a_idx[i]), c_leaf, c_path)  {
            return false
        }



    }




    
    return true
}
