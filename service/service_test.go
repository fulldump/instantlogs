package service

import (
	"bytes"
	"fmt"
	. "github.com/fulldump/biff"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestService_write_HappyPath(t *testing.T) {

	s := NewService()

	s.write([]byte("Line 1\n"))
	s.write([]byte("Line 2\n"))
	s.write([]byte("Line 3\n"))

	AssertEqual(string(s.Data[:s.Size]), "Line 1\nLine 2\nLine 3\n")
}

func TestService_newReader_HappyPath(t *testing.T) {

	s := NewService()

	s.Data = []byte("Line 1\nLine 2\nLine 3\n")
	s.Size = len(s.Data)

	readerData, readerErr := io.ReadAll(s.newReader())
	AssertEqual(string(readerData), "Line 1\nLine 2\nLine 3\n")
	AssertNil(readerErr)
}

func TestService_Ingest_HappyPath(t *testing.T) {

	// Setup
	s := NewService()
	r, w := io.Pipe()
	go func() {
		w.Write([]byte("hello\n"))
		w.Write([]byte("world\n"))
		w.Close()
	}()

	// Run
	s.Ingest(r)

	// Check
	AssertEqual(string(s.Data[:s.Size]), "hello\nworld\n")

}

func TestService_Ingest_LongLines(t *testing.T) {

	// Setup
	s := NewService()

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
	all := s.Data[:s.Size]
	AssertEqual(string(append(fullBLog, fullALog...)), string(all))

}

func TestService_Ingest_LongLivedRequest(t *testing.T) {

	// setup
	s := NewService()

	pr, pw := io.Pipe()
	go s.Ingest(pr)

	// Sending log 1 + delay
	pw.Write([]byte("Log 1\n"))
	time.Sleep(10 * time.Millisecond)
	AssertEqual(string(s.Data[:s.Size]), "Log 1\n")

	// Sending log 2 + delay
	pw.Write([]byte("Log 2\n"))
	time.Sleep(10 * time.Millisecond)
	AssertEqual(string(s.Data[:s.Size]), "Log 1\nLog 2\n")

	// Finish request
	pw.Close()
}

func TestService_Filter_HappyPath(t *testing.T) {
	// setup
	s := NewService()
	s.Data = []byte("Line 1\nLine 2\nLine 3\n")
	s.Size = len(s.Data)

	// run
	buff := &bytes.Buffer{}
	filterErr := s.Filter(buff, []string{"2"}, false)
	AssertNil(filterErr)
	AssertEqual(buff.String(), "Line 2\n")
}

func TestService_ConcurrentWriters(t *testing.T) {

	service := NewService()

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

	AssertEqual(service.Size, len(logsSample)*concurrentWriters)
}

func TestService_Filter_Benchmark(t *testing.T) {

	if os.Getenv("BENCHMARK") == "" {
		t.SkipNow()
	}

	// Setup
	s := NewService()
	s.Data = make([]byte, 6*1024*1024*1024)

	t0 := time.Now()
	line := []byte(strings.Repeat("a", 1023) + "\n")
	maxLines := cap(s.Data)/len(line) - 1
	for i := 0; i < maxLines; i++ {
		s.write(line)
	}
	s.write([]byte("Hello world!\n"))
	fmt.Println("writing lines took:", time.Since(t0))

	// run
	t1 := time.Now()
	output := &bytes.Buffer{}
	filterErr := s.Filter(output, []string{"world"}, false)
	elapsed := time.Since(t1)
	fmt.Println("filter took:", elapsed)
	fmt.Println("lines:", maxLines)
	fmt.Println("throughput (rows per second):", int(float64(maxLines)/elapsed.Seconds()))

	// check
	AssertNil(filterErr)
	AssertEqual(output.String(), "Hello world!\n")
}
