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

// A Selector defines the seed for the sortition.
type Selector struct {
	_struct struct{} `codec:""` // not omitempty

	Seed   committee.Seed `codec:"seed"`
	Round  basics.Round   `codec:"rnd"`

	// TODO: need field for listener too
}

func (sel Selector) ToBeHashed() (protocol.HashID, []byte) {
	return protocol.ExecutionSelector, protocol.EncodeReflect(&sel)
}

func (sel Selector) CommitteeSize(proto config.ConsensusParams) uint64 {
	return 140
}

// CurrentSelector returns the current execution committee selector used for sortition.
func CurrentSelector() Selector {
	// TODO: This needs to be pulled from the contract committee app.
	return Selector{Round: 1}
}

// ComputeWeightOnCommittee determines the weight of the credential on the selector.
func (sel Selector) ComputeWeightOnCommittee(cred committee.UnauthenticatedCredential, vrfPub crypto.VrfPubkey, verifier crypto.SignatureVerifier) (uint64, error) {
	// TODO: These need to be determined from the user account data. This will be USDC eventually?
	var voterMoney uint64 = 200
	var totalMoney uint64 = 10000000

	proto := config.Consensus[protocol.ConsensusCurrentVersion]

	ok, vrfOut := vrfPub.Verify(cred.Proof, sel)
	if !ok {
		err := fmt.Errorf("UnauthenticatedCredential.Verify: could not verify VRF Proof with %v (parameters = %+v, proof = %#v)", vrfPub, sel, cred.Proof)
		return 0, err
	}

	h := crypto.Hash(append(vrfOut[:], verifier[:]...))
	weight := sortition.Select(voterMoney, totalMoney, float64(sel.CommitteeSize(proto)), h)
	return weight, nil
}