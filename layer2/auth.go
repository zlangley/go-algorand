package layer2

import (
	"fmt"
	"github.com/algorand/go-algorand/crypto"
	"strings"
	"text/template"

	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/data/transactions/logic"
)

var (
	ForwardDeclareAppIndex basics.AppIndex = 20189847
	VersionNumberAppIndex  basics.AppIndex = 21328402
)

func GetContractAddress(contractPreID crypto.Digest) basics.Address {
	addr, _, err := logicSigFromTemplateFile("layer2/committee-defer-logicsig.teal.template", contractPreID)
	if err != nil {
		panic(err)
	}
	return addr
}

func logicSigFromTemplateFile(filename string, contractPreID crypto.Digest) (basics.Address, []byte, error) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		return basics.Address{}, nil, fmt.Errorf("could not parse template file: %v", err)
	}
	b := &strings.Builder{}
	err = t.Execute(b, contractPreID)
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
