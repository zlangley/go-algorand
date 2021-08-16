package layer2

import (
	"database/sql"
	"testing"

	"github.com/algorand/go-algorand/crypto"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	store, err := NewStableStore(true)
	require.NoError(t, err)

	cid := ContractID(crypto.Hash([]byte("test")))
	ret, err := store.Select(cid)
	require.NoError(t, err)
	require.Empty(t, ret)

	store.Write(cid, []byte("bbbk"), []byte("bbbv"))
	val, err := store.Get(cid, []byte("bbbk"))
	require.NoError(t, err)
	require.Equal(t, val, []byte("bbbv"))

	store.Write(cid, []byte("ccck"), []byte("cccv"))
	val, err = store.Get(cid, []byte("ccck"))
	require.NoError(t, err)
	require.Equal(t, val, []byte("cccv"))

	require.NoError(t, store.Write(cid, []byte("dddk"), []byte("dddv")))
	require.NoError(t, store.Write(cid, []byte("dddk"), nil))
	val, err = store.Get(cid, []byte("dddk"))
	require.Equal(t, sql.ErrNoRows, err)
	require.Nil(t, val)

	store.Write(cid, []byte("aaak"), []byte("aaav"))
	val, err = store.Get(cid, []byte("aaak"))
	require.NoError(t, err)
	require.Equal(t, val, []byte("aaav"))

	cid2 := ContractID(crypto.Hash([]byte("test2")))
	store.Write(cid2, []byte("foo"), []byte("bar2"))

	ret, err = store.Select(cid)
	require.NoError(t, err)
	require.Equal(t, 3, len(ret))
	require.Equal(t, []byte("aaak"), ret[0].Key)
	require.Equal(t, []byte("aaav"), ret[0].Value)
	require.Equal(t, []byte("bbbk"), ret[1].Key)
	require.Equal(t, []byte("bbbv"), ret[1].Value)
	require.Equal(t, []byte("ccck"), ret[2].Key)
	require.Equal(t, []byte("cccv"), ret[2].Value)
}

func TestCache(t *testing.T) {
	store, err := NewStableStore(true)
	require.NoError(t, err)

	cache, err := store.Speculation()
	require.NoError(t, err)

	cid := ContractID(crypto.Hash([]byte("test")))

	store.Write(cid, []byte("store-only-key"), []byte("store-only-value"))
	store.Write(cid, []byte("both-key"), []byte("both-store-value"))

	commit1, err := cache.Commitment(cid)
	require.NoError(t, err)

	// Test absence of key.
	val, err := cache.Get(cid, []byte("cache-only-key"))
	require.NoError(t, err)
	require.Nil(t, val)

	// Test read-your-writes.
	cache.Write(cid, []byte("cache-only-key"), []byte("cache-only-value"), 0)
	val, err = cache.Get(cid, []byte("cache-only-key"))
	require.NoError(t, err)
	require.Equal(t, []byte("cache-only-value"), val)

	// Test commitment changed.
	commit2, err := cache.Commitment(cid)
	require.NoError(t, err)
	require.NotEqual(t, commit1, commit2)

	// Test delegation to underlying store for missing cache key.
	val, err = cache.Get(cid, []byte("both-key"))
	require.NoError(t, err)
	require.Equal(t, []byte("both-store-value"), val)

	// Test overriding underlying store key.
	cache.Write(cid, []byte("both-key"), []byte("both-cache-value"), 1)
	val, err = cache.Get(cid, []byte("both-key"))
	require.NoError(t, err)
	require.Equal(t, []byte("both-cache-value"), val)

	// Test delete store key in cache.
	cache.Write(cid, []byte("both-key"), nil, 1)
	val, err = cache.Get(cid, []byte("both-key"))
	require.NoError(t, err)
	require.Nil(t, val)

	// Test value remains in persistent store.
	val, err = store.Get(cid, []byte("both-key"))
	require.NoError(t, err)
	require.Equal(t, []byte("both-store-value"), val)

	// Test commitment changed again.
	commit3, err := cache.Commitment(cid)
	require.NoError(t, err)
	require.NotEqual(t, commit2, commit3)

	// Set batch.
	groupID0 := crypto.Hash([]byte{1, 2, 3})
	cache.SetBatchIndexGroup(0, groupID0)
	groupID1 := crypto.Hash([]byte{4, 5, 6})
	cache.SetBatchIndexGroup(1, groupID1)

	// Persist write set for index 0.
	err = cache.PersistGroupState(groupID0)
	require.NoError(t, err)
	val, err = store.Get(cid, []byte("both-key"))
	require.NoError(t, err)
	require.Equal(t, []byte("both-store-value"), val)
	val, err = store.Get(cid, []byte("cache-only-key"))
	require.NoError(t, err)
	require.Equal(t, []byte("cache-only-value"), val)

	// Persist write set for index 1.
	err = cache.PersistGroupState(groupID1)
	require.NoError(t, err)
	val, err = store.Get(cid, []byte("both-key"))
	require.Equal(t, sql.ErrNoRows, err)
	require.Nil(t, val)
	val, err = store.Get(cid, []byte("cache-only-key"))
	require.NoError(t, err)
	require.Equal(t, []byte("cache-only-value"), val)
}
