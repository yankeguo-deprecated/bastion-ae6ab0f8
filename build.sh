#!/bin/bash

COMMAND="$1"

if [ -z "${COMMAND}" ]
then
  COMMAND="help"
fi

set -e
set -u

SRC_ROOT=$(pwd)/$(dirname $0)

build_grpc () {
    protoc -I ${SRC_ROOT}/types --go_out=plugins=grpc:${SRC_ROOT}/types ${SRC_ROOT}/types/daemon.proto
}

build_binfs () {
    cd ${SRC_ROOT}/web/ui && yarn build
    cd ${SRC_ROOT}/web && PKG=web ${GOPATH}/bin/binfs public > public.bfs.go
}

build_cmd () {
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
}

case "${COMMAND}" in
        grpc)
            build_grpc
            ;;

        binfs)
            build_binfs
            ;;

        cmd)
            build_cmd
            ;;

        all)
            build_grpc
            build_binfs
            build_cmd
            ;;

        *)
            echo $"Usage: $0 {grpc|binfs|cmd|all}"
            exit 1
esac
