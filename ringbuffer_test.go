package ringbuffer

import "testing"

func testBasic(t *testing.T, rb *RingBuffer, expectLen uint32, expectLeftSpace uint32) {
	if rb.Len() != expectLen {
		t.Fatal(rb.Len())
	}
	if rb.LeftSpace() != expectLeftSpace {
		t.Fatal(rb.LeftSpace())
	}
}

func TestInit(t *testing.T) {
	rb := New(DefaultCacheMaxSize)
	if rb.maxsize != DefaultCacheMaxSize {
		t.Fatal(rb.maxsize, DefaultCacheMaxSize)
	}
	testBasic(t, rb, 0, DefaultCacheMaxSize)

	in := []byte{'a', 'b'}
	ret := rb.Push(in)
	if ret != 2 {
		t.Fatal(ret)
	}
	testBasic(t, rb, 2, DefaultCacheMaxSize-2)
	if rb.buffer[0] != 'a' || rb.buffer[1] != 'b' {
		t.Fatal(rb.buffer)
	}

	out := make([]byte, 2)
	ret = rb.Get(out)
	if ret != 2 {
		t.Fatal(ret)
	}
	testBasic(t, rb, 0, DefaultCacheMaxSize)
	if out[0] != 'a' || out[1] != 'b' {
		t.Fatal(out)
	}
}
