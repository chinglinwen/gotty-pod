#!/bin/sh
# build image
suffix="$1"
suffix=${suffix:=v1}

sh build.sh

image="gotty-pod:$suffix"
tag="harbor.haodai.net/ops/$image"
echo -e "building image: $tag\n"
docker build --no-cache -t $tag .
docker push $tag
