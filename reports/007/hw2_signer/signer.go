package main

import (
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

}

func MultiHash(in, out chan interface{}) {

}

func CombineResults(in, out chan interface{}) {

}
