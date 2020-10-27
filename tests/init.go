package tests

import (
	"fmt"
	"math/big"

	"gbchain-org/go-gbchain/params"
)

// Forks table defines supported forks and their chain config.
var Forks = map[string]*params.ChainConfig{
	"Singularity": {
		ChainID:          big.NewInt(1),
		SingularityBlock: big.NewInt(0),
	},
}

// UnsupportedForkError is returned when a test requests a fork that isn't implemented.
type UnsupportedForkError struct {
	Name string
}

func (e UnsupportedForkError) Error() string {
	return fmt.Sprintf("unsupported fork %q", e.Name)
}
