package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fulldump/biff"
	"github.com/fulldump/box"

	"instantlogs/service"
)

// TODO: add concurrency tests here...

func TestNewApi_HappyPath(t *testing.T) {

	// setup
	api := NewApi(service.NewService(), "")
	mockserver := httptest.NewServer(box.Box2Http(api))
	defer mockserver.Close()

	// run
	http.Post(mockserver.URL+"/ingest", "text/plain", strings.NewReader("log1\nlog2\nlog3\n"))

	// check
	resp, _ := http.Get(mockserver.URL + "/filter")
	respBody, _ := io.ReadAll(resp.Body)
	biff.AssertEqual(string(respBody), "log1\nlog2\nlog3\n")
}
