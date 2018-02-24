package main

import (
	"os"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "piio"
	app.Usage = "supply digits of Pi via a RESTful API"

	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{
	// 		Name:  "config, c",
	// 		Usage: "Load configuration from `FILE`",
	// 	},
	// }

	app.Commands = []cli.Command{
		{
			Name:  "compress",
			Usage: "compresses a text file of digits of pi",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	app.Run(os.Args)
}
