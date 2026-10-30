package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
)

const maxInput = 100

// 1
// func ExecutePipeline(workers ...job) {
// 	in := make(chan interface{}, 1)
// 	out := make(chan interface{}, 1)

// 	for _, f := range workers {
// 		// if i&1 == 0 {
// 		// 	go f(in, out)
// 		// } else {
// 		// 	go f(out, in)
// 		// }
// 		go f(in, out)
// 		in, out = out, in
// 	}
// }

// 2
// попытка обеспечить последовательную передачу каналов между воркерами
// func PipeSwap(src, dst chan interface{}) {
// 	mu := &sync.Mutex{}
// 	mu.Lock()
// 	for srcVal := range src {
// 		dst <- srcVal
// 	}
// 	mu.Unlock()
// }

// func StartPipe(worker job, bindingCh chan interface{}, wg *sync.WaitGroup) {
// 	fmt.Println("Pipe has started") //
// 	quotaCh <- struct{}{}
// 	localIn := make(chan interface{}, maxInput)
// 	localOut := make(chan interface{}, maxInput)

// 	PipeSwap(bindingCh, localIn)
// 	close(localIn)
// 	fmt.Println("PipeSwap occured") //

// 	worker(localIn, localOut)

// 	PipeSwap(localOut, bindingCh)
// 	close(localOut)

// 	wg.Done()
// 	<-quotaCh
// }

// 3
func PipeWrapper(in chan interface{}, worker job, wg *sync.WaitGroup) chan interface{} {
	out := make(chan interface{}, maxInput)

	go func() {
		worker(in, out)
		close(out)
		wg.Done()
	}()

	return out
}

func ExecutePipeline(workers ...job) {
	out := make(chan interface{}, maxInput)
	wg := &sync.WaitGroup{}

	for _, f := range workers {
		wg.Add(1)
		out = PipeWrapper(out, f, wg)
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	fmt.Println("SingleHash started") //

	for input := range in {
		md5Data := DataSignerMd5(strconv.Itoa(input.(int)))
		crc32md5Data := DataSignerCrc32(md5Data)
		crc32Data := DataSignerCrc32(strconv.Itoa(input.(int)))
		result := crc32Data + "~" + crc32md5Data

		out <- result
		// runtime.Gosched()
	}
}

func MultiHash(in, out chan interface{}) {
	fmt.Println("MultiHash started") //

	for input := range in {
		hashes := make([]string, 6)
		for th := range hashes {
			crc32Data := DataSignerCrc32(strconv.Itoa(th) + input.(string))
			hashes[th] = crc32Data
		}
		result := strings.Join(hashes, "")

		out <- result
		// runtime.Gosched()
	}
}

func CombineResults(in, out chan interface{}) {
	fmt.Println("CombineResults started") //
	var result []string

	for input := range in {
		fmt.Println("CombineResults: ", input) //
		result = append(result, input.(string))
	}

	slices.Sort(result)
	out <- strings.Join(result, "_")
}
