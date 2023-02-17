package service

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/fulldump/instantlogs/blocks"
)

type Service struct {
	block blocks.Blocker

	// stats
	totalBytesSent     int64
	totalBytesReceived int64
	totalBytesFiltered int64
	lastRegexps        []string
	concurrentFilters  int
	concurrentIngests  int
	totalFilters       int
	totalIngests       int
}

func NewService(block blocks.Blocker) *Service {
	return &Service{
		block: block,
	}
}

func (s *Service) Filter(w io.Writer, regexps []string, follow *bool) error {

	s.lastRegexps = append(s.lastRegexps, strings.Join(regexps, " AND "))
	for len(s.lastRegexps) > 10 {
		s.lastRegexps = s.lastRegexps[1:]
	}

	s.totalFilters++
	s.concurrentFilters++
	defer func() { s.concurrentFilters-- }()

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

		s.totalBytesFiltered += int64(len(line))

		match := true
		for _, r := range compiledRegexps {
			if !r.Match(line) {
				match = false
				break
			}
		}
		if match {
			n, err := w.Write(line) // todo: handle error
			if err != nil {
				fmt.Println("ERR Filtering:", err.Error())
			}
			s.totalBytesSent += int64(n)
		}

		if readErr == io.EOF || len(line) == 0 {
			if *follow {
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

	s.totalIngests++
	s.concurrentIngests++
	defer func() { s.concurrentIngests-- }()

	breader := bufio.NewReader(reader)
	for {
		// todo: use scanner?
		data, readErr := breader.ReadBytes('\n') // TODO: split also by \r and \r\n to support mac and win log styles
		if readErr == io.EOF {
			n, err := s.block.Write(data) // Write line...
			totaln += n
			s.totalBytesReceived += int64(n)
			return totaln, err
		}
		if readErr != nil {
			err = fmt.Errorf("read: %w", readErr)
			return
		}
		n, err := s.block.Write(data) // Write line...
		totaln += n
		s.totalBytesReceived += int64(n)
		if err != nil {
			return totaln, err
		}
	}

	return
}

func (s *Service) Stats() map[string]interface{} {
	result := map[string]interface{}{
		"total_bytes_sent":     s.totalBytesSent,
		"total_bytes_received": s.totalBytesReceived,
		"total_bytes_filtered": s.totalBytesFiltered,
		"last_regexps":         s.lastRegexps,
		"concurrent_filters":   s.concurrentFilters,
		"concurrent_ingests":   s.concurrentIngests,
		"total_filters":        s.totalFilters,
		"total_ingests":        s.totalIngests,
	}

	if stats, ok := s.block.(blocks.Stater); ok {
		result["block"] = stats.Stats()
	}

	return result
}
