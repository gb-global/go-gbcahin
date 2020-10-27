package scrypt

import (
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"gbchain-org/go-gbchain/common/math"
	"gbchain-org/go-gbchain/core/types"
	"gbchain-org/go-gbchain/params"
)

type diffTest struct {
	ParentTimestamp    uint64
	ParentDifficulty   *big.Int
	CurrentTimestamp   uint64
	CurrentBlocknumber *big.Int
	CurrentDifficulty  *big.Int
}

func (d *diffTest) UnmarshalJSON(b []byte) (err error) {
	var ext struct {
		ParentTimestamp    string
		ParentDifficulty   string
		CurrentTimestamp   string
		CurrentBlocknumber string
		CurrentDifficulty  string
	}
	if err := json.Unmarshal(b, &ext); err != nil {
		return err
	}

	d.ParentTimestamp = math.MustParseUint64(ext.ParentTimestamp)
	d.ParentDifficulty = math.MustParseBig256(ext.ParentDifficulty)
	d.CurrentTimestamp = math.MustParseUint64(ext.CurrentTimestamp)
	d.CurrentBlocknumber = math.MustParseBig256(ext.CurrentBlocknumber)
	d.CurrentDifficulty = math.MustParseBig256(ext.CurrentDifficulty)

	return nil
}

var testData string = `
{
    "preExpDiffIncrease" : {
        "parentTimestamp" : "42",
        "parentDifficulty" : "1000000",
        "currentTimestamp" : "43",
        "currentBlockNumber" : "42",
        "currentDifficulty" : "1001920",
	"parentUncles" : "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"
    },
    "preExpDiffDecrease" : {
        "parentTimestamp" : "42",
        "parentDifficulty" : "1000000",
        "currentTimestamp" : "60",
        "currentBlockNumber" : "42",
        "currentDifficulty" : "1000569",
	"parentUncles" : "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"
    }
}`

func TestCalcDifficulty(t *testing.T) {
	tests := make(map[string]diffTest)
	strRead := strings.NewReader(testData)
	err := json.NewDecoder(strRead).Decode(&tests)
	if err != nil {
		t.Fatal(err)
	}

	config := &params.ChainConfig{SingularityBlock: big.NewInt(1150000)}
	for name, test := range tests {
		number := new(big.Int).Sub(test.CurrentBlocknumber, big.NewInt(1))
		diff := CalcDifficulty(config, test.CurrentTimestamp, &types.Header{
			Number:     number,
			Time:       test.ParentTimestamp,
			Difficulty: test.ParentDifficulty,
		})
		if diff.Cmp(test.CurrentDifficulty) != 0 {
			t.Error(name, "failed. Expected", test.CurrentDifficulty, "and calculated", diff)
		}
	}
}
