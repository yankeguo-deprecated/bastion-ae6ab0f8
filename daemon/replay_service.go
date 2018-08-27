package daemon

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/yankeguo/bastion/types"
	"github.com/yankeguo/bastion/utils"
)

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
			filename := FilenameForSessionID(f.SessionId, d.opts.ReplayDir)
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
		if err = utils.WriteReplayFrame(f, zw); err != nil {
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
	filename := FilenameForSessionID(req.SessionId, d.opts.ReplayDir)
	var r *os.File
	if r, err = os.Open(filename); err != nil {
		return
	}
	defer r.Close()
	var zr *gzip.Reader
	if zr, err = gzip.NewReader(r); err != nil {
		return
	}
	defer zr.Close()
	for {
		var f types.ReplayFrame
		if err = utils.ReadReplayFrame(&f, zr); err != nil {
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
	return
}
