#!/bin/bash

PKG=static_test go run ../binfs/cmd/binfs/main.go testdata > staticdata_test.go
