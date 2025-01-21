#!/bin/sh
tmux kill-session
ENCRYPTION_DIR=/home/$USER
tmux new-session -d bash
tmux split-window -h bash
tmux send -t "{ENCRYPTION_DIR}:0" -t 0:0.0 "python3 ${ENCRYPTION_DIR}/encryption/symmetric/discord.py listen" C-m
tmux send -t "{ENCRYPTION_DIR}:0.1" -t 0:0.1 "python3 ${ENCRYPTION_DIR}/encryption/symmetric/discord.py send" C-m
tmux -2 attach-session -d
