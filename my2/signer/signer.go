package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
)

const maxInput = 100
const maxMultiHashThreads = 6

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
	}
}

func MultiHash(in, out chan interface{}) {
	fmt.Println("MultiHash started") //

	for input := range in {
		// возможно стоит вынести эти 2 инициализации выше
		hashesCh := make(chan string, maxMultiHashThreads)
		hashes := make([]string, 0, maxMultiHashThreads)

		for th := range maxMultiHashThreads {
			go func(th int) {
				crcData := DataSignerCrc32(strconv.Itoa(th) + input.(string))
				// возможно хуевая идея
				for len(hashes) != th {
				}
				hashesCh <- crcData
			}(th)
		}

		// возможно стоит избавить от этого цикла путем избавления от hashesCh
		for hash := range hashesCh {
			hashes = append(hashes, hash)
			if len(hashes) == maxMultiHashThreads {
				close(hashesCh)
			}
		}

		fmt.Println("Multihash hashes", hashes)
		result := strings.Join(hashes, "")

		out <- result
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
