package data

import (
	"github.com/algorand/go-algorand/agreement"
	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/data/bookkeeping"
	"testing"

	"github.com/algorand/go-algorand/config"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/data/transactions"
	"github.com/algorand/go-algorand/ledger"
	"github.com/algorand/go-algorand/logging"
	"github.com/algorand/go-algorand/protocol"

	"github.com/stretchr/testify/require"
)


func TestSquashSpeculativeTransactions(t *testing.T) {
	genesisInitState, keys := testGenerateInitState(t, protocol.ConsensusCurrentVersion)
	genesisHash := genesisInitState.Block.GenesisHash()

	src := findAccountWithStatus(t, genesisInitState.Accounts, basics.Online)
	srcKey := keys[src]
	require.NotNil(t, srcKey)

	dst := findAccountWithStatus(t, genesisInitState.Accounts, basics.NotParticipating)

	l := createLedger(t, genesisInitState)
	defer l.Close()

	sl, err := NewSpeculationLedger(l, l.Latest())
	require.NoError(t, err)

	srcData, err := sl.LookupLatest(src)
	require.NoError(t, err)

	require.NoError(t, sl.Apply(makeSignedTxn(1, genesisHash, src, dst, srcKey)))
	require.NoError(t, sl.Apply(makeSignedTxn(2, genesisHash, src, dst, srcKey)))
	require.NoError(t, sl.Apply(makeSignedTxn(3, genesisHash, src, dst, srcKey)))
	require.Equal(t, 3, len(sl.workingTxGroups))
	srcData2, err := sl.LookupLatest(src)
	require.NoError(t, err)
	require.Equal(t, srcData.MicroAlgos.Raw - 30000 - 3, srcData2.MicroAlgos.Raw)
	require.NoError(t, sl.StageWorking())
	require.NoError(t, err)
	require.Equal(t, 1, len(sl.stagedTxGroups))
	require.Equal(t, 0, len(sl.workingTxGroups))

	require.NoError(t, sl.Apply(makeSignedTxn(4, genesisHash, src, dst, srcKey)))
	require.NoError(t, sl.Apply(makeSignedTxn(5, genesisHash, src, dst, srcKey)))
	require.Equal(t, 2, len(sl.workingTxGroups))
	require.NoError(t, sl.StageWorking())
	require.NoError(t, err)
	require.Equal(t, 2, len(sl.stagedTxGroups))
	require.Equal(t, 0, len(sl.workingTxGroups))
}


// creates a ledger with an initial block.
func createLedger(t *testing.T, genesisInitState ledger.InitState) *Ledger {
	const inMem = true
	cfg := config.GetDefaultLocal()
	cfg.Archival = true
	log := logging.TestingLog(t)
	log.SetLevel(logging.Warn)
	realLedger, err := ledger.OpenLedger(log, t.Name(), inMem, genesisInitState, cfg)
	require.NoError(t, err, "could not open ledger")

	l := &Ledger{Ledger: realLedger}
	require.NotNil(t, &l)

	blk := bookkeeping.MakeBlock(genesisInitState.Block.BlockHeader)
	err = l.AddBlock(blk, agreement.Certificate{})
	require.NoError(t, err)

	l.WaitForCommit(1)
	return l
}


func findAccountWithStatus(t *testing.T, accts map[basics.Address]basics.AccountData, status basics.Status) basics.Address {
	var ret basics.Address
	for addr, acctData := range accts {
		if addr != testPoolAddr && addr != testSinkAddr && acctData.Status == status {
			ret = addr
			break
		}
	}
	require.False(t, ret.IsZero())
	return ret
}


func makeSignedTxn(id int, genesisHash crypto.Digest, src, dst basics.Address, srcKey *crypto.SignatureSecrets) []transactions.SignedTxn {
	tx := transactions.Transaction{
		Type: protocol.PaymentTx,
		Header: transactions.Header{
			Sender:      src,
			Fee:         basics.MicroAlgos{Raw: 10000},
			FirstValid:  0,
			LastValid:   1000,
			GenesisHash: genesisHash,
			Note:        []byte{byte(id)},
		},
		PaymentTxnFields: transactions.PaymentTxnFields{
			Receiver: dst,
			Amount:   basics.MicroAlgos{Raw: 1},
		},
	}
	return []transactions.SignedTxn{tx.Sign(srcKey)}
}

