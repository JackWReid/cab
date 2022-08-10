package main

import (
	"errors"
	"fmt"
	"time"
)

type bookRecord struct {
	Title   string
	Author  *string
	OkuGuid *string
}

type bookEvent struct {
	Title  string  `json:"title"`
	Author *string `json:"author"`
	Date   *string `json:"date_updated"`
}

type bookSearchResults struct {
	Count   int
	BookIds []string
	Books   []bookRecord
}

func getBooksByStatus(bookStatus string) ([]bookEvent, error) {
	var qStr string
	var bookResults []bookEvent

	switch bookStatus {
	case "doing":
		qStr = "SELECT title, author, date_updated FROM book INNER JOIN book_collection ON oku_guid = book_id WHERE collection = reading ORDER BY date DESC"
	case "done":
		qStr = "SELECT title, author, date_updated FROM book INNER JOIN book_collection ON oku_guid = book_id WHERE collection = read ORDER BY date DESC"
	case "todo":
		qStr = "SELECT title, author, date_updated FROM book INNER JOIN book_collection ON oku_guid = book_id WHERE collection = toread ORDER BY date DESC"
	default:
		return bookResults, errors.New("Invalid book status")
	}

	rows, err := DB.Query(qStr)

	for rows.Next() {
		var row bookEvent
		err = rows.Scan(&row.Title, &row.Author, &row.Date)
		checkErr(err)
		bookResults = append(bookResults, row)
	}

	rows.Close()

	if err != nil {
		return bookResults, err
	}

	return bookResults, nil
}

func searchBooks(titleQuery string) (bookSearchResults, error) {
	var bookResults []bookRecord
	bookQuery := `SELECT title, author FROM book WHERE title LIKE '%' || $1 || '%'`
	rows, err := DB.Query(bookQuery, titleQuery)
	checkErr(err)

	var bookIds []string
	for rows.Next() {
		var row bookRecord
		err = rows.Scan(&row.Title, &row.Author)
		checkErr(err)
		bookIds = append(bookIds, *row.OkuGuid)
		bookResults = append(bookResults, row)
	}

	rows.Close()

	return bookSearchResults{
		Count:   len(bookResults),
		BookIds: bookIds,
		Books:   bookResults,
	}, nil
}

func searchBooksByAuthor(authorQuery string) (bookSearchResults, error) {
	var bookResults []bookRecord
	bookQuery := `SELECT title, author, isbn FROM book WHERE author LIKE '%' || $1 || '%'`
	rows, err := DB.Query(bookQuery, authorQuery)
	checkErr(err)

	var bookIds []string
	for rows.Next() {
		var row bookRecord
		err = rows.Scan(&row.Title, &row.Author)
		checkErr(err)
		bookIds = append(bookIds, *row.OkuGuid)
		bookResults = append(bookResults, row)
	}

	rows.Close()

	return bookSearchResults{
		Count:   len(bookResults),
		BookIds: bookIds,
		Books:   bookResults,
	}, nil
}

func getBookByGuid(guid string) (bookRecord, error) {
	var book bookRecord
	bookQuery := `SELECT title, author, oku_guid FROM book WHERE oku_guid = $1`
	row := DB.QueryRow(bookQuery, guid)
	err := row.Scan(&book.Title, &book.Author, &book.OkuGuid)

	if err != nil {
		return book, err
	}

	return book, err
}

func addBookToCollection(book bookRecord, collection string, date *time.Time) {
	insQ := `INSERT OR IGNORE INTO book_collection(book_id, collection, date_updated) VALUES (?,?,?)`
	_, err := DB.Exec(insQ, book.OkuGuid, collection, date)

	if err != nil {
		fmt.Println(err)
	}
}
