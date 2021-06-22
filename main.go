package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "kafktail",
		Usage: "tail kafka logs of any wire format",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "bootstrap_servers",
				Usage:    "list of kafka `bootstrap_servers` separated by comma",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "topic",
				Usage:    "`topic` whose message you want to tail",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Hello World")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
