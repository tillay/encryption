#!/bin/sh
tmux new-session -d bash
tmux split-window -h bash
tmux send -t 0:0.0 "python3 /home/semblanceofsense/bin/encryption/symmetric/discord.py listen" C-m
tmux send -t 0:0.1 "python3 /home/semblanceofsense/bin/encryption/symmetric/discord.py send" C-m
tmux -2 attach-session -d
