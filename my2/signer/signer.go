package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func ExecutePipeline(workers ...job) {
	in := make(chan interface{}, 100)
	out := make(chan interface{}, 100)

	for i, f := range workers {
		if i&1 == 0 {
			f(in, out)
		} else {
			f(out, in)
		}
	}
}

func SingleHash(in, out chan interface{}) {
	for input := range in {
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

	for input := range in {
		fmt.Println(input)
		result = append(result, input.(string))
	}

	slices.Sort(result)
	out <- strings.Join(result, "_")
}
