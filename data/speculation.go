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
	"time"

	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/data/bookkeeping"
	"github.com/algorand/go-algorand/data/transactions"
	"github.com/algorand/go-algorand/ledger"
	"github.com/algorand/go-algorand/protocol"
)

// A SpeculationLedger adapts a BlockEvaluator to the Ledger interface
// (and provides access to the BlockEvalutor's ability to execute
// trasnactions) This means we can write code that expects a Ledger to
// report on balances and such as we go.

type SpeculationLedger struct {
	baseLedger      *Ledger
	baseRound       basics.Round
	workingTxGroups [][]transactions.SignedTxn
	stagedTxGroups  [][]transactions.SignedTxn
	Checkpoints     []uint64

	Evaluator *ledger.BlockEvaluator
	Version   protocol.ConsensusVersion
}

func NewSpeculationLedger(l *Ledger, rnd basics.Round) (*SpeculationLedger, error) {
	sl := &SpeculationLedger{baseLedger: l, baseRound: rnd}
	err := sl.start()
	return sl, err
}

// Note that start() does not manipulate workingTxGroups or checkpoints, so
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
	for _, txgroup := range sl.stagedTxGroups {
		err = evaluator.TransactionGroup(txgroup)
		if err != nil {
			return err
		}
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
func (sl *SpeculationLedger) Apply(txgroup []transactions.SignedTxn) error {
	err := sl.Evaluator.TransactionGroup(txgroup)
	if err != nil {
		return err
	}
	sl.workingTxGroups = append(sl.workingTxGroups, txgroup)
	return nil
}

func (sl *SpeculationLedger) Checkpoint() error {
	sl.Checkpoints = append(sl.Checkpoints, uint64(len(sl.workingTxGroups)))
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
	replays := sl.workingTxGroups[:sl.Checkpoints[last]]
	sl.workingTxGroups = nil
	for _, txgroup := range replays {
		err := sl.Apply(txgroup)
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

func (sl *SpeculationLedger) StageWorking() error {
	err := sl.start()
	if err != nil {
		return err
	}

	var group transactions.TxGroup
	stxns := bookkeeping.SignedTxnGroupsFlatten(sl.workingTxGroups)
	for _, stxn := range stxns {
		group.TxGroupHashes = append(group.TxGroupHashes, crypto.HashObj(stxn.Txn))
	}
	groupHash := crypto.HashObj(group)
	for i := range stxns {
		stxns[i].Txn.Group = groupHash
	}
	err = sl.Evaluator.TransactionGroup(stxns)
	if err != nil {
		return err
	}
	sl.stagedTxGroups = append(sl.stagedTxGroups, stxns)
	sl.workingTxGroups = nil
	return nil
}

type Profiler struct {
	curr    string
	watches map[string]*Stopwatch
}

func NewProfiler() *Profiler {
	return &Profiler{
		curr: "default",
		watches: map[string]*Stopwatch{
			"default": &Stopwatch{},
		},
	}
}

func (prof *Profiler) Start(key string) {
	if key == prof.curr {
		return
	}
	prof.Stop()
	prof.curr = key
	if _, ok := prof.watches[key]; !ok {
		prof.watches[key] = &Stopwatch{}
	}
	prof.watches[key].Start()
}

func (prof Profiler) Stop() {
	prof.watches[prof.curr].Stop()
}

func (prof Profiler) Elapsed(key string) uint64 {
	if key == prof.curr {
		prof.Stop()
		prof.Start(key)
	}
	sw, ok := prof.watches[key]
	if !ok {
		return 0
	}
	return uint64(sw.Elapsed().Nanoseconds())
}

type Stopwatch struct {
	startTime time.Time
	elapsed   time.Duration
	running   bool
}

func (sw *Stopwatch) Start() {
	sw.running = true
	sw.startTime = time.Now()
}

func (sw *Stopwatch) Stop() {
	elapsed := time.Since(sw.startTime)
	if !sw.running {
		return
	}
	sw.elapsed += elapsed
	sw.running = false
}

func (sw *Stopwatch) Elapsed() time.Duration {
	return sw.elapsed
}
