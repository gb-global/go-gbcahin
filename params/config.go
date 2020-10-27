package params

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/crypto"
)

// Genesis hashes to enforce below configs on.
var (
	MainnetGenesisHash = common.HexToHash("0xda24a722f45ce6cb3aa7cb735a14b830869142be54ada9f8e2c27c877def648b")
	TestnetGenesisHash = common.HexToHash("0xfc26e8a96571fab7c359830a0651bcb442cdd63af5927fb1c74e4dddb7b759ae")
	FoundationAddress  = common.HexToAddress("0x843c5ce134ef3924625b03b800d43e591d502838")
)

// TrustedCheckpoints associates each known checkpoint with the genesis hash of
// the chain it belongs to.
var TrustedCheckpoints = map[common.Hash]*TrustedCheckpoint{
	MainnetGenesisHash: MainnetTrustedCheckpoint,
	TestnetGenesisHash: TestnetTrustedCheckpoint,
}

// CheckpointOracles associates each known checkpoint oracles with the genesis hash of
// the chain it belongs to.
var CheckpointOracles = map[common.Hash]*CheckpointOracleConfig{
	MainnetGenesisHash: MainnetCheckpointOracle,
	TestnetGenesisHash: TestnetCheckpointOracle,
}

var (
	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &ChainConfig{
		ChainID:          big.NewInt(1),
		SingularityBlock: big.NewInt(3966693),
		Scrypt:           new(ScryptConfig),
	}

	// MainnetTrustedCheckpoint contains the light client trusted checkpoint for the main network.
	MainnetTrustedCheckpoint = &TrustedCheckpoint{
		SectionIndex: 71,
		SectionHead:  common.HexToHash("0xe7189072f0e7114186efc591b41bcc9bbc13627e5adf73390663dc0c529c3fe8"),
		CHTRoot:      common.HexToHash("0x4297109e69fa6d8f5f79fb407cd62aa70a17632ac8173dec27954cf39d6fae53"),
		BloomRoot:    common.HexToHash("0xd819b965ebe9ad0cb8a572575f9b9b2823dc00897cb1a978f5fd8ef7e305e575"),
	}

	// MainnetCheckpointOracle contains a set of configs for the main network oracle.
	MainnetCheckpointOracle = &CheckpointOracleConfig{
		Address: common.HexToAddress("0x843c5ce134ef3924625b03b800d43e591d502838"), //合约地址
		Signers: []common.Address{
			common.HexToAddress("0x843c5ce134ef3924625b03b800d43e591d502838"), // sky  地址
		},
		Threshold: 3,
	}

	// TestnetChainConfig contains the chain parameters to run a node on the Ropsten test network.
	TestnetChainConfig = &ChainConfig{
		ChainID:          big.NewInt(3),
		SingularityBlock: big.NewInt(2690000),
		Scrypt:           new(ScryptConfig),
	}

	// TestnetTrustedCheckpoint contains the light client trusted checkpoint for the Ropsten test network.
	TestnetTrustedCheckpoint = &TrustedCheckpoint{
		SectionIndex: 62,
		SectionHead:  common.HexToHash("0x9c1b119ed1cf80123838b66a6feb380eea5907985d4a824a05f950e21650808a"),
		CHTRoot:      common.HexToHash("0x52fe75523f8e2864bed4a988e2330df3bdb5278477315c987a5b82d2dd03f8e9"),
		BloomRoot:    common.HexToHash("0x6969fdea930386b669be52b82e5eb6a827b79dd22daf815c8f1dd16517c8c6f8"),
	}

	// TestnetCheckpointOracle contains a set of configs for the Ropsten test network oracle.
	TestnetCheckpointOracle = &CheckpointOracleConfig{
		Address: common.HexToAddress("0x0A37FE53Fa04Db73028440084c01fA66Ea32123d"),
		Signers: []common.Address{
			common.HexToAddress("0x843c5ce134ef3924625b03b800d43e591d502838"), // sky  地址
		},
		Threshold: 2,
	}

	// AllCliqueProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Clique consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.

	AllCliqueProtocolChanges = &ChainConfig{big.NewInt(1337), big.NewInt(0), nil, nil, &CliqueConfig{Period: 0, Epoch: 30000}, nil, nil, false, nil}

	AllDPoSProtocolChanges = &ChainConfig{big.NewInt(1337), big.NewInt(0), nil, nil, nil, nil, &DPoSConfig{Period: 3, Epoch: 30000, MaxSignerCount: 21, MinVoterBalance: new(big.Int).Mul(big.NewInt(10000), big.NewInt(1000000000000000000))}, false, nil}

	// AllScryptProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Scrypt consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.

	AllScryptProtocolChanges = &ChainConfig{big.NewInt(1337), big.NewInt(0), nil, nil, nil, new(ScryptConfig), nil, false, nil}

	TestChainConfig = &ChainConfig{big.NewInt(1), big.NewInt(0), nil, new(EthashConfig), nil, nil, nil, false, nil}

	TestRules = TestChainConfig.Rules(new(big.Int))

	RaftChainConfig = &ChainConfig{big.NewInt(1337), big.NewInt(0), nil, nil, nil, nil, nil, true, nil}
)

// TrustedCheckpoint represents a set of post-processed trie roots (CHT and
// BloomTrie) associated with the appropriate section index and head hash. It is
// used to start light syncing from this checkpoint and avoid downloading the
// entire header chain while still being able to securely access old headers/logs.
type TrustedCheckpoint struct {
	SectionIndex uint64      `json:"sectionIndex"`
	SectionHead  common.Hash `json:"sectionHead"`
	CHTRoot      common.Hash `json:"chtRoot"`
	BloomRoot    common.Hash `json:"bloomRoot"`
}

// HashEqual returns an indicator comparing the itself hash with given one.
func (c *TrustedCheckpoint) HashEqual(hash common.Hash) bool {
	if c.Empty() {
		return hash == common.Hash{}
	}
	return c.Hash() == hash
}

// Hash returns the hash of checkpoint's four key fields(index, sectionHead, chtRoot and bloomTrieRoot).
func (c *TrustedCheckpoint) Hash() common.Hash {
	buf := make([]byte, 8+3*common.HashLength)
	binary.BigEndian.PutUint64(buf, c.SectionIndex)
	copy(buf[8:], c.SectionHead.Bytes())
	copy(buf[8+common.HashLength:], c.CHTRoot.Bytes())
	copy(buf[8+2*common.HashLength:], c.BloomRoot.Bytes())
	return crypto.Keccak256Hash(buf)
}

// Empty returns an indicator whether the checkpoint is regarded as empty.
func (c *TrustedCheckpoint) Empty() bool {
	return c.SectionHead == (common.Hash{}) || c.CHTRoot == (common.Hash{}) || c.BloomRoot == (common.Hash{})
}

// CheckpointOracleConfig represents a set of checkpoint contract(which acts as an oracle)
// config which used for light client checkpoint syncing.
type CheckpointOracleConfig struct {
	Address   common.Address   `json:"address"`
	Signers   []common.Address `json:"signers"`
	Threshold uint64           `json:"threshold"`
}

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	ChainID *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection

	SingularityBlock *big.Int `json:"singularityBlock,omitempty"` // Singularity switch block (nil = no fork, 0 = already on singularity)
	EWASMBlock       *big.Int `json:"ewasmBlock,omitempty"`       // EWASM switch block (nil = no fork, 0 = already activated)

	// Various consensus engines
	Ethash   *EthashConfig   `json:"ethash,omitempty"`
	Clique   *CliqueConfig   `json:"clique,omitempty"`
	Scrypt   *ScryptConfig   `json:"scrypt,omitempty"`
	DPoS     *DPoSConfig     `json:"dpos,omitempty"`
	Raft     bool            `json:"raft,omitempty"`
	Istanbul *IstanbulConfig `json:"istanbul,omitempty"`
}

// EthashConfig is the consensus engine configs for proof-of-work based sealing.
type EthashConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c *EthashConfig) String() string {
	return "ethash"
}

// CliqueConfig is the consensus engine configs for proof-of-authority based sealing.
type CliqueConfig struct {
	Period uint64 `json:"period"` // Number of seconds between blocks to enforce
	Epoch  uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
}

// String implements the stringer interface, returning the consensus engine details.
func (c *CliqueConfig) String() string {
	return "clique"
}

// IstanbulConfig is the consensus engine configs for Istanbul based sealing.
type IstanbulConfig struct {
	Epoch          uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
	ProposerPolicy uint64 `json:"policy"` // The policy for proposer selection
}

type RaftConfig struct {
	BlockTime uint64 `json:"blockTime"`
}

// String implements the stringer interface, returning the consensus engine details.
func (c *IstanbulConfig) String() string {
	return "istanbul"
}

type GenesisAccount struct {
	Balance string `json:"balance"`
}

// DPoSLightConfig is the config for light node of dpos
type DPoSLightConfig struct {
	Alloc map[common.UnprefixedAddress]GenesisAccount `json:"alloc"`
}

// DPoSConfig is the consensus engine configs for delegated-proof-of-stake based sealing.
type DPoSConfig struct {
	Period           uint64                     `json:"period"`           // Number of seconds between blocks to enforce
	Epoch            uint64                     `json:"epoch"`            // Epoch length to reset votes and checkpoint
	MaxSignerCount   uint64                     `json:"maxSignersCount"`  // Max count of signers
	MinVoterBalance  *big.Int                   `json:"minVoterBalance"`  // Min voter balance to valid this vote
	GenesisTimestamp uint64                     `json:"genesisTimestamp"` // The LoopStartTime of first Block
	SelfVoteSigners  []common.UnprefixedAddress `json:"signers"`          // Signers vote by themselves to seal the block, make sure the signer accounts are pre-funded
	PBFTEnable       bool                       `json:"pbft"`             //
	VoterReward      bool                       `json:"voterReward"`
	LightConfig      *DPoSLightConfig           `json:"lightConfig,omitempty"`
}

// String implements the stringer interface, returning the consensus engine details.
func (a *DPoSConfig) String() string {
	return "dpos"
}

// ScryptConfig is the consensus engine configs for proof-of-work based sealing.
type ScryptConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c *ScryptConfig) String() string {
	return "scrypt"
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	var engine interface{}
	switch {
	case c.Ethash != nil:
		engine = c.Ethash
	case c.Clique != nil:
		engine = c.Clique
	case c.Scrypt != nil:
		engine = c.Scrypt
	case c.DPoS != nil:
		engine = c.DPoS
	case c.Istanbul != nil:
		engine = c.Istanbul
	case c.Raft:
		engine = "raft"
	default:
		engine = "unknown"
	}
	return fmt.Sprintf("{ChainID: %v Singularity: %v, Engine: %v}",
		c.ChainID,
		c.SingularityBlock,
		engine,
	)
}

// IsSingularity returns whether num is either equal to the Istanbul fork block or greater.
func (c *ChainConfig) IsSingularity(num *big.Int) bool {
	return isForked(c.SingularityBlock, num)
}

// IsEWASM returns whether num represents a block number after the EWASM fork
func (c *ChainConfig) IsEWASM(num *big.Int) bool {
	return isForked(c.EWASMBlock, num)
}

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height uint64) *ConfigCompatError {
	bhead := new(big.Int).SetUint64(height)

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bhead)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bhead.SetUint64(err.RewindTo)
	}
	return lasterr
}

// CheckConfigForkOrder checks that we don't "skip" any forks, geth isn't pluggable enough
// to guarantee that forks can be implemented in a different order than on official networks
func (c *ChainConfig) CheckConfigForkOrder() error {
	type fork struct {
		name  string
		block *big.Int
	}
	var lastFork fork
	for _, cur := range []fork{
		{"singularityBlock", c.SingularityBlock},
	} {
		if lastFork.name != "" {
			// Next one must be higher number
			if lastFork.block == nil && cur.block != nil {
				return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at %v",
					lastFork.name, cur.name, cur.block)
			}
			if lastFork.block != nil && cur.block != nil {
				if lastFork.block.Cmp(cur.block) > 0 {
					return fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
						lastFork.name, lastFork.block, cur.name, cur.block)
				}
			}
		}
		lastFork = cur
	}
	return nil
}

func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, head *big.Int) *ConfigCompatError {
	if isForkIncompatible(c.SingularityBlock, newcfg.SingularityBlock, head) {
		return newCompatError("SingularityBlock fork block", c.SingularityBlock, newcfg.SingularityBlock)
	}
	if isForkIncompatible(c.EWASMBlock, newcfg.EWASMBlock, head) {
		return newCompatError("ewasm fork block", c.EWASMBlock, newcfg.EWASMBlock)
	}
	return nil
}

// isForkIncompatible returns true if a fork scheduled at s1 cannot be rescheduled to
// block s2 because head is already past the fork.
func isForkIncompatible(s1, s2, head *big.Int) bool {
	return (isForked(s1, head) || isForked(s2, head)) && !configNumEqual(s1, s2)
}

// isForked returns whether a fork scheduled at block s is active at the given head block.
func isForked(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) <= 0
}

func configNumEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

// Rules wraps ChainConfig and is merely syntactic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID       *big.Int
	IsSingularity bool
}

// Rules ensures c's ChainID is not nil.
func (c *ChainConfig) Rules(num *big.Int) Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}
	return Rules{
		ChainID:       new(big.Int).Set(chainID),
		IsSingularity: c.IsSingularity(num),
	}
}
