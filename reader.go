package piio

import (
	"io"
	"os"
)

type compressedToStringReader struct {
	file         io.ReadCloser
	currentIndex int64
}

func NewCompressedToStringReader(filename string) (io.ReadCloser, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return &compressedToStringReader{
		file:         f,
		currentIndex: 0,
	}, nil
}

func (r *compressedToStringReader) Read(p []byte) (n int, err error) {
	chunk, err := ReadNextCompressedChunk(r.file, len(p))
	if err != nil {
		return
	}
	for i := 0; i < chunk.Length(); i++ {
		d, _ := chunk.Digit(int64(i) + r.currentIndex)
		p[i] = byte('0') + d
	}
	n = chunk.Length()
	r.currentIndex += int64(n)

	return
}

func (r *compressedToStringReader) Close() error {
	return r.file.Close()
}
