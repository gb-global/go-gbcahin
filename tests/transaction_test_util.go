package tests

import (
	"fmt"

	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/common/hexutil"
	"gbchain-org/go-gbchain/core"
	"gbchain-org/go-gbchain/core/types"
	"gbchain-org/go-gbchain/params"
	"gbchain-org/go-gbchain/rlp"
)

// TransactionTest checks RLP decoding and sender derivation of transactions.
type TransactionTest struct {
	RLP         hexutil.Bytes `json:"rlp"`
	Singularity ttFork
}

type ttFork struct {
	Sender common.UnprefixedAddress `json:"sender"`
	Hash   common.UnprefixedHash    `json:"hash"`
}

func (tt *TransactionTest) Run(config *params.ChainConfig) error {

	validateTx := func(rlpData hexutil.Bytes, signer types.Signer, singularity bool) (*common.Address, *common.Hash, error) {
		tx := new(types.Transaction)
		if err := rlp.DecodeBytes(rlpData, tx); err != nil {
			return nil, nil, err
		}
		sender, err := types.Sender(signer, tx)
		if err != nil {
			return nil, nil, err
		}
		// Intrinsic gas
		requiredGas, err := core.IntrinsicGas(tx.Data(), tx.To() == nil, singularity)
		if err != nil {
			return nil, nil, err
		}
		if requiredGas > tx.Gas() {
			return nil, nil, fmt.Errorf("insufficient gas ( %d < %d )", tx.Gas(), requiredGas)
		}
		h := tx.Hash()
		return &sender, &h, nil
	}

	for _, testcase := range []struct {
		name        string
		signer      types.Signer
		fork        ttFork
		singularity bool
	}{
		{"Singularity", types.NewEIP155Signer(config.ChainID), tt.Singularity, true},
	} {
		sender, txhash, err := validateTx(tt.RLP, testcase.signer, testcase.singularity)

		if testcase.fork.Sender == (common.UnprefixedAddress{}) {
			if err == nil {
				return fmt.Errorf("expected error, got none (address %v)[%v]", sender.String(), testcase.name)
			}
			continue
		}
		// Should resolve the right address
		if err != nil {
			return fmt.Errorf("got error, expected none: %v", err)
		}
		if sender == nil {
			return fmt.Errorf("sender was nil, should be %x", common.Address(testcase.fork.Sender))
		}
		if *sender != common.Address(testcase.fork.Sender) {
			return fmt.Errorf("sender mismatch: got %x, want %x", sender, testcase.fork.Sender)
		}
		if txhash == nil {
			return fmt.Errorf("txhash was nil, should be %x", common.Hash(testcase.fork.Hash))
		}
		if *txhash != common.Hash(testcase.fork.Hash) {
			return fmt.Errorf("hash mismatch: got %x, want %x", *txhash, testcase.fork.Hash)
		}
	}
	return nil
}
