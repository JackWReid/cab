package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"
)

func addBook(db *sql.DB, bookTitle string) {
	var existingBooks []bookRecord
	fmt.Println("Check for existing books called", bookTitle)
	bookQuery := `SELECT id, title, author, isbn FROM book WHERE title LIKE '%' || $1 || '%'`
	rows, err := db.Query(bookQuery, bookTitle)
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
		gbResults := searchGoogleBooks(bookTitle)
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

func addMovie(db *sql.DB, movieTitle string) {
	var existingMovies []movieRecord
	fmt.Println("Check for existing movies called", movieTitle)
	movieQuery := `SELECT id, title, year, letterboxd_uri FROM movie WHERE title LIKE '%' || $1 || '%'`
	rows, err := db.Query(movieQuery, movieTitle)
	checkErr(err)

	var exMovieIds []string
	for rows.Next() {
		var row movieRecord
		err = rows.Scan(&row.Id, &row.Title, &row.Year, &row.Slug)
		checkErr(err)
		exMovieIds = append(exMovieIds, row.Id)
		existingMovies = append(existingMovies, row)
	}

	rows.Close()

	tableMovieRecord(existingMovies)
	var exMovieRes string
	fmt.Println("Enter existing movie ID or [x] to search Letterboxd:")
	fmt.Scanln(&exMovieRes)

	if slices.Contains(exMovieIds, exMovieRes) {
		fmt.Println("selected ex movie", exMovieRes)
		f := func(m movieRecord) bool {
			return exMovieRes == m.Id
		}
		selectedMovieIdx := slices.IndexFunc(existingMovies, f)
		selectedMovie := existingMovies[selectedMovieIdx]
		fmt.Println(selectedMovie)
	} else {
		lbResults := searchLb(movieTitle)
		tableLbResults(lbResults)

		var lbMovieRes string
		fmt.Println("Enter # of movie from Letterboxd:")
		fmt.Scanln(&lbMovieRes)

		fmt.Println("selected lb movie", lbMovieRes)
		idx, _ := strconv.Atoi(lbMovieRes)
		sm := lbResults[idx]

		insQ := `INSERT INTO movie(title, year, letterboxd_uri) VALUES(?, ?, ?)`
		_, err := db.Exec(insQ, sm.Title, sm.Year, sm.Slug)
		if err != nil {
			fmt.Println("db err", err)
		}
	}
}
