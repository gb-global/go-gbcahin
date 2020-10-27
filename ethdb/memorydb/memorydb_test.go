package memorydb

import (
	"testing"

	"gbchain-org/go-gbchain/ethdb"
	"gbchain-org/go-gbchain/ethdb/dbtest"
)

func TestMemoryDB(t *testing.T) {
	t.Run("DatabaseSuite", func(t *testing.T) {
		dbtest.TestDatabaseSuite(t, func() ethdb.KeyValueStore {
			return New()
		})
	})
}
