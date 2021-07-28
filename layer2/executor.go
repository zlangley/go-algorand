package layer2

import (
	"fmt"
	"github.com/algorand/go-algorand/data/transactions"
	"github.com/algorand/go-algorand/data/transactions/logic"
	"github.com/algorand/go-algorand/protocol"
	"github.com/algorand/go-algorand/util"
	"path"
	"strconv"
	"strings"
	"text/template"

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
	kind       batchItemKind
	runner     kalgo.Runner
	contractID basics.Address
	sender     basics.Address
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
		contractID: addr,
		sender:     sender,
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
		contractID: addr,
		sender:     sender,
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

	// TODO: Run() should probably just return a *kalgo.Output, but for now
	// they are separate to make profiling more contained.
	prof.Start("node")
	out, err := kalgo.ParseOutput(rawout)
	if err != nil {
		return fmt.Errorf("parsing kalgo: %w", err)
	}

	for _, commit := range out.Commitments {
		if commit.PreviousCommitment == "" {
			err = ex.addCommitmentAppOptInTxn([]byte(commit.NewCommitment), item.contractID)
		} else {
			err = ex.addCommitmentAppCallTxn([]byte(commit.PreviousCommitment), []byte(commit.NewCommitment), item.contractID)
		}
		if err != nil {
			return err
		}

		prof.Start("cow")
		if err = ex.copyContract(commit.Contract); err != nil {
			return err
		}
		prof.Start("node")
	}
	ex.ledger.CommitStack()
	return nil
}

func (ex *Executor) copyContract(contract string) error {
	src := path.Join(ex.kenv.SourcePrefix, contract)
	dst := path.Join(ex.kenv.SourcePrefix, "..", strconv.Itoa(len(ex.ledger.TransactionBatch())), contract)
	return util.CopyFolder(src, dst)
}

// associated with every l2 contract, we have:
//   - contract ID
//   - escrow account (logicsig)

func logicSigFromTemplateFile(filename string, contractID basics.Address) (basics.Address, []byte, error) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		return basics.Address{}, nil, fmt.Errorf("could not parse template file: %v", err)
	}
	b := &strings.Builder{}
	err = t.Execute(b, contractID)
	if err != nil {
		return basics.Address{}, nil, fmt.Errorf("could not execute template: %v", err)
	}
	ops, err := logic.AssembleString(b.String())
	if err != nil {
		return basics.Address{}, nil, fmt.Errorf("could not assemble TEAL: %v", b)
	}
	addr := basics.Address(logic.HashProgram(ops.Program))
	return addr, ops.Program, err
}

var ForwardDeclareAppIndex basics.AppIndex = 20189847
var VersionNumberAppIndex basics.AppIndex = 19841517

//var CommitmentMapAddress = "Z72YSHBKWQMW66SGB6WHH6VOK46KDSTJZ2NZJPRXQG4RK6U2Y62ZNR6FJY"

func (ex *Executor) addCommitmentAppOptInTxn(commitment []byte, contractID basics.Address) error {
	addr, prog, err := logicSigFromTemplateFile("layer2/optin-logicsig.teal.template", contractID)
	if err != nil {
		return err
	}
	tx := transactions.Transaction{
		Type: protocol.ApplicationCallTx,
		Header: transactions.Header{
			Sender:      addr,
			Fee:         basics.MicroAlgos{Raw: 1000},
			FirstValid:  ex.ledger.Latest(),
			LastValid:   ex.ledger.Latest() + 1000,
			GenesisHash: ex.ledger.GenesisHash(),
			Note:        commitment,
		},
		ApplicationCallTxnFields: transactions.ApplicationCallTxnFields{
			ApplicationID:   VersionNumberAppIndex,
			ApplicationArgs: [][]byte{commitment},
			OnCompletion:    transactions.OptInOC,
		},
	}
	stx := transactions.SignedTxn{
		Txn:  tx,
		Lsig: transactions.LogicSig{Logic: prog},
	}
	err = ex.ledger.Apply([]transactions.SignedTxn{stx})
	if err != nil {
		return fmt.Errorf("could not apply transaction: %v", err)
	}
	return nil
}

func (ex *Executor) addCommitmentAppCallTxn(prev, new []byte, contractID basics.Address) error {
	addr, committeeDeferProg, err := logicSigFromTemplateFile("layer2/committee-defer-logicsic.teal.template", contractID)
	if err != nil {
		return err
	}
	updateCommitmentTxn := transactions.Transaction{
		Type: protocol.ApplicationCallTx,
		Header: transactions.Header{
			Sender:      addr,
			Fee:         basics.MicroAlgos{Raw: 1000},
			FirstValid:  ex.ledger.Latest(),
			LastValid:   ex.ledger.Latest() + 1000,
			GenesisHash: ex.ledger.GenesisHash(),
		},
		ApplicationCallTxnFields: transactions.ApplicationCallTxnFields{
			ApplicationID:   VersionNumberAppIndex,
			ApplicationArgs: [][]byte{prev, new},
			Accounts:        []basics.Address{addr},
		},
	}
	updateCommitmentStxn := transactions.SignedTxn{
		Txn:  updateCommitmentTxn,
		Lsig: transactions.LogicSig{Logic: committeeDeferProg},
	}
	authCheckOps, err := logic.AssembleString(`#pragma version 4
bytecblock b32 L5SP2W7BEXYQEHIMNQHCKOEGXYBAG56T5VVAFYOR3QDG5XHQ5D6B57BRKA
int 1
`)
	authCheckSender := basics.Address(logic.HashProgram(authCheckOps.Program))
	authCheckStx := transactions.SignedTxn{
		Txn: transactions.Transaction{
			Type: protocol.ApplicationCallTx,
			Header: transactions.Header{
				Sender:      authCheckSender,
				Fee:         basics.MicroAlgos{Raw: 1000},
				FirstValid:  ex.ledger.Latest(),
				LastValid:   ex.ledger.Latest() + 1000,
				GenesisHash: ex.ledger.GenesisHash(),
			},
			ApplicationCallTxnFields: transactions.ApplicationCallTxnFields{
				ApplicationID:   ForwardDeclareAppIndex,
				ApplicationArgs: [][]byte{[]byte("extract")},
			},
		},
		Lsig: transactions.LogicSig{Logic: authCheckOps.Program},
	}
	err = ex.ledger.Apply([]transactions.SignedTxn{authCheckStx, updateCommitmentStxn})
	if err != nil {
		return err //fmt.Errorf("could not apply transaction (%v -> %v): %v: %v", base64.StdEncoding.EncodeToString(prev), base64.StdEncoding.EncodeToString(new), tx, err)
	}
	return nil
}
