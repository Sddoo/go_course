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

	for _, worker := range workers {
		wg.Add(1)
		out = PipeWrapper(out, worker, wg)
	}
	wg.Wait()
}

// надо распараллелить это
func SingleHash(in, out chan interface{}) {
	var crc32md5Data, crc32Data string
	wg := &sync.WaitGroup{}

	for input := range in {
		fmt.Println("SingleHash goes in", input)
		md5Data := DataSignerMd5(strconv.Itoa(input.(int)))
		wg.Add(2)

		go func() {
			crc32md5Data = DataSignerCrc32(md5Data)
			wg.Done()
		}()
		go func() {
			crc32Data = DataSignerCrc32(strconv.Itoa(input.(int)))
			wg.Done()
		}()

		wg.Wait()
		result := crc32Data + "~" + crc32md5Data

		fmt.Println("SingleHash goes out", result)
		out <- result
	}
}

func ComputeMultiHash(input interface{}, out chan interface{}, outWg *sync.WaitGroup) {
	wg := &sync.WaitGroup{}
	mu := sync.Mutex{}
	hashes := make([]string, maxMultiHashThreads)

	for th := range maxMultiHashThreads {
		wg.Add(1)
		go func(th int) {
			crcData := DataSignerCrc32(strconv.Itoa(th) + input.(string))
			mu.Lock()
			hashes[th] = crcData
			mu.Unlock()
			wg.Done()
		}(th)
	}

	wg.Wait()
	result := strings.Join(hashes, "")

	// fmt.Println("ComputeMultiHash goes out", result)
	out <- result
	outWg.Done()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for input := range in {
		wg.Add(1)
		go ComputeMultiHash(input, out, wg)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	result := make([]string, 0, maxInput)

	for input := range in {
		result = append(result, input.(string))
	}

	slices.Sort(result)
	out <- strings.Join(result, "_")
}
