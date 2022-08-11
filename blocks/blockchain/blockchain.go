package blockchain

import (
	"io"
	"sync"

	"instantlogs/blocks"
)

type BlockChain struct {
	firstEntry   *blockNode
	lastEntry    *blockNode
	blocksMutex  sync.Mutex
	blockFactory func() blocks.Blocker
}

type blockNode struct {
	block blocks.Blocker
	next  *blockNode
	// todo: add timestamp
}

func New(blockFactory func() blocks.Blocker) *BlockChain {
	initialEntry := &blockNode{
		block: blockFactory(),
	}
	return &BlockChain{
		blockFactory: blockFactory,
		firstEntry:   initialEntry,
		lastEntry:    initialEntry,
	}
}

func (b *BlockChain) Write(p []byte) (n int, err error) {

	b.blocksMutex.Lock()
	defer b.blocksMutex.Unlock()

	n, err = b.lastEntry.block.Write(p)
	if err == io.EOF {
		newEntry := &blockNode{
			block: b.blockFactory(),
		}
		b.lastEntry.next = newEntry
		b.lastEntry = newEntry
		return newEntry.block.Write(p)
	}

	return
}

func (b *BlockChain) NewReader() io.Reader {
	return &blockChainReader{
		currentEntry: b.firstEntry,
		blockReader:  b.firstEntry.block.NewReader(),
		blockChain:   b,
	}
}

type blockChainReader struct {
	currentEntry *blockNode
	blockReader  io.Reader
	blockChain   *BlockChain
}

func (b *blockChainReader) Read(p []byte) (n int, err error) {

	n, err = b.blockReader.Read(p)
	if err == io.EOF {
		if b.currentEntry.next != nil {
			b.currentEntry = b.currentEntry.next
			b.blockReader = b.currentEntry.block.NewReader()
			return b.blockReader.Read(p)
		}
	}

	return
}
