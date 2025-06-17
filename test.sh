#!/bin/bash

# 配置
VM_USER=lcx
VM_IP=lcx-Virtual-Machine
PROJECT_DIR=/home/lcx/feast
NAME=feast
SSH_KEY=./.ssh

# 在远程主机上创建目录（作为 root）
ssh -i "$SSH_KEY" "$VM_USER@$VM_IP" "mkdir -p $PROJECT_DIR"

# 构建 Go Linux 版本（在本地）
GOOS=linux GOARCH=amd64 go build -o ./bin/$NAME

# 上传二进制文件到远程机器
scp -i "$SSH_KEY" ./bin/$NAME "$VM_USER@$VM_IP:$PROJECT_DIR"

# 在远程机器上赋权并执行（在当前终端中运行）
ssh -i "$SSH_KEY" -t "$VM_USER@$VM_IP" "cd $PROJECT_DIR && chmod +x $NAME && ./$NAME"

