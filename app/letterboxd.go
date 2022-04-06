package main

import (
	"fmt"
	"os"

	"github.com/anaskhan96/soup"
)

type letterboxdMovie struct {
	Slug  string
	Title string
	Year  string
}

func searchLb(query string) []letterboxdMovie {
	var resList []letterboxdMovie
	sUrl := fmt.Sprintf("https://letterboxd.com/search/%s", query)
	resp, err := soup.Get(sUrl)
	if err != nil {
		os.Exit(1)
	}
	doc := soup.HTMLParse(resp)
	posters := doc.Find("ul", "class", "results").FindAll("div", "class", "linked-film-poster")
	for _, poster := range posters {
		meta := poster.FindNextElementSibling()
		metaHead := meta.Find("span", "class", "film-title-wrapper")
		metaLinks := metaHead.FindAll("a")
		slug := poster.Attrs()["data-film-slug"]
		title := metaLinks[0].Text()
		year := metaLinks[1].Text()
		result := letterboxdMovie{
			Slug:  slug,
			Title: title,
			Year:  year,
		}
		resList = append(resList, result)
	}

	return resList
}
