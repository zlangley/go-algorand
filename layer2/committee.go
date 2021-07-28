package layer2

import (
	"github.com/algorand/go-algorand/config"
	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/data/committee"
	"github.com/algorand/go-algorand/data/committee/sortition"
	"github.com/algorand/go-algorand/protocol"
	"math/rand"
)

type (
	// step is a sequence number denoting distinct stages in Algorand
	period uint64

	// period is used to track progress with a given round in the protocol
	step uint64
)

type selector struct {
	_struct struct{} `codec:""` // not omitempty

	Seed   committee.Seed `codec:"seed"`
	Round  basics.Round   `codec:"rnd"`
	Period period         `codec:"per"`
}

func (sel selector) ToBeHashed() (protocol.HashID, []byte) {
	return protocol.ExecutionSelector, protocol.EncodeReflect(&sel)
}

func (sel selector) CommitteeSize(proto config.ConsensusParams) uint64 {
	return 140
}

func ComputeWeight() (uint64, error) {
	var voterMoney uint64 = 200
	var totalMoney uint64 = 10000000

	proto := config.Consensus[protocol.ConsensusCurrentVersion]

	gen := rand.New(rand.NewSource(2))
	var seed crypto.Seed
	gen.Read(seed[:])
	s := crypto.GenerateSignatureSecrets(seed)
	vrfPub, vrfSec := crypto.VrfKeygenFromSeed(seed)

	sel := selector{
		Round: 1,
	}
	// This should happen once per round.
	cred := committee.MakeCredential(&vrfSec, sel)

	_, vrfOut := vrfPub.Verify(cred.Proof, sel)

	h := crypto.Hash(append(vrfOut[:], s.SignatureVerifier[:]...))
	weight := sortition.Select(voterMoney, totalMoney, float64(sel.CommitteeSize(proto)), h)
	return weight, nil
}