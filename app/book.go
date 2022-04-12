package main

import (
	"errors"
	"fmt"
)

type bookRecord struct {
	Id     string
	Title  string
	Author *string
	Isbn   *string
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
		qStr = "SELECT title, author, date FROM view_book_reading ORDER BY date DESC"
	case "done":
		qStr = "SELECT title, author, date FROM view_book_read ORDER BY date DESC"
	case "todo":
		qStr = "SELECT title, author, date FROM view_book_toread ORDER BY date DESC"
	default:
		return bookResults, errors.New("Invalid bookStatus")
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

func insertBook(gb condensedGoogleResult) (bookRecord, error) {
	var insertedBook bookRecord
	insQ := `INSERT INTO book(title, author, isbn) VALUES(?, ?, ?)`
	result, err := DB.Exec(insQ, gb.title, gb.author, gb.isbn)

	if err != nil {
		return insertedBook, err
	}

	id, _ := result.LastInsertId()

	reQ := `SELECT id, title, author FROM book WHERE id = $1`
	row := DB.QueryRow(reQ, id)

	row.Scan(&insertedBook.Id, &insertedBook.Title, &insertedBook.Author)

	return insertedBook, nil
}

func readingBook(book bookRecord) error {
	insQ := `INSERT INTO book_reading(book_id) VALUES(?)`
	_, err := DB.Exec(insQ, book.Id)

	if err != nil {
		return err
	}

	return nil
}

func readBook(book bookRecord) error {
	cheQ := `SELECT COUNT(*) FROM book_log WHERE book_id = $1`
	row := DB.QueryRow(cheQ, book.Id)
	var cheCount int
	row.Scan(&cheCount)

	if cheCount > 0 {
		fmt.Println(book.Title, "already marked as read")
	} else {
		insQ := `INSERT INTO book_log(book_id) VALUES(?)`
		_, err := DB.Exec(insQ, book.Id)

		if err != nil {
			fmt.Println("insert to log err")
			return err
		}
	}

	unreadingQ := `DELETE FROM book_reading WHERE book_id = $1`
	_, err := DB.Exec(unreadingQ, book.Id)

	untoreadQ := `DELETE FROM book_to_book_collection WHERE book_id = $1 AND collection_id = 1`
	_, err = DB.Exec(untoreadQ, book.Id)

	if err != nil {
		return err
	}

	return nil
}

func toreadBook(book bookRecord) error {
	insQ := `INSERT INTO book_to_book_collection(book_id, collection_id) VALUES(?, 1)`
	_, err := DB.Exec(insQ, book.Id)

	if err != nil {
		return err
	}

	return nil
}

func searchBooks(query string) (bookSearchResults, error) {
	var bookResults []bookRecord
	bookQuery := `SELECT id, title, author, isbn FROM book WHERE title LIKE '%' || $1 || '%'`
	rows, err := DB.Query(bookQuery, query)
	checkErr(err)

	var bookIds []string
	for rows.Next() {
		var row bookRecord
		err = rows.Scan(&row.Id, &row.Title, &row.Author, &row.Isbn)
		checkErr(err)
		bookIds = append(bookIds, row.Id)
		bookResults = append(bookResults, row)
	}

	rows.Close()

	return bookSearchResults{
		Count:   len(bookResults),
		BookIds: bookIds,
		Books:   bookResults,
	}, nil
}
