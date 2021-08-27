#!/bin/sh

ALGOD_ADDRESS=$(cat ~/testnet-data/algod.net)
export ALGOD_ADDRESS

ALGOD_TOKEN=$(cat ~/testnet-data/algod.token)
export ALGOD_TOKEN

SPECULATION_TOKEN=$(curl -s -X POST "${ALGOD_ADDRESS}/v2/blocks/0/speculation" -H  "accept: application/json" -H "X-Algo-API-Token: $ALGOD_TOKEN" | jq -r .token)
export SPECULATION_TOKEN

#echo "Created speculative state, SPECULATION_TOKEN=$SPECULATION_TOKEN"

curl -X POST "${ALGOD_ADDRESS}/v2/contracts/batch?speculation=${SPECULATION_TOKEN}" -H "Content-Type:application/json" -H "accept:application/json" -H "X-Algo-API-Token: $ALGOD_TOKEN" -d '[
  {
    "id": "do-nothing",
    "command": "init",
    "source": "
      (define-public (do-nothing) (ok 1))"
  },
  {
    "id": "do-nothing5",
    "command": "init",
    "source": "
      (define-public (do-nothing)
        (begin
          (contract-call? .do-nothing do-nothing)
          (contract-call? .do-nothing do-nothing)
          (contract-call? .do-nothing do-nothing)
          (contract-call? .do-nothing do-nothing)
          (contract-call? .do-nothing do-nothing)))"
  },
  {
    "id": "do-nothing5",
    "command": "call",
    "function": "do-nothing",
    "args": ""
  }
]'

#echo "Cleaning up speculation state"
curl -s -X POST "${ALGOD_ADDRESS}/v2/speculation/$SPECULATION_TOKEN/delete" -H  "accept: application/json" -H "X-Algo-API-Token: $ALGOD_TOKEN"
