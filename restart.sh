#!/bin/bash


make

./server.sh restart admin


./server.sh restart node
