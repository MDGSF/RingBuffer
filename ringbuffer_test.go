package ringbuffer

import (
	"fmt"
	"testing"
)

func init() {
	fmt.Printf("a = %d\n", 'a')
	fmt.Printf("b = %d\n", 'b')
	fmt.Printf("c = %d\n", 'c')
	fmt.Printf("d = %d\n", 'd')
}

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

func TestOverflow(t *testing.T) {
	rb := New(4)
	testBasic(t, rb, 0, 4)

	in := []byte{'a', 'b', 'c', 'd', 'e'}
	ret := rb.Push(in)
	if ret != 4 {
		t.Fatal(ret)
	}
	testBasic(t, rb, 4, 0)
	if rb.buffer[0] != 'a' ||
		rb.buffer[1] != 'b' ||
		rb.buffer[2] != 'c' ||
		rb.buffer[3] != 'd' {
		t.Fatal(rb.buffer)
	}

	out := make([]byte, 1)
	ret = rb.Get(out)
	if ret != 1 {
		t.Fatal(ret)
	}
	testBasic(t, rb, 3, 1)
	if out[0] != 'a' {
		t.Fatal(out)
	}

	in2 := []byte{'x', 'y'}
	ret = rb.Push(in2)
	if ret != 1 {
		t.Fatal(ret)
	}
	testBasic(t, rb, 4, 0)
	if rb.buffer[0] != 'x' ||
		rb.buffer[1] != 'b' ||
		rb.buffer[2] != 'c' ||
		rb.buffer[3] != 'd' {
		t.Fatal(rb.buffer)
	}

	out2 := make([]byte, 5)
	ret = rb.Get(out2)
	if ret != 4 {
		t.Fatal(ret)
	}
	testBasic(t, rb, 0, 4)
	if out2[0] != 'b' ||
		out2[1] != 'c' ||
		out2[2] != 'd' ||
		out2[3] != 'x' {
		t.Fatal(out2)
	}
}
