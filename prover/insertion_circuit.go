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

	NumOfUTXOs int
	Depth      int
}

func (circuit *InsertionCircuit) Define(api frontend.API) error {
	//api.AssertIsEqual(circuit.Leaf, circuit.Leaf)
	//api.AssertIsEqual(circuit.InPathIndices, circuit.InPathIndices)
	//api.AssertIsEqual(circuit.InPathElements, circuit.InPathElements)

	// Actual batch merkle proof verification.
	root := abstractor.Call(api, InsertionProof{
		Root:           circuit.Root,
		Leaf:           circuit.Leaf,
		InPathIndices:  circuit.InPathIndices,
		InPathElements: circuit.InPathElements,

		NumOfUTXOs: circuit.NumOfUTXOs,
		Depth:      circuit.Depth,
	})

	// Final root needs to match.
	api.AssertIsEqual(root, circuit.Root)

	return nil
}

func ImportInsertionSetup(treeDepth uint32, batchSize uint32, pkPath string, vkPath string) (*ProvingSystem, error) {
	proofs := make([][]frontend.Variable, batchSize)
	for i := 0; i < int(batchSize); i++ {
		proofs[i] = make([]frontend.Variable, treeDepth)
	}

	circuit := InsertionCircuit{
		Root:           make([]frontend.Variable, batchSize),
		Leaf:           make([]frontend.Variable, batchSize),
		InPathIndices:  make([]frontend.Variable, batchSize),
		InPathElements: proofs,
		NumOfUTXOs:     int(batchSize),
		Depth:          int(treeDepth),
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

	return &ProvingSystem{treeDepth, batchSize, pk, vk, ccs}, nil
}
