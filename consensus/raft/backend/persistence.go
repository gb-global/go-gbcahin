package backend

import (
	"encoding/binary"

	"gbchain-org/go-gbchain/consensus/raft"
	"gbchain-org/go-gbchain/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	noFsync = &opt.WriteOptions{
		NoWriteMerge: false,
		Sync:         false,
	}
)

func openRaftDb(path string) (db *leveldb.DB, err error) {
	// Open the db and recover any potential corruptions
	db, err = leveldb.OpenFile(path, &opt.Options{
		OpenFilesCacheCapacity: -1, // -1 means 0??
		BlockCacheCapacity:     -1,
	})
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(path, nil)
	}
	return
}

func (pm *ProtocolManager) loadAppliedIndex() uint64 {
	dat, err := pm.raftDb.Get(raft.AppliedDbKey, nil)
	var lastAppliedIndex uint64
	if err == errors.ErrNotFound {
		lastAppliedIndex = 0
	} else if err != nil {
		raft.Fatalf("loadAppliedIndex error: %s", err)
	} else {
		lastAppliedIndex = binary.LittleEndian.Uint64(dat)
	}

	pm.mu.Lock()
	pm.appliedIndex = lastAppliedIndex
	pm.mu.Unlock()

	log.Info("loaded the latest applied index", "lastAppliedIndex", lastAppliedIndex)

	return lastAppliedIndex
}

func (pm *ProtocolManager) writeAppliedIndex(index uint64) {
	log.Info("persisted the latest applied index", "index", index)
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, index)
	pm.raftDb.Put(raft.AppliedDbKey, buf, noFsync)
}
