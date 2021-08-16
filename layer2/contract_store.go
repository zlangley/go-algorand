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
	CREATE TABLE IF NOT EXISTS contract_key_value_pairs(
		contract_id CHAR(58) NOT NULL,
		key CHAR(256) NOT NULL,
		value BLOB,
		PRIMARY KEY (contract_id, key)
	);
`
var speculationSchema = `
	CREATE TABLE IF NOT EXISTS speculative_contract_key_value_pairs(
		contract_id CHAR(58) NOT NULL,
		key CHAR(256) NOT NULL,
		value BLOB,
		seqno INT,
		PRIMARY KEY (contract_id, key, seqno)
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

type SpeculationStore struct {
	backingStore *StableStore
	db           db.Accessor
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
		    contract_key_value_pairs
		WHERE
		    contract_id = $1 AND key = $2
		LIMIT
			1
	`, cid.String(), key)

	if err := row.Err(); err != nil {
		return nil, err
	}

	var value []byte
	err := row.Scan(&value)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return value, nil
}

func (s *StableStore) Write(cid ContractID, key []byte, val []byte) error {
	if val == nil {
		_, err := s.db.Handle.Exec(`
			DELETE FROM
		        contract_key_value_pairs
    		WHERE
    		    contract_id = $1 AND key = $2
		`, cid.String(), key)
		return err
	}
	_, err := s.db.Handle.Exec(`
		INSERT INTO contract_key_value_pairs(contract_id, key, value)
		VALUES ($1, $2, $3)
		ON CONFLICT(contract_id, key)
		DO UPDATE SET value=$3;
	`, cid.String(), key, val)
	return err
}

func (s *StableStore) Select(cid ContractID) ([]KeyValuePair, error) {
	rows, err := s.db.Handle.Query(`
		SELECT
			key, value
		FROM
			contract_key_value_pairs
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

func (s *StableStore) Speculation() (*SpeculationStore, error) {
	dba, err := db.MakeAccessor("layer2speculation.sqlite", false, true)
	if err != nil {
		return nil, err
	}
	_, err = dba.Handle.Exec(speculationSchema)
	if err != nil {
		return nil, err
	}
	spec := &SpeculationStore{backingStore: s, db: dba}
	spec.Reset()
	return spec, nil
}

func (s *SpeculationStore) Get(cid ContractID, key []byte) ([]byte, error) {
	row := s.db.Handle.QueryRow(`
		SELECT
		    value
		FROM
		    speculative_contract_key_value_pairs
		WHERE
		    contract_id = $1 AND key = $2
		ORDER BY
			seqno DESC
		LIMIT
			1
	`, cid.String(), key)

	if err := row.Err(); err != nil {
		panic(err)
	}

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

func (s *SpeculationStore) Write(cid ContractID, key []byte, val []byte, seqno int) {
	_, err := s.db.Handle.Exec(`
		INSERT INTO speculative_contract_key_value_pairs(contract_id, key, value, seqno)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT(contract_id, key, seqno)
		DO UPDATE SET value=$3;
	`, cid.String(), key, val, seqno)
	if err != nil {
		panic(err)
	}
}

func (s *SpeculationStore) Commitment(cid ContractID) (crypto.Digest, error) {
	rows, err := s.db.Handle.Query(`
		SELECT
			a.key, a.value
		FROM
			speculative_contract_key_value_pairs a
			LEFT JOIN
				speculative_contract_key_value_pairs b
			ON
				a.key = b.key AND 
				a.seqno < b.seqno
		WHERE
			a.contract_id = $1 AND b.seqno IS NULL
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

func (s *SpeculationStore) Reset() {
	_, err := s.db.Handle.Exec("DELETE FROM speculative_contract_key_value_pairs")
	if err != nil {
		panic(err)
	}
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
