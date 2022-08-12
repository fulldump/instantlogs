package blocks

import "io"

// Blocker allows multiple writers and multiple readers
type Blocker interface {
	io.Writer
	NewReader() io.Reader // todo: review the name for a better understanding
}

// Future
//type Liner interface {
//	WriteLine(line []byte) (err error)
//	NewLineReader() LineReader
//}
//
//type LineReader interface {
//	ReadLine() (line []byte, err error)
//}

// Service interface
type Service interface {
	Ingest(r io.Reader) (n int, err error)
	Filter(w io.Writer, regexps []string, follow bool) (err error) // todo: explore composition on this
}
