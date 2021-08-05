// Copyright (C) 2019-2021 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

package data

import (
	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/data/transactions"
	"github.com/algorand/go-algorand/ledger"
	"github.com/algorand/go-algorand/protocol"
)

// A SpeculationLedger provides speculative execution on top of a Ledger.
//
// Effectively, a SpeculationLedger adapts a ledger.BlockEvaluator to the
// Ledger interface.
type SpeculationLedger struct {
	baseLedger  *Ledger
	baseRound   basics.Round
	stack       [][]transactions.Transaction
	Checkpoints []uint64

	Evaluator *ledger.BlockEvaluator
	Version   protocol.ConsensusVersion
}

func NewSpeculationLedger(l *Ledger, rnd basics.Round) (*SpeculationLedger, error) {
	sl := &SpeculationLedger{baseLedger: l, baseRound: rnd}
	err := sl.start()
	return sl, err
}

// Note that start() does not manipulate stack or checkpoints, so
// it can be used at construction time and during rollback.
func (sl *SpeculationLedger) start() error {
	hdr, err := sl.baseLedger.BlockHdr(sl.baseRound)
	if err != nil {
		return err
	}
	evaluator, err := sl.baseLedger.StartEvaluator(hdr, 0)
	if err != nil {
		return err
	}
	sl.Evaluator = evaluator
	sl.Version = hdr.CurrentProtocol
	return nil
}

func (sl *SpeculationLedger) GetCreator(cidx basics.CreatableIndex, ctype basics.CreatableType) (basics.Address, bool, error) {
	return sl.Evaluator.State().GetCreator(cidx, ctype)
}

func (sl *SpeculationLedger) Latest() basics.Round {
	return sl.baseRound // or +1 per group? The speculative txns are certainly not in the ledger's round.
}

func (sl *SpeculationLedger) LookupLatest(addr basics.Address) (basics.AccountData, error) {
	return sl.Evaluator.State().Get(addr, true)
}

func (sl *SpeculationLedger) LookupLatestWithoutRewards(addr basics.Address) (basics.AccountData, basics.Round, error) {
	acct, err := sl.Evaluator.State().Get(addr, false)
	// Need to understand what the "validThrough" round returned here should mean
	return acct, basics.Round(0), err
}

func (sl *SpeculationLedger) Apply(txgroup []transactions.Transaction) error {
	var stxgroup []transactions.SignedTxn
	for _, txn := range txgroup {
		stxgroup = append(stxgroup, transactions.SignedTxn{Txn: txn})
	}
	err := sl.Evaluator.TransactionGroup(stxgroup)
	if err != nil {
		return err
	}
	sl.stack = append(sl.stack, txgroup)
	return nil
}

func (sl *SpeculationLedger) ApplyIgnoringSignatures(stxgroup []transactions.SignedTxn) error {
	var txgroup []transactions.Transaction
	for _, txn := range stxgroup {
		txgroup = append(txgroup, txn.Txn)
	}
	return sl.Apply(txgroup)
}

func (sl *SpeculationLedger) Checkpoint() error {
	sl.Checkpoints = append(sl.Checkpoints, uint64(len(sl.stack)))
	return nil
}

func (sl *SpeculationLedger) Rollback() error {
	// Start the evaluator over again from the beginning
	err := sl.start()
	if err != nil {
		return err
	}

	// Replay the txns up until the last checkpoint
	last := len(sl.Checkpoints) - 1
	replays := sl.stack[:sl.Checkpoints[last]]
	sl.stack = nil
	for _, txgroup := range replays {
		err = sl.Apply(txgroup)
		if err != nil {
			return err
		}
	}

	// Discard that checkpoint
	sl.Checkpoints = sl.Checkpoints[:last]
	return nil
}

func (sl *SpeculationLedger) Commit() error {
	last := len(sl.Checkpoints) - 1
	sl.Checkpoints = sl.Checkpoints[:last]
	return nil
}

func (sl *SpeculationLedger) TxStack() [][]transactions.Transaction {
	return sl.stack
}

func (sl *SpeculationLedger) GenesisHash() crypto.Digest {
	return sl.baseLedger.GenesisHash()
}
