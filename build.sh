#!/bin/sh
# build binary

echo "start compiling..."
go build
curl fs.devops.haodai.net/soft/uploadapi -F file=@gotty-pod -F truncate=yes
cksum gotty-pod
