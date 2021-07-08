package data

import (
	"github.com/algorand/go-algorand/agreement"
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

	const inMem = true
	cfg := config.GetDefaultLocal()
	cfg.Archival = true
	log := logging.TestingLog(t)
	log.SetLevel(logging.Warn)
	realLedger, err := ledger.OpenLedger(log, t.Name(), inMem, genesisInitState, cfg)
	require.NoError(t, err, "could not open ledger")
	defer realLedger.Close()

	l := &Ledger{Ledger: realLedger}
	require.NotNil(t, &l)

	var sourceAccount basics.Address
	var destAccount basics.Address
	for addr, acctData := range genesisInitState.Accounts {
		if addr == testPoolAddr || addr == testSinkAddr {
			continue
		}
		if acctData.Status == basics.Online {
			sourceAccount = addr
			break
		}
	}
	for addr, acctData := range genesisInitState.Accounts {
		if addr == testPoolAddr || addr == testSinkAddr {
			continue
		}
		if acctData.Status == basics.NotParticipating {
			destAccount = addr
			break
		}
	}
	require.False(t, sourceAccount.IsZero())
	require.False(t, destAccount.IsZero())

	srcAccountKey := keys[sourceAccount]
	require.NotNil(t, srcAccountKey)

	blk := bookkeeping.MakeBlock(genesisInitState.Block.BlockHeader)
	err = l.AddBlock(blk, agreement.Certificate{})
	require.NoError(t, err)

	l.WaitForCommit(1)

	signedTx := func (id int) []transactions.SignedTxn {
		tx := transactions.Transaction{
			Type: protocol.PaymentTx,
			Header: transactions.Header{
				Sender: sourceAccount,
				Fee: basics.MicroAlgos{Raw: 10000},
				FirstValid: 0,
				LastValid: 1000,
				GenesisHash: genesisInitState.Block.BlockHeader.GenesisHash,
				Note: []byte{byte(id)},
			},
			PaymentTxnFields: transactions.PaymentTxnFields{
				Receiver: destAccount,
				Amount: basics.MicroAlgos{Raw: 1},
			},
		}
		return []transactions.SignedTxn{tx.Sign(srcAccountKey)}
	}

	sl, err := NewSpeculationLedger(l, l.Latest())
	require.NoError(t, err)

	require.NoError(t, sl.Apply(signedTx(1)))
	require.NoError(t, sl.Apply(signedTx(2)))
	require.NoError(t, sl.Apply(signedTx(3)))
	require.Equal(t, 3, len(sl.workingTxGroups))
	require.NoError(t, sl.StageWorking())
	require.NoError(t, err)
	require.Equal(t, 1, len(sl.stagedTxGroups))
	require.Equal(t, 0, len(sl.workingTxGroups))

	require.NoError(t, sl.Apply(signedTx(4)))
	require.NoError(t, sl.Apply(signedTx(5)))
	require.Equal(t, 2, len(sl.workingTxGroups))
	require.NoError(t, sl.StageWorking())
	require.NoError(t, err)
	require.Equal(t, 2, len(sl.stagedTxGroups))
	require.Equal(t, 0, len(sl.workingTxGroups))
}
