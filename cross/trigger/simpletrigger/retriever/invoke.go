package retriever

import (
	"context"
	"math/big"
	"time"

	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/common/math"
	"gbchain-org/go-gbchain/core"
	"gbchain-org/go-gbchain/core/state"
	"gbchain-org/go-gbchain/core/types"
	"gbchain-org/go-gbchain/core/vm"
	"gbchain-org/go-gbchain/params"

	"gbchain-org/go-gbchain/cross/trigger"
	"gbchain-org/go-gbchain/cross/trigger/simpletrigger"
)

// ChainInvoke invoke blockchain interfaces to get block and transaction states
type ChainInvoke struct {
	bc simpletrigger.BlockChain
}

func NewChainInvoke(chain simpletrigger.BlockChain) *ChainInvoke {
	return &ChainInvoke{bc: chain}
}

func (c ChainInvoke) CurrentBlockNumber() uint64 {
	if blk := c.bc.CurrentBlock(); blk != nil {
		return blk.NumberU64()
	}
	return 0
}

func (c ChainInvoke) GetTransactionNumberOnChain(tx trigger.Transaction) uint64 {
	if num := c.bc.GetBlockNumber(tx.BlockHash()); num != nil {
		return *num
	}
	//TODO return current for invisible block?
	return c.bc.CurrentBlock().NumberU64()
}

func (c ChainInvoke) GetConfirmedTransactionNumberOnChain(tx trigger.Transaction) uint64 {
	if num := c.bc.GetBlockNumber(tx.BlockHash()); num != nil {
		return *num + uint64(simpletrigger.DefaultConfirmDepth)
	}
	//TODO return current for invisible block?
	return c.bc.CurrentBlock().NumberU64()
}

func (c ChainInvoke) GetTransactionTimeOnChain(tx trigger.Transaction) uint64 {
	if header := c.bc.GetHeaderByHash(tx.BlockHash()); header != nil {
		return header.Time
	}
	return 0
}

func (c ChainInvoke) IsTransactionExpired(tx trigger.Transaction, expiredHeight uint64) bool {
	return c.bc.CurrentBlock().NumberU64()-c.GetTransactionNumberOnChain(tx) > expiredHeight
}

// EvmInvoke invoke evm interfaces to call contract function
type EvmInvoke struct {
	bc      core.ChainContext
	header  *types.Header
	stateDB *state.StateDB

	chainConfig *params.ChainConfig
	vmConfig    vm.Config
}

func NewEvmInvoke(bc core.ChainContext, header *types.Header, stateDB *state.StateDB, config *params.ChainConfig,
	vmCfg vm.Config) *EvmInvoke {
	return &EvmInvoke{bc: bc, header: header, stateDB: stateDB, chainConfig: config, vmConfig: vmCfg}
}

func (e EvmInvoke) CallContract(from common.Address, to *common.Address, function []byte, inputs ...[]byte) ([]byte, error) {
	var data []byte
	data = append(data, function...)
	for _, input := range inputs {
		data = append(data, input...)
	}

	checkMsg := types.NewMessage(from, to, 0, big.NewInt(0), math.MaxUint64/2,
		big.NewInt(params.GWei), data, false)
	var cancel context.CancelFunc
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

	// Make sure the context is cancelled when the call has completed
	// this makes sure resources are cleaned up.
	defer cancel()

	// Get a new instance of the EVM.
	// Create a new context to be used in the EVM environment
	vmContext := core.NewEVMContext(checkMsg, e.header, e.bc, nil)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	tempState := e.stateDB.Copy()
	tempState.SetBalance(checkMsg.From(), math.MaxBig256)
	evm := vm.NewEVM(vmContext, tempState, e.chainConfig, e.vmConfig)
	// Wait for the context to be done and cancel the evm. Even if the
	// EVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		evm.Cancel()
	}()

	// Setup the gas pool (also for unmetered requests)
	// and apply the messages
	gp := new(core.GasPool).AddGas(math.MaxUint64)
	res, _, _, err := core.ApplyMessage(evm, checkMsg, gp)
	if err != nil {
		return nil, err
	}
	return res, nil
}
