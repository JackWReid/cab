package main

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

type OkuBookEvent struct {
	Date    *time.Time
	Title   string
	Author  string
	OkuGuid string
}

type OkuFeed struct {
	Date  *time.Time
	Items []OkuBookEvent
}

var OkuReadUrl string = "https://oku.club/rss/collection/zQtTo"
var OkuToreadUrl string = "https://oku.club/rss/collection/JSKHS"
var OkuReadingUrl string = "https://oku.club/rss/collection/2f67M"

func getOkuFetcher(feedUrl string) func() OkuFeed {
	return func() OkuFeed {
		var resFeed OkuFeed
		var resItems []OkuBookEvent
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(feedUrl)
		for _, feedEvent := range feed.Items {
			okuEvent := OkuBookEvent{
				Date:    feedEvent.PublishedParsed,
				Title:   feedEvent.Title,
				Author:  feedEvent.Author.Name,
				OkuGuid: feedEvent.GUID,
			}
			resItems = append(resItems, okuEvent)
		}

		resFeed.Date = feed.UpdatedParsed
		resFeed.Items = resItems

		resolveOkuFeedBooks(resFeed.Items)

		return resFeed
	}
}

func resolveOkuFeeds() {

}

func giveBookOkuGuid(book bookRecord, guid string) error {
	upQ := `UPDATE book SET oku_guid = $1 WHERE id = $2`
	_, err := DB.Exec(upQ, guid, book.Id)

	if err != nil {
		fmt.Println("Failed to update Oku GUID for", book.Title)
	}

	return err
}

func resolveOkuFeedBooks(be []OkuBookEvent) {
	var mBooks []bookRecord
	var unBooks []OkuBookEvent
	for _, e := range be {
		sr, _ := searchBooks(e.Title)
		if len(sr.Books) > 0 {
			mBooks = append(mBooks, sr.Books[0])
			giveBookOkuGuid(sr.Books[0], e.OkuGuid)
		} else {
			unBooks = append(unBooks, e)
		}
	}

	var auUnBooks []OkuBookEvent

	for _, e2 := range unBooks {
		sr, _ := searchBooksByAuthor(e2.Author)
		if len(sr.Books) == 1 {
			mBooks = append(mBooks, sr.Books[0])
			giveBookOkuGuid(sr.Books[0], e2.OkuGuid)
		} else {
			auUnBooks = append(auUnBooks, e2)
		}
	}

	for _, um := range auUnBooks {
		fmt.Println(um.Title)
	}
}
