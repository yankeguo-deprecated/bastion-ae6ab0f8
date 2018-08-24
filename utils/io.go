package utils

import (
	"io"
	"sync"
)

var (
	DummyCloser io.Closer = dummyCloser(0)
)

type dummyCloser int

func (d dummyCloser) Close() error {
	return nil
}

func DualCopy(c1 io.ReadWriteCloser, c2 io.ReadWriteCloser) (err error) {
	wr := &sync.WaitGroup{}
	wr.Add(2)
	go func() {
		if _, err1 := io.Copy(c1, c2); err1 != nil {
			err = err1
		}
		wr.Done()
	}()
	go func() {
		if _, err1 := io.Copy(c2, c1); err1 != nil {
			err = err1
		}
		wr.Done()
	}()
	wr.Wait()
	return
}
