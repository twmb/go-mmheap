package mmheap

import (
	"math/rand"
	"testing"
)

// The tests below are ripped straight from stdlib's container/heap
// with minor modifications where necessary.
//
// All benchmarks but BenchmarkDup are new.

/*
Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

type intHeap []int

func (h *intHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *intHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *intHeap) Len() int {
	return len(*h)
}

func (h *intHeap) Pop() (v interface{}) {
	*h, v = (*h)[:h.Len()-1], (*h)[h.Len()-1]
	return
}

func (h *intHeap) Push(v interface{}) {
	*h = append(*h, v.(int))
}

func (h intHeap) verify(t *testing.T, i int) {
	t.Helper()
	n := h.Len()
	j1 := 2*i + 1
	j2 := 2*i + 2
	badDirection := ">"
	if !isMinLevel(i) {
		badDirection = "<"
	}
	if j1 < n {
		if isMinLevel(i) == h.Less(j1, i) && h[j1] != h[i] {
			t.Errorf("heap invariant invalidated [%d] = %d %s [%d] = %d",
				i, h[i],
				badDirection,
				j1, h[j1],
			)
			return
		}
		h.verify(t, j1)
	}
	if j2 < n {
		if isMinLevel(i) == h.Less(j2, i) && h[j2] != h[i] {
			t.Errorf("heap invariant invalidated [%d] = %d %s [%d] = %d",
				i, h[i],
				badDirection,
				j1, h[j2],
			)
			return
		}
		h.verify(t, j2)
	}
}

func TestInit0(t *testing.T) {
	h := new(intHeap)
	for i := 20; i > 0; i-- {
		h.Push(0) // all elements are the same
	}
	Init(h)
	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		x := Pop(h).(int)
		h.verify(t, 0)
		if x != 0 {
			t.Errorf("%d.th pop got %d; want %d", i, x, 0)
		}
	}
}

func TestInit1(t *testing.T) {
	h := new(intHeap)
	for i := 20; i > 0; i-- {
		h.Push(i) // all elements are different
	}
	Init(h)
	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		x := Pop(h).(int)
		h.verify(t, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func Test(t *testing.T) {
	h := new(intHeap)
	h.verify(t, 0)

	for i := 20; i > 10; i-- {
		h.Push(i)
	}
	Init(h)
	h.verify(t, 0)

	for i := 10; i > 0; i-- {
		Push(h, i)
		h.verify(t, 0)
	}

	for i := 1; h.Len() > 0; i++ {
		x := Pop(h).(int)
		if i < 20 {
			Push(h, 20+i)
		}
		h.verify(t, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

// The Go container/heap Remove tests relied on left to right heap ordering and
// did not initialize the order of the heap.
//
// We init the heap after pushing the first 10 elements as well as do a little
// bit more to determine the max (which was always removed in Remove0).

func TestRemove0(t *testing.T) {
	h := new(intHeap)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	Init(h)
	h.verify(t, 0)

	for h.Len() > 0 {
		i := MaxIndex(h)
		exp := h.Len() - 1
		x := Remove(h, i).(int)
		if x != exp {
			t.Errorf("Remove(%d) got %d; want %d", i, x, exp)
		}
		h.verify(t, 0)
	}
}

func TestRemove1(t *testing.T) {
	h := new(intHeap)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	Init(h)
	h.verify(t, 0)

	for i := 0; h.Len() > 0; i++ {
		x := Remove(h, 0).(int)
		if x != i {
			t.Errorf("Remove(0) got %d; want %d", x, i)
		}
		h.verify(t, 0)
	}
}

func TestRemove2(t *testing.T) {
	N := 10

	h := new(intHeap)
	for i := 0; i < N; i++ {
		h.Push(i)
	}
	Init(h)
	h.verify(t, 0)

	m := make(map[int]bool)
	for h.Len() > 0 {
		m[Remove(h, (h.Len()-1)/2).(int)] = true
		h.verify(t, 0)
	}

	if len(m) != N {
		t.Errorf("len(m) = %d; want %d", len(m), N)
	}
	for i := 0; i < len(m); i++ {
		if !m[i] {
			t.Errorf("m[%d] doesn't exist", i)
		}
	}
}

func TestFix(t *testing.T) {
	h := new(intHeap)
	h.verify(t, 0)

	for i := 200; i > 0; i -= 10 {
		Push(h, i)
	}
	h.verify(t, 0)

	if (*h)[0] != 10 {
		t.Fatalf("Expected head to be 10, was %d", (*h)[0])
	}
	(*h)[0] = 210
	Fix(h, 0)
	h.verify(t, 0)

	for i := 100; i > 0; i-- {
		elem := rand.Intn(h.Len())
		if i&1 == 0 {
			(*h)[elem] *= 2
		} else {
			(*h)[elem] /= 2
		}
		Fix(h, elem)
		h.verify(t, 0)
	}
}

func BenchmarkDup(b *testing.B) {
	const n = 10000
	h := make(intHeap, 0, n)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			Push(&h, 0) // all elements are the same
		}
		for h.Len() > 0 {
			Pop(&h)
		}
	}
}

func BenchmarkOrdered(b *testing.B) {
	const n = 1000
	h := make(intHeap, 0, n)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			Push(&h, j)
		}
		for h.Len() > 0 {
			Pop(&h)
		}
	}
}

func BenchmarkRandom(b *testing.B) {
	rng := rand.New(rand.NewSource(0))
	const n = 1000
	h := make(intHeap, 0, n)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			if rng.Intn(10) == 1 && h.Len() > 0 {
				h.Pop()
			} else {
				Push(&h, rng.Intn(n))
			}
		}
		for h.Len() > 0 {
			Pop(&h)
		}
	}
}

func BenchmarkOrderedPushOnly(b *testing.B) {
	const n = 1000
	for i := 0; i < b.N; i++ {
		h := make(intHeap, 0, n)
		for j := 0; j < n; j++ {
			Push(&h, j)
		}
	}
}
