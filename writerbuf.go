// User: szh
// Date: 2020/11/3
// Time: 11:02

package glog

import (
	"io"
	"sync"
	"sync/atomic"
)
const (
	defaultBufSize = 4096

)

// buffered output

// Writer implements buffering for an io.Writer object.
// If an error occurs writing to a Writer, no more data will be
// accepted and all subsequent writes, and Flush, will return the error.
// After all data has been written, the client should call the
// Flush method to guarantee all data has been forwarded to
// the underlying io.Writer.
type Writer struct {
	err     error
	buf     []byte
	n       int
	wr      io.Writer
	mu      sync.Mutex
	muw     sync.Mutex
	flushing int32
}

// NewWriterSize returns a new Writer whose buffer has at least the specified
// size. If the argument io.Writer is already a Writer with large enough
// size, it returns the underlying Writer.
func NewWriterSize(w io.Writer, size int) *Writer {
	// Is it already a Writer?
	b, ok := w.(*Writer)
	if ok && len(b.buf) >= size {
		return b
	}
	if size <= 0 {
		size = defaultBufSize
	}
	return &Writer{
		buf: make([]byte, size),
		wr:  w,
	}
}

// NewWriter returns a new Writer whose buffer has the default size.
func NewWriter(w io.Writer) *Writer {
	return NewWriterSize(w, defaultBufSize)
}

// Size returns the size of the underlying buffer in bytes.
func (b *Writer) Size() int { return len(b.buf) }

// Reset discards any unflushed buffered data, clears any error, and
// resets b to write its output to w.
func (b *Writer) Reset(w io.Writer) {
	b.err = nil
	b.n = 0
	b.wr = w
}

func (b *Writer) ResetWriter(w io.Writer) {
	b.wr = w
}

// Flush writes any buffered data to the underlying io.Writer.
func (b *Writer) syncFlush() error {
	if b.err != nil {
		return b.err
	}
	nn := b.n
	if nn == 0 {
		return nil
	}
	//
	n, err := b.wr.Write(b.buf[0:nn])
	if err != nil{
		b.err = err
		if n <= 0{
			return err
		}
	}
	if n < nn{
		err = io.ErrShortWrite
	}
	b.mu.Lock()
	b.err = err
	if n > 0{
		copy(b.buf[0:], b.buf[n:])
		b.n -= n
		if b.n < 0{
			b.n = 0
		}
	}
	b.mu.Unlock()
	return err
}

//
func (b *Writer) Flush() error {
	if atomic.AddInt32(&b.flushing,1) == 1{
		go func(){
			for {
				b.mu.Lock()
				if b.n <= 0 || b.err != nil{
					b.mu.Unlock()
					break
				}
				b.mu.Unlock()
				b.syncFlush()
			}
			atomic.SwapInt32(&b.flushing,0)
		}()
	}
	return nil
}

// Available returns how many bytes are unused in the buffer.
func (b *Writer) Available() int {
	b.mu.Lock()
	n := len(b.buf) - b.n
	b.mu.Unlock()
	return n
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *Writer) Buffered() int {
	b.mu.Lock()
	n := b.n
	b.mu.Unlock()
	return n
}

// Write writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *Writer) Write(p []byte) (nn int, err error) {
	//b.muw.Lock()
	for len(p) > 0 && b.err == nil {
		var n int
		if b.n == 0 {
			// Large write, empty buffer.
			// Write directly from p to avoid copy.
			n, b.err = b.wr.Write(p)
		} else if b.n >= len(b.buf) {
			b.err = b.syncFlush()
		}else {
			//b.Flush()
			b.mu.Lock()
			n = copy(b.buf[b.n:], p)
			b.n += n
			b.mu.Unlock()
		}
		nn += n
		p = p[n:]
	}
	//
	if b.err != nil {
		//b.muw.Unlock()
		return nn, b.err
	}
	//
	if len(p) > 0 {
		b.mu.Lock()
		n := copy(b.buf[b.n:], p)
		b.n += n
		b.mu.Unlock()
		nn += n
	}
	//
	//b.muw.Unlock()
	if b.n > 0{
		b.Flush()
	}
	return nn, nil
}

