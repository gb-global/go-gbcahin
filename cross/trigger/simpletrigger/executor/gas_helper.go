package executor

import (
	"context"
	"math/big"
	"time"

	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/common/hexutil"
	"gbchain-org/go-gbchain/common/math"
	"gbchain-org/go-gbchain/core"
	"gbchain-org/go-gbchain/core/types"
	"gbchain-org/go-gbchain/core/vm"
	"gbchain-org/go-gbchain/log"
	"gbchain-org/go-gbchain/params"
	"gbchain-org/go-gbchain/rpc"

	"gbchain-org/go-gbchain/cross/trigger/simpletrigger"
)

var (
	defaultGasPrice = big.NewInt(params.GWei)
	MaxGasPrice     = big.NewInt(500 * params.GWei)
)

type CallArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      hexutil.Uint64  `json:"gas"`
	GasPrice hexutil.Big     `json:"gasPrice"`
	Value    hexutil.Big     `json:"value"`
	Data     hexutil.Bytes   `json:"data"`
}

type GasHelper struct {
	blockchain *core.BlockChain
	chain      simpletrigger.GBChain
}

func NewGasHelper(blockchain *core.BlockChain, chain simpletrigger.GBChain) *GasHelper {
	return &GasHelper{blockchain: blockchain, chain: chain}
}

func (h *GasHelper) GetBalance(addr common.Address) (*big.Int, error) {
	state, err := h.blockchain.State()
	if err != nil {
		return nil, err
	}
	return state.GetBalance(addr), nil
}

func (this *GasHelper) doCall(ctx context.Context, args CallArgs, blockNr rpc.BlockNumber, vmCfg vm.Config,
	timeout time.Duration) ([]byte, uint64, bool, error) {

	defer func(start time.Time) {
		log.Trace("Executing EVM call finished", "runtime", time.Since(start))
	}(time.Now())

	state, header, err := this.chain.StateAndHeaderByNumber(ctx, blockNr)
	if state == nil || err != nil {
		return nil, 0, false, err
	}
	// Set sender address or use a default if none specified
	addr := args.From
	// Set default gas & gas price if none were set
	gas, gasPrice := uint64(args.Gas), args.GasPrice.ToInt()
	if gas == 0 {
		gas = math.MaxUint64 / 2
	}
	if gasPrice.Sign() == 0 {
		gasPrice.Set(defaultGasPrice)
	}

	// Create new call message
	msg := types.NewMessage(addr, args.To, 0, args.Value.ToInt(), gas, gasPrice, args.Data, false)

	// Setup context so it may be cancelled the call has completed
	// or, in case of unmetered gas, setup a context with a timeout.
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	// Make sure the context is cancelled when the call has completed
	// this makes sure resources are cleaned up.
	defer cancel()

	// Get a new instance of the EVM.
	evm, vmError, err := this.chain.GetEVM(ctx, msg, state, header, vmCfg)
	if err != nil {
		return nil, 0, false, err
	}
	// Wait for the context to be done and cancel the evm. Even if the
	// EVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		evm.Cancel()
	}()
	// Setup the gas pool (also for unmetered requests)
	// and apply the message.
	gp := new(core.GasPool).AddGas(math.MaxUint64)
	res, gas, failed, err := core.ApplyMessage(evm, msg, gp)
	if err := vmError(); err != nil {
		return nil, 0, false, err
	}
	return res, gas, failed, err
}

func (this *GasHelper) checkExec(ctx context.Context, args CallArgs) (bool, error) {
	_, _, failed, err := this.doCall(ctx, args, rpc.LatestBlockNumber, vm.Config{}, 0)
	if err != nil || failed {
		return false, err
	}
	return true, nil
}
