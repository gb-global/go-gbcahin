package tests

import (
	"math/big"
	"testing"

	//"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/params"
)

var (
	mainnetChainConfig = params.ChainConfig{
		ChainID:          big.NewInt(1),
		SingularityBlock: big.NewInt(3000000),
	}
)

func TestDifficulty(t *testing.T) {
	t.Parallel()

	dt := new(testMatcher)
	// Not difficulty-tests
	dt.skipLoad("hexencodetest.*")
	dt.skipLoad("crypto.*")
	dt.skipLoad("blockgenesistest\\.json")
	dt.skipLoad("genesishashestest\\.json")
	dt.skipLoad("keyaddrtest\\.json")
	dt.skipLoad("txtest\\.json")

	// files are 2 years old, contains strange values
	dt.skipLoad("difficultyCustomHomestead\\.json")
	dt.skipLoad("difficultyMorden\\.json")
	dt.skipLoad("difficultyOlimpic\\.json")

	dt.config("Ropsten", *params.TestnetChainConfig)
	dt.config("Morden", *params.TestnetChainConfig)
	dt.config("Frontier", params.ChainConfig{})
	dt.config("Singularity", params.ChainConfig{
		SingularityBlock: big.NewInt(0),
	})

	dt.config("Frontier", *params.TestnetChainConfig)
	dt.config("MainNetwork", mainnetChainConfig)
	dt.config("CustomMainNetwork", mainnetChainConfig)
	dt.config("difficulty.json", mainnetChainConfig)

	dt.walk(t, difficultyTestDir, func(t *testing.T, name string, test *DifficultyTest) {
		cfg := dt.findConfig(name)
		if test.ParentDifficulty.Cmp(params.MinimumDifficulty) < 0 {
			t.Skip("difficulty below minimum")
			return
		}
		if err := dt.checkFailure(t, name, test.Run(cfg)); err != nil {
			t.Error(err)
		}
	})
}
