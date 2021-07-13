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

package v2

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/algorand/go-algorand/config"
	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/daemon/algod/api/server/v2/generated"
	"github.com/algorand/go-algorand/daemon/algod/api/server/v2/generated/private"
	"github.com/algorand/go-algorand/data"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/data/bookkeeping"
	"github.com/algorand/go-algorand/data/transactions"
	"github.com/algorand/go-algorand/data/transactions/logic"
	"github.com/algorand/go-algorand/kalgo"
	"github.com/algorand/go-algorand/ledger/ledgercore"
	"github.com/algorand/go-algorand/logging"
	"github.com/algorand/go-algorand/node"
	"github.com/algorand/go-algorand/protocol"
	"github.com/algorand/go-algorand/rpcs"
	"github.com/algorand/go-algorand/util"
)

const maxTealSourceBytes = 1e5
const maxTealDryrunBytes = 1e5
const maxAlgoClaritySourceBytes = 1e5
const maxAlgoClarityBatchBytes = 1e6

var prof *util.Profiler

// Handlers is an implementation to the V2 route handler interface defined by the generated code.
type Handlers struct {
	Node     NodeInterface
	Log      logging.Logger
	Shutdown <-chan struct{}
}

// ledgerForApiHandlers represents the subset of Ledger functionality
// needed by the handler functions. It exists so that a
// SpeculationLedger can implement it and substitute in.
type ledgerForApiHandlers interface {
	GetCreator(cidx basics.CreatableIndex, ctype basics.CreatableType) (basics.Address, bool, error)
	Latest() basics.Round
	LookupLatest(addr basics.Address) (basics.AccountData, error)
	LookupLatestWithoutRewards(addr basics.Address) (basics.AccountData, basics.Round, error)
}

// NodeInterface represents node fns used by the handlers.
type NodeInterface interface {
	Ledger() *data.Ledger
	NewSpeculationLedger(rnd basics.Round) (string, error)
	SpeculationLedger(token string) (*data.SpeculationLedger, error)
	DestroySpeculationLedger(token string)
	Status() (s node.StatusReport, err error)
	GenesisID() string
	GenesisHash() crypto.Digest
	BroadcastSignedTxGroup(txgroup []transactions.SignedTxn) error
	GetPendingTransaction(txID transactions.Txid) (res node.TxnWithStatus, found bool)
	GetPendingTxnsFromPool() ([]transactions.SignedTxn, error)
	SuggestedFee() basics.MicroAlgos
	StartCatchup(catchpoint string) error
	AbortCatchup(catchpoint string) error
	Config() config.Local
}

// RegisterParticipationKeys registers participation keys.
// (POST /v2/register-participation-keys/{address})
func (v2 *Handlers) RegisterParticipationKeys(ctx echo.Context, address string, params private.RegisterParticipationKeysParams) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	// TODO: register participation keys endpoint
	return ctx.String(http.StatusNotImplemented, "Endpoint not implemented.")
}

// ShutdownNode shuts down the node.
// (POST /v2/shutdown)
func (v2 *Handlers) ShutdownNode(ctx echo.Context, params private.ShutdownNodeParams) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	// TODO: shutdown endpoint
	return ctx.String(http.StatusNotImplemented, "Endpoint not implemented.")
}

// AccountInformation gets account information for a given account.
// (GET /v2/accounts/{address})
func (v2 *Handlers) AccountInformation(ctx echo.Context, address string, params generated.AccountInformationParams) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)

	handle, contentType, err := getCodecHandle(params.Format)
	if err != nil {
		return badRequest(ctx, err, errFailedParsingFormatOption, v2.Log)
	}

	addr, err := basics.UnmarshalChecksumAddress(address)
	if err != nil {
		return badRequest(ctx, err, errFailedToParseAddress, v2.Log)
	}

	var ledger ledgerForApiHandlers
	if params.Speculation == nil {
		ledger = v2.Node.Ledger()
	} else {
		ledger, err = v2.Node.SpeculationLedger(*params.Speculation)
		if err != nil {
			return badRequest(ctx, err, errFailedLookingUpLedger, v2.Log)
		}
	}
	record, err := ledger.LookupLatest(addr)
	if err != nil {
		return internalError(ctx, err, errFailedLookingUpLedger, v2.Log)
	}

	if handle == protocol.CodecHandle {
		data, err := encode(handle, record)
		if err != nil {
			return internalError(ctx, err, errFailedToEncodeResponse, v2.Log)
		}
		return ctx.Blob(http.StatusOK, contentType, data)
	}

	recordWithoutPendingRewards, _, err := ledger.LookupLatestWithoutRewards(addr)
	if err != nil {
		return internalError(ctx, err, errFailedLookingUpLedger, v2.Log)
	}
	amountWithoutPendingRewards := recordWithoutPendingRewards.MicroAlgos

	assetsCreators := make(map[basics.AssetIndex]string, len(record.Assets))
	if len(record.Assets) > 0 {
		//assets = make(map[uint64]v1.AssetHolding)
		for curid := range record.Assets {
			var creator string
			creatorAddr, ok, err := ledger.GetCreator(basics.CreatableIndex(curid), basics.AssetCreatable)
			if err == nil && ok {
				creator = creatorAddr.String()
			} else {
				// Asset may have been deleted, so we can no
				// longer fetch the creator
				creator = ""
			}
			assetsCreators[curid] = creator
		}
	}

	account, err := AccountDataToAccount(address, &record, assetsCreators, ledger.Latest(), amountWithoutPendingRewards)
	if err != nil {
		return internalError(ctx, err, errInternalFailure, v2.Log)
	}

	response := generated.AccountResponse(account)
	return ctx.JSON(http.StatusOK, response)
}

// GetBlock gets the block for the given round.
// (GET /v2/blocks/{round})
func (v2 *Handlers) GetBlock(ctx echo.Context, round uint64, params generated.GetBlockParams) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	handle, contentType, err := getCodecHandle(params.Format)
	if err != nil {
		return badRequest(ctx, err, errFailedParsingFormatOption, v2.Log)
	}

	// msgpack format uses 'RawBlockBytes' and attaches a custom header.
	if handle == protocol.CodecHandle {
		blockbytes, err := rpcs.RawBlockBytes(v2.Node.Ledger(), basics.Round(round))
		if err != nil {
			return internalError(ctx, err, err.Error(), v2.Log)
		}

		ctx.Response().Writer.Header().Add("X-Algorand-Struct", "block-v1")
		return ctx.Blob(http.StatusOK, contentType, blockbytes)
	}

	ledger := v2.Node.Ledger()
	block, _, err := ledger.BlockCert(basics.Round(round))
	if err != nil {
		return internalError(ctx, err, errFailedLookingUpLedger, v2.Log)
	}

	// Encoding wasn't working well without embedding "real" objects.
	response := struct {
		Block bookkeeping.Block `codec:"block"`
	}{
		Block: block,
	}

	data, err := encode(handle, response)
	if err != nil {
		return internalError(ctx, err, errFailedToEncodeResponse, v2.Log)
	}

	return ctx.Blob(http.StatusOK, contentType, data)
}

// Create a speculation context starting at the given block.
// (POST /v2/blocks/{round}/speculation)
func (v2 *Handlers) CreateSpeculation(ctx echo.Context, round uint64) error {
	prof = util.NewProfiler()
	prof.Start(kNode)
	defer prof.Stop()
	if round == 0 {
		round = uint64(v2.Node.Ledger().Latest())
	}
	token, err := v2.Node.NewSpeculationLedger(basics.Round(round))
	if err != nil {
		return badRequest(ctx, err, fmt.Sprintf("%v", err), v2.Log)
	}

	prof.Start(kCopyOnWrite)
	current := filepath.Join(os.Getenv("ALGO_CLARITY_PREFIX"), "current")
	if _, err = os.Stat(current); os.IsNotExist(err) {
		err = os.MkdirAll(current, 0777)
		if err != nil {
			return internalError(ctx, err, err.Error(), v2.Log)
		}
	}
	kenv := kalgoEnv(ctx.Request(), token)
	if err = dircopy(current, kenv.SourcePrefix); err != nil {
		return internalError(ctx, err, err.Error(), v2.Log)
	}
	return ctx.JSON(http.StatusOK, generated.SpeculationResponse{
		Base:  round,
		Token: token,
	})
}

func kalgoEnv(req *http.Request, speculation string) kalgo.Env {
	return kalgo.Env{
		AlgodAddress:     req.Context().Value(http.LocalAddrContextKey).(net.Addr).String(),
		AlgodToken:       req.Header.Get("X-Algo-API-Token"),
		SpeculationToken: speculation,
		SourcePrefix:     filepath.Join(os.Getenv("ALGO_CLARITY_PREFIX"), speculation, "current"),
	}
}

var (
	kNode        = "node"
	kKalgoTotal  = "kalgo"
	kCopyOnWrite = "cow"
)

func createContract(name, source string, ledger *data.SpeculationLedger, kenv kalgo.Env, sender basics.Address) (*kalgo.Output, basics.Address, error) {
	addr := basics.Address(crypto.Hash([]byte(source)))

	cmd := &kalgo.InitCmd{
		Cmd: kalgo.Cmd{
//			Address: addr,
			Name:    name,
			Sender:  sender,
		},
		Source: source,
	}

	out, err := executeVM(kenv, cmd, ledger)
	if err != nil {
		return nil, addr, err
	}

	/*
		tx := transactions.Transaction{
			Type: protocol.ApplicationCallTx,
			Header: transactions.Header{
				Sender:      addr,
				Fee:         basics.MicroAlgos{Raw: 10000},
				FirstValid:  ledger.Latest(),
				LastValid:   ledger.Latest() + 1000,
				//GenesisHash: ledger.GenesisHash,
			},
			ApplicationCallTxnFields: transactions.ApplicationCallTxnFields{
				ApplicationID: basics.AppIndex(-1), // todo
				ApplicationArgs: [][]byte{[]byte(""), []byte(out.Commitments[0].NewCommitment)},
			},
		}
		// sign from sender
		ledger.Apply()
	*/

	return out, addr, nil
}

func callContract(addr basics.Address, name, function, args string, sender basics.Address, kenv kalgo.Env, ledger *data.SpeculationLedger) (*kalgo.Output, error) {
	cmd := &kalgo.CallCmd{
		Cmd: kalgo.Cmd{
//			Address: addr,
			Name:    name,
			Sender:  sender,
		},
		Function: function,
		Args:     args,
	}
	return executeVM(kenv, cmd, ledger)
}

// CreateContract creates an AlgoClarity contract.
// (POST /v2/contract/{name})
func (v2 *Handlers) CreateContract(ctx echo.Context, name string, params generated.CreateContractParams) error {
	prof.Start(kNode)
	defer prof.Stop()

	if params.Speculation == nil {
		err := errors.New("speculation token required (for now)")
		return badRequest(ctx, err, err.Error(), v2.Log)
	}
	ledger, err := v2.Node.SpeculationLedger(*params.Speculation)
	if err != nil {
		return badRequest(ctx, err, errFailedLookingUpLedger, v2.Log)
	}
	buf := new(bytes.Buffer)
	ctx.Request().Body = http.MaxBytesReader(nil, ctx.Request().Body, maxAlgoClaritySourceBytes)
	buf.ReadFrom(ctx.Request().Body)

	var sender basics.Address
	if params.Sender != nil {
		err = sender.UnmarshalText([]byte(*params.Sender))
		if err != nil {
			return badRequest(ctx, err, errFailedToParseAddress, v2.Log)
		}
	}

	kenv := kalgoEnv(ctx.Request(), *params.Speculation)
	out, addr, err := createContract(name, buf.String(), ledger, kenv, sender)
	if err != nil {
		return internalError(ctx, err, err.Error(), v2.Log)
	}
	return ctx.JSON(http.StatusOK, generated.SpeculationResponse{
		Base:        uint64(ledger.Latest()),
		Checkpoints: &ledger.Checkpoints,
		Token:       *params.Speculation,
		Commitments: out.Commitments,
		Address:     addr.GetUserAddress(),
		TotalTime:   prof.ElapsedTotal(),
		NodeTime:    prof.Elapsed(kNode),
		KalgoTime:   prof.Elapsed(kKalgoTotal),
		DbTime:      prof.Elapsed(kCopyOnWrite),
	})
}

// Calls a function on a previously initialized contract.
// (POST /v2/contracts/{name}/call/{function})
func (v2 *Handlers) CallContract(ctx echo.Context, name string, function string, params generated.CallContractParams) error {
	prof.Start(kNode)
	defer prof.Stop()
	if params.Speculation == nil {
		err := errors.New("speculation token required (for now)")
		return badRequest(ctx, err, err.Error(), v2.Log)
	}
	speculation := *params.Speculation
	ledger, err := v2.Node.SpeculationLedger(speculation)
	if err != nil {
		return badRequest(ctx, err, errFailedLookingUpLedger, v2.Log)
	}

	var addr basics.Address
	if params.Address != nil {
		err = addr.UnmarshalText([]byte(*params.Address))
		if err != nil {
			return badRequest(ctx, err, errFailedToParseAddress, v2.Log)
		}
	}

	var sender basics.Address
	if params.Sender != nil {
		err = sender.UnmarshalText([]byte(*params.Sender))
		if err != nil {
			return badRequest(ctx, err, errFailedToParseAddress, v2.Log)
		}
	}

	args := ""
	if params.Args != nil {
		args = *params.Args
	}
	kenv := kalgoEnv(ctx.Request(), speculation)
	out, err := callContract(addr, name, function, args, sender, kenv, ledger)
	if err != nil {
		return internalError(ctx, err, err.Error(), v2.Log)
	}
	return ctx.JSON(http.StatusOK, generated.SpeculationResponse{
		Base:        uint64(ledger.Latest()),
		Checkpoints: &ledger.Checkpoints,
		Token:       speculation,
		Commitments: out.Commitments,
		TotalTime:   prof.ElapsedTotal(),
		NodeTime:    prof.Elapsed(kNode),
		KalgoTime:   prof.Elapsed(kKalgoTotal),
		DbTime:      prof.Elapsed(kCopyOnWrite),
	})
}

type cmdDiscriminator struct {
	Command string `json:"command"`
}

func decodeBatch(data []byte) ([]interface{}, error) {
	// FIXME[zach]: I haven't been able to get go-swagger to generate the right interface for a
	// heterogeneous list (perhaps we are using an old version?). The parsing here is kind of gnarly,
	// but is intended to somewhat mirror what go-swagger supposedly expects. We decode the data twice.
	// First we decode using the `cmdDiscriminator` helper, which just extracts the "command" discriminator
	// field. This should be a bit better than just decoding to map[interface{}]interface{}, which would
	// not enforce the presence of the "command" key. The second decoding just decodes the array, and leaves
	// the items as json.RawMessage which are then decoded dynamically based on the "command" discriminator.
	var discrims []cmdDiscriminator
	if err := decode(protocol.JSONUnstrictHandle, data, &discrims); err != nil {
		return nil, err
	}
	var rawCmds []json.RawMessage
	if err := decode(protocol.JSONHandle, data, &rawCmds); err != nil {
		return nil, err
	}
	ret := make([]interface{}, len(rawCmds))
	for i, discrim := range discrims {
		switch discrim.Command {
		case "init":
			var init generated.ContractInit
			if err := decode(protocol.JSONHandle, rawCmds[i], &init); err != nil {
				return nil, err
			}
			ret[i] = init
		case "call":
			var call generated.ContractCall
			if err := decode(protocol.JSONHandle, rawCmds[i], &call); err != nil {
				return nil, err
			}
			ret[i] = call
		default:
			return nil, errors.New(fmt.Sprintf("item in batch missing command descriminator (index %d)", i))
		}
	}
	return ret, nil
}

// Calls a function on a previously initialized contract.
// (POST /v2/contracts/batch)
func (v2 *Handlers) ContractBatchExecute(ctx echo.Context, params generated.ContractBatchExecuteParams) error {
	prof = util.NewProfiler()
	prof.Start(kNode)
	defer prof.Stop()
	if params.Speculation == nil {
		err := errors.New("speculation token required (for now)")
		return badRequest(ctx, err, err.Error(), v2.Log)
	}
	speculation := *params.Speculation
	ledger, err := v2.Node.SpeculationLedger(speculation)
	if err != nil {
		return badRequest(ctx, err, errFailedLookingUpLedger, v2.Log)
	}

	req := ctx.Request()
	buf := new(bytes.Buffer)
	req.Body = http.MaxBytesReader(nil, req.Body, maxAlgoClarityBatchBytes)
	buf.ReadFrom(req.Body)
	data := buf.Bytes()

	gcmds, err := decodeBatch(data)
	if err != nil {
		return badRequest(ctx, err, err.Error(), v2.Log)
	}

	kenv := kalgoEnv(ctx.Request(), speculation)
	var commits []generated.ContractCommitment
	for _, gcmd := range gcmds {
		prof.Start(kNode)
		var out *kalgo.Output
		if init, ok := gcmd.(generated.ContractInit); ok {
			var sender basics.Address
			if init.Sender != nil {
				err = sender.UnmarshalText([]byte(*init.Sender))
				if err != nil {
					return badRequest(ctx, err, errFailedToParseAddress, v2.Log)
				}
			}
			out, _, err = createContract(init.Id, init.Source, ledger, kenv, sender)
		} else if call, ok := gcmd.(generated.ContractCall); ok {
			var addr, sender basics.Address
			if call.Sender != nil {
				sender, err = basics.UnmarshalChecksumAddress(*call.Sender)
				if err != nil {
					return badRequest(ctx, err, errFailedToParseAddress, v2.Log)
				}
			}
			if call.Address != nil {
				addr, err = basics.UnmarshalChecksumAddress(*call.Address)
				if err != nil {
					return badRequest(ctx, err, errFailedToParseAddress, v2.Log)
				}
			}
			out, err = callContract(addr, call.Id, call.Function, call.Args, sender, kenv, ledger)
		} else {
			err = errors.New("unknown command type returned from decodeBatch()")
			return internalError(ctx, err, err.Error(), v2.Log)
		}
		if err != nil {
			return internalError(ctx, err, err.Error(), v2.Log)
		}
		commits = append(commits, out.Commitments...)
	}
	return ctx.JSON(http.StatusOK, generated.SpeculationResponse{
		Base:        uint64(ledger.Latest()),
		Checkpoints: &ledger.Checkpoints,
		Commitments: commits,
		Token:       speculation,
		TotalTime:   prof.ElapsedTotal(),
		NodeTime:    prof.Elapsed(kNode),
		KalgoTime:   prof.Elapsed(kKalgoTotal),
		DbTime:      prof.Elapsed(kCopyOnWrite),
	})
}

func executeVM(kenv kalgo.Env, cmd kalgo.Runner, ledger *data.SpeculationLedger) (*kalgo.Output, error) {
	prof.Start(kKalgoTotal)
	rawout, err := cmd.Run(kenv)
	if err != nil {
		return nil, err
	}

	prof.Start(kNode)
	out, err := kalgo.ParseOutput(rawout)
	if err != nil {
		return nil, err
	}

	if err = ledger.Checkpoint(); err != nil {
		return nil, err
	}

	prof.Start(kCopyOnWrite)
	for _, commit := range out.Commitments {
		if commit.PreviousCommitment == commit.NewCommitment {
			continue
		}
		src := path.Join(kenv.SourcePrefix, commit.Contract)
		dst := path.Join(kenv.SourcePrefix, "..", strconv.Itoa(len(ledger.Checkpoints)), commit.Contract)
		if err = dircopy(src, dst); err != nil {
			return nil, err
		}
	}
	prof.Start(kNode)
	return out, nil
}

// Perform operations on a speculation object.
// (POST /v2/speculation/{token}/{operation})
func (v2 *Handlers) SpeculationOperation(ctx echo.Context, speculation string, operation string) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	if operation == "delete" {
		v2.Node.DestroySpeculationLedger(speculation)
		os.RemoveAll(filepath.Join(os.Getenv("ALGO_CLARITY_PREFIX"), speculation))
		// XXX: return something more reasonable
		return ctx.JSON(http.StatusOK, generated.SpeculationResponse{
			Base:  0,
			Token: speculation,
		})
	}
	ledger, err := v2.Node.SpeculationLedger(speculation)
	if err != nil {
		return badRequest(ctx, err, errFailedLookingUpLedger, v2.Log)
	}
	if operation == "checkpoint" {
		err := ledger.Checkpoint()
		if err != nil {
			return badRequest(ctx, err, err.Error(), v2.Log)
		}
		return ctx.JSON(http.StatusOK, generated.SpeculationResponse{
			Base:        uint64(ledger.Latest()),
			Checkpoints: &ledger.Checkpoints,
			Token:       speculation,
		})
	}
	if operation == "rollback" {
		err := ledger.Rollback()
		if err != nil {
			return badRequest(ctx, err, err.Error(), v2.Log)
		}
		return ctx.JSON(http.StatusOK, generated.SpeculationResponse{
			Base:        uint64(ledger.Latest()),
			Checkpoints: &ledger.Checkpoints,
			Token:       speculation,
		})
	}
	if operation == "commit" {
		ledger.Commit()
		if err != nil {
			return badRequest(ctx, err, err.Error(), v2.Log)
		}
		return ctx.JSON(http.StatusOK, generated.SpeculationResponse{
			Base:        uint64(ledger.Latest()),
			Checkpoints: &ledger.Checkpoints,
			Token:       speculation,
		})
	}
	message := fmt.Sprintf("unknown operation '%s'", operation)
	return badRequest(ctx, errors.New(message), message, v2.Log)
}

// GetProof generates a Merkle proof for a transaction in a block.
// (GET /v2/blocks/{round}/transactions/{txid}/proof)
func (v2 *Handlers) GetProof(ctx echo.Context, round uint64, txid string, params generated.GetProofParams) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	var txID transactions.Txid
	err := txID.UnmarshalText([]byte(txid))
	if err != nil {
		return badRequest(ctx, err, errNoTxnSpecified, v2.Log)
	}

	ledger := v2.Node.Ledger()
	block, _, err := ledger.BlockCert(basics.Round(round))
	if err != nil {
		return internalError(ctx, err, errFailedLookingUpLedger, v2.Log)
	}

	proto := config.Consensus[block.CurrentProtocol]
	if proto.PaysetCommit != config.PaysetCommitMerkle {
		return notFound(ctx, err, "protocol does not support Merkle proofs", v2.Log)
	}

	txns, err := block.DecodePaysetFlat()
	if err != nil {
		return internalError(ctx, err, "decoding transactions", v2.Log)
	}

	for idx := range txns {
		if txns[idx].Txn.ID() == txID {
			tree, err := block.TxnMerkleTree()
			if err != nil {
				return internalError(ctx, err, "building Merkle tree", v2.Log)
			}

			proof, err := tree.Prove([]uint64{uint64(idx)})
			if err != nil {
				return internalError(ctx, err, "generating proof", v2.Log)
			}

			proofconcat := make([]byte, 0)
			for _, proofelem := range proof {
				proofconcat = append(proofconcat, proofelem[:]...)
			}

			stibhash := block.Payset[idx].Hash()

			response := generated.ProofResponse{
				Proof:    proofconcat,
				Stibhash: stibhash[:],
				Idx:      uint64(idx),
			}

			return ctx.JSON(http.StatusOK, response)
		}
	}

	err = errors.New(errTransactionNotFound)
	return notFound(ctx, err, err.Error(), v2.Log)
}

// GetSupply gets the current supply reported by the ledger.
// (GET /v2/ledger/supply)
func (v2 *Handlers) GetSupply(ctx echo.Context) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	latest := v2.Node.Ledger().Latest()
	totals, err := v2.Node.Ledger().Totals(latest)
	if err != nil {
		err = fmt.Errorf("GetSupply(): round %d, failed: %v", latest, err)
		return internalError(ctx, err, errInternalFailure, v2.Log)
	}

	supply := generated.SupplyResponse{
		CurrentRound: uint64(latest),
		TotalMoney:   totals.Participating().Raw,
		OnlineMoney:  totals.Online.Money.Raw,
	}

	return ctx.JSON(http.StatusOK, supply)
}

// GetStatus gets the current node status.
// (GET /v2/status)
func (v2 *Handlers) GetStatus(ctx echo.Context) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	stat, err := v2.Node.Status()
	if err != nil {
		return internalError(ctx, err, errFailedRetrievingNodeStatus, v2.Log)
	}

	response := generated.NodeStatusResponse{
		LastRound:                   uint64(stat.LastRound),
		LastVersion:                 string(stat.LastVersion),
		NextVersion:                 string(stat.NextVersion),
		NextVersionRound:            uint64(stat.NextVersionRound),
		NextVersionSupported:        stat.NextVersionSupported,
		TimeSinceLastRound:          uint64(stat.TimeSinceLastRound().Nanoseconds()),
		CatchupTime:                 uint64(stat.CatchupTime.Nanoseconds()),
		StoppedAtUnsupportedRound:   stat.StoppedAtUnsupportedRound,
		LastCatchpoint:              &stat.LastCatchpoint,
		Catchpoint:                  &stat.Catchpoint,
		CatchpointTotalAccounts:     &stat.CatchpointCatchupTotalAccounts,
		CatchpointProcessedAccounts: &stat.CatchpointCatchupProcessedAccounts,
		CatchpointVerifiedAccounts:  &stat.CatchpointCatchupVerifiedAccounts,
		CatchpointTotalBlocks:       &stat.CatchpointCatchupTotalBlocks,
		CatchpointAcquiredBlocks:    &stat.CatchpointCatchupAcquiredBlocks,
	}

	return ctx.JSON(http.StatusOK, response)
}

// WaitForBlock returns the node status after waiting for the given round.
// (GET /v2/status/wait-for-block-after/{round}/)
func (v2 *Handlers) WaitForBlock(ctx echo.Context, round uint64) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	ledger := v2.Node.Ledger()

	stat, err := v2.Node.Status()
	if err != nil {
		return internalError(ctx, err, errFailedRetrievingNodeStatus, v2.Log)
	}
	if stat.StoppedAtUnsupportedRound {
		return badRequest(ctx, err, errRequestedRoundInUnsupportedRound, v2.Log)
	}
	if stat.Catchpoint != "" {
		// node is currently catching up to the requested catchpoint.
		return serviceUnavailable(ctx, fmt.Errorf("WaitForBlock failed as the node was catchpoint catchuping"), errOperationNotAvailableDuringCatchup, v2.Log)
	}

	latestBlkHdr, err := ledger.BlockHdr(ledger.Latest())
	if err != nil {
		return internalError(ctx, err, errFailedRetrievingLatestBlockHeaderStatus, v2.Log)
	}
	if latestBlkHdr.NextProtocol != "" {
		if _, nextProtocolSupported := config.Consensus[latestBlkHdr.NextProtocol]; !nextProtocolSupported {
			// see if the desired protocol switch is expect to happen before or after the above point.
			if latestBlkHdr.NextProtocolSwitchOn <= basics.Round(round+1) {
				// we would never reach to this round, since this round would happen after the (unsupported) protocol upgrade.
				return badRequest(ctx, err, errRequestedRoundInUnsupportedRound, v2.Log)
			}
		}
	}

	// Wait
	select {
	case <-v2.Shutdown:
		return internalError(ctx, err, errServiceShuttingDown, v2.Log)
	case <-time.After(1 * time.Minute):
	case <-ledger.Wait(basics.Round(round + 1)):
	}

	// Return status after the wait
	return v2.GetStatus(ctx)
}

// RawTransaction broadcasts a raw transaction to the network.
// (POST /v2/transactions)
func (v2 *Handlers) RawTransaction(ctx echo.Context, params generated.RawTransactionParams) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	stat, err := v2.Node.Status()
	if err != nil {
		return internalError(ctx, err, errFailedRetrievingNodeStatus, v2.Log)
	}
	if stat.Catchpoint != "" {
		// node is currently catching up to the requested catchpoint.
		return serviceUnavailable(ctx, fmt.Errorf("RawTransaction failed as the node was catchpoint catchuping"), errOperationNotAvailableDuringCatchup, v2.Log)
	}
	proto := config.Consensus[stat.LastVersion]

	var txgroup []transactions.SignedTxn
	dec := protocol.NewDecoder(ctx.Request().Body)
	for {
		var st transactions.SignedTxn
		err := dec.Decode(&st)
		if err == io.EOF {
			break
		}
		if err != nil {
			return badRequest(ctx, err, err.Error(), v2.Log)
		}
		txgroup = append(txgroup, st)

		if len(txgroup) > proto.MaxTxGroupSize {
			err := fmt.Errorf("max group size is %d", proto.MaxTxGroupSize)
			return badRequest(ctx, err, err.Error(), v2.Log)
		}
	}

	if len(txgroup) == 0 {
		err := errors.New("empty txgroup")
		return badRequest(ctx, err, err.Error(), v2.Log)
	}

	if params.Speculation != nil {
		// Rather than broadcast the txgroup, apply it to the speculation ledger
		ledger, err := v2.Node.SpeculationLedger(*params.Speculation)
		if err != nil {
			return badRequest(ctx, err, errFailedLookingUpLedger, v2.Log)
		}
		err = ledger.Apply(txgroup)
		if err != nil {
			return badRequest(ctx, err, fmt.Sprintf("%v", err), v2.Log)
		}
		return ctx.JSON(http.StatusOK, generated.PostTransactionsResponse{TxId: "notImPleMENted"})
	}

	err = v2.Node.BroadcastSignedTxGroup(txgroup)
	if err != nil {
		return badRequest(ctx, err, err.Error(), v2.Log)
	}

	// For backwards compatibility, return txid of first tx in group
	txid := txgroup[0].ID()
	return ctx.JSON(http.StatusOK, generated.PostTransactionsResponse{TxId: txid.String()})
}

// TealDryrun takes transactions and additional simulated ledger state and returns debugging information.
// (POST /v2/teal/dryrun)
func (v2 *Handlers) TealDryrun(ctx echo.Context) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	if !v2.Node.Config().EnableDeveloperAPI {
		return ctx.String(http.StatusNotFound, "/teal/dryrun was not enabled in the configuration file by setting the EnableDeveloperAPI to true")
	}
	req := ctx.Request()
	buf := new(bytes.Buffer)
	req.Body = http.MaxBytesReader(nil, req.Body, maxTealDryrunBytes)
	buf.ReadFrom(req.Body)
	data := buf.Bytes()

	var dr DryrunRequest
	var gdr generated.DryrunRequest
	err := decode(protocol.JSONHandle, data, &gdr)
	if err == nil {
		dr, err = DryrunRequestFromGenerated(&gdr)
		if err != nil {
			return badRequest(ctx, err, err.Error(), v2.Log)
		}
	} else {
		err = decode(protocol.CodecHandle, data, &dr)
		if err != nil {
			return badRequest(ctx, err, err.Error(), v2.Log)
		}
	}

	// fetch previous block header just once to prevent racing with network
	var hdr bookkeeping.BlockHeader
	if dr.ProtocolVersion == "" || dr.Round == 0 || dr.LatestTimestamp == 0 {
		actualLedger := v2.Node.Ledger()
		hdr, err = actualLedger.BlockHdr(actualLedger.Latest())
		if err != nil {
			return internalError(ctx, err, "current block error", v2.Log)
		}
	}

	var response generated.DryrunResponse

	var protocolVersion protocol.ConsensusVersion
	if dr.ProtocolVersion != "" {
		var ok bool
		_, ok = config.Consensus[protocol.ConsensusVersion(dr.ProtocolVersion)]
		if !ok {
			return badRequest(ctx, nil, "unsupported protocol version", v2.Log)
		}
		protocolVersion = protocol.ConsensusVersion(dr.ProtocolVersion)
	} else {
		protocolVersion = hdr.CurrentProtocol
	}
	dr.ProtocolVersion = string(protocolVersion)

	if dr.Round == 0 {
		dr.Round = uint64(hdr.Round + 1)
	}

	if dr.LatestTimestamp == 0 {
		dr.LatestTimestamp = hdr.TimeStamp
	}

	doDryrunRequest(&dr, &response)
	response.ProtocolVersion = string(protocolVersion)
	return ctx.JSON(http.StatusOK, response)
}

// TransactionParams returns the suggested parameters for constructing a new transaction.
// (GET /v2/transactions/params)
func (v2 *Handlers) TransactionParams(ctx echo.Context, params generated.TransactionParamsParams) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	stat, err := v2.Node.Status()
	if err != nil {
		return internalError(ctx, err, errFailedRetrievingNodeStatus, v2.Log)
	}
	if stat.Catchpoint != "" {
		// node is currently catching up to the requested catchpoint.
		return serviceUnavailable(ctx, fmt.Errorf("TransactionParams failed as the node was catchpoint catchuping"), errOperationNotAvailableDuringCatchup, v2.Log)
	}

	version := stat.LastVersion
	last := stat.LastRound

	if params.Speculation != nil {
		ledger, err := v2.Node.SpeculationLedger(*params.Speculation)
		if err != nil {
			return badRequest(ctx, err, errFailedLookingUpLedger, v2.Log)
		}
		version = ledger.Version
		last = ledger.Latest()
	}

	gh := v2.Node.GenesisHash()
	proto := config.Consensus[version]

	response := generated.TransactionParametersResponse{
		ConsensusVersion: string(version),
		Fee:              v2.Node.SuggestedFee().Raw,
		GenesisHash:      gh[:],
		GenesisId:        v2.Node.GenesisID(),
		LastRound:        uint64(last),
		MinFee:           proto.MinTxnFee,
	}

	return ctx.JSON(http.StatusOK, response)
}

// PendingTransactionInformation returns a transaction with the specified txID
// from the transaction pool. If not found looks for the transaction in the
// last proto.MaxTxnLife rounds
// (GET /v2/transactions/pending/{txid})
func (v2 *Handlers) PendingTransactionInformation(ctx echo.Context, txid string, params generated.PendingTransactionInformationParams) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)

	stat, err := v2.Node.Status()
	if err != nil {
		return internalError(ctx, err, errFailedRetrievingNodeStatus, v2.Log)
	}
	if stat.Catchpoint != "" {
		// node is currently catching up to the requested catchpoint.
		return serviceUnavailable(ctx, fmt.Errorf("PendingTransactionInformation failed as the node was catchpoint catchuping"), errOperationNotAvailableDuringCatchup, v2.Log)
	}

	txID := transactions.Txid{}
	if err := txID.UnmarshalText([]byte(txid)); err != nil {
		return badRequest(ctx, err, errNoTxnSpecified, v2.Log)
	}

	txn, ok := v2.Node.GetPendingTransaction(txID)

	// We didn't find it, return a failure
	if !ok {
		err := errors.New(errTransactionNotFound)
		return notFound(ctx, err, err.Error(), v2.Log)
	}

	// Encoding wasn't working well without embedding "real" objects.
	response := struct {
		AssetIndex         *uint64                        `codec:"asset-index,omitempty"`
		AssetClosingAmount *uint64                        `codec:"asset-closing-amount,omitempty"`
		ApplicationIndex   *uint64                        `codec:"application-index,omitempty"`
		CloseRewards       *uint64                        `codec:"close-rewards,omitempty"`
		ClosingAmount      *uint64                        `codec:"closing-amount,omitempty"`
		ConfirmedRound     *uint64                        `codec:"confirmed-round,omitempty"`
		GlobalStateDelta   *generated.StateDelta          `codec:"global-state-delta,omitempty"`
		LocalStateDelta    *[]generated.AccountStateDelta `codec:"local-state-delta,omitempty"`
		PoolError          string                         `codec:"pool-error"`
		ReceiverRewards    *uint64                        `codec:"receiver-rewards,omitempty"`
		SenderRewards      *uint64                        `codec:"sender-rewards,omitempty"`
		Txn                transactions.SignedTxn         `codec:"txn"`
	}{
		Txn: txn.Txn,
	}

	handle, contentType, err := getCodecHandle(params.Format)
	if err != nil {
		return badRequest(ctx, err, errFailedParsingFormatOption, v2.Log)
	}

	if txn.ConfirmedRound != 0 {
		r := uint64(txn.ConfirmedRound)
		response.ConfirmedRound = &r

		response.ClosingAmount = &txn.ApplyData.ClosingAmount.Raw
		response.AssetClosingAmount = &txn.ApplyData.AssetClosingAmount
		response.SenderRewards = &txn.ApplyData.SenderRewards.Raw
		response.ReceiverRewards = &txn.ApplyData.ReceiverRewards.Raw
		response.CloseRewards = &txn.ApplyData.CloseRewards.Raw

		response.AssetIndex = computeAssetIndexFromTxn(txn, v2.Node.Ledger())
		response.ApplicationIndex = computeAppIndexFromTxn(txn, v2.Node.Ledger())

		response.LocalStateDelta, response.GlobalStateDelta = convertToDeltas(txn)
	}

	data, err := encode(handle, response)
	if err != nil {
		return internalError(ctx, err, errFailedToEncodeResponse, v2.Log)
	}

	return ctx.Blob(http.StatusOK, contentType, data)
}

// getPendingTransactions returns to the provided context a list of uncomfirmed transactions currently in the transaction pool with optional Max/Address filters.
func (v2 *Handlers) getPendingTransactions(ctx echo.Context, max *uint64, format *string, addrFilter *string) error {

	stat, err := v2.Node.Status()
	if err != nil {
		return internalError(ctx, err, errFailedRetrievingNodeStatus, v2.Log)
	}
	if stat.Catchpoint != "" {
		// node is currently catching up to the requested catchpoint.
		return serviceUnavailable(ctx, fmt.Errorf("PendingTransactionInformation failed as the node was catchpoint catchuping"), errOperationNotAvailableDuringCatchup, v2.Log)
	}

	var addrPtr *basics.Address

	if addrFilter != nil {
		addr, err := basics.UnmarshalChecksumAddress(*addrFilter)
		if err != nil {
			return badRequest(ctx, err, errFailedToParseAddress, v2.Log)
		}
		addrPtr = &addr
	}

	handle, contentType, err := getCodecHandle(format)
	if err != nil {
		return badRequest(ctx, err, errFailedParsingFormatOption, v2.Log)
	}

	txnPool, err := v2.Node.GetPendingTxnsFromPool()
	if err != nil {
		return internalError(ctx, err, errFailedLookingUpTransactionPool, v2.Log)
	}

	// MatchAddress uses this to check FeeSink, we don't care about that here.
	spec := transactions.SpecialAddresses{
		FeeSink:     basics.Address{},
		RewardsPool: basics.Address{},
	}

	txnLimit := uint64(math.MaxUint64)
	if max != nil && *max != 0 {
		txnLimit = *max
	}

	// Convert transactions to msgp / json strings
	topTxns := make([]transactions.SignedTxn, 0)
	for _, txn := range txnPool {
		// break out if we've reached the max number of transactions
		if uint64(len(topTxns)) >= txnLimit {
			break
		}

		// continue if we have an address filter and the address doesn't match the transaction.
		if addrPtr != nil && !txn.Txn.MatchAddress(*addrPtr, spec) {
			continue
		}

		topTxns = append(topTxns, txn)
	}

	// Encoding wasn't working well without embedding "real" objects.
	response := struct {
		TopTransactions   []transactions.SignedTxn `json:"top-transactions"`
		TotalTransactions uint64                   `json:"total-transactions"`
	}{
		TopTransactions:   topTxns,
		TotalTransactions: uint64(len(txnPool)),
	}

	data, err := encode(handle, response)
	if err != nil {
		return internalError(ctx, err, errFailedToEncodeResponse, v2.Log)
	}

	return ctx.Blob(http.StatusOK, contentType, data)
}

// startCatchup Given a catchpoint, it starts catching up to this catchpoint
func (v2 *Handlers) startCatchup(ctx echo.Context, catchpoint string) error {
	_, _, err := ledgercore.ParseCatchpointLabel(catchpoint)
	if err != nil {
		return badRequest(ctx, err, errFailedToParseCatchpoint, v2.Log)
	}

	// Select 200/201, or return an error
	var code int
	err = v2.Node.StartCatchup(catchpoint)
	switch err.(type) {
	case nil:
		code = http.StatusCreated
	case *node.CatchpointAlreadyInProgressError:
		code = http.StatusOK
	case *node.CatchpointUnableToStartError:
		return badRequest(ctx, err, err.Error(), v2.Log)
	default:
		return internalError(ctx, err, fmt.Sprintf(errFailedToStartCatchup, err), v2.Log)
	}

	return ctx.JSON(code, private.CatchpointStartResponse{
		CatchupMessage: catchpoint,
	})
}

// abortCatchup Given a catchpoint, it aborts catching up to this catchpoint
func (v2 *Handlers) abortCatchup(ctx echo.Context, catchpoint string) error {
	_, _, err := ledgercore.ParseCatchpointLabel(catchpoint)
	if err != nil {
		return badRequest(ctx, err, errFailedToParseCatchpoint, v2.Log)
	}

	err = v2.Node.AbortCatchup(catchpoint)
	if err != nil {
		return internalError(ctx, err, fmt.Sprintf(errFailedToAbortCatchup, err), v2.Log)
	}

	return ctx.JSON(http.StatusOK, private.CatchpointAbortResponse{
		CatchupMessage: catchpoint,
	})
}

// GetPendingTransactions returns the list of unconfirmed transactions currently in the transaction pool.
// (GET /v2/transactions/pending)
func (v2 *Handlers) GetPendingTransactions(ctx echo.Context, params generated.GetPendingTransactionsParams) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	return v2.getPendingTransactions(ctx, params.Max, params.Format, nil)
}

// GetApplicationByID returns application information by app idx.
// (GET /v2/applications/{application-id})
func (v2 *Handlers) GetApplicationByID(ctx echo.Context, applicationId uint64, params generated.GetApplicationByIDParams) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	appIdx := basics.AppIndex(applicationId)
	ledger := v2.Node.Ledger()
	creator, ok, err := ledger.GetCreator(basics.CreatableIndex(appIdx), basics.AppCreatable)
	if err != nil {
		return internalError(ctx, err, errFailedLookingUpLedger, v2.Log)
	}
	if !ok {
		return notFound(ctx, errors.New(errAppDoesNotExist), errAppDoesNotExist, v2.Log)
	}

	record, _, err := ledger.LookupLatestWithoutRewards(creator)
	if err != nil {
		return internalError(ctx, err, errFailedLookingUpLedger, v2.Log)
	}

	appParams, ok := record.AppParams[appIdx]
	if !ok {
		return notFound(ctx, errors.New(errAppDoesNotExist), errAppDoesNotExist, v2.Log)
	}
	app := AppParamsToApplication(creator.String(), appIdx, &appParams)
	response := generated.ApplicationResponse(app)
	return ctx.JSON(http.StatusOK, response)
}

// GetAssetByID returns application information by app idx.
// (GET /v2/assets/{asset-id})
func (v2 *Handlers) GetAssetByID(ctx echo.Context, assetId uint64, params generated.GetAssetByIDParams) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	assetIdx := basics.AssetIndex(assetId)
	var ledger ledgerForApiHandlers
	var err error
	if params.Speculation == nil {
		ledger = v2.Node.Ledger()
	} else {
		ledger, err = v2.Node.SpeculationLedger(*params.Speculation)
		if err != nil {
			return badRequest(ctx, err, errFailedLookingUpLedger, v2.Log)
		}
	}
	creator, ok, err := ledger.GetCreator(basics.CreatableIndex(assetIdx), basics.AssetCreatable)
	if err != nil {
		return internalError(ctx, err, errFailedLookingUpLedger, v2.Log)
	}
	if !ok {
		return notFound(ctx, errors.New(errAssetDoesNotExist), errAssetDoesNotExist, v2.Log)
	}

	record, err := ledger.LookupLatest(creator)
	if err != nil {
		return internalError(ctx, err, errFailedLookingUpLedger, v2.Log)
	}

	assetParams, ok := record.AssetParams[assetIdx]
	if !ok {
		return notFound(ctx, errors.New(errAssetDoesNotExist), errAssetDoesNotExist, v2.Log)
	}

	asset := AssetParamsToAsset(creator.String(), assetIdx, &assetParams)
	response := generated.AssetResponse(asset)
	return ctx.JSON(http.StatusOK, response)
}

// GetPendingTransactionsByAddress takes an Algorand address and returns its associated list of unconfirmed transactions currently in the transaction pool.
// (GET /v2/accounts/{address}/transactions/pending)
func (v2 *Handlers) GetPendingTransactionsByAddress(ctx echo.Context, addr string, params generated.GetPendingTransactionsByAddressParams) error {
	return v2.getPendingTransactions(ctx, params.Max, params.Format, &addr)
}

// StartCatchup Given a catchpoint, it starts catching up to this catchpoint
// (POST /v2/catchup/{catchpoint})
func (v2 *Handlers) StartCatchup(ctx echo.Context, catchpoint string) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	return v2.startCatchup(ctx, catchpoint)
}

// AbortCatchup Given a catchpoint, it aborts catching up to this catchpoint
// (DELETE /v2/catchup/{catchpoint})
func (v2 *Handlers) AbortCatchup(ctx echo.Context, catchpoint string) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	return v2.abortCatchup(ctx, catchpoint)
}

// TealCompile compiles TEAL code to binary, return both binary and hash
// (POST /v2/teal/compile)
func (v2 *Handlers) TealCompile(ctx echo.Context) error {
	prof.Start(kNode)
	defer prof.Start(kKalgoTotal)
	// return early if teal compile is not allowed in node config
	if !v2.Node.Config().EnableDeveloperAPI {
		return ctx.String(http.StatusNotFound, "/teal/compile was not enabled in the configuration file by setting the EnableDeveloperAPI to true")
	}
	buf := new(bytes.Buffer)
	ctx.Request().Body = http.MaxBytesReader(nil, ctx.Request().Body, maxTealSourceBytes)
	buf.ReadFrom(ctx.Request().Body)
	source := buf.String()
	ops, err := logic.AssembleString(source)
	if err != nil {
		return badRequest(ctx, err, err.Error(), v2.Log)
	}
	pd := logic.HashProgram(ops.Program)
	addr := basics.Address(pd)
	response := generated.CompileResponse{
		Hash:   addr.String(),
		Result: base64.StdEncoding.EncodeToString(ops.Program),
	}
	return ctx.JSON(http.StatusOK, response)
}
