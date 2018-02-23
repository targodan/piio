package piio

import (
	"fmt"
	"os"
)

// ChunkSource represents a source of chunks.
// It has to be thread safe.
type ChunkSource interface {
	// GetChunk returns the requested chunk.
	GetChunk(firstIndex int64, size int) (Chunk, error)
}

type uncachedChunkSource struct {
	filename string
	maxSize  int
}

// NewUncachedChunkSource creates a new uncached ChunkSource.
func NewUncachedChunkSource(filename string, maxSize int) ChunkSource {
	return &uncachedChunkSource{
		filename: filename,
		maxSize:  maxSize,
	}
}

func (cs *uncachedChunkSource) GetChunk(firstIndex int64, size int) (Chunk, error) {
	if size > cs.maxSize {
		return nil, fmt.Errorf("requested chunk of size %d but only supporting chunks of size up to %d", size, cs.maxSize)
	}

	file, err := os.Open(cs.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ReadCompressedChunkFile(file, firstIndex, size)
}
