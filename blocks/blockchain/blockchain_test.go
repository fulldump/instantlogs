package blockchain

import (
	"github.com/fulldump/biff"
	"instantlogs/blocks"
	"instantlogs/blocks/bigblock"
	"testing"
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

	biff.AssertEqual(len(b.blocks), 4)

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

	biff.AssertEqual(len(b.blocks), 1)

	//data, _ := io.ReadAll(b.blocks[0].(bigblock.BigBlock).Buffer)

	biff.AssertEqual(string(b.blocks[0].(*bigblock.BigBlock).Buffer[0:36]),
		"hello\nworld\nwhatever\nwhatever22\nzzz\n")

}
