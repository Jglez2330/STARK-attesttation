package proof

import (
	"crypto/sha256"
	"encoding/json"
)

type Fiat_Shamir_Proof struct {
	Objects []any `json:"objects"`
	Index   int   `json:"index"`
}

func NewFiat_Shamir_Proof(objects []any) Fiat_Shamir_Proof {
	return Fiat_Shamir_Proof{objects, 0}
}

func (f *Fiat_Shamir_Proof) Push(object any) {
	f.Objects = append(f.Objects, object)
}

func (f *Fiat_Shamir_Proof) Pop() any {
	if f.Index < len(f.Objects) {
		f.Index++
		return f.Objects[f.Index-1]
	}
	return nil
}

func (f *Fiat_Shamir_Proof) Serialize() []byte {
	data, err := json.Marshal(f)
	if err != nil {
		panic(err)
	}
	return data
}

func Deserialize(data []byte) Fiat_Shamir_Proof {
	var res Fiat_Shamir_Proof
	err := json.Unmarshal(data, &res)
	if err != nil {
		panic(err)
	}
	return res
}

func (f *Fiat_Shamir_Proof) Prover_FS() []byte {
	h := sha256.New()
	h.Write(f.Serialize())
	return h.Sum(nil)
}

func (f *Fiat_Shamir_Proof) Verifier_FS() []byte {
	data, err := json.Marshal(f.Objects[:f.Index])
	if err != nil {
		panic(err)
	}
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}
