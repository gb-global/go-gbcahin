package validator

import (
	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/consensus/istanbul"
)

func New(addr common.Address) istanbul.Validator {
	return &defaultValidator{
		address: addr,
	}
}

func NewSet(addrs []common.Address, policy istanbul.ProposerPolicy) istanbul.ValidatorSet {
	return newDefaultSet(addrs, policy)
}

func ExtractValidators(extraData []byte) []common.Address {
	// get the validator addresses
	addrs := make([]common.Address, (len(extraData) / common.AddressLength))
	for i := 0; i < len(addrs); i++ {
		copy(addrs[i][:], extraData[i*common.AddressLength:])
	}

	return addrs
}

// Check whether the extraData is presented in prescribed form
func ValidExtraData(extraData []byte) bool {
	return len(extraData)%common.AddressLength == 0
}
