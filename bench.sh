#!/bin/bash

for ((batch_size = 1; batch_size <= 10; batch_size++)); do
    echo "init $batch_size..." >> log.txt
    echo "Batch size: $batch_size"
    ./gnark-mbu setup --batch-size $batch_size --mode insertion --tree-depth 22 --output circuit_${batch_size}_22
    ./gnark-mbu start --keys-file circuit_${batch_size}_22 --mode insertion --json-logging >> log.txt &
    # MBU_PID=$!  # Store the PID of gnark-mbu process
    sleep 3
    ./gnark-mbu gen-test-params --mode insertion --batch-size $batch_size --tree-depth 22 > ${batch_size}_22_test.json
    curl -X POST -d @${batch_size}_22_test.json http://localhost:3001/prove
    # Kill gnark-mbu process
    # kill -9 $MBU_PID
    killall gnark-mbu
    sleep 2
done
