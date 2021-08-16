package layer2

import (
	"bytes"
	"database/sql"
	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/protocol"
	"github.com/algorand/go-algorand/util/db"
)

var stableSchema = `
	CREATE TABLE IF NOT EXISTS contract_kv_pairs(
		contract_id CHAR(58) NOT NULL,
		key CHAR(256) NOT NULL,
		value BLOB NOT NULL,
		PRIMARY KEY (contract_id, key)
	);
`
var speculationSchema = `
	CREATE TABLE speculation.sequenced_contract_kv_pairs(
		contract_id CHAR(58) NOT NULL,
		key CHAR(256) NOT NULL,
		value BLOB,
		batch_idx INT,
		PRIMARY KEY (contract_id, key, batch_idx)
	);

	CREATE INDEX speculation.sequenced_contact_kv_pairs__batch_idx__idx ON sequenced_contract_kv_pairs(batch_idx);

	CREATE TABLE speculation.txgroups(
	    group_id CHAR(58) NOT NULL PRIMARY KEY,
	    batch_idx INT
    );
`

type ContractID basics.Address

func (cid ContractID) String() string {
	return basics.Address(cid).String()
}

type KeyValuePair struct {
	Key   []byte `json:"key""`
	Value []byte `json:"value"`
}

type StableStore struct {
	db db.Accessor
}

func NewStableStore(inMemory bool) (*StableStore, error) {
	dba, err := db.MakeAccessor("contract.sqlite", false, inMemory)
	if err != nil {
		return nil, err
	}
	_, err = dba.Handle.Exec(stableSchema)
	if err != nil {
		return nil, err
	}
	return &StableStore{db: dba}, nil
}

func (s *StableStore) Get(cid ContractID, key []byte) ([]byte, error) {
	row := s.db.Handle.QueryRow(`
		SELECT
		    value
		FROM
		    contract_kv_pairs
		WHERE
		    contract_id = $1 AND key = $2
	`, cid.String(), key)

	var value []byte
	err := row.Scan(&value)
	return value, err
}

func (s *StableStore) Select(cid ContractID) ([]KeyValuePair, error) {
	rows, err := s.db.Handle.Query(`
		SELECT
			key, value
		FROM
			contract_kv_pairs
		WHERE
			contract_id = $1
		ORDER BY
			key ASC
    `, cid.String())

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return resultSetToKVPairs(rows), nil
}

type SpeculationStore struct {
	backingStore *StableStore
	db           db.Accessor
}

func (s *StableStore) Speculation() (*SpeculationStore, error) {
	_, err := s.db.Handle.Exec("ATTACH DATABASE ':memory:' AS speculation")
	dba := s.db
	if err != nil {
		return nil, err
	}
	_, err = dba.Handle.Exec(speculationSchema)
	if err != nil {
		return nil, err
	}
	spec := &SpeculationStore{backingStore: s, db: dba}
	return spec, nil
}

func (s *SpeculationStore) Get(cid ContractID, key []byte) ([]byte, error) {
	row := s.db.Handle.QueryRow(`
		SELECT
		    value
		FROM
		    sequenced_contract_kv_pairs
		WHERE
		    contract_id = $1 AND key = $2
		ORDER BY
			batch_idx DESC
	`, cid.String(), key)

	var value []byte
	err := row.Scan(&value)
	if err == sql.ErrNoRows {
		value, err = s.backingStore.Get(cid, key)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
	} else if err != nil {
		// In-memory store error should never happen.
		panic(err)
	}
	return value, nil
}

func (s *SpeculationStore) Write(cid ContractID, key []byte, val []byte, batch_idx int) {
	_, err := s.db.Handle.Exec(`
		INSERT INTO sequenced_contract_kv_pairs(contract_id, key, value, batch_idx)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT(contract_id, key, batch_idx)
		DO UPDATE SET value=$3;
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
	_, err := s.db.Handle.Exec(`
		DELETE FROM contract_kv_pairs
		WHERE (contract_id, key) IN (
      		SELECT contract_id, key FROM
				txgroups t
			JOIN sequenced_contract_kv_pairs s ON
				t.batch_idx = s.batch_idx
			WHERE
				t.group_id = $1 AND s.value IS NULL
		);

		INSERT OR REPLACE INTO contract_kv_pairs(contract_id, key, value)
		SELECT contract_id, key, value FROM
			txgroups t
		JOIN sequenced_contract_kv_pairs s ON
			t.batch_idx = s.batch_idx
		WHERE
			t.group_id = $1 AND s.value IS NOT NULL;
	`, groupID.String(), groupID.String())
	if err != nil {
		return err
	}
	return nil
}

func (s *SpeculationStore) Commitment(cid ContractID) (crypto.Digest, error) {
	rows, err := s.db.Handle.Query(`
		SELECT
			a.key, a.value
		FROM
			sequenced_contract_kv_pairs a
			LEFT JOIN
				sequenced_contract_kv_pairs b
			ON
				a.key = b.key AND 
				a.batch_idx < b.batch_idx
		WHERE
			a.contract_id = $1 AND b.batch_idx IS NULL
		ORDER BY
			a.key ASC
    `, cid.String())
	if err != nil {
		panic(err)
	}

	cachePairs := resultSetToKVPairs(rows)
	storePairs, err := s.backingStore.Select(cid)
	if err != nil {
		return crypto.Digest{}, err
	}
	merged := mergeKeyValuePairs(storePairs, cachePairs)
	encoded := protocol.EncodeJSON(merged)
	return crypto.Hash(encoded), nil
}

func mergeKeyValuePairs(storePairs, cachePairs []KeyValuePair) []KeyValuePair {
	var cidx, sidx int
	var merged []KeyValuePair
	for cidx < len(cachePairs) && sidx < len(storePairs) {
		ckv := cachePairs[cidx]
		skv := storePairs[sidx]
		if bytes.Compare(skv.Key, ckv.Key) < 0 {
			merged = append(merged, storePairs[sidx])
			sidx++
		} else {
			merged = append(merged, cachePairs[cidx])
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

func resultSetToKVPairs(rows *sql.Rows) []KeyValuePair {
	var kvs []KeyValuePair
	for rows.Next() {
		var key, value []byte
		rows.Scan(&key, &value)
		kvs = append(kvs, KeyValuePair{key, value})
	}
	return kvs
}
