package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"io"
	"net/http"
	"os"
	"strings"
)

var Version = "dev"

func main() {
	cmd := &cli.Command{
		Name:    "gurl",
		Usage:   "A simple HTTP client",
		Suggest: true,

		ArgsUsage: "<url>",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "url",
			},
		},

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "request",
				Value:   "GET",
				Usage:   "HTTP method to use",
				Aliases: []string{"X"},
			},
			&cli.StringSliceFlag{
				Name:    "header",
				Usage:   "Custom header(s) to include in the request",
				Aliases: []string{"H"},
			},
			&cli.StringFlag{
				Name:    "data",
				Usage:   "HTTP POST data",
				Aliases: []string{"d"},
			},
			&cli.StringFlag{
				Name:  "data-raw",
				Usage: "HTTP POST data, '@' allowed, '@-' reads from stdin",
			},
			&cli.StringFlag{
				Name:    "user",
				Usage:   "Basic auth user and password (e.g., user:password)",
				Aliases: []string{"u"},
			},
			&cli.BoolFlag{
				Name:    "head",
				Usage:   "Show response headers only",
				Aliases: []string{"I"},
			},
			&cli.BoolFlag{
				Name:    "fail",
				Usage:   "Fail fast with no output on HTTP errors",
				Aliases: []string{"f"},
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "Show verbose output",
				Aliases: []string{"v"},
			},
			&cli.BoolFlag{
				Name:    "version",
				Usage:   "Show version number and quit",
				Aliases: []string{"V"},
			},
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Bool("version") {
				fmt.Println("gurl version", Version)
				return nil
			}

			url := cmd.StringArg("url")
			if cmd.StringArg("url") == "" {
				return cli.DefaultShowRootCommandHelp(cmd)
			}

			method := cmd.String("request")
			headers := make(map[string]string)

			for _, h := range cmd.StringSlice("header") {
				nameAndValue := strings.SplitN(h, ":", 2)
				if len(nameAndValue) != 2 {
					continue
				}
				headers[strings.ToLower(nameAndValue[0])] = strings.TrimSpace(nameAndValue[1])
			}
			if _, ok := headers["user-agent"]; !ok {
				headers["user-agent"] = "gurl/" + Version
			}
			data := cmd.String("data")
			dataRaw := cmd.String("data-raw")
			var body io.Reader
			if dataRaw != "" {
				if strings.HasPrefix(dataRaw, "@") {
					filePath := strings.TrimPrefix(dataRaw, "@")
					if filePath == "-" {
						body = os.Stdin
					} else {
						file, err := os.Open(filePath)
						if err != nil {
							return err
						}
						defer file.Close()
						body = file
					}
				} else {
					data = dataRaw
				}
			}
			if data != "" {
				body = io.NopCloser(strings.NewReader(data))
			}

			verbose := cmd.Bool("verbose")
			headOnly := cmd.Bool("head")
			failFast := cmd.Bool("fail")

			client := &http.Client{}
			req, err := http.NewRequest(method, url, body)
			if err != nil {
				return err
			}

			userAndPassword := strings.Split(cmd.String("user"), ":")
			user := userAndPassword[0]
			password := ""
			if len(userAndPassword) > 1 {
				password = userAndPassword[1]
			}
			if user != "" && password != "" {
				req.SetBasicAuth(user, password)
			}

			for k, v := range headers {
				req.Header.Set(k, v)
			}

			resp, err := client.Do(req)
			if err != nil {
				return err
			}

			if failFast {
				if resp.StatusCode >= 400 {
					return fmt.Errorf("the requested URL returned error: %s", resp.Status)
				}
			}

			if verbose || headOnly {
				fmt.Println(headers)
				fmt.Println(resp.Header)
				fmt.Println(resp.Status)
			}

			if !headOnly {
				defer resp.Body.Close()
				bodyText, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}

				fmt.Print(string(bodyText))
			}

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Println("gurl:", err)
	}
}
