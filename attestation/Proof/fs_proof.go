package proof

import (
	"crypto/sha256"
	"encoding/json"
)

type Fiat_Shamir_Proof struct{
    objects []any
    index int
}

func NewFiat_Shamir_Proof(objects []any) Fiat_Shamir_Proof{
    return Fiat_Shamir_Proof{objects, 0}
}

func (f Fiat_Shamir_Proof) Push(object any){
    f.objects = append(f.objects, object)
}

func (f Fiat_Shamir_Proof) Pop() any{
    if f.index < len(f.objects){
        f.index++
        return f.objects[f.index-1]
    }
    return nil
}

func (f Fiat_Shamir_Proof) serialize_json() []byte{
    json, err := json.Marshal(f)
    if err != nil{
        return nil
    }
    return json
}


func (f Fiat_Shamir_Proof) deserialize_json(json_s string) Fiat_Shamir_Proof{
    var obj Fiat_Shamir_Proof
    err := json.Unmarshal([]byte(json_s), &obj)
    if err != nil{
        return Fiat_Shamir_Proof{}
    }
    return obj
}

func (f Fiat_Shamir_Proof)  Prover_fiat_shamir(num_bytes int) []byte{
    sha256_hash := sha256.Sum256(f.serialize_json())
    return sha256_hash[:num_bytes]
}

func (f Fiat_Shamir_Proof)  Verifier_fiat_shamir(num_bytes int) []byte{
    json_dump, err := json.Marshal(f.objects[:f.index])
    if err != nil{
        return nil
    }
    sha256_hash := sha256.Sum256(json_dump)
    return sha256_hash[:num_bytes]


}

