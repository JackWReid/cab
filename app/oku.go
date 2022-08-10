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

type BookIdPair struct {
	Id   int64
	Guid string
}

func getOkuFetcher(collection string) func() OkuFeed {
	urlMap := map[string]string{
		"reading": GlobalConfig.OkuReadingURL,
		"toread":  GlobalConfig.OkuToreadURL,
		"read":    GlobalConfig.OkuReadURL,
	}

	return func() OkuFeed {
		var resFeed OkuFeed
		var resItems []OkuBookEvent
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(urlMap[collection])
		fmt.Println(urlMap[collection])

		if err != nil {
			fmt.Println(err)
		}

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

		resolveOkuFeed(collection, resFeed)

		return resFeed
	}
}

func resolveOkuFeed(collection string, feedData OkuFeed) {
	fmt.Println("resolve")
	fmt.Println(collection)
	for _, book := range feedData.Items {
		insQ := `INSERT OR IGNORE INTO book(title, author, oku_guid) VALUES (?, ?, ?)`
		_, err := DB.Exec(insQ, book.Title, book.Author, book.OkuGuid)

		if err != nil {
			fmt.Println(err)
		}

		var resolvedBook bookRecord
		resolvedBook, _ = getBookByGuid(book.OkuGuid)

		addBookToCollection(resolvedBook, collection, book.Date)
	}

	clearCollectionExcept(collection, feedData.Items)
}

func clearCollectionExcept(collection string, nBooks []OkuBookEvent) {
	fmt.Println("clear")
	fmt.Println(collection)
	var existingBookIds []BookIdPair

	q := `SELECT oku_guid FROM book INNER JOIN book_collection ON book_id = oku_guid WHERE collection = $1`
	rows, _ := DB.Query(q, collection)

	for rows.Next() {
		var oPair BookIdPair
		rows.Scan(&oPair.Id, &oPair.Guid)
		existingBookIds = append(existingBookIds, oPair)
	}

	rows.Close()

	for _, existingBook := range existingBookIds {
		var isBookStillInCollection bool = false
		var bookDate *time.Time
		for _, nBook := range nBooks {
			if nBook.OkuGuid == existingBook.Guid {
				isBookStillInCollection = true
				bookDate = nBook.Date
			}
		}

		if isBookStillInCollection {
			uQ := `UPDATE book_collection SET date_updates = $1 WHERE book_id = $2 AND collection = $3`
			_, err := DB.Exec(uQ, bookDate, existingBook.Id, collection)

			if err != nil {
				fmt.Println("Error updating")
				fmt.Println(err)
			}
		} else {
			dQ := `DELETE FROM book_collection WHERE book_id = $1 AND collection = $2`
			_, err := DB.Exec(dQ, existingBook.Id, collection)

			if err != nil {
				fmt.Println("Error deleting")
			}
		}
	}
}
