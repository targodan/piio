package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/targodan/piio"
	"github.com/targodan/piio/rest"
	"gopkg.in/urfave/cli.v1"
)

const defaultChunkSize = 512

type readSeekCloser interface {
	io.ReadSeeker
	io.Closer
}
type writeSeekCloser interface {
	io.WriteSeeker
	io.Closer
}

func main() {
	app := cli.NewApp()
	app.Name = "piio"
	app.Usage = "supply digits of Pi via a RESTful API"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "pi,p",
			Usage: "The file of pi.",
			Value: "pi.bin",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "compress",
			Usage:  "compresses a text file of digits of pi",
			Action: compressAction,
		},
		{
			Name:   "search",
			Usage:  "searches for a string of digits in pi",
			Action: searchAction,
		},
		{
			Name:  "serve",
			Usage: "listen and serve",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "addr,a",
					Usage: "The address and port to listen on.",
					Value: "127.0.0.1:8080",
				},
				cli.IntFlag{
					Name:  "max-chunk-size,c",
					Usage: "The maximum size of a chunk to be served.",
					Value: defaultChunkSize,
				},
			},
			Action: serveAction,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func searchAction(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return cli.NewExitError("expected exactly one argument. Usage: piio search <searchString>", 1)
	}
	searchString := c.Args().Get(0)
	pifile := c.GlobalString("pi")

	index, err := piio.Search(pifile, searchString)
	if err != nil {
		return cli.NewExitError(err, 2)
	}
	if index != -1 {
		fmt.Printf("Found \"%s\" at position %d.\n", searchString, index)
	} else {
		fmt.Printf("Could not find \"%s\".\n", searchString)
	}
	return nil
}

func compressAction(c *cli.Context) error {
	infile := c.Args().Get(0)
	outfile := c.Args().Get(1)

	if infile == "" || outfile == "" {
		return cli.NewExitError("expected exactly two arguments usage: piio compress <infile> <outfile>", 1)
	}

	chunkSize := defaultChunkSize

	var in readSeekCloser
	var out io.WriteCloser
	var err error

	if infile == "-" {
		in = os.Stdin
	} else {
		in, err = os.Open(infile)
		if err != nil {
			return cli.NewExitError(err, 2)
		}
	}
	defer in.Close()
	if outfile == "-" {
		out = os.Stdout
	} else {
		out, err = os.Create(outfile)
		if err != nil {
			return cli.NewExitError(err, 2)
		}
	}
	defer out.Close()

	var chnk piio.Chunk
	for i := int64(0); err == nil; i += int64(chunkSize) {
		chnk, err = piio.ReadTextChunk(in, i, chunkSize)
		if err != nil {
			break
		}
		err = piio.WriteChunk(chnk, piio.FileFormatCompressed, out)
		if err != nil {
			break
		}
	}
	if err != io.EOF {
		return cli.NewExitError(err, 3)
	}

	return nil
}

func serveAction(c *cli.Context) error {
	chunkSource := piio.NewUncachedChunkSource(c.GlobalString("pi"), piio.FileFormatCompressed, c.Int("max-chunk-size"))
	api := rest.NewAPI(chunkSource)

	server := &http.Server{
		Addr:           c.String("addr"),
		MaxHeaderBytes: 512,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		Handler:        api.Handler(),
	}

	err := server.ListenAndServe()

	return err
}
