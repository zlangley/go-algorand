package layer2

import (
	"testing"

	"github.com/algorand/go-algorand/crypto"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	store, err := NewStableStore(true)
	require.NoError(t, err)

	cache := NewSpeculationStore(store)
	require.NoError(t, err)

	cid := ContractID(crypto.Hash([]byte("test")))

	_, err = store.db.Handle.Exec("INSERT INTO contract_key_values(contract_id, key, value) VALUES($1, $2, $3)", cid.String(), []byte("store-only-key"), []byte("store-only-value"))
	require.NoError(t, err)
	_, err = store.db.Handle.Exec("INSERT INTO contract_key_values(contract_id, key, value) VALUES($1, $2, $3)", cid.String(), []byte("both-key"), []byte("both-store-value"))
	require.NoError(t, err)

	commit1, err := cache.Commitment(cid) // {store-only-key: store-only-value, both-key: both-store-value}
	require.NoError(t, err)

	// Test absence of key.
	val, err := cache.Get(cid, []byte("cache-only-key"))
	require.Error(t, ErrNoKey, err)
	require.Nil(t, val)

	// Test read-your-writes.
	cache.Write(cid, []byte("cache-only-key"), []byte("cache-only-value"), 0)
	val, err = cache.Get(cid, []byte("cache-only-key"))
	require.NoError(t, err)
	require.Equal(t, []byte("cache-only-value"), val)

	// Test commitment changed.
	commit2, err := cache.Commitment(cid) // {store-only-key: store-only-value, both-key: both-store-value, cache-only-key: cache-only-value}
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

	// Test commitment reverted.
	commit3, err := cache.Commitment(cid) // {store-only-key: store-only-value, cache-only-key: cache-only-value}
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
	require.NoError(t, err)
	require.Nil(t, val)
	val, err = store.Get(cid, []byte("cache-only-key"))
	require.NoError(t, err)
	require.Equal(t, []byte("cache-only-value"), val)

	// Commitment should stay the same.
	commit4, err := cache.Commitment(cid)
	require.NoError(t, err)
	require.Equal(t, commit3, commit4)
}

func TestNewSpeculationStore(t *testing.T) {
	cid := ContractID(crypto.Hash([]byte("test")))

	store, err := NewStableStore(true)
	require.NoError(t, err)
	spec := NewSpeculationStore(store)
	spec.Write(cid, []byte("foo"), []byte("bar"), 0)
	spec = NewSpeculationStore(store)
	val, err := spec.Get(cid, []byte("foo"))
	require.Equal(t, ErrNoKey, err)
	require.Nil(t, val)
}