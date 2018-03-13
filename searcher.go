package piio

import (
	"bufio"
	"errors"
	"io"
)

const searchChunkSize = 64

func Search(compressedFile string, search string) (int64, error) {
	for _, c := range search {
		if c < '0' || c > '9' {
			return -1, errors.New("you can only search for series of digits")
		}
	}

	file, err := NewCompressedToStringReader(compressedFile)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, searchChunkSize)
	searchIndex := 0
	globalIndex := int64(0)
	matchIndex := int64(-1)
	for {
		_, err = reader.Read(buffer)
		if err != nil {
			break
		}
		for i := range buffer {
			if buffer[i] == search[searchIndex] {
				searchIndex++
			} else {
				searchIndex = 0
			}
			globalIndex++
			if searchIndex == len(search) {
				matchIndex = globalIndex - int64(len(search)) + 1
				break
			}
		}
		if matchIndex != -1 {
			break
		}
	}
	if err != nil && err != io.EOF {
		return -1, err
	}
	return matchIndex, nil
}
