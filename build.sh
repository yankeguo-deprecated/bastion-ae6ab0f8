#!/bin/bash
set -e
set -u

SRC_ROOT=$(pwd)/$(dirname $0)

# protobuf
protoc -I ${SRC_ROOT}/types --go_out=plugins=grpc:${SRC_ROOT}/types ${SRC_ROOT}/types/daemon.proto

# build vue project
cd ${SRC_ROOT}/web/ui && yarn build
# remove *.map files
cd ${SRC_ROOT}/web/public && rm -rf **/*.map
# generate binfs package
cd ${SRC_ROOT}/web && PKG=web binfs public > public.bfs.go

# build executables
rm -rf ${SRC_ROOT}/build
mkdir -p ${SRC_ROOT}/build
for ARCH in "amd64"
do
    for OS in "darwin" "linux"
    do
        for CMD in $(ls -1 ${SRC_ROOT}/cmd)
        do
            echo "building ${CMD}-${OS}-${ARCH}..."
            GOOS=${OS} GOARCH=${ARCH} go build -o ${SRC_ROOT}/build/${CMD}-${OS}-${ARCH} github.com/yankeguo/bastion/cmd/${CMD}
        done
    done
done
