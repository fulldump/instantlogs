package blockchain

import (
	"fmt"
	"io"
	"testing"

	"github.com/fulldump/biff"

	"github.com/fulldump/instantlogs/blocks"
	"github.com/fulldump/instantlogs/blocks/bigblock"
)

func Test_BlockChain_HappyPath(t *testing.T) {

	b := New(func() blocks.Blocker {
		return bigblock.NewWithBuffer(make([]byte, 10))
	})

	b.Write([]byte("hello\n"))
	b.Write([]byte("world\n"))
	b.Write([]byte("whatever\n"))
	b.Write([]byte("whatever22\n"))
	b.Write([]byte("zzz\n"))

	// biff.AssertEqual(len(b.blocks), 4)

	data, err := io.ReadAll(b.NewReader())
	fmt.Println(string(data), err)

}

func Test_BlockChain_MaxLine(t *testing.T) {

	b := New(func() blocks.Blocker {
		return bigblock.NewWithBuffer(make([]byte, 10))
	})

	b.Write([]byte("short\n"))
	b.Write([]byte("exactline\n"))
	b.Write([]byte("looooooooong\n"))

	data, err := io.ReadAll(b.NewReader())
	biff.AssertNil(err)
	biff.AssertEqual(string(data), "short\nexactline\n")
}

func Test_BlockChain_1Kbuffer(t *testing.T) {

	b := New(func() blocks.Blocker {
		return bigblock.NewWithBuffer(make([]byte, 1024))
	})

	b.Write([]byte("hello\n"))
	b.Write([]byte("world\n"))
	b.Write([]byte("whatever\n"))
	b.Write([]byte("whatever22\n"))
	b.Write([]byte("zzz\n"))

	allLogs, err := io.ReadAll(b.NewReader())
	biff.AssertNil(err)
	biff.AssertEqual(string(allLogs), "hello\nworld\nwhatever\nwhatever22\nzzz\n")

}

func Test_BlockChain_OnBlockCompleted(t *testing.T) {

	b := New(func() blocks.Blocker {
		return bigblock.NewWithBuffer(make([]byte, 10))
	})

	expectedBlocks := []string{
		"hello\n",
		"world!!!\n",
		"a\nb\n",
	}

	b.OnBlockCompleted(func(block blocks.Blocker) {
		blockCompleted, blockCompletedErr := io.ReadAll(block.NewReader())
		biff.AssertNil(blockCompletedErr)
		biff.AssertEqual(string(blockCompleted), expectedBlocks[0])
		expectedBlocks = expectedBlocks[1:]
	})

	b.Write([]byte("hello\n"))
	b.Write([]byte("world!!!\n"))
	b.Write([]byte("a\n"))
	b.Write([]byte("b\n"))
	b.Write([]byte("whatever\n"))

}
