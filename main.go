package main

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

type kv struct {
	Key   string
	Value int
}

// KVHeap https://golang.org/pkg/container/heap/
type KVHeap []kv

func (h KVHeap) Less(i, j int) bool {
	return h[i].Value > h[j].Value
}

func (h KVHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h KVHeap) Len() int {
	return len(h)
}

// Push on the heap
func (h *KVHeap) Push(x interface{}) {
	*h = append(*h, x.(kv))
}

// Pop off the heap
func (h *KVHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func getOrderedHeap(m map[string]int) *KVHeap {
	h := &KVHeap{}
	heap.Init(h)
	for k, v := range m {
		heap.Push(h, kv{k, v})
	}
	return h
}

func processLine(currentLine string) (ip string, uri string) {
	logFormatRegexp, err := regexp.Compile("^([0-9/.]+) [^ ]+ [^ ]+ (.*) \"GET (http://example.net)?([^ ]+) .*\"")
	if err != nil {
		log.Fatal(err)
	}
	matches := logFormatRegexp.FindStringSubmatch(currentLine)
	if len(matches) > 0 {
		return matches[1], matches[4]
	}
	return
}

func incrementMapEntry(m map[string]int, entry string) {
	mapEntry := m[entry]
	if mapEntry == 0 {
		m[entry] = 1
	} else {
		m[entry] = mapEntry + 1
	}
}

func scanLogFile(handle io.Reader) (map[string]int, map[string]int) {
	ipList := make(map[string]int)
	urlList := make(map[string]int)
	scanner := bufio.NewScanner(handle)

	for scanner.Scan() {
		ip, url := processLine(scanner.Text())
		incrementMapEntry(ipList, ip)
		incrementMapEntry(urlList, url)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ipList, urlList
}

func main() {
	fptr := flag.String("f", "programming-task-example-data.log", "File to read logs from")
	flag.Parse()
	fileName := *fptr
	handle, err := os.Open(string(fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	ipList, urlList := scanLogFile(handle)

	ipListHeap := getOrderedHeap(ipList)
	urlListHeap := getOrderedHeap(urlList)

	// Output results
	fmt.Println("Number of IP's:")
	fmt.Println(len(ipList))
	fmt.Println("Top 3 IP's:")
	for i := 1; i <= 3; i++ {
		fmt.Println(heap.Pop(ipListHeap))
	}
	fmt.Println("Top 3 URL's")
	for i := 1; i <= 3; i++ {
		fmt.Println(heap.Pop(urlListHeap))
	}
}
