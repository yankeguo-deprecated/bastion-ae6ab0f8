package daemon

import "time"

func now() int64 {
	return time.Now().Unix()
}
