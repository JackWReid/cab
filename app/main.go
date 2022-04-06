package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func queryGenerate(queryType string) (queryString string) {
	switch queryType {
	case "reading":
		return "SELECT title, author, date FROM view_book_reading ORDER BY date DESC"
	case "read":
		return "SELECT title, author, date FROM view_book_read ORDER BY date DESC"
	case "toread":
		return "SELECT title, author, date FROM view_book_toread ORDER BY date DESC"
	case "watched":
		return "SELECT title, date FROM view_movie_watched ORDER BY date DESC"
	case "towatch":
		return "SELECT title, date FROM view_movie_towatch ORDER BY date DESC"
	default:
		return ""
	}
}

func mediaList(db *sql.DB, listType string, jsonFlag bool) {
	bookEvents := []bookEvent{}
	var rows *sql.Rows
	var err error
	queryString := queryGenerate(listType)
	rows, err = db.Query(queryString)
	checkErr(err)

	for rows.Next() {
		var row bookEvent
		err = rows.Scan(&row.Title, &row.Author, &row.Date)
		checkErr(err)
		bookEvents = append(bookEvents, row)
	}

	if jsonFlag {
		jsonBookChrono(bookEvents)
	} else {
		tableBookChrono(bookEvents, 1000)
	}

	rows.Close()
}

func main() {
	db, connErr := sql.Open("sqlite3", "./media.db")
	checkErr(connErr)

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listJson := listCmd.Bool("json", false, "Output lists in JSON, default is table")
	listMedia := listCmd.String("media", "reading", "Which media and status to list")

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addMediaType := addCmd.String("type", "book", "Book or movie")
	addTitle := addCmd.String("title", "", "Name of the book or movie to add")

	switch os.Args[1] {
	case "list":
		listCmd.Parse(os.Args[2:])
		mediaList(db, *listMedia, *listJson)
	case "backfill":
		backfillBooks(db)
	case "add":
		addCmd.Parse(os.Args[2:])
		switch *addMediaType {
		case "book":
			addBook(db, *addTitle)
		case "movie":
			addMovie(db, *addTitle)
		default:
			fmt.Println("Invalid media type to add: book or movie")
		}

	default:
		fmt.Println("\ncot (1) help")
		fmt.Println("\nSUBCOMMANDS")
		fmt.Println("\nlist\nPrint a list of media that's been marked as done, doing, or to do")
		listCmd.PrintDefaults()
		fmt.Println("\nbackfill\nUse Google Books to fill in missing author and ISBN data\n\n")
		os.Exit(0)
	}

	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
