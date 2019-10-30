#!/bin/bash

# 删除数据库
rm -rf *.db

# build 项目
go build -o bc *.go

# 执行文件
./bc 

