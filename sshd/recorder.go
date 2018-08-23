package sshd

import (
	"context"
	"encoding/binary"
	"github.com/yankeguo/bastion/sshd/sandbox"
	"github.com/yankeguo/bastion/types"
	"io"
	"log"
	"time"
)

func timestamp(start time.Time) uint32 {
	return uint32(time.Now().Sub(start) / time.Millisecond)
}

type Recorder struct {
	wch chan sandbox.Window
	c   types.ReplayService_WriteReplayClient
}

func StartRecording(opts *sandbox.ExecAttachOptions, sessionID int64, rs types.ReplayServiceClient) (rec *Recorder) {
	rec = &Recorder{}
	var err error

	// starting time
	start := time.Now()

	// build replay write client
	if rec.c, err = rs.WriteReplay(context.Background()); err != nil {
		return
	}

	// if opts.WindowChan is not nil, replace it
	if opts.WindowChan != nil {
		wch := opts.WindowChan
		rec.wch = make(chan sandbox.Window, 4)
		go func() {
			// iterate original chan
			for w := range wch {
				// send frame
				buf := make([]byte, 8, 8)
				binary.BigEndian.PutUint32(buf, uint32(w.Width))
				binary.BigEndian.PutUint32(buf[4:], uint32(w.Height))
				if err = rec.c.Send(&types.ReplayFrame{
					SessionId: sessionID,
					Timestamp: timestamp(start),
					Type:      types.ReplayFrameTypeWindowSize,
					Payload:   buf,
				}); err != nil {
					log.Println("failed to send window-size record frame:", err)
				}
				// proxy to chan
				rec.wch <- w
			}
			close(rec.wch)
		}()
		opts.WindowChan = rec.wch
	}
	// if opts.Stdout is not nil, replace it
	if opts.Stdout != nil {
		opts.Stdout = NewRecorderWriter(opts.Stdout, sessionID, types.ReplayFrameTypeStdout, start, rec.c)
	}
	// if opts.Stderr is not nil, replace it
	if opts.Stderr != nil {
		opts.Stderr = NewRecorderWriter(opts.Stderr, sessionID, types.ReplayFrameTypeStderr, start, rec.c)
	}
	return
}

func (r *Recorder) Stop() {
	if r.c != nil {
		r.c.CloseAndRecv()
	}
}

type RecorderWriter struct {
	w         io.Writer
	sessionId int64
	typ       uint32
	start     time.Time
	client    types.ReplayService_WriteReplayClient
}

func (w *RecorderWriter) Write(p []byte) (int, error) {
	var err error
	if err = w.client.Send(&types.ReplayFrame{
		SessionId: w.sessionId,
		Timestamp: timestamp(w.start),
		Type:      w.typ,
		Payload:   p,
	}); err != nil {
		log.Println("failed to write frame:", err)
	}
	return w.w.Write(p)
}

func (w *RecorderWriter) Close() error {
	if c, ok := w.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func NewRecorderWriter(w io.Writer, sessionId int64, typ uint32, start time.Time, client types.ReplayService_WriteReplayClient) io.Writer {
	return &RecorderWriter{
		w:         w,
		sessionId: sessionId,
		typ:       typ,
		start:     start,
		client:    client,
	}
}
