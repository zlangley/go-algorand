package layer2

import (
	"encoding/base64"
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
	contractID basics.Address
	sender     basics.Address
}

func GetContractPreID(sender basics.Address, source string) crypto.Digest {
	sourceHash := crypto.Hash([]byte(source))
	return crypto.Hash([]byte(sender.GetUserAddress() + sourceHash.String()))
}

func NewInitInvocation(name, source string, contractAddr, sender basics.Address) *Invocation {
	return &Invocation{
		kind: Init,
		name: name,
		runner: &kalgo.InitCmd{
			Cmd: kalgo.Cmd{
				Address: contractAddr,
				Name:    name,
				Sender:  sender,
			},
			Source: source,
		},
		contractID: contractAddr,
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
}

func NewExecutor(ledger *data.SpeculationLedger, kenv kalgo.Env) *Executor {
	return &Executor{
		ledger: ledger,
		kenv:   kenv,
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

	// TODO: Run() should probably just return a *kalgo.Output, but for now, running
	// the VM and parsing its output are separate to contain the profiling logic.
	out, err := kalgo.ParseOutput(rawout)
	if err != nil {
		return fmt.Errorf("parsing kalgo: %w", err)
	}

	prof.Start("effects")
	// Add version updates.
	for _, commit := range out.Commitments {
		// TODO: use acutal ContractID once we get it from kalgo output.
		contractID := crypto.Hash([]byte(commit.Contract))
		contractAddr := GetContractAddress(contractID)

		var txn transactions.Transaction
		if item.kind == Init {
			txn, err = CommitmentAppOptInTxn(contractAddr, commit.NewCommitment, ex.ledger.Latest(), ex.ledger.Latest() + 1000, ex.ledger.GenesisHash())
		} else {
			txn, err = CommitmentAppCallTxn(contractAddr, commit.PreviousCommitment, commit.NewCommitment, ex.ledger.Latest(), ex.ledger.Latest()+1000, ex.ledger.GenesisHash())
		}
		if err != nil {
			return fmt.Errorf("could not create commitment swap contract call: %v", err)
		}
		err = ex.ledger.Apply([]transactions.Transaction{txn})
		if err != nil {
			return fmt.Errorf("could not apply version swap call (%v -> %v) for contract %v: %v", base64.StdEncoding.EncodeToString(commit.PreviousCommitment), base64.StdEncoding.EncodeToString(commit.NewCommitment), commit.Contract, err)
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

func CommitmentAppOptInTxn(contractAddr basics.Address, commitment []byte, fv, lv basics.Round, genesisHash crypto.Digest) (transactions.Transaction, error) {
	return transactions.Transaction{
		Type: protocol.ApplicationCallTx,
		Header: transactions.Header{
			Sender:      contractAddr,
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

func CommitmentAppCallTxn(contractAddr basics.Address, prev, new []byte, fv, lv basics.Round, genesisHash crypto.Digest) (transactions.Transaction, error) {
	return transactions.Transaction{
		Type: protocol.ApplicationCallTx,
		Header: transactions.Header{
			Sender:      contractAddr,
			Fee:         basics.MicroAlgos{Raw: 1000},
			FirstValid:  fv,
			LastValid:   lv,
			GenesisHash: genesisHash,
		},
		ApplicationCallTxnFields: transactions.ApplicationCallTxnFields{
			ApplicationID:   VersionNumberAppIndex,
			ApplicationArgs: [][]byte{prev, new},
			Accounts:        []basics.Address{contractAddr},
		},
	}, nil
}
