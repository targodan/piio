package main

import (
	"io"
	"os"

	"github.com/targodan/piio"
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

	app.Commands = []cli.Command{
		{
			Name:   "compress",
			Usage:  "compresses a text file of digits of pi",
			Action: compressAction,
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
			},
			Action: serveAction,
		},
	}

	app.Run(os.Args)
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
	return nil
}
