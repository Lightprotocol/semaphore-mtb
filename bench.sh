#!/bin/bash

for ((batch_size = 1; batch_size <= 10; batch_size++)); do
    echo "init $batch_size..." >> log.txt
    echo "Batch size: $batch_size"
    if [ ! -f "circuit_$batch_size_22" ]; then
        ./gnark-mbu setup --batch-size $batch_size --mode insertion --tree-depth 22 --output circuit_${batch_size}_22
    fi
    ./gnark-mbu start --keys-file circuit_${batch_size}_22 --mode insertion --json-logging >> log.txt &
    sleep $((batch_size))
    ./gnark-mbu gen-test-params --mode insertion --batch-size $batch_size --tree-depth 22 > ${batch_size}_22_test.json
    curl -X POST -d @${batch_size}_22_test.json http://localhost:3001/prove
    killall gnark-mbu
    sleep 2
done
