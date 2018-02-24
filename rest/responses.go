package rest

type DigitResponse struct {
	Index int64   `json:"index"`
	Digit byte    `json:"digit"`
	Error *string `json:"error"`
}

type ChunkResponse struct {
	FirstIndex int64   `json:"firstIndex"`
	Digits     []int   `json:"digits"`
	Error      *string `json:"error"`
}
