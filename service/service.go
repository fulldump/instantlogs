package service

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"instantlogs/blocks"
)

type Service struct {
	block blocks.Blocker
}

func NewService(block blocks.Blocker) *Service {
	return &Service{
		block: block,
	}
}

func (s *Service) Filter(w io.Writer, regexps []string, follow bool) error {

	// Multiple regexps
	compiledRegexps := make([]*regexp.Regexp, len(regexps))
	for i, e := range regexps {
		r, err := regexp.Compile(e)
		if err != nil {
			return fmt.Errorf("bad regexp '%s': %w", e, err)
		}
		compiledRegexps[i] = r
	}

	// Workaround, embed into io.Writer?
	flusher, flusherOk := w.(http.Flusher)

	lines := bufio.NewReader(s.block.NewReader())

	for {
		line, readErr := lines.ReadBytes('\n')

		match := true
		for _, r := range compiledRegexps {
			if !r.Match(line) {
				match = false
				break
			}
		}
		if match {
			w.Write(line) // todo: handle error
		}

		if readErr == io.EOF || len(line) == 0 {
			if follow {
				// Workaround part2
				if flusherOk {
					flusher.Flush()
				}
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return nil
		}

	}

	return nil
}

func (s *Service) Ingest(reader io.Reader) (totaln int, err error) {

	breader := bufio.NewReader(reader)
	for {
		// todo: use scanner?
		data, readErr := breader.ReadBytes('\n') // TODO: split also by \r and \r\n to support mac and win log styles
		if readErr == io.EOF {
			n, err := s.block.Write(data) // Write line...
			totaln += n
			return totaln, err
		}
		if readErr != nil {
			err = fmt.Errorf("read: %w", readErr)
			return
		}
		n, err := s.block.Write(data) // Write line...
		totaln += n
		if err != nil {
			return totaln, err
		}
	}

	return
}
