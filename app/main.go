package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func bookList(bookStatus string, jsonFlag bool) {
	bookEvents, err := getBooksByStatus(bookStatus)

	if err != nil {
		fmt.Println("Failed to get books by status")
		return
	}

	if jsonFlag {
		jsonBookEvents(bookEvents)
	} else {
		tableBookEvents(bookEvents, 1000)
	}
}

func movieList(movieStatus string, jsonFlag bool) {
	movieEvents, err := getMoviesByStatus(movieStatus)

	if err != nil {
		fmt.Println("Failed to get movies by status")
		fmt.Println(err)
		return
	}

	if jsonFlag {
		jsonMovieEvents(movieEvents)
	} else {
		tableMovieEvents(movieEvents, 1000)
	}
}

func printHelp(listCmd *flag.FlagSet, addCmd *flag.FlagSet) {
	fmt.Println("===\ncab")
	fmt.Println("cab is a way to manage logs and lists for books and movies.\n===\n")
	fmt.Println("SUBCOMMANDS")
	fmt.Println("\nlist\nPrint a list of media that's been marked as done, doing, or to do")
	listCmd.PrintDefaults()
	fmt.Println("\nadd\nAdd books and movies from Google Books and Letterboxd")
	addCmd.PrintDefaults()
	fmt.Println("\nbackfill\nUse Google Books to fill in missing author and ISBN data\n")
	os.Exit(0)
}

func main() {
	db, connErr := sql.Open("sqlite", "./media.db")

	if connErr != nil {
		fmt.Println("Failed to connect to DB")
		fmt.Println(connErr)
		os.Exit(1)
	}

	DB = db

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listJson := listCmd.Bool("json", false, "Output lists in JSON, default is table")
	listMediaType := listCmd.String("type", "book", "Book or movie")
	listMediaStatus := listCmd.String("status", "done", "Done, doing, todo")

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addMediaType := addCmd.String("type", "book", "Book or movie")
	addTitle := addCmd.String("title", "", "Name of the book or movie to add")

	if len(os.Args) == 1 {
		printHelp(listCmd, addCmd)
		db.Close()
		return
	}

	switch os.Args[1] {
	case "list":
		listCmd.Parse(os.Args[2:])
		switch *listMediaType {
		case "book":
			bookList(*listMediaStatus, *listJson)
		case "movie":
			movieList(*listMediaStatus, *listJson)
		default:
			fmt.Println("Invalid media type to add: book or movie")
		}
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
		printHelp(listCmd, addCmd)
	}

	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
