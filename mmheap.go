// Package mmheap provides a drop-in min-max heap for any type that implements
// heap.Interface.
//
// In a min-max heap, the minimum element is at the root at index 0, and the
// maximum element is in one of the two root children.
//
// Thus, a min-max heap provides an efficient way to access either the minimum
// or maximum elements in a set in O(1) time.
//
// Due to the increased number of comparisons to implement a min-heap, its
// runtime is slower than a general heap (in this implementation, just under
// 2x). Because of this, unless you need to continually access both the minimum
// and maximum elements, it may be beneficial to use a different data
// structure. Even if you do need to continually access the minimum or maximum,
// a different data structure may be better.
//
// This package exists for anybody looking, but I recommend benchmarking this
// against github.com/google/btree. Generally, a btree is a good and fast
// data structure.
//
// For more information about a min-max heap, read the following paper:
//
//     http://www.cs.otago.ac.nz/staffpriv/mike/Papers/MinMaxHeaps/MinMaxHeaps.pdf
//
// Since this package is identical to the stdlib's container/heap
// documentation, it is elided here.
package mmheap

import (
	"container/heap"
	"math/bits"
)

type Interface interface {
	heap.Interface
}

func Init(h Interface) {
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		down(h, i, n)
	}
}

func Push(h Interface, x interface{}) {
	h.Push(x)
	up(h, h.Len()-1)
}

func Pop(h Interface) interface{} {
	n := h.Len() - 1
	h.Swap(0, n)
	down(h, 0, n)
	return h.Pop()
}

func Remove(h Interface, i int) interface{} {
	n := h.Len() - 1
	if n != i {
		h.Swap(i, n)
		if !down(h, i, n) {
			up(h, i)
		}
	}
	return h.Pop()
}

func Fix(h Interface, i int) {
	if !down(h, i, h.Len()) {
		up(h, i)
	}
}

// Max returns the index of the maximum element of the heap.
//
// This is a convenience function that always returns either 0, 1, or 2.
// This will panic if the heap has no elements.
func Max(h Interface) int {
	switch h.Len() {
	case 1:
		return 0
	case 2:
		return 1
	default:
		if h.Less(1, 2) {
			return 2
		}
		return 1
	}
}

func up(h Interface, on int) {
	onMinLevel := isMinLevel(on)

	// On min level:
	// If we have a parent, if our parent is less than us, then we swap
	// with the parent. Our parent should not be less than us.
	//
	// On max level:
	// If we have a parent, if our parent is more than us, then we swap
	// with the parent. Our parent should not be more than us.
	//
	// If we swap, our level changed by one, and we need to swap onMinLevel.
	parent := parent(on)
	if hasParent(on) {
		if onMinLevel == h.Less(parent, on) {
			h.Swap(on, parent)
			on = parent
			onMinLevel = !onMinLevel
		}
	}

	// On min level:
	// While we have a grandparent, if our grandparent is less than us,
	// then we swap with our grandparent.
	//
	// On max level:
	// While we have a grandparent, if our grandparent is more than us,
	// same.
	for hasGrandparent(on) {
		grandparent := grandparent(on)
		if onMinLevel == h.Less(on, grandparent) {
			h.Swap(on, grandparent)
			on = grandparent
			continue
		}
		break
	}
}

// min levels are odd levels, following a log pattern, so the odd expression
// below works out.
func isMinLevel(index int) bool {
	return bits.LeadingZeros(uint(index+1))&1 == 1
}

func hasParent(index int) bool {
	return index > 0
}

func parent(index int) int {
	return (index - 1) / 2
}

func hasGrandparent(index int) bool {
	return index > 2
}

func grandparent(index int) int {
	return parent(parent(index))
}

func down(h Interface, i0, n int) bool {
	on := i0
	onMinLevel := isMinLevel(i0)

	for {
		l := 2*on + 1
		r := l + 1

		ll := 2*l + 1
		lr := ll + 1

		rl := 2*r + 1
		rr := rl + 1

		type relation uint8
		const (
			self relation = iota
			child
			grandchild
		)

		type progeny struct {
			index    int
			relation relation
		}

		smallest := progeny{on, self}
		for _, progeny := range &[...]progeny{
			{l, child},
			{r, child},
			{ll, grandchild},
			{lr, grandchild},
			{rl, grandchild},
			{rr, grandchild},
		} {
			if progeny.index >= n || progeny.index < 0 {
				break
			}
			if onMinLevel == h.Less(progeny.index, smallest.index) {
				smallest = progeny
			}
		}

		if smallest.relation == self {
			break
		}
		h.Swap(on, smallest.index)
		on = smallest.index
		if smallest.relation == child {
			break
		}
		p := parent(smallest.index)
		if onMinLevel == h.Less(p, smallest.index) {
			h.Swap(smallest.index, p)
		}
	}
	return on > i0
}
