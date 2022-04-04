package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

type backfillBook struct {
	id    string
	title string
}

func backwriteBook(db *sql.DB, id string, author string, isbn string) {
	backwriteQuery := `UPDATE book SET author = $1, isbn = $2 WHERE id = $3`
	result, err := db.Exec(backwriteQuery, author, isbn, id)
	if err != nil {
		panic(err)
	}

	affectedRows, _ := result.RowsAffected()

	if affectedRows == 0 {
		log.Fatal("Warning: zero rows written: ", id, author, isbn)
	}
}

func backfillBooks(db *sql.DB) {
	queueQuery := "SELECT id, title FROM book WHERE author IS NULL OR isbn IS NULL ORDER BY id DESC"
	backfillQueue := []backfillBook{}

	rows, err := db.Query(queueQuery)
	checkErr(err)

	for rows.Next() {
		var row backfillBook
		err = rows.Scan(&row.id, &row.title)
		checkErr(err)
		backfillQueue = append(backfillQueue, row)
	}

	rows.Close()

	for _, i := range backfillQueue {
		fmt.Println("\nSearching Google Books for", i.title)
		bookResults := searchGoogleBooks(i.title)
		tableGoogleResults(bookResults)

		var selectedRawId string
		fmt.Println("Select ID to backfill, default [0] or [x] to skip:")
		fmt.Scanln(&selectedRawId)

		if selectedRawId == "x" {
			fmt.Println("Skipping...")
		} else {
			var selectedId int
			if len(selectedRawId) == 0 {
				selectedId = 0
			}
			selectedId, _ = strconv.Atoi(selectedRawId)
			sb := bookResults[selectedId]
			backwriteBook(db, i.id, sb.author, sb.isbn)
		}
	}
}
