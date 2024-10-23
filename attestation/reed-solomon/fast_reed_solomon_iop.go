package reedsolomon

import (
	"bytes"
	"encoding/binary"
	proof "jglez2330/stark-attestation/Proof"
	field "jglez2330/stark-attestation/field"
	merkle "jglez2330/stark-attestation/merkle"
	"jglez2330/stark-attestation/polynomial"
	"math/big"
)

type FastReedSolomonIOP struct {
	domain_size int
	Field       field.Field
	omega       field.FieldElement
    exp_fac     int
}

func NewFastReedSolomonIOP(domain_size int, field field.Field) FastReedSolomonIOP {
	return FastReedSolomonIOP{domain_size, field}
}

func (r FastReedSolomonIOP) Rounds() int {
	count := 0
	domain := r.domain_size
	for domain > 1 {
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

func (r FastReedSolomonIOP) next_polynomial(poly polynomial.Polynomial, alpha field.FieldElement) polynomial.Polynomial {

	// To create the next polynomial we need to cut the degree by 2 and generate a new polynomial
	odd_coefficients := make([]field.FieldElement, len(poly.Coefficients())/2)
	even_coefficients := make([]field.FieldElement, len(poly.Coefficients())/2)
	for i := 0; i < len(poly.Coefficients()); i++ {
		if i%2 == 0 {
			even_coefficients = append(even_coefficients, poly.Coefficients()[i])
		} else {
			odd_coefficients = append(odd_coefficients, poly.Coefficients()[i])
		}
	}
	alpha_poly := polynomial.NewPolynomial([]field.FieldElement{alpha})
	polynomial_odd := polynomial.NewPolynomial(odd_coefficients).Multiply(alpha_poly)
	polynomial_even := polynomial.NewPolynomial(even_coefficients)

	return polynomial_odd.Add(polynomial_even)
}

func (r FastReedSolomonIOP) next_domain(domain []field.FieldElement) []field.FieldElement {
	next_domain := make([]field.FieldElement, 0)
	two := field.FieldElement{big.NewInt(2), r.Field}
	for i := 0; i < len(domain)/2; i++ {
		next_domain = append(next_domain, domain[i].Exp(two))
	}
	return next_domain
}

func (r FastReedSolomonIOP) next_layer(poly polynomial.Polynomial, domain []field.FieldElement) []field.FieldElement {
	return poly.Evaluate_Domain(domain)
}

func (r FastReedSolomonIOP) next_fri_layer(poly polynomial.Polynomial, domain []field.FieldElement, alpha field.FieldElement) (next_poly polynomial.Polynomial, next_domain []field.FieldElement, next_layer []field.FieldElement) {
	next_poly = r.next_polynomial(poly, alpha)
	next_domain = r.next_domain(domain)
	next_layer = r.next_layer(next_poly, next_domain)
	return next_poly, next_domain, next_layer
}

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

func (r FastReedSolomonIOP) Commit(codeword []field.FieldElement, fs_proof proof.Fiat_Shamir_Proof) [][]field.FieldElement {
	codewords := make([][]field.FieldElement, 0)
	offset := field.FieldElement{big.NewInt(1), r.Field}
	omega := r.omega.Multiply(offset)
	for i := 0; i < r.Rounds(); i++ {

		mk := merkle.NewMerkle(codeword)
		layer_root := mk.Root
		fs_proof.Push(layer_root)
		codewords = append(codewords, codeword)

		alpha_raw := binary.BigEndian.Uint64(fs_proof.Prover_fiat_shamir(32))
		alpha := field.FieldElement{big.NewInt(int64(alpha_raw)), r.Field}
		codeword = r.split_fold(codeword, alpha, omega)
		offset = offset.Exp(field.FieldElement{big.NewInt(2), r.Field})
		omega = omega.Exp(field.FieldElement{big.NewInt(2), r.Field}).Multiply(offset)

	}
    fs_proof.Push(codeword)
	codewords = append(codewords, codeword)

	return codewords

}

func (r FastReedSolomonIOP) split_fold(codeword []field.FieldElement, alpha, omega field.FieldElement) (codewords []field.FieldElement) {
	new_codeword := make([]field.FieldElement, 0)
	two := field.FieldElement{big.NewInt(2), r.Field}
	two_inv := two.Inverse()
	for i := 0; i < len(codeword)/2; i += 1 {
		// (alpha +1) * codeword[i]
		lower_idx_num := alpha.Add(r.Field.One()).Multiply(codeword[i])
		// omega^i
		lower_idx_den := omega.Exp(field.FieldElement{big.NewInt(int64(i)), r.Field})
		// (alpha + 1)*codeword[i] / omega^i
		lower_idx := lower_idx_num.Multiply(lower_idx_den.Inverse())
		// (1- alpha) * codeword[i + len(codeword)/2]
		upper_idx_num := r.Field.One().Subtract(lower_idx).Multiply(codeword[len(codeword)/2+i])
		// omega^i
		upper_idx_den := omega.Exp(field.FieldElement{big.NewInt(int64(i)), r.Field})
		// (1-alpha)*codeword[i + len(codeword)/2] / omega^i
		upper_idx := upper_idx_num.Multiply(upper_idx_den.Inverse())

		// (alpha + 1)*codeword[i] /2*omega^i + (1-alpha)*codeword[i + len(codeword)/2] / 2*omega^i
		element := lower_idx.Add(upper_idx).Multiply(two_inv)
		new_codeword = append(new_codeword, element)
	}
	return new_codeword
}

func (r FastReedSolomonIOP) Query(current_word, next_word []field.FieldElement, index int, fs_proof proof.Fiat_Shamir_Proof) int {
	sib_idx := index + (len(current_word) / 2)

	//Tell the verifier the leaves we are querying
	fs_proof.Push(current_word[index].Value.Bytes())
	fs_proof.Push(current_word[sib_idx].Value.Bytes())
	fs_proof.Push(next_word[index].Value.Bytes())

	//Now we need to prove that the leaves are in the merkle tree and the path to the root
	merkle_idx := merkle.NewMerkle(current_word)
	path := merkle_idx.Open_path(index)
	for i := 0; i < len(path); i++ {
		fs_proof.Push(path[i])
	}
	path_sib := merkle_idx.Open_path(sib_idx)
	for i := 0; i < len(path_sib); i++ {
		fs_proof.Push(path_sib[i])
	}
	merkle_next := merkle.NewMerkle(next_word)
	path_next := merkle_next.Open_path(index)
	for i := 0; i < len(path_next); i++ {
		fs_proof.Push(path_next[i])
	}

	return index + sib_idx

}

func (r FastReedSolomonIOP) Prove(codeword []field.FieldElement, fs_proof proof.Fiat_Shamir_Proof) []int {

	codewords := r.Commit(codeword, fs_proof)
	indeces := make([]int, 0)
	for i := 0; i < len(codewords); i++ {
		//indeces = append(indeces, r.random_index(fs_proof))
        indeces = append(indeces, 1)
	}

	//Query phase

	new_indeces := indeces
	for i := 0; i < len(codewords)-1; i++ {
		for j := 0; j < len(codewords[i]); j++ {
			new_indeces[j] = indeces[j] % len(codewords[i]) / 2
		}
		r.Query(codewords[i], codewords[i+1], new_indeces[0], fs_proof)
	}
	return indeces
}

func (r FastReedSolomonIOP) random_index(fs_proof proof.Fiat_Shamir_Proof) int {
	idx_raw := binary.BigEndian.Uint64(fs_proof.Prover_fiat_shamir(32))
	idx := int(idx_raw)
	return idx
}

func (r FastReedSolomonIOP) random_index_ver(fs_proof proof.Fiat_Shamir_Proof) int {
	idx_raw := binary.BigEndian.Uint64(fs_proof.Verifier_fiat_shamir(32))
	idx := int(idx_raw) 
	return idx
}
func (r FastReedSolomonIOP) Verify(fs_proof proof.Fiat_Shamir_Proof) bool {
	offset := field.FieldElement{big.NewInt(1), r.Field}
	omega := r.omega.Multiply(offset)

	//First we need to get the roots and alphas used on the commitment phase
	roots := make([][]byte, 0)
	alphas := make([]int, 0)

	for i := 0; i < r.Rounds(); i++ {
		roots = append(roots, fs_proof.Pop().([]byte))

        alphas = append(alphas, r.random_index_ver(fs_proof))
	}

    //Extract the last codeword
    last_codeword := fs_proof.Pop().([]field.FieldElement)

    //Convwer into field elements
    mk := merkle.NewMerkle(last_codeword)

    //Check if last code word is correct
    if bytes.Compare(roots[len(roots)-1], mk.Root) != 0{
        return false
    }

    //Now we need to verify that the codewords belong to a low degree polynomial
    degree := len(last_codeword)/r.exp_fac - 1

    last_domain := make([]field.FieldElement, 0)
    h := r.Field.Generator_subgroup_order(len(last_codeword)) 
    for i := 0; i < len(last_codeword); i++ {
        h_i := h.Exp(field.FieldElement{big.NewInt(int64(i)),r.Field})
        last_domain = append(last_domain, h_i)
    }

    poly_new := polynomial.Polynomial{}
    poly := poly_new.Interpolate_domain(last_domain, last_codeword)

    if poly.Degree() > degree{
        return false
    }

    for i := 0; i < r.Rounds()-1; i++ {
        indeces := make([]int, 0) 
        for j := 0; j < r.domain_size; j++ {
            //indeces = append(indeces, r.random_index_ver(fs_proof))
            indeces = append(indeces, 1)
        }

        //Query phase
        path_next_word := fs_proof.Pop().([]field.FieldElement)
        Merkle_next := merkle.NewMerkle(path_next_word)
        if !Merkle_next.Verify(roots[i], indeces[0], last_codeword[indeces[0]], path_next_word){
            return false
        }
        path_sib_word := fs_proof.Pop().([]field.FieldElement)
        Merkle_sib := merkle.NewMerkle(path_sib_word)
        if !Merkle_sib.Verify(roots[i], indeces[1], last_codeword[indeces[1]], path_sib_word){
            return false
        }
        path_word := fs_proof.Pop().([]field.FieldElement)
        Merkle := merkle.NewMerkle(path_word)
        if !Merkle.Verify(roots[i], indeces[0], last_codeword[indeces[0]], path_word){
            return false
        }
        omega = omega.Exp(field.FieldElement{big.NewInt(2), r.Field}).Multiply(offset)      
        offset = offset.Exp(field.FieldElement{big.NewInt(2), r.Field})



    }


	return true
}
