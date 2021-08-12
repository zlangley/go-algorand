package layer2

import (
	"database/sql"
	"fmt"

	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/protocol"
	"github.com/algorand/go-algorand/util/db"
)

var schema = `
	CREATE TABLE IF NOT EXISTS contract_key_value_pairs(
		contract_id CHAR(58) NOT NULL,
		key CHAR(256) NOT NULL,
		value BLOB,
        dirty INTEGER,
		PRIMARY KEY (contract_id, key)
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

type Store struct {
	db db.Accessor
}

type SpeculationStore struct {
	db db.Accessor
}

func NewStore() (*Store, *SpeculationStore, error) {
	db, err := db.MakeAccessor("contract.sqlite", false, true)
	if err != nil {
		return nil, nil, err
	}
	_, err = db.Handle.Exec(schema)
	if err != nil {
		return nil, nil, err
	}
	return &Store{db: db}, &SpeculationStore{db}, nil
}

func (s *Store) Get(cid ContractID, key []byte) []byte {
	return get(s.db, cid, key, false)
}

func (s *Store) Write(cid ContractID, key []byte, val []byte) {
	if val == nil {
		s.db.Handle.Exec(`
			DELETE FROM
		        contract_key_value_pairs
    		WHERE
    		    contract_id = $1 AND key = $2 AND dirty = 0
		`, cid.String(), key)
		return
	}
	write(s.db, cid, key, val, false)
}

func (s *Store) Select(cid ContractID) []KeyValuePair {
	return selectPairs(s.db, cid, false)
}

func (s *SpeculationStore) Get(cid ContractID, key []byte) []byte {
	return get(s.db, cid, key, true)
}

func (s *SpeculationStore) Write(cid ContractID, key []byte, val []byte) {
	write(s.db, cid, key, val, true)
}

func (s *SpeculationStore) Commitment(cid ContractID) crypto.Digest {
	kvs := selectPairs(s.db, cid, true)
	encoded := protocol.EncodeJSON(kvs)
	return crypto.Hash(encoded)
}

func get(db db.Accessor, cid ContractID, key []byte, includeDirty bool) []byte {
	dirty := 0
	if includeDirty {
		dirty = 1
	}
	row := db.Handle.QueryRow(`
		SELECT
		    value
		FROM
		    contract_key_value_pairs
		WHERE
		    contract_id = $1 AND key = $2 AND dirty <= $3
		LIMIT
			1
	`, cid.String(), key, dirty)
	var value []byte
	err := row.Scan(&value)
	if err == sql.ErrNoRows {
		return nil
	}
	return value
}

func selectPairs(db db.Accessor, cid ContractID, includeDirty bool) []KeyValuePair {
	dirtyBound := 0
	if includeDirty {
		dirtyBound = 1
	}
	rows, err := db.Handle.Query(`
		SELECT
			key, value
		FROM
			contract_key_value_pairs
		WHERE
			contract_id = $1 AND dirty <= $2
		ORDER BY
			key ASC
    `, cid.String(), dirtyBound)
	fmt.Println(err)

	var kvs []KeyValuePair
	for rows.Next() {
		var key, value []byte
		rows.Scan(&key, &value)
		kvs = append(kvs, KeyValuePair{key, value})
	}
	return kvs
}

func write(db db.Accessor, cid ContractID, key []byte, val []byte, isDirty bool) {
	dirtyBit := 0
	if isDirty {
		dirtyBit = 1
	}
	// TODO: handle nil values
	_, err := db.Handle.Exec(`
		INSERT INTO contract_key_value_pairs(contract_id, key, value, dirty)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT(contract_id, key)
		DO UPDATE SET
			value=$3,
		    dirty=$4;
	`, cid.String(), key, val, dirtyBit)
	fmt.Println(err)
}
