package crypto

import "errors"

type AlgUsage int

const (
	USAGEUNKNOWN AlgUsage = iota
	SIGNATURE
	ENCRYPTION
	KEM
	KEX
)

type Algorithm interface {
	GenerateKeys() (publicKey interface{}, privateKey interface{}, err error)
	Usage() []AlgUsage
	Name() string
	Version() string
}

type Algorithms []Algorithm

type AlgorithmsRepositoryItem struct {
	Name    string
	Version string
	Use     []AlgUsage
}

type AlgRepo struct {
	Items []AlgorithmsRepositoryItem
	index int
}

var indexes map[string]bool

func (ar *AlgRepo) Register(name, version string, usage AlgUsage) error {

	if checkAlgExists(name) {
		return errors.New("this algorithm is already registered")
	}

	algItem := AlgorithmsRepositoryItem{
		Name:    name,
		Version: version,
		Use:     make([]AlgUsage, 0),
	}
	algItem.Use = append(algItem.Use, usage)
	currIndex := ar.index + 1
	ar.Items[currIndex] = algItem
	indexes[name] = true
	return nil

}

func checkAlgExists(name string) bool {
	return indexes[name] == true
}
