package circularbuffer

import (
	"testing"

	. "github.com/fulldump/biff"
)

func TestCircularBuffer_Write_1_line_fits(t *testing.T) {

	circularBuffer := NewCircularBuffer(16)

	n, err := circularBuffer.Write([]byte("hello\n"))
	AssertNil(err)
	AssertEqual(n, 6)
	AssertEqual(string(circularBuffer.data[0:6]), "hello\n")
	AssertEqual(circularBuffer.free, 10)
	AssertEqual(circularBuffer.read, 0)
	AssertEqual(circularBuffer.write, 6)
}

func TestCircularBuffer_Write_2_lines_fit(t *testing.T) {

	circularBuffer := NewCircularBuffer(16)

	circularBuffer.Write([]byte("hello\n"))
	circularBuffer.Write([]byte("world!\n"))
	AssertEqual(string(circularBuffer.data[0:13]), "hello\nworld!\n")
	AssertEqual(circularBuffer.free, 3)
	AssertEqual(circularBuffer.read, 0)
	AssertEqual(circularBuffer.write, 13)
}

func TestCircularBuffer_Write_2_lines_fit_1_line_dont_fit(t *testing.T) {

	circularBuffer := NewCircularBuffer(16)

	circularBuffer.Write([]byte("hello\n"))  // 6 chars
	circularBuffer.Write([]byte("world!\n")) // 7 chars
	circularBuffer.Write([]byte("good\n"))   // 5 chars
	AssertEqual(string(circularBuffer.data), "d\nllo\nworld!\ngoo")
	AssertEqual(circularBuffer.free, 4)
	AssertEqual(circularBuffer.read, 6)
	AssertEqual(circularBuffer.write, 2)
}

func TestCircularBuffer_FreeLineBrokenAtEnd(t *testing.T) {

	//circularBuffer := &CircularBuffer{
	//	data:  []byte("world!\nXXXXXXXLine1\nLine2\nHello "),
	//	mutex: &sync.Mutex{},
	//	read: 13,
	//	write: 25,
	//	free:  0,
	//	size:  0,
	//}

	circularBuffer := NewCircularBuffer(32)

	// first pass
	circularBuffer.Write([]byte("Line1\n"))
	circularBuffer.Write([]byte("Line2\n"))
	circularBuffer.Write([]byte("Line3\n"))
	circularBuffer.Write([]byte("Line4\n"))
	circularBuffer.Write([]byte("Line5\n"))

	// broken line
	circularBuffer.Write([]byte("Hello world\n"))

	// write until broken line
	circularBuffer.Write([]byte("Line1\n"))
	circularBuffer.Write([]byte("Line2\n"))
	circularBuffer.Write([]byte("Line3\n"))
	circularBuffer.Write([]byte("Line4\n"))
	circularBuffer.Write([]byte("Line5\n"))

	AssertEqual(string(circularBuffer.data), "4\nLine5\nd\nLine1\nLine2\nLine3\nLine")
	AssertEqual(circularBuffer.free, 2)
	AssertEqual(circularBuffer.read, 10)
	AssertEqual(circularBuffer.write, 8)
}
