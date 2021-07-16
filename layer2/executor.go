package layer2

import (
	"fmt"
	"github.com/algorand/go-algorand/data/transactions"
	"github.com/algorand/go-algorand/protocol"
	"github.com/algorand/go-algorand/util"
	"path"
	"strconv"

	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/data"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/layer2/kalgo"
)

type batchItemKind int8

const (
	Init batchItemKind = iota
	Call
)

type BatchItem struct {
	kind   batchItemKind
	runner kalgo.Runner
	sender basics.Address
}

func NewInitBatchItem(name, source string, sender basics.Address) *BatchItem {
	sourceHash := crypto.Hash([]byte(source))
	addr := basics.Address(crypto.Hash([]byte(sender.GetUserAddress() + sourceHash.String())))

	return &BatchItem{
		kind: Init,
		runner: &kalgo.InitCmd{
			Cmd: kalgo.Cmd{
				Address: addr,
				Name:    name,
				Sender:  sender,
			},
			Source: source,
		},
		sender: sender,
	}
}

func NewCallBatchItem(name, function, args string, addr, sender basics.Address) *BatchItem {
	return &BatchItem{
		kind: Call,
		runner: &kalgo.CallCmd{
			Cmd: kalgo.Cmd{
				Address: addr,
				Name:    name,
				Sender:  sender,
			},
			Function: function,
			Args:     args,
		},
		sender: sender,
	}
}

type ContractCommitment struct {
	Contract   string
	PrevCommit string
	NewCommit  string
}

type Executor struct {
	ledger *data.SpeculationLedger
	kenv   kalgo.Env
}

func NewExecutor(ledger *data.SpeculationLedger, kenv kalgo.Env) *Executor {
	return &Executor{
		ledger: ledger,
		kenv:   kenv,
	}
}

func (ex *Executor) Execute(item *BatchItem, prof *util.Profiler) error {
	prof.Start("kalgo")
	rawout, err := item.runner.Run(ex.kenv)
	if err != nil {
		return fmt.Errorf("running kalgo: %w", err)
	}

	prof.Start("node")
	out, err := kalgo.ParseOutput(rawout)
	if err != nil {
		return fmt.Errorf("parsing kalgo: %w", err)
	}

	if err = ex.ledger.Checkpoint(); err != nil {
		return err
	}

	for _, commit := range out.Commitments {
		if commit.PreviousCommitment == commit.NewCommitment {
			continue
		}

		//ex.addCommitmentCheckTx([]byte(commit.PreviousCommitment), []byte(commit.NewCommitment), item.sender)

		prof.Start("cow")
		if err = ex.copyContract(commit.Contract); err != nil {
			return err
		}
		prof.Start("node")
	}
	return nil
}

func (ex *Executor) copyContract(contract string) error {
	src := path.Join(ex.kenv.SourcePrefix, contract)
	dst := path.Join(ex.kenv.SourcePrefix, "..", strconv.Itoa(len(ex.ledger.Checkpoints)), contract)
	return util.CopyFolder(src, dst)
}

var CommitmentAppIndex basics.AppIndex = 19841517
//var CommitmentMapAddress = "Z72YSHBKWQMW66SGB6WHH6VOK46KDSTJZ2NZJPRXQG4RK6U2Y62ZNR6FJY"

func (ex *Executor) addCommitmentAppOptInTxn(commitment []byte, sender basics.Address) {
	tx := transactions.Transaction{
		Type: protocol.ApplicationCallTx,
		Header: transactions.Header{
			Sender:      sender,
			Fee:         basics.MicroAlgos{Raw: 10000},
			FirstValid:  ex.ledger.Latest(),
			LastValid:   ex.ledger.Latest() + 1000,
			//GenesisHash: ledger.GenesisHash,
		},
		ApplicationCallTxnFields: transactions.ApplicationCallTxnFields{
			ApplicationID: CommitmentAppIndex,
			ApplicationArgs: [][]byte{commitment},
			OnCompletion: transactions.OptInOC,
		},
	}
/*	stx := transactions.SignedTxn{
		Txn: tx,
		Lsig: transactions.LogicSig{
			Logic:
		}
	}*/
	// TODO: sign from execution committee
	stx := tx.Sign(nil)
	ex.ledger.Apply([]transactions.SignedTxn{stx})
}

func (ex *Executor) addCommitmentAppCallTxn(prev, new []byte, sender basics.Address) {
	tx := transactions.Transaction{
		Type: protocol.ApplicationCallTx,
		Header: transactions.Header{
			Sender:      sender,
			Fee:         basics.MicroAlgos{Raw: 10000},
			FirstValid:  ex.ledger.Latest(),
			LastValid:   ex.ledger.Latest() + 1000,
			//GenesisHash: ledger.GenesisHash,
		},
		ApplicationCallTxnFields: transactions.ApplicationCallTxnFields{
			ApplicationID: CommitmentAppIndex,
			ApplicationArgs: [][]byte{prev, new},
		},
	}
	// TODO: sign from execution committee
	stx := tx.Sign(nil)
	ex.ledger.Apply([]transactions.SignedTxn{stx})
}
