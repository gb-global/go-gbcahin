package core

import (
	"fmt"
	"math/big"

	"gbchain-org/go-gbchain/accounts/abi"
	"gbchain-org/go-gbchain/common"

	"github.com/syndtr/goleveldb/leveldb/errors"
)

var (
	ErrInvalidRecept    = errors.New("invalid recept transaction")
	ErrChainIdMissMatch = fmt.Errorf("[%w]: recept chainId miss match", ErrInvalidRecept)
	ErrToMissMatch      = fmt.Errorf("[%w]: recept to address miss match", ErrInvalidRecept)
	ErrFromMissMatch    = fmt.Errorf("[%w]: recept from address miss match", ErrInvalidRecept)
)

type ReceptTransaction struct {
	CTxId         common.Hash    `json:"ctxId" gencodec:"required"`         //cross_transaction ID
	TxHash        common.Hash    `json:"txHash" gencodec:"required"`        //taker txHash
	From          common.Address `json:"from" gencodec:"required"`          //Token seller
	To            common.Address `json:"to" gencodec:"required"`            //Token buyer
	DestinationId *big.Int       `json:"destinationId" gencodec:"required"` //Message destination networkId
	ChainId       *big.Int       `json:"chainId" gencodec:"required"`
}

func NewReceptTransaction(id, txHash common.Hash, from, to common.Address, remoteChainId, chainId *big.Int) *ReceptTransaction {
	return &ReceptTransaction{
		CTxId:         id,
		TxHash:        txHash,
		From:          from,
		To:            to,
		DestinationId: remoteChainId,
		ChainId:       chainId,
	}
}

func (rtx ReceptTransaction) Check(maker *CrossTransactionWithSignatures) error {
	if maker == nil {
		return ErrInvalidRecept
	}
	if maker.DestinationId().Cmp(rtx.ChainId) != 0 {
		return ErrChainIdMissMatch
	}
	if maker.Data.From != rtx.From {
		return ErrFromMissMatch
	}
	if maker.Data.To != (common.Address{}) && maker.Data.To != rtx.To {
		return ErrToMissMatch
	}
	return nil
}

type Recept struct {
	TxId   common.Hash
	TxHash common.Hash
	From   common.Address
	To     common.Address
	//Input  []byte //TODO delete
}

func (rtx *ReceptTransaction) ConstructData(crossContract abi.ABI) ([]byte, error) {
	rep := Recept{
		TxId:   rtx.CTxId,
		TxHash: rtx.TxHash,
		From:   rtx.From,
		To:     rtx.To,
	}
	return crossContract.Pack("makerFinish", rep, rtx.ChainId)
}
