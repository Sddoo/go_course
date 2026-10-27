package main

import (
	"strconv"
	"strings"
)

// func ExecutePipeline() {

// }

func SingleHash(data string) string {

	md5Data := DataSignerMd5(data)

	crc32md5Data := DataSignerCrc32(md5Data)

	crc32Data := DataSignerCrc32(data)

	result := crc32Data + "~" + crc32md5Data

	return result
}

func MultiHash(singleHashData string) string {
	hashes := make([]string, 6, 6)

	for th := range hashes {
		crc32Data := DataSignerCrc32(strconv.Itoa(th) + singleHashData)
		hashes[th] = crc32Data
	}

	result := strings.Join(hashes, "")

	return result
}

func CombineResults(hashes []string) string {
	return strings.Join(hashes, "_")
}
