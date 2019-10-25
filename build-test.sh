#!/bin/sh
# build binary
set -e

echo "start compiling..."

go build -o gotty-pod1
curl fs.haodai.net/soft/uploadapi -F file=@gotty-pod1 -F truncate=yes
cksum gotty-pod1
