package recorder

import (
	"github.com/rs/zerolog/log"
	"github.com/yankeguo/bastion/types"
	"io"
	"time"
)

func diffUint32(v1, v2 uint32) uint32 {
	if v1 > v2 {
		return v1 - v2
	}
	return v2 - v1
}

func cloneFrame(f *types.ReplayFrame) *types.ReplayFrame {
	nf := &types.ReplayFrame{}
	*nf = *f
	nf.Payload = make([]byte, len(f.Payload), len(f.Payload))
	copy(nf.Payload, f.Payload)
	return nf
}

func concatPayload(p1, p2 []byte) []byte {
	buf := make([]byte, len(p1)+len(p2), len(p1)+len(p2))
	copy(buf, p1)
	copy(buf[len(p1):], p2)
	return buf
}

type FrameWriter struct {
	client types.ReplayService_WriteReplayClient
	last   *types.ReplayFrame
}

func NewFrameWriter(client types.ReplayService_WriteReplayClient) *FrameWriter {
	return &FrameWriter{
		client: client,
	}
}

func (fw *FrameWriter) WriteFrame(f *types.ReplayFrame) (err error) {
	// if no cached frame, just cache it
	if fw.last == nil {
		fw.last = f
		return
	}
	// if cached frame is 100ms ago, or different frame type, send the cached frame and cache the new frame
	if diffUint32(f.Timestamp, fw.last.Timestamp) > 100 || f.Type != fw.last.Type {
		err = fw.client.Send(fw.last)
		fw.last = cloneFrame(f)
		return
	}
	// concat frames
	if f.Type == types.ReplayFrameTypeWindowSize {
		fw.last = cloneFrame(f)
	} else {
		fw.last.Payload = concatPayload(fw.last.Payload, f.Payload)
	}
	return
}

func (fw *FrameWriter) Close() (err error) {
	if fw.last != nil {
		fw.client.Send(fw.last)
	}
	_, err = fw.client.CloseAndRecv()
	return
}

type RecordedWriter struct {
	w         io.Writer
	sessionId int64
	typ       uint32
	start     time.Time
	fr        *FrameWriter
}

func (w *RecordedWriter) Write(p []byte) (int, error) {
	var err error
	if err = w.fr.WriteFrame(&types.ReplayFrame{
		SessionId: w.sessionId,
		Timestamp: timestamp(w.start),
		Type:      w.typ,
		Payload:   p,
	}); err != nil {
		log.Error().Err(err).Msg("failed to write replay frame")
	}
	return w.w.Write(p)
}

func (w *RecordedWriter) Close() error {
	if c, ok := w.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func NewRecordedWriter(w io.Writer, sessionId int64, typ uint32, start time.Time, client *FrameWriter) io.Writer {
	return &RecordedWriter{
		w:         w,
		sessionId: sessionId,
		typ:       typ,
		start:     start,
		fr:        client,
	}
}
