package service

import (
	"io"
	"testing"

	. "github.com/fulldump/biff"
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

	s := NewService()

	s.Data = []byte("Line 1\nLine 2\nLine 3\n")
	s.Size = len(s.Data)

	readerData, readerErr := io.ReadAll(s.newReader())
	AssertEqual(string(readerData), "Line 1\nLine 2\nLine 3\n")
	AssertNil(readerErr)
}
