package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type course struct {
	title string
	link  string
}

func main() {
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	courses := []course{}

	c.OnHTML("td[class=vrtx-course-description-name]", func(e *colly.HTMLElement) {
		nodeContent := e.Text
		text := strings.ReplaceAll(strings.Trim(nodeContent, "\t\n "), "\n", "")
		link := e.ChildAttr("a", "href")

		courses = append(courses, course{text, link})
	})

	c.OnHTML("a[class=vrtx-next]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Following link %v\n", link)
		e.Request.Visit(link)
	})

	c.OnScraped(func(_ *colly.Response) {
		for _, course := range courses {
			fmt.Printf("%v -> %v\n", course.title, course.link)
		}
	})

	c.Visit("https://www.uio.no/studier/emner/alle/?filter.semester=v22")
}
