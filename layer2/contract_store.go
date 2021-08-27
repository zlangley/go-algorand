package layer2

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"

	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/protocol"
	"github.com/algorand/go-algorand/util/db"
)

// Stable storage consists of a single table in which each entry is a triple
// (contract_id, key, value).
const stableSchema = `
	CREATE TABLE IF NOT EXISTS contract_key_values(
		contract_id CHAR(58) NOT NULL,
		key CHAR(256) NOT NULL,
		value BLOB NOT NULL,
		PRIMARY KEY (contract_id, key)
	);
`

// In the speculation cache we have two tables: `contract_key_value_writes` contains
// speculative updates to the underlying storage along with the batch index they came
// from, and `txgroups` links batch indices to effects transaction group ids.
//
// The speculation cache lives in-memory. Since we use SQLite for both the stable
// storage and the speculation cache, we also attach them. Attaching them allows us
// to, for example, perform an INSERT INTO ... SELECT statement when updating stable
// storage from the cache.
const speculationSchema = `
	ATTACH DATABASE ':memory:' AS speculation;

	CREATE TABLE speculation.contract_key_values_writes(
		contract_id CHAR(58) NOT NULL,
		batch_idx INT,
		key CHAR(256) NOT NULL,
		value BLOB,
		PRIMARY KEY (contract_id, key, batch_idx)
	);

	CREATE INDEX speculation.contract_key_value_writes__batch_idx__idx ON contract_key_values_writes(batch_idx);

	CREATE TABLE speculation.txgroups(
		group_id CHAR(58) NOT NULL PRIMARY KEY,
		batch_idx INT
	);
`

// A ContractID uniquely specifies a contract.
type ContractID crypto.Digest

func (cid ContractID) String() string {
	return crypto.Digest(cid).String()
}

type StableStore struct {
	db db.Accessor
}

func NewStableStore(inMemory bool) (*StableStore, error) {
	dba, err := db.MakeAccessor("layer2_contract_state.sqlite", false, inMemory)
	if err != nil {
		return nil, err
	}
	_, err = dba.Handle.Exec(stableSchema)
	if err != nil {
		return nil, err
	}
	return &StableStore{db: dba}, nil
}

// Get returns the value corresponding to the given contract ID and key.
//
// If the key does not exist in the store, nil is returned.
func (s *StableStore) Get(cid ContractID, key []byte) ([]byte, error) {
	row := s.db.Handle.QueryRow(`
		SELECT
			value
		FROM
			contract_key_values
		WHERE
			contract_id = $1 AND key = $2
	`, cid.String(), key)

	var value []byte
	err := row.Scan(&value)
	// The value is non-nullable in the database, so we don't need to use the error to
	// differentiate a missing key from a nil value; a nil value always means the key
	// was absent.
	if err == sql.ErrNoRows {
		err = nil
	}
	return value, err
}

type KeyValue struct {
	Key   []byte `json:"key""`
	Value []byte `json:"value"`
}

func (s *StableStore) getWithPrefix(cid ContractID, keyPrefix []byte) ([]KeyValue, error) {
	rows, err := s.db.Handle.Query(`
		SELECT
			key, value
		FROM
			contract_key_values
		WHERE
			contract_id = $1 AND key LIKE $2
		ORDER BY
			key ASC
	`, cid.String(), string(keyPrefix) + "%")

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return parseKeyValues(rows), nil
}

type SpeculationStore struct {
	backingStore *StableStore
	db           db.Accessor
}

func NewSpeculationStore(backingStore *StableStore) *SpeculationStore {
	_, err := backingStore.db.Handle.Exec(speculationSchema)
	if err != nil {
		_, err = backingStore.db.Handle.Exec("DETACH DATABASE speculation;" + speculationSchema)
		if err != nil {
			panic(err)
		}
	}
	return &SpeculationStore{backingStore: backingStore, db: backingStore.db}
}

var ErrKeyDeleted = errors.New("speculation cache: key has been deleted")

// Get returns the value corresponding to the given contract ID and key.
//
// The speculation cache is checked first, and then stable storage is checked.
// There are two reasons the returned bytes can be nil: either there is no such
// key in speculation or stable storage, or the key is in stable storage but has
// been deleted in speculation. To help differentiate between these two cases, we
// return ErrKeyDeleted when the latter case occurs.
func (s *SpeculationStore) Get(cid ContractID, key []byte) ([]byte, error) {
	row := s.db.Handle.QueryRow(`
		SELECT
			value
		FROM
			contract_key_values_writes
		WHERE
			contract_id = $1 AND key = $2
		ORDER BY
			batch_idx DESC
		LIMIT 1
	`, cid.String(), key)

	var value []byte
	err := row.Scan(&value)
	if err == sql.ErrNoRows {
		// Not in speculation cache? Check store.
		value, err = s.backingStore.Get(cid, key)
		if err != nil {
			return nil, err
		}
		return value, nil
		err = nil
	} else if err != nil {
		// In-memory cache error should never happen.
		panic(err)
	} else if value == nil {
		return value, ErrKeyDeleted
	}
	return value, nil
}

func (s *SpeculationStore) Write(cid ContractID, key []byte, val []byte, batch_idx int) {
	_, err := s.db.Handle.Exec(`
		INSERT OR REPLACE INTO contract_key_values_writes(contract_id, key, value, batch_idx)
		VALUES ($1, $2, $3, $4);
	`, cid.String(), key, val, batch_idx)
	if err != nil {
		panic(err)
	}
}

func (s *SpeculationStore) SetBatchIndexGroup(batch_idx int, groupID crypto.Digest) {
	_, err := s.db.Handle.Exec("INSERT INTO txgroups(group_id, batch_idx) VALUES ($1, $2)", groupID.String(), batch_idx)
	if err != nil {
		panic(err)
	}
}

func (s *SpeculationStore) PersistGroupState(groupID crypto.Digest) error {
	// NULL values need to be handled separately; a NULL value in the cache
	// represents a deletion when we persist, while any other values correspond to upserts.
	//
	// Here we leverage that the speculation database is attached to the persistent one;
	// rather than fetch from the in-memory database, parse the result set, and then issue
	// upserts/deletes to the persistent database, we can just write two simple SQL statements.
	// If we wanted to use some other key-value store implementation for the persistent
	// database, this would need to be changed.
	_, err := s.db.Handle.Exec(`
		DELETE FROM contract_key_values
		WHERE (contract_id, key) IN (
			SELECT contract_id, key FROM
				txgroups t
			JOIN contract_key_values_writes s ON
				t.batch_idx = s.batch_idx
			WHERE
				t.group_id = $1 AND s.value IS NULL
		);

		INSERT OR REPLACE INTO contract_key_values(contract_id, key, value)
		SELECT contract_id, key, value FROM
			txgroups t
		JOIN contract_key_values_writes s ON
			t.batch_idx = s.batch_idx
		WHERE
			t.group_id = $2 AND s.value IS NOT NULL;
	`, groupID.String(), groupID.String())
	if err != nil {
		return err
	}
	return nil
}

func (s *SpeculationStore) GetWithPrefix(cid ContractID, keyPrefix []byte) ([]KeyValue, error) {
	// Since the databases are attached, we could get the latest key-value pairs
	// with a JOIN, but it seems more complicated than just manually merging two
	// separate queries (especially since sqlite3 does not support full outer joins).
	rows, err := s.db.Handle.Query(`
		SELECT
			a.key, a.value
		FROM
			contract_key_values_writes a
			LEFT JOIN
				contract_key_values_writes b
			ON
				a.contract_id = b.contract_id AND
				a.key = b.key AND 
				a.batch_idx < b.batch_idx
		WHERE
			a.contract_id = $1 AND b.batch_idx IS NULL AND a.key LIKE $2
		ORDER BY
			a.key ASC
	`, cid.String(), string(keyPrefix) + "%")

	// In-memory db should never fail.
	if err != nil {
		panic(err)
	}

	cachePairs := parseKeyValues(rows)
	storePairs, err := s.backingStore.getWithPrefix(cid, keyPrefix)
	if err != nil {
		return nil, err
	}
	return mergeKeyValuePairs(storePairs, cachePairs), nil
}

func (s *SpeculationStore) Commitment(cid ContractID) (crypto.Digest, error) {
	kvs, err := s.GetWithPrefix(cid, []byte{})
	fmt.Println(kvs)
	if err != nil {
		return crypto.Digest{}, err
	}
	encoded := protocol.EncodeJSON(kvs)
	return crypto.Hash(encoded), nil
}

func mergeKeyValuePairs(storePairs, cachePairs []KeyValue) []KeyValue {
	var cidx, sidx int
	var merged []KeyValue
	for cidx < len(cachePairs) && sidx < len(storePairs) {
		ckv := cachePairs[cidx]
		skv := storePairs[sidx]
		if bytes.Compare(skv.Key, ckv.Key) < 0 {
			merged = append(merged, storePairs[sidx])
			sidx++
		} else {
			if cachePairs[cidx].Value != nil {
				merged = append(merged, cachePairs[cidx])
			}
			if bytes.Compare(cachePairs[cidx].Key, storePairs[sidx].Key) == 0 {
				sidx++
			}
			cidx++
		}
	}
	for cidx < len(cachePairs) {
		merged = append(merged, cachePairs[cidx])
		cidx++
	}
	for sidx < len(storePairs) {
		merged = append(merged, storePairs[sidx])
		sidx++
	}
	return merged
}

func parseKeyValues(rows *sql.Rows) []KeyValue {
	var kvs []KeyValue
	for rows.Next() {
		var key, value []byte
		rows.Scan(&key, &value)
		kvs = append(kvs, KeyValue{key, value})
	}
	return kvs
}
