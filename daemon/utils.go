package daemon

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func now() int64 {
	return time.Now().Unix()
}

func newToken() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}
