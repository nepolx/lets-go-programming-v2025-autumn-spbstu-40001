package main

import (
	"container/heap"
	"fmt"
)

func main() {
	var dishNum int

	_, err := fmt.Scan(&dishNum)
	if err != nil {
		fmt.Println("Invalid input", err)

		return
	}

	intheap := &IntHeap{}
	heap.Init(intheap)

	for range dishNum {
		var temp int

		_, err := fmt.Scan(&temp)
		if err != nil {
			fmt.Println("Invalid input", err)

			return
		}

		heap.Push(intheap, temp)
	}

	var rating int

	_, err = fmt.Scan(&rating)
	if err != nil {
		fmt.Println("Invalid input", err)

		return
	}

	for range rating - 1 {
		if intheap.Len() == 0 {
			fmt.Println("There is no such dish")
		}

		heap.Pop(intheap)
	}

	fmt.Println(heap.Pop(intheap))
}

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}
