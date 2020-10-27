package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"gbchain-org/go-gbchain/core"
)

// Tests the go-gbchain to Aleth chainspec conversion for the Stureby testnet.
func TestAlethSturebyConverter(t *testing.T) {
	blob, err := ioutil.ReadFile("testdata/stureby_geth.json")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	var genesis core.Genesis
	if err := json.Unmarshal(blob, &genesis); err != nil {
		t.Fatalf("failed parsing genesis: %v", err)
	}
	spec, err := newAlethGenesisSpec("stureby", &genesis)
	if err != nil {
		t.Fatalf("failed creating chainspec: %v", err)
	}

	expBlob, err := ioutil.ReadFile("testdata/stureby_aleth.json")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	expspec := &alethGenesisSpec{}
	if err := json.Unmarshal(expBlob, expspec); err != nil {
		t.Fatalf("failed parsing genesis: %v", err)
	}
	if !reflect.DeepEqual(expspec, spec) {
		t.Errorf("chainspec mismatch")
		c := spew.ConfigState{
			DisablePointerAddresses: true,
			SortKeys:                true,
		}
		exp := strings.Split(c.Sdump(expspec), "\n")
		got := strings.Split(c.Sdump(spec), "\n")
		for i := 0; i < len(exp) && i < len(got); i++ {

			if exp[i] != got[i] {
				t.Logf("got: %v\nexp: %v,%d\n", exp[i], got[i], i)
			}
		}
	}
}

// Tests the go-gbchain to Parity chainspec conversion for the Stureby testnet.
func TestParitySturebyConverter(t *testing.T) {
	blob, err := ioutil.ReadFile("testdata/stureby_geth.json")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	var genesis core.Genesis
	if err := json.Unmarshal(blob, &genesis); err != nil {
		t.Fatalf("failed parsing genesis: %v", err)
	}
	spec, err := newParityChainSpec("stureby", &genesis, []string{})
	if err != nil {
		t.Fatalf("failed creating chainspec: %v", err)
	}
	enc, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		t.Fatalf("failed encoding chainspec: %v", err)
	}
	expBlob, err := ioutil.ReadFile("testdata/stureby_parity.json")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	if !bytes.Equal(expBlob, enc) {
		t.Errorf("chainspec mismatch")
		exp := strings.Split(string(expBlob), "\n")
		got := strings.Split(string(enc), "\n")
		for i := 0; i < len(exp) && i < len(got); i++ {
			if exp[i] != got[i] {
				t.Fatalf("got: %v\n exp: %v\n", strings.Trim(exp[i], " "), strings.Trim(got[i], " "))
			}
		}
	}
}
