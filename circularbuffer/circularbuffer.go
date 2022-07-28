package circularbuffer

import (
	"bytes"
	"sync"
)

/*

# first loop
data: [XXXXXXXnXXXXXXXXXXX                         ]
       |                  |
       read              write

# full data
data: [XXXXXXXnXXXXXXXXnXXXXXXXXXnXXXXXXXXXXXXXXX  ]
       |                                         |
       read                                     write

# first overwrite: allocate line
data: [XXXXXXXnXXXXXXXXnXXXXXXXXXnXXXXXXXXXXXXXXX  ]
               |                                 |
               read                             write

# first overwrite: writing new line
data: [YYYYXXXnXXXXXXXXnXXXXXXXXXnXXXXXXXXXXXXXXXYY]
           |   |
           |   read
           write

Algorithm
lock {
    // Allocate space for the next line
	for len(newline) > free {
		// find next '\n' <--- careful, can be 2 operations also!
        read = next(\n) + 1
	}
    copy(data[read:], newline)
}

*/

type CircularBuffer struct {
	data  []byte
	mutex *sync.Mutex // to protect read and write
	read  int         // todo: rename to read
	write int
	free  int
	size  int
}

func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		data:  make([]byte, size),
		mutex: &sync.Mutex{},
		read:  0,
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
		if c.read < c.write {
			nextLineBreak := bytes.IndexByte(c.data[c.read:c.write], '\n')
			if nextLineBreak < 0 {
				panic("Alfonso pays")
			}
			nextLineBreak++
			c.read += nextLineBreak
			c.free += nextLineBreak
			continue
		}
		nextLineBreak := bytes.IndexByte(c.data[c.read:], '\n')
		if nextLineBreak >= 0 {
			nextLineBreak++
			c.read += nextLineBreak
			c.free += nextLineBreak
			continue
		}
		nextLineBreak = bytes.IndexByte(c.data[0:c.write], '\n')
		if nextLineBreak < 0 {
			panic("Alfonso pays")
		}
		//
		nextLineBreak++
		c.read = nextLineBreak               // resets c.read
		c.free = c.size - (c.write - c.read) // todo review this expression
	}

	n = copy(c.data[c.write:], line)
	c.write += n
	c.free -= n

	if n < l {
		c.write = 0 // resets c.write
		n = copy(c.data[c.write:], line[n:])
		c.write += n
		c.free -= n
		n = l
	}

	return
}
