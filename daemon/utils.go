package daemon

import "github.com/asdine/storm"

func ConvertStormError(err error) error {
	if err == storm.ErrNotFound {
		return errRecordNotFound
	} else {
		return errDatabase
	}
}
