#!/bin/bash
MODE="insertion"
DEPTH="22"
URL="http://localhost:3001/prove"
gnarkmbu() {
    local args=("$@")
    ./gnark-mbu "${args[@]}"
}
for ((batch_size = 1; batch_size <= 10; batch_size++)); do
    echo "init $batch_size..." >> log.txt
    echo "Batch size: $batch_size"

    CIRCUIT_FILE="/tmp/circuit_${batch_size}_${DEPTH}"
    TEST_FILE="/tmp/inputs_${batch_size}_${DEPTH}_test.json"

    if [ ! -f "${CIRCUIT_FILE}" ]; then
        gnarkmbu setup --batch-size "$batch_size" --mode "$MODE" --tree-depth "$DEPTH" --output "${CIRCUIT_FILE}"
    fi

    gnarkmbu start --keys-file "${CIRCUIT_FILE}" --mode "$MODE" --json-logging >> log.txt &
    sleep $((batch_size))

    gnarkmbu gen-test-params --mode "$MODE" --batch-size "$batch_size" --tree-depth "$DEPTH" > "${TEST_FILE}"

    curl -X POST -d @"${TEST_FILE}" "$URL"

    killall gnark-mbu
    sleep 2
done