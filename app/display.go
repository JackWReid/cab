package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

type bookEvent struct {
	Title  string
	Author *string
	Date   *string
}

type bookRecord struct {
	Id     string
	Title  string
	Author *string
	Isbn   *string
}

type movieRecord struct {
	Id    string
	Title string
	Year  string
	Slug  string
}

func parseDate(rawString string) (date time.Time) {
	var t time.Time
	if len(rawString) == 10 {
		t, _ = time.Parse("2006-01-02", rawString)
	} else {
		t, _ = time.Parse("2006-01-02 15:04:05 UTC", rawString)
	}

	return t
}

func tableBookEvents(events []bookEvent, limit int) {
	tab := table.NewWriter()
	tab.SetOutputMirror(os.Stdout)
	tab.AppendHeader(table.Row{"Date", "Title", "Author"})

	limitedEvents := events
	if len(events) > limit {
		limitedEvents = events[:limit]
	}

	for _, row := range limitedEvents {
		strDate := ""
		strAuthor := ""
		if row.Date != nil {
			strDate = *row.Date
		}
		if row.Author != nil {
			strAuthor = *row.Author
		}
		date := parseDate(strDate).Format("2006-01-02")
		tab.AppendRow([]interface{}{date, text.Trim(row.Title, 50), strAuthor})
	}

	tab.Render()
}

func jsonBookEvents(events []bookEvent) {
	jsonBytes, err := json.Marshal(events)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonBytes))
}

func tableGoogleResults(results []condensedGoogleResult) {
	tab := table.NewWriter()
	tab.SetOutputMirror(os.Stdout)
	tab.AppendHeader(table.Row{"#", "ISBN", "Title", "Subtitle", "Author"})

	for i, row := range results {
		tab.AppendRow([]interface{}{
			i,
			row.isbn,
			text.Trim(row.title, 30),
			text.Trim(row.subtitle, 50),
			row.author,
		})
	}

	tab.Render()
}

func tableBookRecord(results []bookRecord) {
	tab := table.NewWriter()
	tab.SetOutputMirror(os.Stdout)
	tab.AppendHeader(table.Row{"ID", "Title", "Author"})

	for _, row := range results {
		tab.AppendRow([]interface{}{
			row.Id,
			text.Trim(row.Title, 50),
			text.Trim(*row.Author, 20),
		})
	}

	tab.Render()
}

func tableMovieEvents(events []movieEvent, limit int) {
	tab := table.NewWriter()
	tab.SetOutputMirror(os.Stdout)
	tab.AppendHeader(table.Row{"Date", "Title", "Year"})

	limitedEvents := events
	if len(events) > limit {
		limitedEvents = events[:limit]
	}

	for _, row := range limitedEvents {
		date := parseDate(row.Date).Format("2006-01-02")
		tab.AppendRow([]interface{}{date, text.Trim(row.Title, 50), row.Year})
	}

	tab.Render()
}

func tableMovieRecord(results []movieRecord) {
	tab := table.NewWriter()
	tab.SetOutputMirror(os.Stdout)
	tab.AppendHeader(table.Row{"ID", "Title", "Year"})

	for _, row := range results {
		tab.AppendRow([]interface{}{
			row.Id,
			row.Title,
			row.Year,
		})
	}

	tab.Render()
}

func jsonMovieEvents(events []movieEvent) {
	jsonBytes, err := json.Marshal(events)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonBytes))
}

func tableLbResults(results []letterboxdMovie) {
	tab := table.NewWriter()
	tab.SetOutputMirror(os.Stdout)
	tab.AppendHeader(table.Row{"#", "Slug", "Title", "Year"})

	for i, row := range results {
		tab.AppendRow([]interface{}{
			i,
			row.Slug,
			row.Title,
			row.Year,
		})
	}

	tab.Render()
}
