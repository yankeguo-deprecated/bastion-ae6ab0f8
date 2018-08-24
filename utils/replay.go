package utils

import (
	"encoding/binary"
	"github.com/yankeguo/bastion/types"
	"io"
)

func MarshalReplayFrameWindowSizePayload(width, height uint32) []byte {
	buf := make([]byte, 8, 8)
	binary.BigEndian.PutUint32(buf, width)
	binary.BigEndian.PutUint32(buf[4:], height)
	return buf
}

func WriteReplayFrame(f *types.ReplayFrame, w io.Writer) (err error) {
	// TIMESTAMP (4 bytes) + TYPE (1 byte) + PAYLOAD_LEN (4 bytes) + PAYLOAD
	l := 4 + 1 + 4 + len(f.Payload)
	buf := make([]byte, l, l)
	binary.BigEndian.PutUint32(buf, f.Timestamp)
	buf[4] = byte(f.Type)
	binary.BigEndian.PutUint32(buf[5:], uint32(len(f.Payload)))
	copy(buf[9:], f.Payload)
	_, err = w.Write(buf)
	return
}

func ReadReplayFrame(f *types.ReplayFrame, r io.Reader) (err error) {
	// TIMESTAMP (4 bytes) + TYPE (1 byte) + PAYLOAD_LEN (4 bytes)
	h := make([]byte, 4+1+4, 4+1+4)
	if _, err = r.Read(h); err != nil {
		return
	}
	f.SessionId = 0
	f.Timestamp = binary.BigEndian.Uint32(h)
	f.Type = uint32(h[4])
	l := binary.BigEndian.Uint32(h[5:])
	if l > 0 {
		f.Payload = make([]byte, l, l)
		var i int
		if i, err = r.Read(f.Payload); err != nil {
			if i == int(l) && err == io.EOF {
				err = nil
			}
		}
	} else {
		f.Payload = make([]byte, 0, 0)
	}
	return
}
