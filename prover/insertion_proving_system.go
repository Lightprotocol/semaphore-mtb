package prover

import (
	"fmt"
	"light/light-prover/logging"
	"math/big"
	"strconv"

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

func (p *InsertionParameters) NumberOfUTXOs() uint32 {
	return uint32(len(p.Root))
}

func (p *InsertionParameters) TreeDepth() uint32 {
	if len(p.InPathElements) == 0 {
		return 0
	}
	return uint32(len(p.InPathElements[0]))
}

func (p *InsertionParameters) ValidateShape(treeDepth uint32, numOfUTXOs uint32) error {
	if p.NumberOfUTXOs() != numOfUTXOs {
		return fmt.Errorf("wrong number of utxos: %d", len(p.Root))
	}
	if p.TreeDepth() != treeDepth {
		return fmt.Errorf("wrong size of merkle proof for proof %d: %d", p.NumberOfUTXOs(), p.TreeDepth())
	}
	return nil
}

func R1CSInsertion(treeDepth uint32, numberOfUtxos uint32) (constraint.ConstraintSystem, error) {
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
	return frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
}

func SetupInsertion(treeDepth uint32, numberOfUtxos uint32) (*ProvingSystem, error) {
	ccs, err := R1CSInsertion(treeDepth, numberOfUtxos)
	if err != nil {
		return nil, err
	}
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		return nil, err
	}
	return &ProvingSystem{treeDepth, numberOfUtxos, pk, vk, ccs}, nil
}

func (ps *ProvingSystem) ProveInsertion(params *InsertionParameters) (*Proof, error) {
	if err := params.ValidateShape(ps.TreeDepth, ps.NumberOfUtxos); err != nil {
		return nil, err
	}

	inPathIndices := make([]frontend.Variable, ps.NumberOfUtxos)
	root := make([]frontend.Variable, ps.NumberOfUtxos)
	leaf := make([]frontend.Variable, ps.NumberOfUtxos)
	inPathElements := make([][]frontend.Variable, ps.NumberOfUtxos)

	for i := 0; i < int(ps.NumberOfUtxos); i++ {
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
	}

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, err
	}

	logging.Logger().Info().Msg("Proof " + strconv.Itoa(int(ps.TreeDepth)) + " " + strconv.Itoa(int(ps.NumberOfUtxos)))
	proof, err := groth16.Prove(ps.ConstraintSystem, ps.ProvingKey, witness)
	if err != nil {
		return nil, err
	}

	return &Proof{proof}, nil
}

func (ps *ProvingSystem) VerifyInsertion(root []big.Int, leaf []big.Int, proof *Proof) error {
	leafArray := make([]frontend.Variable, ps.NumberOfUtxos)
	for i, v := range leaf {
		leafArray[i] = v
	}

	rootArray := make([]frontend.Variable, ps.NumberOfUtxos)
	for i, v := range root {
		rootArray[i] = v
	}

	publicAssignment := InsertionCircuit{
		Leaf: leafArray,
		Root: rootArray,
	}
	witness, err := frontend.NewWitness(&publicAssignment, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		return err
	}
	return groth16.Verify(proof.Proof, ps.VerifyingKey, witness)
}
