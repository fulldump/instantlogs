package service

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	. "github.com/fulldump/biff"

	"github.com/fulldump/instantlogs/blocks"
	"github.com/fulldump/instantlogs/blocks/bigblock"
	"github.com/fulldump/instantlogs/blocks/blockchain"
)

func TestService_Ingest_HappyPath(t *testing.T) {

	// Setup
	s := NewService(bigblock.New())
	r, w := io.Pipe()
	go func() {
		w.Write([]byte("hello\n"))
		w.Write([]byte("world\n"))
		w.Close()
	}()

	// Run
	s.Ingest(r)

	// Check
	blockData, blockDataErr := io.ReadAll(s.block.NewReader())
	AssertEqual(string(blockData), "hello\nworld\n")
	AssertNil(blockDataErr)

}

func TestService_Ingest_LongLines(t *testing.T) {

	// Setup
	bb := bigblock.New()
	s := NewService(bb)

	ar, aw := io.Pipe()
	br, bw := io.Pipe()

	go s.Ingest(ar)
	go s.Ingest(br)

	firstALog := bytes.Repeat([]byte{'a'}, 32*1024+100)
	secondALog := []byte("secondpart\n")
	fullALog := append(firstALog, secondALog...)
	fullBLog := []byte("B Log\n")

	aw.Write(firstALog)
	time.Sleep(10 * time.Millisecond)
	bw.Write(fullBLog)
	time.Sleep(10 * time.Millisecond)
	aw.Write(secondALog)
	time.Sleep(10 * time.Millisecond)

	// Finish request
	aw.Close()
	bw.Close()

	// check
	bigBlockData, bigBlockErr := io.ReadAll(bb.NewReader())
	AssertEqual(string(append(fullBLog, fullALog...)), string(bigBlockData))
	AssertNil(bigBlockErr)

}

func TestService_Ingest_LongLivedRequest(t *testing.T) {

	// setup
	bb := bigblock.New()
	s := NewService(bb)

	pr, pw := io.Pipe()
	go s.Ingest(pr)

	// Sending log 1 + delay
	pw.Write([]byte("Log 1\n"))
	time.Sleep(10 * time.Millisecond)
	AssertEqual(string(bb.Buffer[:bb.BufferSize]), "Log 1\n")

	// Sending log 2 + delay
	pw.Write([]byte("Log 2\n"))
	time.Sleep(10 * time.Millisecond)
	AssertEqual(string(bb.Buffer[:bb.BufferSize]), "Log 1\nLog 2\n")

	// Finish request
	pw.Close()
}

func TestService_Filter_HappyPath(t *testing.T) {
	// setup
	bb := bigblock.New()
	s := NewService(bb)
	bb.Write([]byte("Line 1\nLine 2\nLine 3\n"))

	// run
	follow := false
	buff := &bytes.Buffer{}
	filterErr := s.Filter(buff, []string{"2"}, &follow)
	AssertNil(filterErr)
	AssertEqual(buff.String(), "Line 2\n")
}

func TestService_ConcurrentWriters(t *testing.T) {

	bc := blockchain.New(func() blocks.Blocker {
		return bigblock.New()
	})
	bc.MaxBlocks = 100
	service := NewService(bc)

	logLine := strings.Repeat("a", 1024) + "\n"
	logsSample := strings.Repeat(logLine, 10)

	concurrentWriters := 10000
	wg := &sync.WaitGroup{}
	for i := 0; i < concurrentWriters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			service.Ingest(strings.NewReader(logsSample))
		}()
	}
	wg.Wait()

	n, err := io.Copy(io.Discard, bc.NewReader())
	AssertNil(err)
	AssertEqual(int(n), len(logsSample)*concurrentWriters)
}

func TestService_Filter_Benchmark(t *testing.T) {

	if os.Getenv("BENCHMARK") == "" {
		t.SkipNow()
	}

	// Setup
	bb := bigblock.NewWithBuffer(make([]byte, 6*1024*1024*1024))
	s := NewService(bb)

	t0 := time.Now()
	line := []byte(strings.Repeat("a", 1023) + "\n")
	maxLines := cap(bb.Buffer)/len(line) - 1
	for i := 0; i < maxLines; i++ {
		bb.Write(line)
	}
	bb.Write([]byte("Hello world!\n"))
	fmt.Println("writing lines took:", time.Since(t0))

	// run
	t1 := time.Now()
	follow := false
	output := &bytes.Buffer{}
	filterErr := s.Filter(output, []string{"world"}, &follow)
	elapsed := time.Since(t1)
	fmt.Println("filter took:", elapsed)
	fmt.Println("lines:", maxLines)
	fmt.Println("throughput (rows per second):", int(float64(maxLines)/elapsed.Seconds()))

	// check
	AssertNil(filterErr)
	AssertEqual(output.String(), "Hello world!\n")
}
