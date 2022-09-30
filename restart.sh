#!/bin/bash

make clean

make

./server.sh restart admin


./server.sh restart node
