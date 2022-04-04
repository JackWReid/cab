package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"
)

func addMedia(db *sql.DB, mediaType string, mediaTitle string) {
	var existingBooks []bookRecord
	fmt.Println("Check for existing books called", mediaTitle)
	bookQuery := `SELECT id, title, author, isbn FROM book WHERE title LIKE '%' || $1 || '%'`
	rows, err := db.Query(bookQuery, mediaTitle)
	checkErr(err)

	var exBookIds []string
	for rows.Next() {
		var row bookRecord
		err = rows.Scan(&row.Id, &row.Title, &row.Author, &row.Isbn)
		checkErr(err)
		exBookIds = append(exBookIds, row.Id)
		existingBooks = append(existingBooks, row)
	}

	rows.Close()

	tableBookRecord(existingBooks)

	var exBookRes string
	fmt.Println("Enter existing book ID or [x] to search Google Books:")
	fmt.Scanln(&exBookRes)

	if slices.Contains(exBookIds, exBookRes) {
		f := func(b bookRecord) bool {
			return exBookRes == b.Id
		}
		selectedBookIdx := slices.IndexFunc(existingBooks, f)
		selectedBook := existingBooks[selectedBookIdx]
		fmt.Println(selectedBook)
	} else {
		gbResults := searchGoogleBooks(mediaTitle)
		tableGoogleResults(gbResults)

		var gbBookRes string
		fmt.Println("Enter # of book from Google Books:")
		fmt.Scanln(&gbBookRes)

		idx, _ := strconv.Atoi(gbBookRes)
		sb := gbResults[idx]

		insQ := `INSERT INTO book(title, author, isbn) VALUES(?, ?, ?)`
		db.Exec(insQ, sb.title, sb.author, sb.isbn)

		// var logRes string
		// fmt.Println("Is this book:\n1. Read already\n2. Being read right now\n3. To read")
		// fmt.Scanln(&logRes)

		// switch logRes {
		// case "1":

		// case "2":
		// case "3":
		// }
	}
}
