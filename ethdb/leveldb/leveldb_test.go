package leveldb

import (
	"testing"

	"gbchain-org/go-gbchain/ethdb"
	"gbchain-org/go-gbchain/ethdb/dbtest"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

func TestLevelDB(t *testing.T) {
	t.Run("DatabaseSuite", func(t *testing.T) {
		dbtest.TestDatabaseSuite(t, func() ethdb.KeyValueStore {
			db, err := leveldb.Open(storage.NewMemStorage(), nil)
			if err != nil {
				t.Fatal(err)
			}
			return &Database{
				db: db,
			}
		})
	})
}
