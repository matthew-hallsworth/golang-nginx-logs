package main

import (
	"container/heap"
	"reflect"
	"strings"
	"testing"
)

// TestScanLogFile This tests scanning of the log file and it's main method
func TestScanLogFile(t *testing.T) {
	t.Run("Test that the log file scanning works with a single line", func(t *testing.T) {
		sample := strings.NewReader("79.125.00.21 - - [10/Jul/2018:20:03:40 +0200] \"GET /newsletter/ HTTP/1.1\" 200 3574 \"-\" \"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/5.0)\"")
		ipList, urlList := scanLogFile(sample)
		if len(ipList) != 1 {
			t.Errorf("Number of IP in single line is not 1 - is actually %d", len(ipList))
		}
		keys := reflect.ValueOf(ipList).MapKeys()
		if keys[0].Interface().(string) != "79.125.00.21" {
			t.Errorf("IP fetched is incorrect - is actually %s", keys[0])
		}
		if len(urlList) != 1 {
			t.Errorf("Number of URL in single line is not 1 - is actually %d", len(urlList))
		}
		keys = reflect.ValueOf(urlList).MapKeys()
		if keys[0].Interface().(string) != "/newsletter/" {
			t.Errorf("URL fetched is incorrect - is actually %s", keys[0])
		}
	})

	t.Run("Test parsing a multiline string", func(t *testing.T) {
		sample := strings.NewReader(`79.125.00.22 - - [10/Jul/2018:20:03:40 +0200] "GET /newsletter/abc HTTP/1.1" 200 3574 "-" "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/5.0)"
79.125.00.21 - - [10/Jul/2018:20:05:40 +0200] "GET /newsletter/ HTTP/1.1" 200 3574 "-" "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/5.0)"		
79.125.00.22 - - [10/Jul/2018:20:05:40 +0200] "GET /newsletter/abc HTTP/1.1" 200 3574 "-" "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/5.0)"`)
		ipList, urlList := scanLogFile(sample)
		if len(ipList) != 2 {
			t.Errorf("Number of IP in single line is not 2 - is actually %d", len(ipList))
		}
		keys := reflect.ValueOf(ipList).MapKeys()
		if keys[0].Interface().(string) != "79.125.00.22" {
			t.Errorf("IP fetched is incorrect - is actually %s", keys[0])
		}
		if len(urlList) != 2 {
			t.Errorf("Number of URL in single line is not 2 - is actually %d", len(urlList))
		}
		keys = reflect.ValueOf(urlList).MapKeys()
		if keys[0].Interface().(string) != "/newsletter/abc" {
			t.Errorf("URL fetched is incorrect - is actually %s", keys[0])
		}
	})
}

func TestOrderingMap(t *testing.T) {
	t.Run("Test that ordering a dictionary map occurs correctly", func(t *testing.T) {
		sampleMap := map[string]int{
			"a": 1,
			"b": 4,
			"c": 6,
			"d": 12,
		}
		sampleHeap := getOrderedHeap(sampleMap)
		topValue := heap.Pop(sampleHeap).(kv).Value
		if topValue != 12 {
			t.Errorf("Did not order map correctly to get 12, value returned is %d", topValue)
		}
	})
}
