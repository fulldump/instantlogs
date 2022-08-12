package blockchain

import (
	"fmt"
	"io"
	"runtime"
	"sync"

	"instantlogs/blocks"
)

type BlockChain struct {
	MaxBlocks               int
	NumBlocks               int
	firstEntry              *blockNode
	lastEntry               *blockNode
	blocksMutex             sync.Mutex
	blockFactory            func() blocks.Blocker
	callbacksBlockCompleted []func(block blocks.Blocker)
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
		MaxBlocks:    10,
		NumBlocks:    1,
		blockFactory: blockFactory,
		firstEntry:   initialEntry,
		lastEntry:    initialEntry,
	}
}

func (b *BlockChain) Write(p []byte) (n int, err error) {

	// todo: hint: when a block is dropped, all its readers should be moved
	// to the next block

	b.blocksMutex.Lock()
	defer b.blocksMutex.Unlock()

	n, err = b.lastEntry.block.Write(p)
	if err == io.EOF {
		newEntry := &blockNode{
			block: b.blockFactory(),
		}
		b.NumBlocks++
		b.lastEntry.next = newEntry
		for _, callback := range b.callbacksBlockCompleted {
			callback(b.lastEntry.block)
		}
		if b.NumBlocks > b.MaxBlocks { // todo: add callbackBlockDropped
			// Drop first block
			b.NumBlocks--
			b.firstEntry = b.firstEntry.next
			go runtime.GC() // todo: just in case, review the impact of this!
			fmt.Println("Block dropped")
		}
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

func (b *BlockChain) OnBlockCompleted(f func(block blocks.Blocker)) {
	b.callbacksBlockCompleted = append(b.callbacksBlockCompleted, f)
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
