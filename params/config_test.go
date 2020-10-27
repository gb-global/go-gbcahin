package params

import (
	"math/big"
	"reflect"
	"testing"
)

func TestCheckCompatible(t *testing.T) {
	type test struct {
		stored, new *ChainConfig
		head        uint64
		wantErr     *ConfigCompatError
	}
	tests := []test{
		{stored: AllScryptProtocolChanges, new: AllScryptProtocolChanges, head: 0, wantErr: nil},
		{stored: AllScryptProtocolChanges, new: AllScryptProtocolChanges, head: 100, wantErr: nil},
		{
			stored:  &ChainConfig{SingularityBlock: big.NewInt(10)},
			new:     &ChainConfig{SingularityBlock: big.NewInt(20)},
			head:    9,
			wantErr: nil,
		},
		{
			stored: AllScryptProtocolChanges,
			new:    &ChainConfig{SingularityBlock: nil},
			head:   3,
			wantErr: &ConfigCompatError{
				What:         "SingularityBlock fork block",
				StoredConfig: big.NewInt(0),
				NewConfig:    nil,
				RewindTo:     0,
			},
		},
		{
			stored: AllScryptProtocolChanges,
			new:    &ChainConfig{SingularityBlock: big.NewInt(1)},
			head:   3,
			wantErr: &ConfigCompatError{
				What:         "SingularityBlock fork block",
				StoredConfig: big.NewInt(0),
				NewConfig:    big.NewInt(1),
				RewindTo:     0,
			},
		},
	}

	for _, test := range tests {
		err := test.stored.CheckCompatible(test.new, test.head)
		if !reflect.DeepEqual(err, test.wantErr) {
			t.Errorf("error mismatch:\nstored: %v\nnew: %v\nhead: %v\nerr: %v\nwant: %v", test.stored, test.new, test.head, err, test.wantErr)
		}
	}
}
