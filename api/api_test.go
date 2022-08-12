package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fulldump/biff"
	"github.com/fulldump/box"

	"instantlogs/blocks"
	"instantlogs/blocks/bigblock"
	"instantlogs/blocks/blockchain"
	"instantlogs/service"
)

// TODO: add concurrency tests here...

func TestNewApi_HappyPath(t *testing.T) {

	// setup
	bc := blockchain.New(func() blocks.Blocker {
		return bigblock.NewWithBuffer(make([]byte, 1*1024*1024))
	})
	api := NewApi(service.NewService(bc), "")
	mockserver := httptest.NewServer(box.Box2Http(api))
	defer mockserver.Close()

	// run
	http.Post(mockserver.URL+"/ingest", "text/plain", strings.NewReader("log1\nlog2\nlog3\n"))

	// check
	resp, _ := http.Get(mockserver.URL + "/filter")
	respBody, _ := io.ReadAll(resp.Body)
	biff.AssertEqual(string(respBody), "log1\nlog2\nlog3\n")
}
