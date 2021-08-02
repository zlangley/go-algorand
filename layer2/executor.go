package layer2

import (
	"bytes"
	"fmt"
	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/data"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/data/transactions"
	"github.com/algorand/go-algorand/layer2/kalgo"
	"github.com/algorand/go-algorand/protocol"
	"github.com/algorand/go-algorand/util"
	"path"
	"strconv"
)

type batchItemKind int8

const (
	Init batchItemKind = iota
	Call
)

type BatchItem struct {
	kind       batchItemKind
	name       string
	runner     kalgo.Runner
	contractID crypto.Digest
	sender     basics.Address
}

func NewInitBatchItem(name, source string, sender basics.Address) *BatchItem {
	sourceHash := crypto.Hash([]byte(source))
	contractID := crypto.Hash([]byte(sender.GetUserAddress() + sourceHash.String()))
	addr, _, _ := logicSigFromTemplateFile("layer2/committee-defer-logicsic.teal.template", contractID)

	return &BatchItem{
		kind: Init,
		name: name,
		runner: &kalgo.InitCmd{
			Cmd: kalgo.Cmd{
				Address: addr,
				Name:    name,
				Sender:  sender,
			},
			Source: source,
		},
		contractID: contractID,
		sender:     sender,
	}
}

func NewCallBatchItem(name, function, args string, addr, sender basics.Address) *BatchItem {
	return &BatchItem{
		kind: Call,
		name: name,
		runner: &kalgo.CallCmd{
			Cmd: kalgo.Cmd{
				Address: addr,
				Name:    name,
				Sender:  sender,
			},
			Function: function,
			Args:     args,
		},
		sender:     sender,
	}
}

type executionResult struct {
	effectsTxns []transactions.Transaction
	kalgoOutput *kalgo.Output
}

type Executor struct {
	ledger *data.SpeculationLedger
	kenv   kalgo.Env

	lastLedgerIdx int
	results       []*executionResult

	// TODO: kalgo should return contract ids eventually
	contractIDs map[string]crypto.Digest
}

func NewExecutor(ledger *data.SpeculationLedger, kenv kalgo.Env) *Executor {
	return &Executor{
		ledger: ledger,
		kenv:   kenv,
		contractIDs: make(map[string]crypto.Digest),
	}
}

// Execute executes the given BatchItem and commits to the speculative ledger.
//
// If successful, the effects transactions from the execution are applied to
// the speculative ledger.
func (ex *Executor) Execute(item *BatchItem, prof *util.Profiler) error {
	prof.Start("kalgo")
	rawout, err := item.runner.Run(ex.kenv)
	if err != nil {
		return fmt.Errorf("running kalgo: %w", err)
	}

	prof.Start("node")

	// Store the contract ID away in case we call this again later by name.
	if item.kind == Init {
		ex.contractIDs[item.name] = item.contractID
	}

	// TODO: Run() should probably just return a *kalgo.Output, but for now, running
	// the VM and parsing its output are separate to contain the profiling logic.
	out, err := kalgo.ParseOutput(rawout)
	if err != nil {
		return fmt.Errorf("parsing kalgo: %w", err)
	}

	prof.Start("effects")
	// Add version updates.
	for _, commit := range out.Commitments {
		contractID := ex.contractIDs[commit.Contract]
		var txn transactions.Transaction
		if len(commit.PreviousCommitment) == 0 {
			txn, err = CommitmentAppOptInTxn(contractID, commit.NewCommitment, ex.ledger.Latest(), ex.ledger.Latest() + 1000, ex.ledger.GenesisHash())
		} else if bytes.Compare(commit.PreviousCommitment, commit.NewCommitment) == 0 {
			txn, err = CommitmentAppCallTxn(contractID, commit.PreviousCommitment, commit.NewCommitment, ex.ledger.Latest(), ex.ledger.Latest()+1000, ex.ledger.GenesisHash())
		}
		if err != nil {
			return fmt.Errorf("could not create commitment swap contract call: %v", err)
		}
		err = ex.ledger.Apply([]transactions.SignedTxn{{Txn: txn}})
		if err != nil {
			return fmt.Errorf("could not apply version swap call (%v -> %v): %v", commit.PreviousCommitment, commit.NewCommitment, err)
		}
	}

	stxnGroups := ex.ledger.TxStack()[ex.lastLedgerIdx:]
	var effects []transactions.Transaction
	for _, stxnGroup := range stxnGroups {
		for _, stx := range stxnGroup {
			effects = append(effects, stx.Txn)
		}
	}
	result := &executionResult{
		effectsTxns: effects,
		kalgoOutput: out,
	}
	ex.results = append(ex.results, result)
	ex.lastLedgerIdx = len(ex.ledger.TxStack())

	prof.Start("cow")
	for _, commit := range out.Commitments {
		if err = ex.copyContract(commit.Contract); err != nil {
			return err
		}
	}
	prof.Start("node")
	return nil
}

func (ex *Executor) copyContract(contract string) error {
	src := path.Join(ex.kenv.SourcePrefix, contract)
	dst := path.Join(ex.kenv.SourcePrefix, "..", strconv.Itoa(len(ex.results)), contract)
	return util.CopyFolder(src, dst)
}

func (ex *Executor) EffectsTxns() [][]transactions.Transaction {
	var txns [][]transactions.Transaction
	for _, res := range ex.results {
		txns = append(txns, res.effectsTxns)
	}
	return txns
}

// associated with every l2 contract, we have:
//   - contract ID
//   - escrow account (logicsig)

//var CommitmentMapAddress = "Z72YSHBKWQMW66SGB6WHH6VOK46KDSTJZ2NZJPRXQG4RK6U2Y62ZNR6FJY"

func CommitmentAppOptInTxn(contractID crypto.Digest, commitment []byte, fv, lv basics.Round, genesisHash crypto.Digest) (transactions.Transaction, error) {
	addr, _, err := logicSigFromTemplateFile("layer2/committee-defer-logicsic.teal.template", contractID)
	if err != nil {
		return transactions.Transaction{}, err
	}
	return transactions.Transaction{
		Type: protocol.ApplicationCallTx,
		Header: transactions.Header{
			Sender:      addr,
			Fee:         basics.MicroAlgos{Raw: 1000},
			FirstValid:  fv,
			LastValid:   lv,
			GenesisHash: genesisHash,
		},
		ApplicationCallTxnFields: transactions.ApplicationCallTxnFields{
			ApplicationID:   VersionNumberAppIndex,
			ApplicationArgs: [][]byte{commitment},
			OnCompletion:    transactions.OptInOC,
		},
	}, nil
}

func CommitmentAppCallTxn(contractID crypto.Digest, prev, new []byte, fv, lv basics.Round, genesisHash crypto.Digest) (transactions.Transaction, error) {
	addr, _, err := logicSigFromTemplateFile("layer2/committee-defer-logicsic.teal.template", contractID)
	if err != nil {
		return transactions.Transaction{}, err
	}
	return transactions.Transaction{
		Type: protocol.ApplicationCallTx,
		Header: transactions.Header{
			Sender:      addr,
			Fee:         basics.MicroAlgos{Raw: 1000},
			FirstValid:  fv,
			LastValid:   lv,
			GenesisHash: genesisHash,
		},
		ApplicationCallTxnFields: transactions.ApplicationCallTxnFields{
			ApplicationID:   VersionNumberAppIndex,
			ApplicationArgs: [][]byte{prev, new},
			Accounts:        []basics.Address{addr},
		},
	}, nil
}
