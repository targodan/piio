package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/targodan/piio"

	"github.com/julienschmidt/httprouter"
	errors "github.com/targodan/go-errors"
)

type API struct {
	router      *httprouter.Router
	chunkSource piio.ChunkSource
}

func writeJson(w http.ResponseWriter, data interface{}) {
	jw := json.NewEncoder(w)
	jw.Encode(data)
}

func NewAPI(chunkSource piio.ChunkSource) *API {
	versionRouter := httprouter.New()
	api := &API{
		router:      versionRouter,
		chunkSource: chunkSource,
	}

	v1Router := httprouter.New()

	v1Router.GET("/digit/:index", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		index, err := strconv.ParseInt(p.ByName("index"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJson(w, &DigitResponse{Error: "the index must be a number, got " + p.ByName("index")})
		}
		d, err := api.GetDigit(index)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJson(w, &DigitResponse{Error: err.Error()})
		}
		writeJson(w, &DigitResponse{
			Index: index,
			Digit: d,
		})
	})

	v1Router.GET("/chunk/:startIndex/:size", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		index, err := strconv.ParseInt(p.ByName("startIndex"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJson(w, &ChunkResponse{Error: "the start index must be a number, got " + p.ByName("startIndex")})
		}
		size, err := strconv.ParseInt(p.ByName("size"), 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJson(w, &ChunkResponse{Error: "the size must be a number, got " + p.ByName("size")})
		}
		chnk, err := api.GetChunk(index, int(size))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJson(w, &ChunkResponse{Error: err.Error()})
		}
		unChnk := piio.AsUncompressedChunk(chnk)
		writeJson(w, &ChunkResponse{
			FirstIndex: unChnk.FirstIndex(),
			Digits:     unChnk.Digits,
		})
	})

	versionRouter.Handler("GET", "/v1", v1Router)

	return api
}

func (api *API) GetDigit(index int64) (byte, error) {
	chnk, err := api.chunkSource.GetChunk(index, 1)
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
