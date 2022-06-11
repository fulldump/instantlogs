package main

import (
	"fmt"
	"github.com/fulldump/goconfig"
	"io"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Server string // sever ingesting endpoint
	File   string // log file to read
}

func main() {

	c := &Config{
		Server: "http://localhost:8080/ingest",
	}
	goconfig.Read(&c)

	if stdinHasData() {
		fmt.Println("INFO: gathering data from stdin")
		sendStream(io.NopCloser(os.Stdin), c.Server)
		return
	}

	if c.File != "" {
		fmt.Println("INFO: gathering data from file")
		for {
			f, err := os.Open(c.File)
			if err != nil {
				fmt.Println("ERROR: Open:", err.Error())
				time.Sleep(1 * time.Second) // todo: hardcoded magic number
				continue
			}
			sendStream(io.NopCloser(NewFollowRead(f)), c.Server)
		}
		return
	}

	fmt.Println("ERROR: bad configured")
	os.Exit(-9)
}

func stdinHasData() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}

func sendStream(r io.Reader, endpoint string) error {
	for {
		fmt.Printf("INFO: sending data to %s...\n", endpoint)

		// TODO: put buffer to avoid blocking

		req, err := http.NewRequest(http.MethodPost, endpoint, r)
		if err != nil {
			return fmt.Errorf("NewRequest: %w", err)
		}

		resp, err := http.DefaultClient.Do(req) // todo: optimize http client for this use
		if err != nil {
			fmt.Errorf("ERROR: SendRequest: %w", err)
			time.Sleep(1 * time.Second) // todo: magic number hardcoded
			continue
		}
		fmt.Println("INFO: response status:", resp.Status)
	}
	return nil
}

type FollowRead struct {
	file *os.File
}

func NewFollowRead(f *os.File) *FollowRead {
	return &FollowRead{f}
}

func (f *FollowRead) Read(p []byte) (n int, err error) {

	n, err = f.file.Read(p)
	if err == nil {
		return
	}

	if err == io.EOF {
		time.Sleep(1 * time.Second) // todo: hardcoded magic number
		return n, nil
	}

	return
}
