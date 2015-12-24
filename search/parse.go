package search

import (
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

func dump(s *goquery.Selection) {
	h, err := s.Html()
	if err == nil {
		fmt.Println("##", "--", "before", "------------------------------------------------")
		fmt.Println(h)
		fmt.Println("##", "--", "after ", "------------------------------------------------")
	}
}

type Result struct {
	Kind string
	Url  string
}

type parser struct {
	Results []*Result
	Err     error
}

func (p *parser) append(kind, url string) {
	if p.Results == nil {
		p.Results = []*Result{}
	}

	p.Results = append(p.Results, &Result{
		Kind: kind,
		Url:  url,
	})
}

func (p *parser) parse(s *goquery.Selection) {
	s.Each(func(i int, content *goquery.Selection) {
		kind := content.Find(".position-sec strong").First().Text()
		url := content.Find(".name-sec a").First().AttrOr("href", "")
		p.append(kind, url)
	})
}

func Parse(r io.Reader) ([]*Result, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	p := &parser{}
	p.parse(doc.Find("table.search-results .result"))
	return p.Results, p.Err
}
