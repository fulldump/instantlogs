package service

import (
	"bufio"
	"fmt"
	"io"
	"sync"
)

type Service struct {
	Data  []byte
	Size  int
	Mutex *sync.Mutex // to protect Size
}

func NewService() *Service {
	return &Service{
		Data:  make([]byte, 0, 100*1024*1024), // TODO: make this configurable
		Size:  0,                              // Amount of used bytes
		Mutex: &sync.Mutex{},
	}
}

// write is internal helper to handle concurrent writes
func (s *Service) write(p []byte) (int, error) { // info: implement io.Writer interface

	l := len(p)

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	n := copy(s.Data[s.Size:s.Size+l], p) // todo: handle returning value?
	s.Size += l

	return n, nil
}

// newReader is internal helper to handler concurrent and independent reads
func (s *Service) newReader() *storageReader { // return io.Reader
	return &storageReader{
		service: s,
	}
}

func (s *Service) Filter(w io.Writer, regexps []string, follow bool) error {

	return nil
}

func (s *Service) Ingest(reader io.Reader) (totaln int, err error) {

	breader := bufio.NewReader(reader)
	for {
		// todo: use scanner?
		data, readErr := breader.ReadBytes('\n')
		if readErr == io.EOF {
			n, err := s.write(data)
			totaln += n
			return totaln, err
		}
		if readErr != nil {
			err = fmt.Errorf("read: %w", readErr)
			return
		}
		n, err := s.write(data)
		totaln += n
		if err != nil {
			return totaln, err
		}
	}

	return
}

type storageReader struct {
	nextByte int
	service  *Service
}

func (r *storageReader) Read(p []byte) (n int, err error) { // info: implement io.Reader interface

	pending := r.service.Size - r.nextByte
	if pending == 0 {
		return 0, io.EOF
	}

	len := len(p)
	if len > pending {
		len = pending
	}

	n = copy(p, r.service.Data[r.nextByte:r.nextByte+len])
	r.nextByte += n

	return
}
