package main

import (
	"container/heap"
	"fmt"
)

func main() {
	var (
		dishNum, k int
	)

	_, err := fmt.Scan(&dishNum)
	if err != nil {
		fmt.Println("Invalid input", err)

		return
	}

	h := &IntHeap{}
	heap.Init(h)

	for range dishNum {
		var temp int
		_, err := fmt.Scan(&temp)
		if err != nil {
			fmt.Println("Invalid input", err)

			return
		}
		heap.Push(h, temp)
	}

	_, err = fmt.Scan(&k)
	if err != nil {
		fmt.Println("Invalid input", err)

		return
	}

	for range k - 1 {
		if h.Len() == 0 {
			fmt.Println("There is no such dish")
		}
		heap.Pop(h)
	}

	fmt.Println(heap.Pop(h))
}

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}
