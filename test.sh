#!/bin/bash

VM_USER=lcx
VM_IP=lcx-Virtual-Machine
PROJECT_DIR=/home/lcx/feast
NAME=feast
SSH_KEY=./.ssh

ssh -i "$SSH_KEY" "$VM_USER@$VM_IP" "mkdir -p $PROJECT_DIR"

GOOS=linux GOARCH=amd64 go build -o ./bin/$NAME

scp -i "$SSH_KEY" ./bin/$NAME "$VM_USER@$VM_IP:$PROJECT_DIR"

ssh -i "$SSH_KEY" -t "$VM_USER@$VM_IP" "cd $PROJECT_DIR && chmod +x $NAME && ./$NAME"

