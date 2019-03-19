package ringbuffer

import (
	"math"
	"testing"
	"time"
)

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

func TestInOutOverFlow1(t *testing.T) {
	/*
		0 1 2 3 ..... max-2, max-1, max
	*/
	rb := New(8)
	rb.in = 1
	rb.out = math.MaxUint32 - 1
	testBasic(t, rb, 3, 5)

	in := []byte{'a'}
	ret := rb.Push(in)
	if ret != 1 {
		t.Fatal(ret)
	}
	testBasic(t, rb, 4, 4)

	out := make([]byte, 5)
	ret = rb.Get(out)
	if ret != 4 {
		t.Fatal(ret)
	}
	testBasic(t, rb, 0, 8)
}

func TestInOutOverFlow2(t *testing.T) {
	rb := New(8)
	rb.in = math.MaxUint32
	rb.out = math.MaxUint32 - 1
	testBasic(t, rb, 1, 7)

	idx := rb.in % 8

	in := []byte{'a'}
	ret := rb.Push(in)
	if ret != 1 {
		t.Fatal(ret)
	}
	testBasic(t, rb, 2, 6)
	if rb.buffer[idx] != 'a' {
		t.Fatal(rb.buffer)
	}
	if rb.in != 0 {
		t.Fatal(rb)
	}
}

func TestCrossPush(t *testing.T) {
	rb := New(4)
	rb.in = 2
	rb.out = 2
	testBasic(t, rb, 0, 4)

	in := []byte{'a', 'b', 'c', 'd', 'e'}
	ret := rb.Push(in)
	if ret != 4 {
		t.Fatal(ret)
	}
	testBasic(t, rb, 4, 0)
	if rb.buffer[0] != 'c' ||
		rb.buffer[1] != 'd' ||
		rb.buffer[2] != 'a' ||
		rb.buffer[3] != 'b' {
		t.Fatal(rb.buffer)
	}
	if rb.in != 6 {
		t.Fatal(rb.in)
	}
}

func TestProducerConsumer(t *testing.T) {
	rb := New(1024)
	data := []byte(string("huangjian"))
	N := 10000
	go func() {
		for i := 0; i < N; {
			if int(rb.LeftSpace()) > len(data) {
				rb.Push(data)
				i++
			}
			time.Sleep(time.Microsecond)
		}
	}()

	time.Sleep(time.Second)

	for i := 0; i < N; i++ {
		out := make([]byte, len(data))
		ret := rb.Get(out)
		if string(out) != "huangjian" {
			t.Fatal(i, ret, out)
		}
		time.Sleep(2 * time.Microsecond)
	}
}
