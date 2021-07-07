package kalgo

import (
	"testing"
)

func TestParseKalgoOutput(t *testing.T) {
	contents := `<clarity>
  <k>
    .
  </k>
  <mode>
    root
  </mode>
  <revert>
    false
  </revert>
  <callStack>
    .List
  </callStack>
  <callDepth>
    0
  </callDepth>
  <eventID>
    2
  </eventID>
  <events>
    ListItem ( [ 1 : commit ( 0 ) ] )
  </events>
  <checkpoints>
    .List
  </checkpoints>
  <asserts>
    .List
  </asserts>
  <logStack>
    .List
  </logStack>
  <costTracker>
    <cTotal>
      ExecCost ( 1 , 37 , 0 , 0 , 38 )
    </cTotal>
    <cLimit>
      ExecCost ( 18446744073709551615 , 18446744073709551615 , 18446744073709551615 , 18446744073709551615 , 18446744073709551615 )
    </cLimit>
    <cTxLimit>
      ExecCost ( 18446744073709551615 , 18446744073709551615 , 18446744073709551615 , 18446744073709551615 , 18446744073709551615 )
    </cTxLimit>
    <cMemory>
      0
    </cMemory>
    <cStmts>
      0
    </cStmts>
    <cMemLimit>
      100000000
    </cMemLimit>
  </costTracker>
  <pprefix>
    "tests/clarity"
  </pprefix>
  <current>
    .do-nothing
  </current>
  <currentFun>
    .Name
  </currentFun>
  <commitments>
    .do-nothing |-> InitialCommitment ( "\x07)\xa7\xce#\xde\x99\xee\xf5\x99\f\xdc}&\xbd/\xa5\xdbJ)\xf8\xf7\xb9\xbb\x8b5\xb8G3\xceJ\x84" )
  </commitments>
  <updated>
    SetItem ( .do-nothing )
  </updated>
  <thrownValue>
    .ValueOrError
  </thrownValue>
  <contracts>
    .ContractCellMap
  </contracts>
  <context>
    <txSender>
      @AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAY5HFKQ
    </txSender>
    <bufferedIO>
      <bio-op>
        .
      </bio-op>
      <bio-file>
        ""
      </bio-file>
      <bio-fd>
        -1
      </bio-fd>
      <bio-buffer>
        ""
      </bio-buffer>
      <bio-error>
        ""
      </bio-error>
    </bufferedIO>
    <txGroup>
      <txGroupID>
        0
      </txGroupID>
      <currentTx>
        0
      </currentTx>
      <transactions>
        .TransactionCellMap
      </transactions>
    </txGroup>
  </context>
  <result>
    .
  </result>
  <stdout>
    .
  </stdout>
  <txnEffects>
    .Map
  </txnEffects>
  <returncode>
    0
  </returncode>
  <returnstatus>
    "Success - program executed to completion successfully"
  </returnstatus>
</clarity>
`
	_, err := ParseOutput([]byte(contents))
	if err != nil {
		t.Errorf("Failed to parse contents (%v)", err.Error())
	}
}
