#!/bin/bash

COMMAND="$1"
VERSION="$2"

if [ -z "${COMMAND}" ]
then
  COMMAND="help"
fi
if [ -z "${VERSION}" ]
then
  VERSION=`date '+%Y%m%d%H%M%S'`
fi

set -e
set -u

SRC_ROOT=$(pwd)/$(dirname $0)

build_grpc () {
    protoc -I ${SRC_ROOT}/types --go_out=plugins=grpc:${SRC_ROOT}/types ${SRC_ROOT}/types/daemon.proto
}

build_ui () {
    cd ${SRC_ROOT}/web/ui && yarn build
    cp -f ${SRC_ROOT}/resources/package.json ${SRC_ROOT}/web/public/package.json
    sed -i "" -e "s/__VERSION__/${VERSION}/g" ${SRC_ROOT}/web/public/package.json
    sed -i "" -e "s/\\=\\/static/\\=\\/\\/unpkg.com\\/bastion-assets@1.0.${VERSION}\\/static/g" ${SRC_ROOT}/web/public/index.html
    cd ${SRC_ROOT}/web/public && yarn publish
}

build_binfs () {
    cd ${SRC_ROOT}/web && PKG=web ${GOPATH}/bin/binfs public views > assets.bfs.go
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

        ui)
            build_ui
            ;;

        binfs)
            build_binfs
            ;;

        cmd)
            build_cmd
            ;;

        all)
            build_grpc
            build_ui
            build_binfs
            build_cmd
            ;;

        *)
            echo $"Usage: $0 {grpc|ui|binfs|cmd|all}"
            exit 1
esac
