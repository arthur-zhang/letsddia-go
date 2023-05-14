package merging_k_sorted_list

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"testing"
)

func createNSortedInDir(dir string) {
	N := 10
	for i := 0; i < N; i++ {
		fileName := fmt.Sprintf("%04d.dat", i)
		filePath := dir + "/" + fileName
		file, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		nums := make([]int, N, N)
		for i := 0; i < N; i++ {
			nums[i] = rand.Intn(100)
		}

		sort.Ints(nums)
		for i := 0; i < N; i++ {
			bytes := []byte(fmt.Sprintf("%d", nums[i]))
			_, _ = file.Write(bytes)
			if i != N-1 {
				_, _ = file.Write([]byte("\n"))
			}
		}
	}

}

// mkdirOrRemove creates a directory if it does not exist,
// or removes all files in the directory if it exists.
func mkdirOrRemove(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	} else {
		files, err := os.ReadDir(dir)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			err = os.Remove(dir + "/" + file.Name())
			if err != nil {
				panic(err)
			}
		}
	}
}

type Block struct {
	r     *bufio.Reader
	isEof bool
	line  string
}

func NewBlock(path string) Block {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	return Block{
		r:     r,
		isEof: false,
	}
}

type Iterator interface {
	Next() bool
	Value() int
}

func (b *Block) Next() bool {
	line, _, err := b.r.ReadLine()
	if err != nil {
		if err != io.EOF {
			panic("unexpected error")
		}
		return false
	}
	b.line = string(line)
	return true
}
func (b *Block) Value() int {
	n, err := strconv.Atoi(b.line)
	if err != nil {
		panic("should not happen")
	}
	return n

}
func listFilesInDir(dir string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		panic("open failed")
	}
	result := make([]string, 0)
	for _, file := range files {
		filePath := dir + "/" + file.Name()
		result = append(result, filePath)
	}
	return result
}

func TestMerging(t *testing.T) {
	dir := "/tmp/k_way_merging"
	mkdirOrRemove(dir)
	createNSortedInDir(dir)
	files := listFilesInDir(dir)
	iters := make([]Block, 0)
	for _, file := range files {
		iters = append(iters, NewBlock(file))
	}
	itersLen := len(iters)
	pq := &MinHeap{}

	heap.Init(pq)
	for i := 0; i < itersLen; i++ {
		if !iters[i].Next() {
			continue
		}
		entry := Entry{
			value:   iters[i].Value(),
			iterIdx: i,
		}
		heap.Push(pq, entry)
	}
	for pq.Len() > 0 {
		item := heap.Pop(pq).(Entry)
		println(item.value)

		if iters[item.iterIdx].Next() {
			newEntry := Entry{
				value:   iters[item.iterIdx].Value(),
				iterIdx: item.iterIdx,
			}
			heap.Push(pq, newEntry)
		}
	}
}

// implement priority queue using min-heap
type Entry struct {
	value   int
	iterIdx int
}
type MinHeap []Entry

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].value < h[j].value }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(Entry))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// end of min-heap
