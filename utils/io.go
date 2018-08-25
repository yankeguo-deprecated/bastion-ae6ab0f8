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
	go CopyWG(c1, c2, wr, &err)
	go CopyWG(c2, c1, wr, &err)
	wr.Wait()
	return
}

func CopyWG(dst io.Writer, src io.Reader, wr *sync.WaitGroup, errOut *error) {
	var err error
	if _, err = io.Copy(dst, src); err != nil {
		if errOut != nil && *errOut == nil {
			*errOut = err
		}
	}
	wr.Done()
	return
}
