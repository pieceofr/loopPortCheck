#!/usr/bin/env sh
$CURR_PUBLIC_IP="118.163.120.180"

docker run -d --name loopPortChec -p 2136:2136  \
    -e HOST=$CURR_PUBLIC_IP \
    -e PORT="2136" \
    pieceofr/loop-port-check

#docker run -d --name loopPortCheck -p 2136:2136 -e HOST="118.163.120.180" -e PORT="2136" pieceofr/loop-port-check
#docker run -it --name loopPortCheck -p 2136:2136 -e HOST="118.163.120.180" -e PORT="2136" pieceofr/loop-port-check