package prover

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"io"
	"math/big"
	"os"
)

func fromHex(i *big.Int, s string) error {
	_, ok := i.SetString(s, 0)
	if !ok {
		return fmt.Errorf("invalid number: %s", s)
	}
	return nil
}

func toHex(i *big.Int) string {
	return fmt.Sprintf("0x%s", i.Text(16))
}

type ProofDataJSON struct {
	Root           []string   `json:"root"`
	InPathIndices  []uint32   `json:"inPathIndices"`
	InPathElements [][]string `json:"inPathElements"`
	Leaf           []string   `json:"leaf"`
}

func ParseInput(inputJSON string) (InsertionParameters, error) {
	var proofData InsertionParameters
	err := json.Unmarshal([]byte(inputJSON), &proofData)
	if err != nil {
		return InsertionParameters{}, fmt.Errorf("error parsing JSON: %v", err)
	}
	return proofData, nil
}

func (p *InsertionParameters) MarshalJSON() ([]byte, error) {
	paramsJson := ProofDataJSON{}

	paramsJson.Root = make([]string, len(p.Root))
	for i := 0; i < len(p.Root); i++ {
		paramsJson.Root[i] = toHex(&p.Root[i])
	}

	paramsJson.InPathIndices = make([]uint32, len(p.InPathIndices))
	for i := 0; i < len(p.InPathIndices); i++ {
		paramsJson.InPathIndices[i] = p.InPathIndices[i]
	}

	paramsJson.InPathElements = make([][]string, len(p.InPathElements))
	for i := 0; i < len(p.InPathElements); i++ {
		for j := 0; j < len(p.InPathElements[i]); j++ {
			paramsJson.InPathElements[i][j] = toHex(&p.InPathElements[i][j])
		}
	}

	paramsJson.Leaf = make([]string, len(p.Leaf))
	for i := 0; i < len(p.Leaf); i++ {
		paramsJson.Leaf[i] = toHex(&p.Leaf[i])
	}

	return json.Marshal(paramsJson)
}

func (p *InsertionParameters) UnmarshalJSON(data []byte) error {

	var params ProofDataJSON

	err := json.Unmarshal(data, &params)
	if err != nil {
		return err
	}

	p.Root = make([]big.Int, len(params.Root))
	for i := 0; i < len(params.Root); i++ {
		err = fromHex(&p.Root[i], params.Root[i])
		if err != nil {
			return err
		}
	}

	p.Leaf = make([]big.Int, len(params.Leaf))
	for i := 0; i < len(params.Leaf); i++ {
		err = fromHex(&p.Leaf[i], params.Leaf[i])
		if err != nil {
			return err
		}
	}

	p.InPathIndices = make([]uint32, len(params.InPathIndices))
	for i := 0; i < len(params.InPathIndices); i++ {
		p.InPathIndices[i] = params.InPathIndices[i]
	}

	p.InPathElements = make([][]big.Int, len(params.InPathElements))
	for i := 0; i < len(params.InPathElements); i++ {
		p.InPathElements[i] = make([]big.Int, len(params.InPathElements[i]))
		for j := 0; j < len(params.InPathElements[i]); j++ {
			err = fromHex(&p.InPathElements[i][j], params.InPathElements[i][j])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ps *ProvingSystem) UnsafeReadFrom(r io.Reader) (int64, error) {
	var totalRead int64 = 0
	var intBuf [4]byte

	read, err := io.ReadFull(r, intBuf[:])
	totalRead += int64(read)
	if err != nil {
		return totalRead, err
	}
	ps.TreeDepth = binary.BigEndian.Uint32(intBuf[:])

	read, err = io.ReadFull(r, intBuf[:])
	totalRead += int64(read)
	if err != nil {
		return totalRead, err
	}
	ps.NumOfUTXOs = binary.BigEndian.Uint32(intBuf[:])

	ps.ProvingKey = groth16.NewProvingKey(ecc.BN254)
	keyRead, err := ps.ProvingKey.UnsafeReadFrom(r)
	totalRead += keyRead
	if err != nil {
		return totalRead, err
	}

	ps.VerifyingKey = groth16.NewVerifyingKey(ecc.BN254)
	keyRead, err = ps.VerifyingKey.UnsafeReadFrom(r)
	totalRead += keyRead
	if err != nil {
		return totalRead, err
	}

	ps.ConstraintSystem = groth16.NewCS(ecc.BN254)
	keyRead, err = ps.ConstraintSystem.ReadFrom(r)
	totalRead += keyRead
	if err != nil {
		return totalRead, err
	}

	return totalRead, nil
}
func ReadSystemFromFile(path string) (ps *ProvingSystem, err error) {
	ps = new(ProvingSystem)
	file, err := os.Open(path)
	if err != nil {
		return
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	_, err = ps.UnsafeReadFrom(file)
	if err != nil {
		return
	}
	return
}

func (ps *ProvingSystem) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64 = 0
	var intBuf [4]byte

	binary.BigEndian.PutUint32(intBuf[:], ps.TreeDepth)
	written, err := w.Write(intBuf[:])
	totalWritten += int64(written)
	if err != nil {
		return totalWritten, err
	}

	binary.BigEndian.PutUint32(intBuf[:], ps.NumOfUTXOs)
	written, err = w.Write(intBuf[:])
	totalWritten += int64(written)
	if err != nil {
		return totalWritten, err
	}

	keyWritten, err := ps.ProvingKey.WriteTo(w)
	totalWritten += keyWritten
	if err != nil {
		return totalWritten, err
	}

	keyWritten, err = ps.VerifyingKey.WriteTo(w)
	totalWritten += keyWritten
	if err != nil {
		return totalWritten, err
	}

	keyWritten, err = ps.ConstraintSystem.WriteTo(w)
	totalWritten += keyWritten
	if err != nil {
		return totalWritten, err
	}

	return totalWritten, nil
}
