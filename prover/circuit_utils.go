package prover

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/reilabs/gnark-lean-extractor/v2/abstractor"
	"light/gnark-merkle/logging"
	"light/gnark-merkle/prover/poseidon"
	"os"
)

type Proof struct {
	Proof groth16.Proof
}

type ProvingSystem struct {
	TreeDepth        uint32
	NumOfUTXOs       uint32
	ProvingKey       groth16.ProvingKey
	VerifyingKey     groth16.VerifyingKey
	ConstraintSystem constraint.ConstraintSystem
}

type ProofRound struct {
	Direction frontend.Variable
	Hash      frontend.Variable
	Sibling   frontend.Variable
}

func (gadget ProofRound) DefineGadget(api frontend.API) interface{} {
	api.AssertIsBoolean(gadget.Direction)
	d1 := api.Select(gadget.Direction, gadget.Hash, gadget.Sibling)
	d2 := api.Select(gadget.Direction, gadget.Sibling, gadget.Hash)
	sum := abstractor.Call(api, poseidon.Poseidon2{In1: d1, In2: d2})
	return sum
}

type VerifyProof struct {
	Proof []frontend.Variable
	Path  []frontend.Variable
}

func (gadget VerifyProof) DefineGadget(api frontend.API) interface{} {
	sum := gadget.Proof[0]
	for i := 1; i < len(gadget.Proof); i++ {
		sum = abstractor.Call(api, ProofRound{Direction: gadget.Path[i-1], Hash: gadget.Proof[i], Sibling: sum})
	}
	return sum
}

type InsertionRound struct {
	Index    frontend.Variable
	Item     frontend.Variable
	PrevRoot frontend.Variable
	Proof    []frontend.Variable

	Depth int
}

func (gadget InsertionRound) DefineGadget(api frontend.API) interface{} {
	currentPath := api.ToBinary(gadget.Index, gadget.Depth)

	proof := append([]frontend.Variable{gadget.Item}, gadget.Proof[:]...)
	root := abstractor.Call(api, VerifyProof{Proof: proof, Path: currentPath})

	return root
}

type InsertionProof struct {
	Root           []frontend.Variable
	Leaf           []frontend.Variable
	InPathIndices  []frontend.Variable
	InPathElements [][]frontend.Variable

	NumOfUTXOs int
	Depth      int
}

func (gadget InsertionProof) DefineGadget(api frontend.API) interface{} {
	currentHash := gadget.Leaf
	fmt.Println("currentHash: ", currentHash)
	fmt.Println("gadget.NumOfUTXOs: ", gadget.NumOfUTXOs)
	fmt.Println("gadget.Depth: ", gadget.Depth)

	nextHash := currentHash[0]
	for i := 0; i < gadget.NumOfUTXOs; i++ {
		for j := 0; j < gadget.Depth; j++ {
			el := gadget.InPathElements[i][j]
			fmt.Println("inPathElements[", i, "][", j, "]: ", el)
			nextHash = abstractor.Call(api, ProofRound{Direction: gadget.InPathIndices[i], Hash: nextHash, Sibling: gadget.InPathElements[i][j]})
			fmt.Println("nextHash: ", nextHash)
		}
		currentHash[i] = nextHash
	}
	return nextHash
}

// Trusted setup utility functions
// Taken from: https://github.com/bnb-chain/zkbnb/blob/master/common/prove/proof_keys.go#L19
func LoadProvingKey(filepath string) (pk groth16.ProvingKey, err error) {
	logging.Logger().Info().Msg("start reading proving key")
	pk = groth16.NewProvingKey(ecc.BN254)
	f, _ := os.Open(filepath)
	_, err = pk.ReadFrom(f)
	if err != nil {
		return pk, fmt.Errorf("read file error")
	}
	f.Close()

	return pk, nil
}

// Taken from: https://github.com/bnb-chain/zkbnb/blob/master/common/prove/proof_keys.go#L32
func LoadVerifyingKey(filepath string) (verifyingKey groth16.VerifyingKey, err error) {
	logging.Logger().Info().Msg("start reading verifying key")
	verifyingKey = groth16.NewVerifyingKey(ecc.BN254)
	f, _ := os.Open(filepath)
	_, err = verifyingKey.ReadFrom(f)
	if err != nil {
		return verifyingKey, fmt.Errorf("read file error")
	}
	f.Close()

	return verifyingKey, nil
}
