package layer2

import (
	"bytes"
	"fmt"
	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/data"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/data/bookkeeping"
	"github.com/algorand/go-algorand/data/transactions"
	"github.com/algorand/go-algorand/layer2/kalgo"
	"github.com/algorand/go-algorand/protocol"
	"github.com/algorand/go-algorand/util"
	"path"
	"strconv"
)

type invocationKind int8

const (
	Init invocationKind = iota
	Call
)

type Invocation struct {
	kind       invocationKind
	name       string
	runner     kalgo.Runner
	contractID crypto.Digest
	sender     basics.Address
}

func NewInitInvocation(name, source string, sender basics.Address) *Invocation {
	sourceHash := crypto.Hash([]byte(source))
	contractID := crypto.Hash([]byte(sender.GetUserAddress() + sourceHash.String()))
	addr, _, _ := logicSigFromTemplateFile("layer2/committee-defer-logicsic.teal.template", contractID)

	return &Invocation{
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

func NewCallInvocation(name, function, args string, addr, sender basics.Address) *Invocation {
	return &Invocation{
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

// An Executor executes Layer 2 invocations, managing the required speculative state
// and ultimately producing the Layer 1 Effects Transactions.
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

// Submit executes the given Invocation and commits to the speculative ledger.
//
// If successful, the effects transactions from the execution are applied to
// the speculative ledger, the speculative database is updated.
func (ex *Executor) Submit(item *Invocation, prof *util.Profiler) error {
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
		err = ex.ledger.Apply([]transactions.Transaction{txn})
		if err != nil {
			return fmt.Errorf("could not apply version swap call (%v -> %v): %v", commit.PreviousCommitment, commit.NewCommitment, err)
		}
	}

	effects := bookkeeping.TxnGroupsFlatten(ex.ledger.TxStack()[ex.lastLedgerIdx:])
	result := &executionResult{
		effectsTxns: effects,
		kalgoOutput: out,
	}
	ex.results = append(ex.results, result)
	ex.lastLedgerIdx = len(ex.ledger.TxStack())

	prof.Start("cow")
	for _, commit := range out.Commitments {
		if err = ex.persistWriteSet(commit.Contract); err != nil {
			return err
		}
	}
	prof.Start("node")
	return nil
}

func (ex *Executor) persistWriteSet(contractName string) error {
	src := path.Join(ex.kenv.SourcePrefix, contractName)
	dst := path.Join(ex.kenv.SourcePrefix, "..", strconv.Itoa(len(ex.results)), contractName)
	return util.CopyFolder(src, dst)
}

// EffectsTxns returns the (unsigned) Effects Transactions produced so far.
func (ex *Executor) EffectsTxns() [][]transactions.Transaction {
	var txns [][]transactions.Transaction
	for _, res := range ex.results {
		txns = append(txns, res.effectsTxns)
	}
	return txns
}

func CommitmentAppOptInTxn(contractID crypto.Digest, commitment []byte, fv, lv basics.Round, genesisHash crypto.Digest) (transactions.Transaction, error) {
	return transactions.Transaction{
		Type: protocol.ApplicationCallTx,
		Header: transactions.Header{
			Sender:      ContractAddress(contractID),
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
	addr := ContractAddress(contractID)
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
