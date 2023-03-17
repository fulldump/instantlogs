package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fulldump/goconfig"
)

var VERSION = "dev"

type Config struct {
	Endpoint  string `usage:"[OPTIONAL] Absolute API ingest endpoint"`
	File      string `usage:"Path to log file"`
	LoggerId  string `usage:"Logger id"`
	ApiKey    string `usage:"Credential Api-Key"`
	ApiSecret string `usage:"Credential Api-Secret"`
	Version   bool   `usage:"Show version and exit"`
}

type StdinReader struct {
	Reader io.Reader
}

func (r StdinReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	if err == io.EOF {
		fmt.Println("EOF")
		os.Exit(3)
	}
	return n, err
}

func main() {

	c := &Config{}
	goconfig.Read(&c)

	if c.Version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if c.Endpoint == "" {
		c.Endpoint = "https://instantlogs.io/v1/loggers/" + c.LoggerId + "/ingest"
	}

	if stdinHasData() {
		fmt.Println("INFO: gathering data from stdin")
		sendStream(StdinReader{os.Stdin}, c.Endpoint, c.ApiKey, c.ApiSecret)
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
			sendStream(io.NopCloser(NewFollowRead(f)), c.Endpoint, c.ApiKey, c.ApiSecret)
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

func sendStream(r io.Reader, endpoint, apikey, apisecret string) error {

	for {
		fmt.Printf("INFO: sending data to %s...\n", endpoint)

		// TODO: put buffer to avoid blocking

		req, err := http.NewRequest(http.MethodPost, endpoint, r)
		if err != nil {
			return fmt.Errorf("NewRequest: %w", err)
		}
		if apikey != "" && apisecret != "" {
			req.Header.Set("Api-Key", apikey)
			req.Header.Set("Api-Secret", apisecret)
		}

		resp, err := http.DefaultClient.Do(req) // todo: optimize http client for this use
		if err != nil {
			fmt.Errorf("ERROR: SendRequest: %w", err)
			time.Sleep(1 * time.Second) // todo: magic number hardcoded
			continue
		}
		fmt.Println("INFO: response status:", resp.Status)

		time.Sleep(1*time.Second)
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
