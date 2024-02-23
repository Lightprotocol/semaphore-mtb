#!/bin/bash

DEPTH="26"
URL="http://localhost:3001/prove"

gnark() {
    local args=("$@")
    ./prover-server "${args[@]}"
}

for ((utxos = 1; utxos <= 8; utxos++)); do
    echo "init $utxos..." >> log.txt
    echo "Number of utxos: $utxos"

    CIRCUIT_FILE="circuit_${DEPTH}_${utxos}"
    TEST_FILE="inputs_${utxos}_${DEPTH}_test.json"

    if [ ! -f "${CIRCUIT_FILE}" ]; then
        echo "Prover setup..."
        gnark setup --utxos "$utxos" --tree-depth "$DEPTH" --output "${CIRCUIT_FILE}"
    fi

    if [ ! -f "${TEST_FILE}" ]; then
      echo "Generating test inputs..."
      gnark gen-test-params --utxos "$utxos" --tree-depth "$DEPTH" > "${TEST_FILE}"
    fi

    curl -X POST -d @"${TEST_FILE}" "$URL"
done