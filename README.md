go-mmheap
=========

Package mmheap provides a drop-in min-max heap for any type that implements
heap.Interface.

In a min-max heap, the minimum element is at the root at index 0, and the
maximum element is in one of the two root children.

Thus, a min-max heap provides an efficient way to access either the minimum
or maximum elements in a set in O(1) time.

Due to the increased number of comparisons to implement a min-heap, its
runtime is slower than a general heap (in this implementation, just under
2x). Because of this, unless you need to continually access both the minimum
and maximum elements, it may be beneficial to use a different data
structure. Even if you do need to continually access the minimum or maximum,
a different data structure may be better.

This package exists for anybody looking, but I recommend benchmarking this
against github.com/google/btree. Generally, a btree is a good and fast
data structure.

For more information about a min-max heap, read the following paper:

    http://www.cs.otago.ac.nz/staffpriv/mike/Papers/MinMaxHeaps/MinMaxHeaps.pdf

Since this package is identical to the stdlib's container/heap
documentation, it is elided here.

All of the tests are taken from the stdlib's `container/heap package and the Go
BSD license is copied into `mmheap_test.go` for completeness. As well, the main
package-level functions are duplicated, but only because there is about one way
to write Push/Pop on a heap. The sift up and sift down functions are new.

Documentation
-------------

[![GoDoc](https://godoc.org/github.com/twmb/go-mmheap?status.svg)](https://godoc.org/github.com/twmb/go-mmheap)
