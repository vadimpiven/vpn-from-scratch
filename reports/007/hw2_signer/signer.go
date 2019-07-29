package main

import (
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// ExecutePipeline creates the pipe chain between given job instances.
// It is assumed that first instance do not use `in chan` and the last instance do not use `out chan`.
// Pipes are buffered, they can hold up to `MaxInputDataLen` values as it's the largest expected input length.
func ExecutePipeline(jobs ...job) {
	var prev, curr chan interface{}
	wg := &sync.WaitGroup{}
	wg.Add(len(jobs))
	for _, j := range jobs {
		curr = make(chan interface{}, MaxInputDataLen)
		go func(in, out chan interface{}, j job) {
			defer wg.Done()
			defer close(out)
			j(in, out)
		}(prev, curr, j)
		prev = curr
	}
	wg.Wait()
}

// SingleHash calculates `crc32(data) + "~" + crc32(md5(data))`.
// It is assumed that data has type `int`.
func SingleHash(in, out chan interface{}) {
	wge := &sync.WaitGroup{}
	q := make(chan struct{}, 1)
	for val := range in {
		wge.Add(1)
		go func(val string) {
			defer wge.Done()
			var crc1, crc2 string
			wgi := &sync.WaitGroup{}
			wgi.Add(2)
			go func(val string) {
				defer wgi.Done()
				crc1 = DataSignerCrc32(val)
			}(val)
			go func(val string) {
				defer wgi.Done()
				q <- struct{}{}
				val = DataSignerMd5(val)
				<-q
				crc2 = DataSignerCrc32(val)
			}(val)
			wgi.Wait()
			out <- crc1 + "~" + crc2
		}(strconv.Itoa(val.(int)))
		runtime.Gosched()
	}
	wge.Wait()
}

// MultiHash calculates `crc32("0"+data) + ... + crc32("5"+data)`.
// It is assumed that data has type `string`.
func MultiHash(in, out chan interface{}) {
	var num = [6]string{"0", "1", "2", "3", "4", "5"}
	wge := &sync.WaitGroup{}
	for val := range in {
		wge.Add(1)
		go func(val string) {
			defer wge.Done()
			var res [6]string
			wgi := &sync.WaitGroup{}
			wgi.Add(6)
			for i, th := range num {
				go func(i int, val string) {
					defer wgi.Done()
					res[i] = DataSignerCrc32(val)
				}(i, th+val)
			}
			wgi.Wait()
			out <- res[0] + res[1] + res[2] + res[3] + res[4] + res[5]
		}(val.(string))
		runtime.Gosched()
	}
	wge.Wait()
}

// CombineResults collects all input data, sorts it and joins with `_` separator.
// It is assumed that data has type `string`.
// To collect all input data `MaxInputDataLen` sized buffer is used as it's the largest expected input length.
func CombineResults(in, out chan interface{}) {
	all := make([]string, 0, MaxInputDataLen)
	for val := range in {
		all = append(all, val.(string))
		runtime.Gosched()
	}
	sort.Strings(all)
	out <- strings.Join(all, "_")
}
