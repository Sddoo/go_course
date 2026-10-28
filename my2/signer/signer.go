package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
)

const maxInput = 2

// func StartPipes(in, out chan interface{}, workers ...job) {
// 	for i, f := range workers {
// 		if i&1 == 0 {
// 			f(in, out)
// 		} else {
// 			f(out, in)
// 		}
// 	}

// }

func ExecutePipeline(workers ...job) {
	in := make(chan interface{}, maxInput)
	out := make(chan interface{}, maxInput)

	// workers[0](in, out)
	// for input := range in {
	// 	StartPipes(in,)
	// }
	for i, f := range workers {
		if i&1 == 0 {
			f(in, out)
		} else {
			f(out, in)
		}
	}
}

func SingleHash(in, out chan interface{}) {
	// for range maxInput {
	for input := range in {
		// input := <-in
		fmt.Println(input)
		md5Data := DataSignerMd5(strconv.Itoa(input.(int)))
		crc32md5Data := DataSignerCrc32(md5Data)
		crc32Data := DataSignerCrc32(strconv.Itoa(input.(int)))
		result := crc32Data + "~" + crc32md5Data

		out <- result
	}
}

func MultiHash(in, out chan interface{}) {
	for input := range in {
		hashes := make([]string, 6)
		for th := range hashes {
			crc32Data := DataSignerCrc32(strconv.Itoa(th) + input.(string))
			hashes[th] = crc32Data
		}
		result := strings.Join(hashes, "")

		out <- result
	}
}

func CombineResults(in, out chan interface{}) {
	var result []string
	wg := &sync.WaitGroup{}

	for input := range in {
		wg.Add(1)
		fmt.Println(input)
		result = append(result, strconv.Itoa(input.(int)))
		wg.Done()
	}

	wg.Wait()
	slices.Sort(result)
	out <- strings.Join(result, "_")
}
