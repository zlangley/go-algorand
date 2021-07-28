package layer2

import (
	"github.com/algorand/go-algorand/data/basics"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTealTemplate(t *testing.T) {
	contractID, err := basics.UnmarshalChecksumAddress("SYY6J5YNCZWLAV5MMKRTLSPEQHB576IQDEEC7AJBLJYDJAWIHHPVYFE7GM")
	require.NoError(t, err)
	_, _, source := logicSigFromTemplateFile("committee-defer-logicsig.teal.template", contractID)
	require.Contains(t, source, contractID.String())
}
