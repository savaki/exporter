package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/savaki/exporter/partner"
	"github.com/savaki/exporter/search"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:   "partner",
			Action: parsePartner,
		},
		{
			Name:   "search",
			Action: parseSearch,
		},
		{
			Name:   "crawl",
			Action: crawl,
			Flags: []cli.Flag{
				cli.StringFlag{"codebase", "", "codebase", ""},
				cli.StringFlag{"key", "PageNum", "page number parameter", ""},
				cli.StringFlag{"dir", "target", "output directory", ""},
				cli.IntFlag{"pages", 1, "number of pages to fetch", ""},
			},
		},
	}
	app.Run(os.Args)
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func parsePartner(c *cli.Context) {
	filename := c.Args().First()
	data, err := ioutil.ReadFile(filename)
	check(err)

	partner, err := partner.Parse(bytes.NewReader(data))
	check(err)

	data, err = json.MarshalIndent(partner, "", "  ")
	check(err)

	fmt.Println(string(data))
}

func parseSearch(c *cli.Context) {
	filename := c.Args().First()
	data, err := ioutil.ReadFile(filename)
	check(err)

	results, err := search.Parse(bytes.NewReader(data))
	check(err)

	data, err = json.MarshalIndent(results, "", "  ")
	check(err)

	fmt.Println(string(data))
}
