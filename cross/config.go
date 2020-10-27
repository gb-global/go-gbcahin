package cross

import (
	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/cross/backend/synchronise"
)

const (
	LogDir   = "crosslog"
	TxLogDir = "crosstxlog"
	DataDir  = "crossdata"
)

type Config struct {
	MainContract common.Address       `json:"mainContract"`
	SubContract  common.Address       `json:"subContract"`
	Signer       common.Address       `json:"signer"`
	Anchors      []common.Address     `json:"anchors"`
	SyncMode     synchronise.SyncMode `json:"syncMode"`
}

var DefaultConfig = Config{
	SyncMode: synchronise.ALL,
}

func (config *Config) Sanitize() Config {
	cfg := Config{
		MainContract: config.MainContract,
		SubContract:  config.SubContract,
		Signer:       config.Signer,
	}
	set := make(map[common.Address]struct{})
	for _, anchor := range config.Anchors {
		if _, ok := set[anchor]; !ok {
			cfg.Anchors = append(cfg.Anchors, anchor)
			set[anchor] = struct{}{}
		}
	}
	return cfg
}
