package daemon

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errRecordNotFound = status.Error(codes.InvalidArgument, "record not found")
	errDatabase       = status.Error(codes.Internal, "database error")
	errInternal       = status.Error(codes.Internal, "internal error")
)

func errFieldMissing(key string) error {
	return status.Errorf(codes.InvalidArgument, "missing field '%s'", key)
}

func errDuplicatedField(key string) error {
	return status.Errorf(codes.InvalidArgument, "duplicated field '%s'", key)
}
