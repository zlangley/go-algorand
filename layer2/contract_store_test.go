package layer2

import (
	"testing"

	"github.com/algorand/go-algorand/crypto"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	store, err := NewStableStore(true)
	require.NoError(t, err)

	cid := ContractID(crypto.Hash([]byte("test")))
	ret := store.Select(cid)
	require.Empty(t, ret)

	store.Write(cid, []byte("bbbk"), []byte("bbbv"))
	val := store.Get(cid, []byte("bbbk"))
	require.Equal(t, val, []byte("bbbv"))

	store.Write(cid, []byte("ccck"), []byte("cccv"))
	val = store.Get(cid, []byte("ccck"))
	require.Equal(t, val, []byte("cccv"))

	store.Write(cid, []byte("dddk"), []byte("dddv"))
	store.Write(cid, []byte("dddk"), nil)
	require.Nil(t, store.Get(cid, []byte("dddk")))

	store.Write(cid, []byte("aaak"), []byte("aaav"))
	val = store.Get(cid, []byte("aaak"))
	require.Equal(t, val, []byte("aaav"))

	cid2 := ContractID(crypto.Hash([]byte("test2")))
	store.Write(cid2, []byte("foo"), []byte("bar2"))

	ret = store.Select(cid)
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

	commit1 := cache.Commitment(cid)

	// Test absence of key.
	require.Nil(t, cache.Get(cid, []byte("cache-only-key")))

	// Test read-your-writes.
	cache.Write(cid, []byte("cache-only-key"), []byte("cache-only-value"), 1)
	require.Equal(t, []byte("cache-only-value"), cache.Get(cid, []byte("cache-only-key")))

	// Test commitment changed.
	commit2 := cache.Commitment(cid)
	require.NotEqual(t, commit1, commit2)

	// Test delegation to underlying store for missing cache key.
	require.Equal(t, []byte("both-store-value"), cache.Get(cid, []byte("both-key")))

	// Test overriding underlying store key.
	cache.Write(cid, []byte("both-key"), []byte("both-cache-value"), 1)
	require.Equal(t, []byte("both-cache-value"), cache.Get(cid, []byte("both-key")))

	// Test delete store key in cache.
	cache.Write(cid, []byte("both-key"), nil, 1)
	require.Nil(t, cache.Get(cid, []byte("both-key")))

	// Test value remains in persistent store.
	require.Equal(t, []byte("both-store-value"), store.Get(cid, []byte("both-key")))

	// Test commitment changed again.
	commit3 := cache.Commitment(cid)
	require.NotEqual(t, commit2, commit3)

	// Test reset.
	cache.Reset()
	require.Nil(t, cache.Get(cid, []byte("cache-only-key")))
	require.Equal(t, commit1, cache.Commitment(cid))
}
