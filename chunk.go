package piio

import (
	"errors"
	"io"
)

// Chunk represents a chunk of digits of pi.
type Chunk interface {
	// FirstIndex returns the index of the first digit of
	// pi contained in this chunk.
	FirstIndex() int64

	// Length returns the amount of digits contained in
	// this chunk.
	Length() int

	// LastIndex returns the index of the last digit of pi
	// contained in this chunk.
	LastIndex() int64

	// Digit returns the index-th digit of pi. It errors if
	// the requested digit is not contained in this chunk.
	Digit(index int64) (byte, error)

	// IsCompressed returns whether or not the chunk is compressed.
	IsCompressed() bool
}

// CompressedChunk represents a compressed chunk
// of digits of pi.
type CompressedChunk struct {
	firstIndex int64
	data       []byte
}

// ReadCompressedChunk reads a certain chunk defined by
// the index of the first requested digit of pi and the
// amount of digits requested.
// Both the first index and the size have to be positive and
// even. The given file has to be seekable.
//
// The expected file format is as follows.
// Binary digits, each digit 4 bits wide with the lower index
// digit in the higher nibble of each byte.
func ReadCompressedChunk(input io.ReadSeeker, firstIndex int64, size int) (Chunk, error) {
	if firstIndex < 0 || firstIndex%2 != 0 {
		return nil, errors.New("only positive even first indexes are supported")
	}
	if size <= 0 || size%2 != 0 {
		return nil, errors.New("only positive even sizes are supported")
	}
	_, err := input.Seek(firstIndex/2, 0)
	if err != nil {
		return nil, err
	}

	c, err := ReadNextCompressedChunk(input, size)
	chunk, _ := c.(*CompressedChunk)
	chunk.firstIndex = firstIndex

	return chunk, nil
}

func ReadNextCompressedChunk(input io.Reader, size int) (Chunk, error) {
	var err error
	chunk := &CompressedChunk{
		firstIndex: -1,
		data:       make([]byte, size/2),
	}

	size, err = input.Read(chunk.data)
	if err != nil {
		return nil, err
	}

	// Trim in case we requested more than the file can give us.
	chunk.data = chunk.data[:size]

	return chunk, err
}

// Compress compresses an uncompressed chunk.
func Compress(chnk Chunk) Chunk {
	if chnk.IsCompressed() {
		return chnk
	}

	c, ok := chnk.(*UncompressedChunk)
	if !ok {
		panic("can only uncompress CompressedChunks")
	}

	chunk := &CompressedChunk{
		firstIndex: c.FirstIndex(),
		data:       make([]byte, len(c.Digits)/2),
	}
	for i := 0; i < len(chunk.data); i++ {
		chunk.data[i] = c.Digits[i*2] << 4
		chunk.data[i] |= c.Digits[i*2+1]
	}
	return chunk
}

// IsCompressed returns true.
func (c *CompressedChunk) IsCompressed() bool {
	return true
}

func (c *CompressedChunk) digitIndexToDataIndex(index int64) (ind int, isHighNibble bool) {
	ind = int((index - c.firstIndex) / 2)
	isHighNibble = (index-c.firstIndex)%2 == 0
	return
}

// FirstIndex returns the index of the first digit of
// pi contained in this chunk.
func (c *CompressedChunk) FirstIndex() int64 {
	return c.firstIndex
}

// Length returns the amount of digits contained in
// this chunk.
func (c *CompressedChunk) Length() int {
	return len(c.data) * 2
}

// LastIndex returns the index of the last digit of pi
// contained in this chunk.
func (c *CompressedChunk) LastIndex() int64 {
	return c.firstIndex + int64(c.Length()) - 1
}

// Digit returns the index-th digit of pi. It errors if
// the requested digit is not contained in this chunk.
func (c *CompressedChunk) Digit(index int64) (byte, error) {
	if index < c.firstIndex || int((index-c.firstIndex)/2) >= len(c.data) {
		return 255, errors.New("index out of range")
	}
	ind, isHighNibble := c.digitIndexToDataIndex(index)

	digit := c.data[ind]
	if isHighNibble {
		digit = digit >> 4
	}
	digit = (digit & 0x0F)

	return digit, nil
}

// UncompressedChunk represents a non-compressed chunk
// of digits of pi.
type UncompressedChunk struct {
	FirstDigitIndex int64  `json:"firstDigitIndex"`
	Digits          []byte `json:"digits"`
}

// ReadTextChunk reads a certain chunk defined by
// the index of the first requested digit of pi and the
// amount of digits requested from a text based input.
// Both the first index and the size have to be positive and
// even. The given file has to be seekable.
//
// The expected file format is one character for each digit
// no decimal point. So the file should start with `314`...
func ReadTextChunk(input io.ReadSeeker, firstIndex int64, size int) (Chunk, error) {
	if firstIndex < 0 || firstIndex%2 != 0 {
		return nil, errors.New("only positive even first indexes are supported")
	}
	if size <= 0 || size%2 != 0 {
		return nil, errors.New("only positive even sizes are supported")
	}
	_, err := input.Seek(firstIndex, 0)
	if err != nil {
		return nil, err
	}

	chunk := &UncompressedChunk{
		FirstDigitIndex: firstIndex,
		Digits:          make([]byte, size),
	}

	size, err = input.Read(chunk.Digits)
	if err != nil {
		return nil, err
	}

	// Trim in case we requested more than the file can give us.
	chunk.Digits = chunk.Digits[:size]

	digits := make([]byte, 0, len(chunk.Digits))
	// Convert from character to number and filter
	for i := range chunk.Digits {
		if chunk.Digits[i] < '0' || '9' < chunk.Digits[i] {
			continue
		}
		digits = append(digits, chunk.Digits[i]-byte('0'))
	}
	chunk.Digits = digits

	return chunk, nil
}

// Decompress decompresses a compressed chunk.
func Decompress(chnk Chunk) Chunk {
	if !chnk.IsCompressed() {
		return chnk
	}

	c, ok := chnk.(*CompressedChunk)
	if !ok {
		panic("can only uncompress CompressedChunks")
	}

	chunk := &UncompressedChunk{
		FirstDigitIndex: c.firstIndex,
		Digits:          make([]byte, len(c.data)*2),
	}
	for i := 0; i < len(c.data); i++ {
		chunk.Digits[i*2] = (c.data[i] >> 4) & 0x0F
		chunk.Digits[i*2+1] = (c.data[i] & 0x0F)
	}
	return chunk
}

// IsCompressed returns false.
func (c *UncompressedChunk) IsCompressed() bool {
	return false
}

// FirstIndex returns the index of the first digit of
// pi contained in this chunk.
func (c *UncompressedChunk) FirstIndex() int64 {
	return c.FirstDigitIndex
}

// Length returns the amount of digits contained in
// this chunk.
func (c *UncompressedChunk) Length() int {
	return len(c.Digits)
}

// LastIndex returns the index of the last digit of pi
// contained in this chunk.
func (c *UncompressedChunk) LastIndex() int64 {
	return c.FirstDigitIndex + int64(c.Length()) - 1
}

// Digit returns the index-th digit of pi. It errors if
// the requested digit is not contained in this chunk.
func (c *UncompressedChunk) Digit(index int64) (byte, error) {
	ind := int(index - c.FirstDigitIndex)
	if index < c.FirstDigitIndex || ind >= len(c.Digits) {
		return 255, errors.New("index out of range")
	}
	return c.Digits[ind], nil
}

// AsUncompressedChunk returns the given Chunk as an
// UncompressedChunk. Panics if wrong type is given.
func AsUncompressedChunk(chnk Chunk) *UncompressedChunk {
	chnk = Decompress(chnk)
	c, ok := chnk.(*UncompressedChunk)
	if !ok {
		panic("only builtin chunk types are supported")
	}
	return c
}

func WriteChunk(chnk Chunk, format FileFormat, w io.Writer) error {
	if format == FileFormatCompressed {
		chnk = Compress(chnk)
		c, ok := chnk.(*CompressedChunk)
		if !ok {
			return errors.New("unknown Chunk type")
		}
		return writeCompressedChunk(c, w)
	} else if format == FileFormatText {
		chnk = Decompress(chnk)
		c, ok := chnk.(*UncompressedChunk)
		if !ok {
			return errors.New("unknown Chunk type")
		}
		return writeUncompressedChunkText(c, w)
	}
	return errors.New("unknown file format")
}

func writeCompressedChunk(chnk *CompressedChunk, w io.Writer) error {
	n, err := w.Write(chnk.data)
	if err != nil {
		return err
	}
	if n != len(chnk.data) {
		return errors.New("not all bytes could be written")
	}
	return nil
}

func writeUncompressedChunkText(chnk *UncompressedChunk, w io.Writer) error {
	data := append([]byte{}, chnk.Digits...)
	for i := range data {
		data[i] = data[i] + byte('0')
	}
	n, err := w.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return errors.New("not all bytes could be written")
	}
	return nil
}
