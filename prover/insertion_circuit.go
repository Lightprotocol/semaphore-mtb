package prover

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/reilabs/gnark-lean-extractor/v2/abstractor"
)

type InsertionMbuCircuit struct {
	// single public input
	InputHash frontend.Variable `gnark:",public"`

	// private inputs, but used as public inputs
	StartIndex frontend.Variable   `gnark:"input"`
	PreRoot    frontend.Variable   `gnark:"input"`
	PostRoot   frontend.Variable   `gnark:"input"`
	IdComms    []frontend.Variable `gnark:"input"`

	// private inputs
	MerkleProofs [][]frontend.Variable `gnark:"input"`

	BatchSize int
	Depth     int
}

func (circuit *InsertionMbuCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(circuit.InputHash, circuit.InputHash)
	api.AssertIsEqual(circuit.StartIndex, circuit.StartIndex)
	api.AssertIsEqual(circuit.PreRoot, circuit.PreRoot)
	api.AssertIsEqual(circuit.StartIndex, circuit.StartIndex)

	// Actual batch merkle proof verification.
	root := abstractor.Call(api, InsertionProof{
		StartIndex: circuit.StartIndex,
		PreRoot:    circuit.PreRoot,
		IdComms:    circuit.IdComms,

		MerkleProofs: circuit.MerkleProofs,

		BatchSize: circuit.BatchSize,
		Depth:     circuit.Depth,
	})

	// Final root needs to match.
	api.AssertIsEqual(root, circuit.PostRoot)

	return nil
}

func ImportInsertionSetup(treeDepth uint32, batchSize uint32, pkPath string, vkPath string) (*ProvingSystem, error) {
	proofs := make([][]frontend.Variable, batchSize)
	for i := 0; i < int(batchSize); i++ {
		proofs[i] = make([]frontend.Variable, treeDepth)
	}
	circuit := InsertionMbuCircuit{
		Depth:        int(treeDepth),
		BatchSize:    int(batchSize),
		IdComms:      make([]frontend.Variable, batchSize),
		MerkleProofs: proofs,
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
