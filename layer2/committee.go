package layer2

import (
	"github.com/algorand/go-algorand/config"
	"github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/data/committee"
	"github.com/algorand/go-algorand/data/committee/sortition"
	"github.com/algorand/go-algorand/protocol"
)

type simulationSelector struct{}

func (sel simulationSelector) ToBeHashed() (protocol.HashID, []byte) {
	return protocol.AgreementSelector, protocol.EncodeReflect(&sel)
}

func (sel simulationSelector) CommitteeSize(proto config.ConsensusParams) uint64 {
	return 1000
}

func ComputeWeight(acct basics.Address, acctData basics.AccountData) uint64 {
	sel := simulationSelector{}
	m := committee.Membership{
		Record:     committee.BalanceRecord{Addr: acct, AccountData: acctData},
		Selector:   sel,
		TotalMoney: basics.MicroAlgos{Raw: 10000000},
	}
	_, secret := crypto.VrfKeygen()
	u := committee.MakeCredential(&secret, sel)
	cred, _ := u.Verify(config.Consensus[protocol.ConsensusCurrentVersion], m)

	_, vrfOut := acctData.SelectionID.Verify(cred.Proof, m.Selector)
	h := crypto.Hash(append(vrfOut[:], acct[:]...))
	return sortition.Select(1000, 1000000000, 10000, h)
}