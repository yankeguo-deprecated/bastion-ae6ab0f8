#!/bin/bash

protoc -I types --go_out=plugins=grpc:types types/daemon.proto
