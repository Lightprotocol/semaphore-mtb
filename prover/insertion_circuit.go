package prover

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/reilabs/gnark-lean-extractor/v2/abstractor"
)

type InsertionCircuit struct {
	// public inputs
	Root []frontend.Variable `gnark:",public"`
	Leaf []frontend.Variable `gnark:",public"`

	// private inputs
	InPathIndices  []frontend.Variable   `gnark:"input"`
	InPathElements [][]frontend.Variable `gnark:"input"`

	NumberOfUtxos int
	Depth         int
}

func (circuit *InsertionCircuit) Define(api frontend.API) error {
	for i := 0; i < circuit.NumberOfUtxos; i++ {
		api.AssertIsEqual(circuit.Root[i], circuit.Root[i])
		api.AssertIsEqual(circuit.Leaf[i], circuit.Leaf[i])
		api.AssertIsEqual(circuit.InPathIndices[i], circuit.InPathIndices[i])
		for j := 0; j < circuit.Depth; j++ {
			api.AssertIsEqual(circuit.InPathElements[i][j], circuit.InPathElements[i][j])
		}
	}

	// Actual merkle proof verification.
	_ = abstractor.Call1(api, InsertionProof{
		Root:           circuit.Root,
		Leaf:           circuit.Leaf,
		InPathElements: circuit.InPathElements,
		InPathIndices:  circuit.InPathIndices,

		NumberOfUtxos: circuit.NumberOfUtxos,
		Depth:         circuit.Depth,
	})

	return nil
}

func ImportInsertionSetup(treeDepth uint32, numberOfUtxos uint32, pkPath string, vkPath string) (*ProvingSystem, error) {
	root := make([]frontend.Variable, numberOfUtxos)
	leaf := make([]frontend.Variable, numberOfUtxos)
	inPathIndices := make([]frontend.Variable, numberOfUtxos)
	inPathElements := make([][]frontend.Variable, numberOfUtxos)

	for i := 0; i < int(numberOfUtxos); i++ {
		inPathElements[i] = make([]frontend.Variable, treeDepth)
	}

	circuit := InsertionCircuit{
		Depth:          int(treeDepth),
		NumberOfUtxos:  int(numberOfUtxos),
		Root:           root,
		Leaf:           leaf,
		InPathIndices:  inPathIndices,
		InPathElements: inPathElements,
	}

	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return nil, err
	}

	pk, err := LoadProvingKey(pkPath)

	if err != nil {
		return nil, err
	}

	vk, err := LoadVerifyingKey(vkPath)

	if err != nil {
		return nil, err
	}

	return &ProvingSystem{treeDepth, numberOfUtxos, pk, vk, ccs}, nil
}
