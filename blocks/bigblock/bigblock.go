package bigblock

import (
	"io"
	"sync"
)

type BigBlock struct {
	Buffer      []byte
	BufferSize  int
	BufferMutex sync.Mutex
}

func (b *BigBlock) Stats() map[string]interface{} {
	return map[string]interface{}{
		"cap":  cap(b.Buffer),
		"len":  b.BufferSize,
		"type": "bigblock",
	}
}

func New() *BigBlock {
	return NewWithBuffer(make([]byte, 1*1024*1024))
}

func NewWithBuffer(buffer []byte) *BigBlock {
	return &BigBlock{
		Buffer: buffer,
	}
}

func (b *BigBlock) Write(p []byte) (n int, err error) {
	l := len(p)

	b.BufferMutex.Lock()
	defer b.BufferMutex.Unlock()

	if cap(b.Buffer) < b.BufferSize+l {
		return 0, io.EOF
	}

	n = copy(b.Buffer[b.BufferSize:b.BufferSize+l], p)
	b.BufferSize += l

	return n, nil
}

func (b *BigBlock) NewReader() io.Reader {
	return &bigBlockReader{
		nextByte: 0,
		bigBlock: b,
	}
}

type bigBlockReader struct {
	nextByte int
	bigBlock *BigBlock
}

func (b *bigBlockReader) Read(p []byte) (n int, err error) {

	pending := b.bigBlock.BufferSize - b.nextByte
	if pending <= 0 {
		return 0, io.EOF
	}

	len := len(p)
	if len > pending { // todo: change to >=
		len = pending
	}

	n = copy(p, b.bigBlock.Buffer[b.nextByte:b.nextByte+len])
	b.nextByte += n

	return
}
