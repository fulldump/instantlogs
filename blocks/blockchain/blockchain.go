package blockchain

import (
	"io"
	"runtime"
	"sync"
	"time"

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
	callbacksBlockDiscarded []func(block blocks.Blocker)
}

func (b *BlockChain) Stats() map[string]interface{} {

	blocksList := []interface{}{}
	for entry := b.firstEntry; entry != nil; entry = entry.next {
		stats := map[string]interface{}{}
		if block, ok := entry.block.(blocks.Stater); ok {
			stats = block.Stats()
		}
		stats["timestamp"] = entry.timestamp
		blocksList = append(blocksList, stats)
	}

	result := map[string]interface{}{
		"type":       "blockchain",
		"num_blocks": b.NumBlocks,
		"max_blocks": b.MaxBlocks,
		"chain":      blocksList,
	}

	return result
}

type blockNode struct {
	block     blocks.Blocker
	stats     blocks.Stater
	next      *blockNode
	timestamp time.Time
}

func New(blockFactory func() blocks.Blocker) *BlockChain {
	initialEntry := &blockNode{
		block:     blockFactory(),
		timestamp: time.Now().UTC(),
	}
	return &BlockChain{
		MaxBlocks:    4,
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
			block:     b.blockFactory(),
			timestamp: time.Now().UTC(),
		}
		b.NumBlocks++
		b.lastEntry.next = newEntry
		for _, callback := range b.callbacksBlockCompleted {
			callback(b.lastEntry.block)
		}
		if b.NumBlocks > b.MaxBlocks {
			nodeDiscarded := b.firstEntry
			for _, callback := range b.callbacksBlockDiscarded {
				callback(nodeDiscarded.block)
			}
			// Drop first block
			b.NumBlocks--
			b.firstEntry = b.firstEntry.next
			go runtime.GC() // todo: just in case, review the impact of this!
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

func (b *BlockChain) OnBlockDiscarded(f func(block blocks.Blocker)) {
	b.callbacksBlockDiscarded = append(b.callbacksBlockDiscarded, f)
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
