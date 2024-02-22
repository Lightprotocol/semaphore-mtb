package prover

import (
	"fmt"
	"light/gnark-merkle/logging"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type InsertionParameters struct {
	Root           []big.Int
	InPathIndices  []uint32
	InPathElements [][]big.Int
	Leaf           []big.Int
}

func (p *InsertionParameters) ValidateShape(treeDepth uint32, numOfUTXOs uint32) error {
	if len(p.Root) != int(numOfUTXOs) || len(p.InPathIndices) != int(numOfUTXOs) || len(p.InPathElements) != int(numOfUTXOs) || len(p.Leaf) != int(numOfUTXOs) {
		return fmt.Errorf("wrong number of utxos: %d", len(p.Root))
	}
	for i, proof := range p.InPathElements {
		if len(proof) != int(treeDepth) {
			return fmt.Errorf("wrong size of merkle proof for proof %d: %d", i, len(proof))
		}
	}
	return nil
}

func BuildR1CSInsertion(treeDepth uint32, numOfUTXOs uint32) (constraint.ConstraintSystem, error) {
	root := make([]frontend.Variable, treeDepth)
	leaf := make([]frontend.Variable, treeDepth)
	inPathIndices := make([]frontend.Variable, treeDepth)
	inPathElements := make([][]frontend.Variable, numOfUTXOs)
	for i := 0; i < int(numOfUTXOs); i++ {
		inPathElements[i] = make([]frontend.Variable, treeDepth)
	}

	circuit := InsertionCircuit{
		Depth:          int(treeDepth),
		NumOfUTXOs:     int(numOfUTXOs),
		InPathIndices:  inPathIndices,
		InPathElements: inPathElements,
		Leaf:           leaf,
		Root:           root,
	}
	return frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
}

func SetupInsertion(treeDepth uint32, batchSize uint32) (*ProvingSystem, error) {
	ccs, err := BuildR1CSInsertion(treeDepth, batchSize)
	if err != nil {
		return nil, err
	}
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		return nil, err
	}
	return &ProvingSystem{treeDepth, batchSize, pk, vk, ccs}, nil
}

func (ps *ProvingSystem) ProveInsertion(params *InsertionParameters) (*Proof, error) {
	fmt.Println("Validating shape")
	if err := params.ValidateShape(ps.TreeDepth, ps.NumOfUTXOs); err != nil {
		return nil, err
	}
	fmt.Println("Validated shape")

	inPathIndices := make([]frontend.Variable, ps.NumOfUTXOs)
	root := make([]frontend.Variable, ps.NumOfUTXOs)
	leaf := make([]frontend.Variable, ps.NumOfUTXOs)
	inPathElements := make([][]frontend.Variable, ps.NumOfUTXOs)

	for i := 0; i < int(ps.NumOfUTXOs); i++ {
		root[i] = params.Root[i]
		leaf[i] = params.Leaf[i]
		inPathIndices[i] = params.InPathIndices[i]
		inPathElements[i] = make([]frontend.Variable, ps.TreeDepth)
		for j := 0; j < int(ps.TreeDepth); j++ {
			inPathElements[i][j] = params.InPathElements[i][j]
		}
	}

	assignment := InsertionCircuit{
		Root:           root,
		Leaf:           leaf,
		InPathIndices:  inPathIndices,
		InPathElements: inPathElements,
		NumOfUTXOs:     int(ps.NumOfUTXOs),
		Depth:          int(ps.TreeDepth),
	}
	logging.Logger().Info().Msg("generating proof")

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	fmt.Println("Witness generated")
	if err != nil {
		return nil, err
	}

	fmt.Println("Proving")
	proof, err := groth16.Prove(ps.ConstraintSystem, ps.ProvingKey, witness)
	fmt.Println("Proved")
	if err != nil {
		return nil, err
	}
	logging.Logger().Info().Msg("proof generated successfully")
	return &Proof{proof}, nil
}

func (ps *ProvingSystem) VerifyInsertion(root big.Int, leaf big.Int, proof *Proof) error {
	//TODO: fix

	roots := make([]frontend.Variable, ps.NumOfUTXOs)
	roots[0] = root

	leafs := make([]frontend.Variable, ps.NumOfUTXOs)
	leafs[0] = leaf

	publicAssignment := InsertionCircuit{
		Root: roots,
		Leaf: leafs,
	}

	witness, err := frontend.NewWitness(&publicAssignment, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		return err
	}
	return groth16.Verify(proof.Proof, ps.VerifyingKey, witness)
}
