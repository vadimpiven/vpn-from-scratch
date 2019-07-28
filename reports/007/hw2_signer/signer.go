package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	var prev, curr chan interface{}
	wg := &sync.WaitGroup{}
	wg.Add(len(jobs))
	for _, j := range jobs {
		curr = make(chan interface{}, MaxInputDataLen)
		go func(in, out chan interface{}, j job) {
			defer wg.Done()
			j(in, out)
			close(out)
		}(prev, curr, j)
		prev = curr
	}
	wg.Wait()
}

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
	}
	wge.Wait()
}

var num = [6]string{"0", "1", "2", "3", "4", "5"}

func MultiHash(in, out chan interface{}) {
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
	}
	wge.Wait()
}

func CombineResults(in, out chan interface{}) {
	all := make([]string, 0, MaxInputDataLen)
	for val := range in {
		all = append(all, val.(string))
	}
	sort.Strings(all)
	out <- strings.Join(all, "_")
}
