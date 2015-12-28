package partner

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Link struct {
	Label string `json:"label,omitempty"`
	Href  string `json:"href,omitempty"`
}

type ContactInfo struct {
	Office   *Link  `json:"office,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Country  string `json:"country,omitempty"`
	LinkedIn *Link  `json:"linked-in,omitempty"`
	VCard    *Link  `json:"vcard,omitempty"`
	TimeZone string `json:"tz,omitempty"`
}

type Partner struct {
	Name            string            `json:"name,omitempty"`
	Image           string            `json:"image,omitempty"`
	Practices       map[string]string `json:"practices,omitempty"`
	IndustrySectors map[string]string `json:"industry-sectors,omitempty"`
	Education       []string          `json:"education,omitempty"`
	ContactInfo     *ContactInfo      `json:"contact-info,omitempty"`
	Summary         string            `json:"summary,omitempty"`
	SummaryRaw      string            `json:"summary-raw,omitempty"`
	Qualifications  []string          `json:"qualifications,omitempty"`
}

type ParseMap map[string]func(*Partner, *goquery.Selection)

func (p ParseMap) Parse(partner *Partner, s *goquery.Selection, selector string) {
	s.Each(func(index int, section *goquery.Selection) {
		section.Find(selector).Each(func(index int, content *goquery.Selection) {
			label := strings.TrimSpace(content.Text())

			for key, fn := range p {
				if key == label {
					fn(partner, section)
					return
				}
			}
		})
	})
}

func dump(s *goquery.Selection) {
	h, err := s.Html()
	if err == nil {
		fmt.Println("##", "--", "before", "------------------------------------------------")
		fmt.Println(h)
		fmt.Println("##", "--", "after ", "------------------------------------------------")
	}
}

// -- parse --------------------------------------------------------------------

func Parse(r io.Reader) (*Partner, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	s := doc.Find("body").First()
	partner := &Partner{}

	err = parseSummary(partner, s)
	if err != nil {
		return nil, err
	}

	parseSidebar(partner, s)
	if err != nil {
		return nil, err
	}

	return partner, nil
}

// -- Summary ------------------------------------------------------------------

func parseSummary(partner *Partner, s *goquery.Selection) error {
	// parse the summary
	summary := s.Find(".view-stories").First()
	h, err := summary.Html()
	if err != nil {
		return err
	}
	partner.Summary = strings.TrimSpace(summary.Text())
	partner.SummaryRaw = strings.TrimSpace(h)

	partner.Image = s.Find(".carousel-img img").First().AttrOr("src", "")
	partner.Name = s.Find("h1#top").First().Text()
	if strings.HasSuffix(partner.Name, "Partner") {
		l := len(partner.Name) - len("Partner")
		partner.Name = strings.TrimSpace(partner.Name[0:l])
	}

	return nil
}

// -- Sidebar ------------------------------------------------------------------

func parseSidebar(partner *Partner, s *goquery.Selection) {
	pm := ParseMap{
		"Areas of focus":                parseAreasOfFocus,
		"Contact information":           parseContactInfo,
		"Education":                     parseEducation,
		"Admissions and qualifications": parseQualifications,
	}
	pm.Parse(partner, s.Find(".BioLeftControl section"), "h4 a")
}

// -- Areas of focus -----------------------------------------------------------

func parsePractices(partner *Partner, s *goquery.Selection) {
	practices := map[string]string{}
	s.Find("a").Each(func(index int, content *goquery.Selection) {
		href := content.AttrOr("href", "")
		text := content.Text()
		practices[text] = href
	})
	partner.Practices = practices
}

func parseIndustrySectors(partner *Partner, s *goquery.Selection) {
	industrySectors := map[string]string{}
	s.Find("a").Each(func(index int, content *goquery.Selection) {
		href := content.AttrOr("href", "")
		text := content.Text()
		industrySectors[text] = href
	})
	partner.IndustrySectors = industrySectors
}

func parseAreasOfFocus(partner *Partner, s *goquery.Selection) {
	pm := ParseMap{
		"Practices":        parsePractices,
		"Industry sectors": parseIndustrySectors,
	}
	pm.Parse(partner, s.Find(".areaOfFocus"), "h5")
}

// -- Contact Information ------------------------------------------------------

func parseContactInfo(partner *Partner, s *goquery.Selection) {
	vCard := s.Find(".vcard .social-links").First().AttrOr("href", "")
	if strings.HasPrefix(vCard, "..") {
		vCard = vCard[len(".."):]
	}

	contactInfo := &ContactInfo{
		Office: &Link{
			Label: s.Find(".bio-contact a").First().Text(),
			Href:  s.Find(".bio-contact a").First().AttrOr("href", ""),
		},
		Email: s.Find(".pemail a").First().Text(),
		LinkedIn: &Link{
			Href: s.Find(".social-links").Last().AttrOr("href", ""),
		},
		VCard: &Link{
			Href: vCard,
		},
		TimeZone: s.Find(".timezone input").First().AttrOr("value", ""),
		Phone:    s.Find(".left_assign").First().Text(),
	}

	partner.ContactInfo = contactInfo
}

func parseEducation(partner *Partner, s *goquery.Selection) {
	education := []string{}
	s.Find("p").Each(func(i int, content *goquery.Selection) {
		education = append(education, content.Text())
	})
	partner.Education = education
}

func parseQualifications(partner *Partner, s *goquery.Selection) {
	qualifications := []string{}
	s.Find("p").Each(func(i int, content *goquery.Selection) {
		qualifications = append(qualifications, content.Text())
	})
	partner.Qualifications = qualifications
}
