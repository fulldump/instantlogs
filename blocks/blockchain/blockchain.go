package blockchain

import (
	"io"
	"sync"

	"instantlogs/blocks"
)

type BlockChain struct {
	blocksMutex sync.Mutex

	blockFactory func() blocks.Blocker

	blocks []blocks.Blocker // TODO: make an intermediate struct to storage timestamp, etc
}

func New(blockFactory func() blocks.Blocker) *BlockChain {
	return &BlockChain{
		blockFactory: blockFactory,
		blocks:       []blocks.Blocker{blockFactory()},
	}
}

//type blockElement struct {
//	blocker blocks.Blocker
//
//}

func (b *BlockChain) Write(p []byte) (n int, err error) {

	b.blocksMutex.Lock()
	defer b.blocksMutex.Unlock()

	lastBlock := b.blocks[len(b.blocks)-1]

	n, err = lastBlock.Write(p)
	if err == io.EOF {
		newBlock := b.blockFactory()
		b.blocks = append(b.blocks, newBlock)
		n, err = newBlock.Write(p)
	}

	return
}

func (b *BlockChain) NewReader() io.Reader {
	return &blockChainReader{
		nextBlock:   1,
		blockReader: b.blocks[0].NewReader(),
		blockChain:  b,
	}
}

type blockChainReader struct {
	nextBlock   int
	blockReader io.Reader
	blockChain  *BlockChain
}

func (b *blockChainReader) Read(p []byte) (n int, err error) {

	n, err = b.blockReader.Read(p)
	if err == io.EOF {
		if len(b.blockChain.blocks) > b.nextBlock {
			b.blockReader = b.blockChain.blocks[b.nextBlock].NewReader()
			b.nextBlock++
			return b.blockReader.Read(p)
		}
	}

	return
}
