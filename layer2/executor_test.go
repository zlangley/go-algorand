package layer2

import (
	"testing"

	"github.com/algorand/go-algorand/crypto"
	"github.com/stretchr/testify/require"
)

func TestTealTemplate(t *testing.T) {
	contractPreID := crypto.Hash([]byte("test1"))
	addr1, prog1, err := logicSigFromTemplateFile("committee-defer-logicsig.teal.template", contractPreID)
	require.NoError(t, err)

	contractPreID = crypto.Hash([]byte("test2"))
	addr2, prog2, err := logicSigFromTemplateFile("committee-defer-logicsig.teal.template", contractPreID)
	require.NoError(t, err)

	require.NotEqual(t, addr1, addr2)
	require.NotEqual(t, prog1, prog2)
}
