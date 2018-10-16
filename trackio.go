package trackio

import (
	"io"
	"sync/atomic"
)

type TrackReader struct {
	r io.Reader
	n int64
}

type TrackWriter struct {
	w io.Writer
	n int64
}

func NewReader(r io.Reader) *TrackReader {
	return &TrackReader{
		r: r,
	}
}

func (r *TrackReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	atomic.AddInt64(&r.n, int64(n))
	return
}

// N gets the number of bytes that have been read
// so far.
func (r *TrackReader) N() int64 {
	return atomic.LoadInt64(&r.n)
}

func NewWriter(w io.Writer) *TrackWriter {
	return &TrackWriter{
		w: w,
	}
}

func (w *TrackWriter) Write(p []byte) (n int, err error) {
	n, err = w.w.Write(p)
	atomic.AddInt64(&w.n, int64(n))
	return
}

// N gets the number of bytes that have been read
// so far.
func (w *TrackWriter) N() int64 {
	return atomic.LoadInt64(&w.n)
}
