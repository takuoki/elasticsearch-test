package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	version = "1.0.0"
)

var (
	subCmdList = []*cli.Command{}
)

type subCmd interface {
	Run(ctx *cli.Context, baseURL string) error
}

func main() {
	app := cli.NewApp()
	app.Name = "elastic-search-cli"
	app.Version = version
	app.Usage = "A tool for elastic search."
	app.Commands = subCmdList
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "host",
			Usage: "Host name of elastic search server.",
			Value: "localhost",
		},
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "Port number of elastic search server.",
			Value:   "9200",
		},
		&cli.BoolFlag{
			Name:    "secure",
			Aliases: []string{"s"},
			Usage:   "Use TLS connection.",
			Value:   false,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err.Error())
	}
}

func action(c *cli.Context, sc subCmd) error {
	var baseURL string
	if c.Bool("secure") {
		baseURL = fmt.Sprintf("https://%s:%s", c.String("host"), c.String("port"))
	} else {
		baseURL = fmt.Sprintf("http://%s:%s", c.String("host"), c.String("port"))
	}

	return sc.Run(c, baseURL)
}
