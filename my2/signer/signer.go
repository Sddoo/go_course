package main

import (
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

func SingleHash(in, out chan interface{}) {
	outWg := &sync.WaitGroup{}

	for input := range in {
		outWg.Add(1)
		md5Data := DataSignerMd5(strconv.Itoa(input.(int)))

		go func(md5Data string, outWg *sync.WaitGroup) {
			wg := &sync.WaitGroup{}
			var crc32md5Data, crc32Data string
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

			out <- result
			outWg.Done()
		}(md5Data, outWg)
	}

	outWg.Wait()
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
