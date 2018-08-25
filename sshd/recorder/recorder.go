package recorder

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/yankeguo/bastion/sshd/sandbox"
	"github.com/yankeguo/bastion/types"
	"github.com/yankeguo/bastion/utils"
	"io"
	"time"
)

func timestamp(start time.Time) uint32 {
	return uint32(time.Now().Sub(start) / time.Millisecond)
}

type Recorder struct {
	c *FrameWriter
}

func StartRecording(opts *sandbox.ExecAttachOptions, sessionID int64, rs types.ReplayServiceClient) io.Closer {
	var err error

	// start time
	start := time.Now()

	// build replay write client
	var rc types.ReplayService_WriteReplayClient
	if rc, err = rs.WriteReplay(context.Background()); err != nil {
		log.Error().Err(err).Int64("sessionId", sessionID).Msg("failed to create replay write stream, dummy closer is returned")
		return utils.DummyCloser
	}

	rec := &Recorder{
		c: NewFrameWriter(rc),
	}

	// if opts.WindowChan is not nil, replace it
	if opts.WindowChan != nil {
		oWch := opts.WindowChan
		nWch := make(chan sandbox.Window, 4)
		go func() {
			// iterate original chan
			for w := range oWch {
				if err = rec.c.WriteFrame(&types.ReplayFrame{
					SessionId: sessionID,
					Timestamp: timestamp(start),
					Type:      types.ReplayFrameTypeWindowSize,
					Payload:   utils.MarshalReplayFrameWindowSizePayload(uint32(w.Width), uint32(w.Height)),
				}); err != nil {
					log.Error().Err(err).Int64("sessionId", sessionID).Msg("failed to send window-size record frame")
				}
				// proxy channel
				nWch <- w
			}
			close(nWch)
		}()
		opts.WindowChan = nWch
	}
	// if opts.Stdout is not nil, replace it
	if opts.Stdout != nil {
		opts.Stdout = NewRecordedWriter(opts.Stdout, sessionID, types.ReplayFrameTypeStdout, start, rec.c)
	}
	// if opts.Stderr is not nil, replace it
	if opts.Stderr != nil {
		opts.Stderr = NewRecordedWriter(opts.Stderr, sessionID, types.ReplayFrameTypeStderr, start, rec.c)
	}

	return rec
}

func (r *Recorder) Close() error {
	return r.c.Close()
}
