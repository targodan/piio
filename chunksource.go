package piio

import (
	"errors"
	"fmt"
	"os"
)

// FileFormat represents the supported file formats.
type FileFormat int

const (
	// FileFormatCompressed represents a compressed binary format. See ReadCompressedChunkFile.
	FileFormatCompressed = iota
	// FileFormatText represents a text format. See ReadChunkFromTextfile.
	FileFormatText
)

// ChunkSource represents a source of chunks.
// It has to be thread safe.
type ChunkSource interface {
	// GetChunk returns the requested chunk.
	GetChunk(firstIndex int64, size int) (Chunk, error)
}

type uncachedChunkSource struct {
	filename   string
	fileFormat FileFormat
	maxSize    int
}

// NewUncachedChunkSource creates a new uncached ChunkSource.
func NewUncachedChunkSource(filename string, fileFormat FileFormat, maxSize int) ChunkSource {
	return &uncachedChunkSource{
		filename:   filename,
		fileFormat: fileFormat,
		maxSize:    maxSize,
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

	switch cs.fileFormat {
	case FileFormatCompressed:
		return ReadCompressedChunkFile(file, firstIndex, size)

	case FileFormatText:
		return ReadChunkFromTextfile(file, firstIndex, size)
	}

	return nil, errors.New("unknown file format")
}
