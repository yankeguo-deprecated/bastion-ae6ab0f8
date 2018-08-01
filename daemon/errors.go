package daemon

import (
	"github.com/asdine/storm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errRecordNotFound = status.Error(codes.InvalidArgument, "record not found")
	errDatabase       = status.Error(codes.Internal, "database error")
	errInternal       = status.Error(codes.Internal, "internal error")
)

func errDuplicatedField(key string) error {
	return status.Errorf(codes.AlreadyExists, "duplicated field '%s'", key)
}

func errFromStorm(err error) error {
	if err == nil {
		return nil
	} else if err == storm.ErrNotFound {
		return errRecordNotFound
	} else {
		return status.Error(codes.Internal, err.Error())
	}
}

func checkDuplicated(db storm.Node, bucket string, keyName string, keyVal interface{}) (err error) {
	var exists bool
	if exists, err = db.KeyExists(bucket, keyVal); err != nil {
		err = errFromStorm(err)
		return
	}
	if exists {
		err = errDuplicatedField(keyName)
		return
	}
	return
}
