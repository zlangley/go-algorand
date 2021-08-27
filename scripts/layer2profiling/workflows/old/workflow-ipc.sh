#!/bin/sh

ALGOD_ADDRESS=$(cat ~/testnet-data/algod.net)
export ALGOD_ADDRESS

ALGOD_TOKEN=$(cat ~/testnet-data/algod.token)
export ALGOD_TOKEN

SPECULATION_TOKEN=$(curl -s -X POST "${ALGOD_ADDRESS}/v2/blocks/0/speculation" -H  "accept: application/json" -H "X-Algo-API-Token: $ALGOD_TOKEN" | jq -r .token)
export SPECULATION_TOKEN

#echo "Created speculative state, SPECULATION_TOKEN=$SPECULATION_TOKEN"

curl -s "${ALGOD_ADDRESS}/v2/contracts/on-chain-assert2?speculation=${SPECULATION_TOKEN}" -H "Content-Type: plain/text" -H "accept: application/json" -H "X-Algo-API-Token: $ALGOD_TOKEN" --data-binary "@tests/clarity/on-chain-assert2/on-chain-assert2.clar"

curl -X POST "${ALGOD_ADDRESS}/v2/contracts/batch?speculation=${SPECULATION_TOKEN}" -H "Content-Type:application/json" -H "accept:application/json" -H "X-Algo-API-Token: $ALGOD_TOKEN" -d '[
  {
    "id": "call-predicate",
    "command": "init",
    "source": "(contract-call? .on-chain-assert3 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)"
  }
]'

curl -X POST "${ALGOD_ADDRESS}/v2/contracts/batch?speculation=${SPECULATION_TOKEN}" -H "Content-Type:application/json" -H "accept:application/json" -H "X-Algo-API-Token: $ALGOD_TOKEN" -d '[
  {
    "id": "call-predicate-2",
    "command": "init",
    "source": "
      (begin
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
      )"
  }
]'

curl -X POST "${ALGOD_ADDRESS}/v2/contracts/batch?speculation=${SPECULATION_TOKEN}" -H "Content-Type:application/json" -H "accept:application/json" -H "X-Algo-API-Token: $ALGOD_TOKEN" -d '[
  {
    "id": "call-predicate-2",
    "command": "init",
    "source": "
      (begin
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
        (contract-call? .on-chain-assert2 check @6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY)
      )"
  }
]'
