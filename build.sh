#!/bin/sh
# build binary
set -e
echo "start compiling..."
go build
curl fs.devops.haodai.net/soft/uploadapi -F file=@gotty-pod -F truncate=yes
cksum gotty-pod
