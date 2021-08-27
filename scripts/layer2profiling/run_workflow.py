#!/usr/bin/env python3

from typing import Sequence
from dataclasses import dataclass
import sys
import os
import hashlib
import base64

import requests
import yaml
import pandas as pd


@dataclass
class SpeculationResponse:
    token: str
    base: int
    timing: dict = None
    checkpoints: str = None
    address: str = None
    commitments: str = None
    txns: str = None


class AlgodClient:
    def __init__(self, algod_address, algod_token):
        self.algod_address = algod_address
        self.algod_token = algod_token

        self._headers = {
            'accept': 'application/json',
            'X-Algo-API-Token': algod_token,
        }

    def create_speculation_token(self):
        response = requests.post(f'http://{algod_address}/v2/blocks/0/speculation', headers=self._headers).json()
        print(response)
        return SpeculationResponse(**response)

    def execute_batch(self, spec_token, payload):
        for item in payload:
            source = item.get('source')
            print(source)

        params = {'speculation': spec_token}
        response = requests.post(f'http://{algod_address}/v2/contracts/batch', headers=self._headers, params=params, json=payload).json()
        print(response)
        return SpeculationResponse(**response)

    def delete_speculation_token(self, spec_token):
        return requests.post(f'http://{algod_address}/v2/speculation/{spec_token}/delete', headers=self._headers).json()


def run_trial(client, workflow_file):
    spec_token = client.create_speculation_token().token

    with open(workflow_file, 'r') as f:
        payloads = yaml.load(f, Loader=yaml.Loader)

    setup_payload = payloads.get('setup')
    if setup_payload:
        client.execute_batch(spec_token, setup_payload)

    response = client.execute_batch(spec_token, payloads['test'])

#    client.delete_speculation_token(spec_token)

    return response


if __name__ == '__main__':
    if len(sys.argv) != 3:
        print(f'usage: {sys.argv[0]} <workflow yml file> <trials>')
        sys.exit(1)

    workflow_file = sys.argv[1]
    num_trials = int(sys.argv[2])

    algod_address = os.environ['ALGOD_ADDRESS']
    algod_token = os.environ['ALGOD_TOKEN']

    client = AlgodClient(algod_address, algod_token)

    rows = []

    for i in range(num_trials):
        print(f'Trial {i+1}')

        response = run_trial(client, workflow_file)
        rows.append([
                response.timing['kalgo'] / 1e6,
                response.timing['db'] / 1e6,
                response.timing['node'] / 1e6,
                response.timing['keygen'] / 1e6,
                response.timing['vrf'] / 1e6,
                response.timing['effects'] / 1e6,
                response.timing['total'] / 1e6,
        ])

    df = pd.DataFrame(rows, columns=['kalgo', 'db', 'node', 'keygen', 'vrf', 'effects', 'total'])
    df.to_csv(sys.stdout, sep='\t')
