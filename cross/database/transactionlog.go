package db

import (
	"math/big"
	"sync"

	"gbchain-org/go-gbchain/common"
	"gbchain-org/go-gbchain/cross/core"
	"gbchain-org/go-gbchain/ethdb"
	"gbchain-org/go-gbchain/log"
	"gbchain-org/go-gbchain/rlp"
	"gbchain-org/go-gbchain/trie"
)

var (
	FinishedRoot = []byte("_FINISHED_ROOT_")
)

type TransactionLogs struct {
	diskDB   ethdb.KeyValueStore
	trieDB   *trie.Database
	finished *trie.Trie

	lock sync.RWMutex
}

type TransactionLog struct {
	*TransactionLogs
	chainID *big.Int
}

func NewTransactionLogs(db ethdb.KeyValueStore) (*TransactionLogs, error) {
	database := trie.NewDatabase(db)
	finishedRoot, _ := db.Get(FinishedRoot)
	finished, err := trie.New(common.BytesToHash(finishedRoot), database)
	if err != nil {
		return nil, err
	}
	return &TransactionLogs{diskDB: db, trieDB: database, finished: finished}, nil
}

func (l *TransactionLogs) Get(chainID *big.Int) *TransactionLog {
	return &TransactionLog{
		TransactionLogs: l,
		chainID:         chainID,
	}
}

func (l *TransactionLogs) Close() {
	if err := l.diskDB.Close(); err != nil {
		log.Warn("transaction logs close failed", "error", err)
	}
}

func getKey(chainID *big.Int, hash common.Hash) []byte {
	return append(chainID.Bytes(), hash.Bytes()...)
}

func (l *TransactionLog) AddFinish(ctx *core.CrossTransactionWithSignatures) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	b, err := rlp.EncodeToBytes(ctx)
	if err != nil {
		return err
	}
	l.finished.Update(getKey(l.chainID, ctx.ID()), b)
	return nil
}

func (l *TransactionLog) GetFinish(hash common.Hash) (*core.CrossTransactionWithSignatures, bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	enc, err := l.finished.TryGet(getKey(l.chainID, hash))
	if err != nil {
		return nil, false
	}
	var ctx core.CrossTransactionWithSignatures
	if err := rlp.DecodeBytes(enc, &ctx); err != nil {
		return nil, false
	}
	return &ctx, true
}

func (l *TransactionLog) IsFinish(hash common.Hash) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	b, err := l.finished.TryGet(getKey(l.chainID, hash))
	return err == nil && len(b) > 0
}

func (l *TransactionLog) Commit() (common.Hash, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	root, err := l.finished.Commit(nil)
	if err != nil {
		return root, err
	}
	if err := l.trieDB.Commit(root, false); err != nil {
		return root, err
	}
	return root, l.diskDB.Put(FinishedRoot, root.Bytes())
}
