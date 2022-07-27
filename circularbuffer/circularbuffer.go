package circularbuffer

import (
	"bytes"
	"sync"
)

/*

# first loop
data: [XXXXXXXnXXXXXXXXXXX                         ]
       |                  |
       start              write

# full data
data: [XXXXXXXnXXXXXXXXnXXXXXXXXXnXXXXXXXXXXXXXXX  ]
       |                                         |
       start                                     write

# first overwrite: allocate line
data: [XXXXXXXnXXXXXXXXnXXXXXXXXXnXXXXXXXXXXXXXXX  ]
               |                                 |
               start                             write

# first overwrite: writing new line
data: [YYYYXXXnXXXXXXXXnXXXXXXXXXnXXXXXXXXXXXXXXXYY]
           |   |
           |   start
           write

Algorithm
lock {
    // Allocate space for the next line
	for len(newline) > free {
		// find next '\n' <--- careful, can be 2 operations also!
        start = next(\n) + 1
	}
    copy(data[start:], newline)
}

*/

type CircularBuffer struct {
	data  []byte
	mutex *sync.Mutex // to protect start and write
	start int         // todo: rename to read
	write int
	free  int
	size  int
}

func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		data:  make([]byte, size),
		mutex: &sync.Mutex{},
		start: 0,
		write: 0,
		free:  size,
		size:  size,
	}
}

func (c *CircularBuffer) Write(line []byte) (n int, err error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	l := len(line)
	for l > c.free {
		if c.start < c.write {
			nextLineBreak := bytes.IndexByte(c.data[c.start:c.write], '\n')
			if nextLineBreak < 0 {
				panic("Alfonso pays")
			}
			nextLineBreak++
			c.start += nextLineBreak
			c.free += nextLineBreak
			continue
		}
		nextLineBreak := bytes.IndexByte(c.data[c.start:], '\n')
		if nextLineBreak >= 0 {
			nextLineBreak++
			c.start += nextLineBreak
			c.free += nextLineBreak
			continue
		}
		nextLineBreak = bytes.IndexByte(c.data[0:c.write], '\n')
		if nextLineBreak < 0 {
			panic("Alfonso pays")
		}
		nextLineBreak++
		c.start = nextLineBreak
		c.free += nextLineBreak + (c.size - c.start)
	}

	n = copy(c.data[c.write:], line)
	c.write += n
	c.free -= n

	if n < l {
		c.write = 0
		n = copy(c.data[c.write:], line[n:])
		c.write += n
		c.free -= n
		n = l
	}

	return
}
