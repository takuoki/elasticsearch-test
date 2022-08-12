package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

func init() {
	subCmdList = append(subCmdList, &cli.Command{
		Name:  "delete",
		Usage: "Delete data.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "index",
				Aliases:  []string{"i"},
				Usage:    "Index name to upload to.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "ID you want to delete.",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			return action(c, &deleteCommand{out: os.Stdout})
		},
	})
}

type deleteCommand struct {
	out io.Writer
}

func (sc *deleteCommand) Run(c *cli.Context, baseURL string) error {

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s/_doc/%s", baseURL, c.String("index"), c.String("id")), nil)
	if err != nil {
		return fmt.Errorf("fail to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fail to send request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("fail to read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("fail to delete data: %s", string(body))
	}

	fmt.Fprintln(sc.out, string(body))

	return nil
}
