package bigblock

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestBigBlock_HappyPath(t *testing.T) {

	bb := New()
	r := bb.NewReader()

	bb.Write([]byte("hello\n"))
	io.Copy(os.Stdout, r)

	bb.Write([]byte("world\n"))
	io.Copy(os.Stdout, r)

	r2 := bb.NewReader()
	io.Copy(os.Stdout, r2)
}

func TestBigBlock_EndOfFile(t *testing.T) {

	bb := NewWithBuffer(make([]byte, 10))

	bb.Write([]byte("hello\n"))

	n, err := bb.Write([]byte("world\n"))

	fmt.Println(n, err)
}
