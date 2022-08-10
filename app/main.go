package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB
var GlobalConfig Config

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

func printHelp(listCmd *flag.FlagSet, importCmd *flag.FlagSet) {
	fmt.Println("===\ncab")
	fmt.Println("cab is a way to manage logs and lists for books and movies.\n===\n")
	fmt.Println("SUBCOMMANDS")
	fmt.Println("\nlist\nPrint a list of media that's been marked as done, doing, or to do")
	listCmd.PrintDefaults()
	fmt.Println("\nimport\nUse Oku and Letterboxd to import media catalog\n")
	importCmd.PrintDefaults()
	os.Exit(0)
}

func main() {
	config, configErr := loadConfig()

	if configErr != nil {
		fmt.Println("Failed to load config")
		fmt.Println(configErr)
		os.Exit(1)
	}

	GlobalConfig = config

	db, connErr := sql.Open("sqlite", GlobalConfig.DbFile)

	if connErr != nil {
		fmt.Println("Failed to open DB")
		fmt.Println(connErr)
		os.Exit(1)
	}

	DB = db

	pingErr := DB.Ping()

	if pingErr != nil {
		fmt.Println("Failed to ping DB")
		fmt.Println(pingErr)
		os.Exit(1)
	}

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listJson := listCmd.Bool("json", false, "Output lists in JSON, default is table")
	listMediaType := listCmd.String("type", "book", "Book or movie")
	listMediaStatus := listCmd.String("status", "done", "Done, doing, todo")

	importCmd := flag.NewFlagSet("import", flag.ExitOnError)
	importMediaType := importCmd.String("type", "book", "Book or movie")

	if len(os.Args) == 1 {
		printHelp(listCmd, importCmd)
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
	case "import":
		importCmd.Parse(os.Args[2:])
		switch *importMediaType {
		case "book":
			getOkuFetcher("reading")()
			getOkuFetcher("read")()
			getOkuFetcher("toread")()
		case "movie":
			fmt.Println("Nothing for now")
		default:
			fmt.Println("Invalid media type to import: book or movie")
		}

	default:
		printHelp(listCmd, importCmd)
	}

	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
