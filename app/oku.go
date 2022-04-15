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

var OkuReadUrl string = "https://oku.club/rss/collection/zQtTo"
var OkuToreadUrl string = "https://oku.club/rss/collection/JSKHS"
var OkuReadingUrl string = "https://oku.club/rss/collection/2f67M"

func getOkuFetcher(collection string) func() OkuFeed {
	urlMap := map[string]string{
		"reading": OkuReadingUrl,
		"toread":  OkuToreadUrl,
		"read":    OkuReadUrl,
	}

	return func() OkuFeed {
		var resFeed OkuFeed
		var resItems []OkuBookEvent
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(urlMap[collection])

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
	idMap := map[string]int64{
		"reading": readingCollectionId,
		"toread":  toreadCollectionId,
		"read":    readCollectionId,
	}

	for _, book := range feedData.Items {
		var insertId int64 = 0
		insQ := `INSERT OR IGNORE INTO book(title, author, oku_guid) VALUES (?, ?, ?)`
		result, err := DB.Exec(insQ, book.Title, book.Author, book.OkuGuid)

		if result != nil {
			insertId, err = result.LastInsertId()
		}

		if err != nil {
			fmt.Println(err)
		}

		var resolvedBook bookRecord
		if insertId == 0 {
			resolvedBook, _ = getBookByGuid(book.OkuGuid)
		} else {
			resolvedBook, _ = getBookById(insertId)
		}

		addBookToCollection(resolvedBook, idMap[collection], book.Date)
	}

	clearCollectionExcept(idMap[collection], feedData.Items)
}

func clearCollectionExcept(collectionId int64, nBooks []OkuBookEvent) {
	var existingBookIds []BookIdPair

	q := `SELECT id, oku_guid FROM book INNER JOIN book_to_book_collection ON book_id = id WHERE collection_id = $1`
	rows, _ := DB.Query(q, collectionId)

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
			uQ := `UPDATE book_to_book_collection SET date = $1 WHERE book_id = $2 AND collection_id = $3`
			_, err := DB.Exec(uQ, bookDate, existingBook.Id, collectionId)

			if err != nil {
				fmt.Println("Error updating")
				fmt.Println(err)
			}
		} else {
			fmt.Println("Deleting", existingBook.Guid, "from", collectionId)
			dQ := `DELETE FROM book_to_book_collection WHERE book_id = $1 AND collection_id = $2`
			_, err := DB.Exec(dQ, existingBook.Id, collectionId)

			if err != nil {
				fmt.Println("Error deleting")
			}
		}
	}
}
