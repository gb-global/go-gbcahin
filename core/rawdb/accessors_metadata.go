package rawdb

import (
	"encoding/json"

	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/ethdb"
	"gbchain-org/go-gbchain/log"
	"gbchain-org/go-gbchain/params"
	"gbchain-org/go-gbchain/rlp"
)

// ReadDatabaseVersion retrieves the version number of the database.
func ReadDatabaseVersion(db ethdb.KeyValueReader) *uint64 {
	var version uint64

	enc, _ := db.Get(databaseVerisionKey)
	if len(enc) == 0 {
		return nil
	}
	if err := rlp.DecodeBytes(enc, &version); err != nil {
		return nil
	}

	return &version
}

// WriteDatabaseVersion stores the version number of the database
func WriteDatabaseVersion(db ethdb.KeyValueWriter, version uint64) {
	enc, err := rlp.EncodeToBytes(version)
	if err != nil {
		log.Crit("Failed to encode database version", "err", err)
	}
	if err = db.Put(databaseVerisionKey, enc); err != nil {
		log.Crit("Failed to store the database version", "err", err)
	}
}

// ReadChainConfig retrieves the consensus settings based on the given genesis hash.
func ReadChainConfig(db ethdb.KeyValueReader, hash common.Hash) *params.ChainConfig {
	data, _ := db.Get(configKey(hash))
	if len(data) == 0 {
		return nil
	}
	var config params.ChainConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Error("Invalid chain config JSON", "hash", hash, "err", err)
		return nil
	}
	return &config
}

// WriteChainConfig writes the chain config settings to the database.
func WriteChainConfig(db ethdb.KeyValueWriter, hash common.Hash, cfg *params.ChainConfig) {
	if cfg == nil {
		return
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		log.Crit("Failed to JSON encode chain config", "err", err)
	}
	if err := db.Put(configKey(hash), data); err != nil {
		log.Crit("Failed to store chain config", "err", err)
	}
}

// ReadPreimage retrieves a single preimage of the provided hash.
func ReadPreimage(db ethdb.KeyValueReader, hash common.Hash) []byte {
	data, _ := db.Get(preimageKey(hash))
	return data
}

// WritePreimages writes the provided set of preimages to the database.
func WritePreimages(db ethdb.KeyValueWriter, preimages map[common.Hash][]byte) {
	for hash, preimage := range preimages {
		if err := db.Put(preimageKey(hash), preimage); err != nil {
			log.Crit("Failed to store trie preimage", "err", err)
		}
	}
	preimageCounter.Inc(int64(len(preimages)))
	preimageHitCounter.Inc(int64(len(preimages)))
}
