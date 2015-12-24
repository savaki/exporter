package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/savaki/exporter/search"
)

type Request struct {
	Id  int
	Url string
	Dir string
}

func crawl(c *cli.Context) {
	dir := c.String("dir")
	pages := c.Int("pages")
	key := c.String("key")

	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		} else {
			log.Fatalln(err)
		}
	}

	codebase := c.String("codebase")
	if codebase == "" {
		log.Fatalln("missing prefix parameter")
	}

	ch := make(chan *Request)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go crawler(wg, ch)

	for i := 1; i <= pages; i++ {
		ch <- &Request{
			Id:  i,
			Url: fmt.Sprintf("%v?%v=%v", codebase, key, i),
			Dir: dir,
		}
	}
	close(ch)

	wg.Wait()
}

func crawler(wg *sync.WaitGroup, ch <-chan *Request) {
	defer wg.Done()

	for req := range ch {
		fmt.Printf("fetching page %d\n", req.Id)
		results, err := fetch(req.Url)
		if err != nil {
			log.Fatalf("unable to fetch contents, %v - %v\n", req.Url, err)
		}

		u, err := url.Parse(req.Url)

		for _, result := range results {
			if result.Kind == "Partner" {
				copy(req.Dir, fmt.Sprintf("%v://%v%v", u.Scheme, u.Host, result.Url))
			}
		}
	}
}

func fetch(url string) ([]*search.Result, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("unable to fetch url, %v - %v\n", url, err)
	}
	defer resp.Body.Close()

	return search.Parse(resp.Body)
}

func copy(dir, url string) {
	fmt.Printf("retrieving partner, %v\n", url)

	base := filepath.Base(url)
	filename := fmt.Sprintf("%v/%v.html", dir, base)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	ioutil.WriteFile(filename, data, 0644)
}
