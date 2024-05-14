#!/bin/bash


for ((i = 1; i <= 10; i++)); do
    echo '$(date +'%H:%M:%S') - Some data $i'
    sleep 1
done
