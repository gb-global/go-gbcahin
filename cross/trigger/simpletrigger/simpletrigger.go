package simpletrigger

import (
	"context"
	"math/big"

	"gbchain-org/go-gbchain/accounts"
	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/core"
	"gbchain-org/go-gbchain/core/state"
	"gbchain-org/go-gbchain/core/types"
	"gbchain-org/go-gbchain/core/vm"
	"gbchain-org/go-gbchain/eth/gasprice"
	"gbchain-org/go-gbchain/params"
	"gbchain-org/go-gbchain/rpc"
)

var DefaultConfirmDepth = 12

type ProtocolManager interface {
	NetworkId() uint64
	GetNonce(address common.Address) uint64
	AddLocals([]*types.Transaction)
	Pending() (map[common.Address]types.Transactions, error)
	CanAcceptTxs() bool
}

type BlockChain interface {
	core.ChainContext
	GetBlockNumber(hash common.Hash) *uint64
	GetHeaderByHash(hash common.Hash) *types.Header
	CurrentBlock() *types.Block
	StateAt(root common.Hash) (*state.StateDB, error)
}

type GasPriceOracle interface {
	SuggestPrice(ctx context.Context) (*big.Int, error)
}

type GBChain interface {
	BlockChain() *core.BlockChain
	ChainConfig() *params.ChainConfig
	GasOracle() *gasprice.Oracle
	ProtocolManager() ProtocolManager
	AccountManager() *accounts.Manager
	RegisterAPIs([]rpc.API)
	GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error)
	StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error)
}

type CreditProtocolChain struct {
	GBChain
}

func NewSimpleProtocolChain(sc GBChain) *CreditProtocolChain {
	return &CreditProtocolChain{sc}
}

func (sc *CreditProtocolChain) ChainID() *big.Int {
	return sc.ChainConfig().ChainID
}

func (sc *CreditProtocolChain) GenesisHash() common.Hash {
	return sc.BlockChain().Genesis().Hash()
}
