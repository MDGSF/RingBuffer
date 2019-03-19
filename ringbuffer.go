package ringbuffer

// DefaultCacheMaxSize default max ring buffer size, must be 2^n.
const DefaultCacheMaxSize = 32 * 1024

// RingBuffer represents a ring buffer.
// Safe in one producer and one consumer.
type RingBuffer struct {
	buffer  []byte // buffer used to store element.
	maxsize uint32 // buffer maximum size, must be uint32.
	in      uint32 // in index, must be uint32.
	out     uint32 // out index, must be uint32.
}

// Init initializes or clear RingBuffer rb.
func (rb *RingBuffer) Init(maxsize uint32) *RingBuffer {
	rb.buffer = make([]byte, maxsize)
	rb.maxsize = maxsize
	rb.in = 0
	rb.out = 0
	return rb
}

// New returns an initialized RingBuffer.
func New(maxsize uint32) *RingBuffer {
	return new(RingBuffer).Init(maxsize)
}

// Len returns the number of elements of RingBuffer rb.
func (rb *RingBuffer) Len() uint32 {
	return rb.in - rb.out
}

// LeftSpace returns the left space in RingBuffer rb.
func (rb *RingBuffer) LeftSpace() uint32 {
	return rb.maxsize - rb.in + rb.out
}

// Push insert data at the front of RingBuffer rb.
func (rb *RingBuffer) Push(data []byte) int {
	dataLen := uint32(len(data))
	m1 := min(dataLen, rb.LeftSpace())
	m2 := min(m1, rb.maxsize-(rb.in&(rb.maxsize-1)))
	copy(rb.buffer[(rb.in&(rb.maxsize-1)):], data[:m2])
	copy(rb.buffer, data[m2:])
	rb.in += m1
	return int(m1)
}

// Get get data from the end of RingBuffer rb.
func (rb *RingBuffer) Get(data []byte) int {
	dataLen := uint32(len(data))
	m1 := min(dataLen, rb.Len())
	m2 := min(m1, rb.maxsize-(rb.out&(rb.maxsize-1)))
	copy(data, rb.buffer[(rb.out&(rb.maxsize-1)):(rb.out&(rb.maxsize-1))+m2])
	copy(data[m2:], rb.buffer[:m1-m2])
	rb.out += m1
	return int(m1)
}

// min returns minimum num.
func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}
