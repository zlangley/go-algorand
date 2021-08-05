package layer2

import (
	"fmt"
	"github.com/algorand/go-algorand/config"
	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/data/committee"
	"github.com/algorand/go-algorand/data/committee/sortition"
	"github.com/algorand/go-algorand/protocol"
)

type (
	// step is a sequence number denoting distinct stages in Algorand
	period uint64

	// period is used to track progress with a given round in the protocol
	step uint64
)

type Selector struct {
	_struct struct{} `codec:""` // not omitempty

	Seed   committee.Seed `codec:"seed"`
	Round  basics.Round   `codec:"rnd"`
	Period period         `codec:"per"`
}

func (sel Selector) ToBeHashed() (protocol.HashID, []byte) {
	return protocol.ExecutionSelector, protocol.EncodeReflect(&sel)
}

func (sel Selector) CommitteeSize(proto config.ConsensusParams) uint64 {
	return 140
}

func CurrentSelector() Selector {
	return Selector{Round: 1}
}

func (sel Selector) ComputeWeightOnCommittee(cred committee.UnauthenticatedCredential, vrfPub crypto.VrfPubkey, verifier crypto.SignatureVerifier) (uint64, error) {
	var voterMoney uint64 = 200
	var totalMoney uint64 = 10000000

	proto := config.Consensus[protocol.ConsensusCurrentVersion]

	// This should happen once per round.
	ok, vrfOut := vrfPub.Verify(cred.Proof, sel)
	if !ok {
		err := fmt.Errorf("UnauthenticatedCredential.Verify: could not verify VRF Proof with %v (parameters = %+v, proof = %#v)", vrfPub, sel, cred.Proof)
		return 0, err
	}

	h := crypto.Hash(append(vrfOut[:], verifier[:]...))
	weight := sortition.Select(voterMoney, totalMoney, float64(sel.CommitteeSize(proto)), h)
	return weight, nil
}