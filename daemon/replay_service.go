package daemon

import (
	"compress/gzip"
	"encoding/binary"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/yankeguo/bastion/types"
)

func filenameForSessionID(id int64, dir string) string {
	buf := make([]byte, 8, 8)
	binary.BigEndian.PutUint64(buf, uint64(id))
	name := hex.EncodeToString(buf)
	ret := make([]string, 0, 5)
	ret = append(ret, dir)
	ret = append(ret, name[:4])
	ret = append(ret, name[4:8])
	ret = append(ret, name[8:12])
	ret = append(ret, name)
	return filepath.Join(ret...)
}

func writeReplayFrame(f *types.ReplayFrame, w io.Writer) (err error) {
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

func readReplayFrame(f *types.ReplayFrame, r io.Reader) (err error) {
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

func (d *Daemon) WriteReplay(s types.ReplayService_WriteReplayServer) (err error) {
	var w *os.File
	var zw *gzip.Writer
	for {
		var f *types.ReplayFrame
		// receive frame
		if f, err = s.Recv(); err != nil {
			if err == io.EOF {
				err = s.SendAndClose(&types.WriteReplayResponse{})
			}
			break
		}
		// ensure rec frame writer
		if zw == nil {
			// create filename
			filename := filenameForSessionID(f.SessionId, d.opts.ReplayDir)
			// ensure directory
			if err = os.MkdirAll(filepath.Dir(filename), 0750); err != nil {
				break
			}
			// open file
			if w, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0640); err != nil {
				break
			}
			// create frame writer with GZIP
			zw = gzip.NewWriter(w)
		}
		// write the frame
		if err = writeReplayFrame(f, zw); err != nil {
			break
		}
	}
	// close GZIP writer
	if zw != nil {
		zw.Close()
	}
	// close the GZIP writer won't close the file, so we have to close it manually
	if w != nil {
		w.Close()
	}
	return
}

func (d *Daemon) ReadReplay(req *types.ReadReplayRequest, s types.ReplayService_ReadReplayServer) (err error) {
	filename := filenameForSessionID(req.SessionId, d.opts.ReplayDir)
	var r *os.File
	if r, err = os.Open(filename); err != nil {
		return
	}
	var zr *gzip.Reader
	if zr, err = gzip.NewReader(r); err != nil {
		return
	}
	for {
		var f types.ReplayFrame
		if err = readReplayFrame(&f, zr); err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		f.SessionId = req.SessionId
		if err = s.Send(&f); err != nil {
			break
		}
	}
	if zr != nil {
		zr.Close()
	}
	if r != nil {
		r.Close()
	}
	return
}
