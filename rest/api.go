package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/targodan/piio"

	"github.com/julienschmidt/httprouter"
	errors "github.com/targodan/go-errors"
)

const BaseURI = "/api/"

type API struct {
	router      *httprouter.Router
	chunkSource piio.ChunkSource
}

func writeJson(w http.ResponseWriter, data interface{}) {
	jw := json.NewEncoder(w)
	jw.Encode(data)
}

func NewAPI(chunkSource piio.ChunkSource) *API {
	router := httprouter.New()
	api := &API{
		router:      router,
		chunkSource: chunkSource,
	}

	router.GET(BaseURI+"v1/digit/:index", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		index, err := strconv.ParseInt(p.ByName("index"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errMsg := "the index must be a number, got " + p.ByName("index")
			writeJson(w, &DigitResponse{Error: &errMsg})
			return
		}
		d, err := api.GetDigit(index)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errMsg := err.Error()
			writeJson(w, &DigitResponse{Error: &errMsg})
			return
		}
		writeJson(w, &DigitResponse{
			Index: index,
			Digit: d,
		})
	})

	router.GET(BaseURI+"v1/chunk/:startIndex/:size", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		index, err := strconv.ParseInt(p.ByName("startIndex"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errMsg := "the start index must be a number, got " + p.ByName("startIndex")
			writeJson(w, &ChunkResponse{Error: &errMsg})
			return
		}
		size, err := strconv.ParseInt(p.ByName("size"), 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errMsg := "the size must be a number, got " + p.ByName("size")
			writeJson(w, &ChunkResponse{Error: &errMsg})
			return
		}
		getSize := size
		if getSize%2 != 0 {
			getSize++
		}
		firstIndex := index
		if firstIndex%2 != 0 {
			firstIndex--
			getSize += 2
		}

		chnk, err := api.GetChunk(firstIndex, int(getSize))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errMsg := err.Error()
			writeJson(w, &ChunkResponse{Error: &errMsg})
			return
		}

		unChnk := piio.AsUncompressedChunk(chnk)

		if firstIndex != index && len(unChnk.Digits) > 1 {
			// We requested one too early
			unChnk.Digits = unChnk.Digits[1:]
			unChnk.FirstDigitIndex++
		}
		if size != getSize && len(unChnk.Digits) > int(size) {
			unChnk.Digits = unChnk.Digits[:size]
		}

		digits := make([]int, len(unChnk.Digits))
		for i, d := range unChnk.Digits {
			digits[i] = int(d)
		}
		writeJson(w, &ChunkResponse{
			FirstIndex: unChnk.FirstIndex(),
			Digits:     digits,
		})
	})
	router.GET(BaseURI+"v1/settings", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		avail, err := chunkSource.AvailableDigits()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errMsg := err.Error()
			writeJson(w, &SettingsResponse{
				Error: &errMsg,
			})
		}
		writeJson(w, &SettingsResponse{
			AvailableDigits:  avail,
			MaximumChunkSize: chunkSource.MaximumChunkSize(),
		})
	})

	return api
}

func (api *API) GetDigit(index int64) (byte, error) {
	firstIndex := index
	if firstIndex%2 != 0 {
		firstIndex--
	}
	chnk, err := api.chunkSource.GetChunk(firstIndex, 2)
	if err != nil {
		return 255, errors.Wrap("could not load digit", err)
	}
	d, err := chnk.Digit(index)
	if err != nil {
		return 255, errors.Wrap("could not load digit", err)
	}
	return d, nil
}

func (api *API) GetChunk(firstIndex int64, size int) (piio.Chunk, error) {
	return api.chunkSource.GetChunk(firstIndex, size)
}

func (api *API) Handler() http.Handler {
	return api.router
}
