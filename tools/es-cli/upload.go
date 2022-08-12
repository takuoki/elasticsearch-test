package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

func init() {
	subCmdList = append(subCmdList, &cli.Command{
		Name:  "upload",
		Usage: "Upload data to elastic search.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "index",
				Aliases:  []string{"i"},
				Usage:    "Index name to upload to.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "Path of the file containing the data you want to upload (Newline-delimited JSON).",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			return action(c, &uploadCommand{out: os.Stdout})
		},
	})
}

type uploadCommand struct {
	out io.Writer
}

func (sc *uploadCommand) Run(c *cli.Context, baseURL string) error {

	const (
		initialBufSize = 10000
		maxBufSize     = math.MaxInt32
	)

	file, err := os.Open(c.String("file"))
	if err != nil {
		return fmt.Errorf("fail to open file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, initialBufSize), maxBufSize)
	cnt := 0
	for scanner.Scan() {
		err := func() error {
			if scanner.Text() == "" {
				return nil
			}

			buf := bytes.NewBuffer([]byte(scanner.Text()))
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/_doc", baseURL, c.String("index")), buf)
			if err != nil {
				return fmt.Errorf("fail to create request (success count: %d): %w", cnt, err)
			}

			req.Header.Add("Content-Type", "application/json")

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("fail to send request (success count: %d): %w", cnt, err)
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				return fmt.Errorf("fail to read body (success count: %d): %w", cnt, err)
			}

			if res.StatusCode != http.StatusCreated {
				return fmt.Errorf("fail to upload data (success count: %d): %s", cnt, string(body))
			}

			fmt.Fprintln(sc.out, string(body))
			cnt++

			return nil
		}()
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("fail to scan (success count: %d): %w", cnt, err)
	}

	fmt.Fprintf(sc.out, "successfully uploaded (count: %d)\n", cnt)

	return nil
}
