#!/bin/sh
# build binary

echo "start compiling..."
cd cmd
go build -o gotty-cmd1
curl fs.devops.haodai.net/soft/uploadapi -F file=@gotty-cmd1 -F truncate=yes
cksum gotty-cmd1
cd ..

go build -o gotty-pod1
curl fs.devops.haodai.net/soft/uploadapi -F file=@gotty-pod1 -F truncate=yes
cksum gotty-pod1
