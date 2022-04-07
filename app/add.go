package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"
)

func addBook(db *sql.DB, bookTitle string) (bookRecord, error) {
	if len(bookTitle) == 0 {
		panic("No title passed")
	}

	var exBookResults []bookRecord
	var addedBook bookRecord
	var bookAdded bool
	searchRes, err := searchBooks(bookTitle)

	if err != nil {
		fmt.Println("Failed to search books for", bookTitle)
		return addedBook, err
	}

	if searchRes.Count > 0 {
		exBookResults = searchRes.Books
		tableBookRecord(exBookResults)
		var exBookRes string
		fmt.Println("Enter existing book ID, or [x] to search Google Books:")
		fmt.Scanln(&exBookRes)

		if slices.Contains(searchRes.BookIds, exBookRes) {
			f := func(b bookRecord) bool {
				return exBookRes == b.Id
			}
			selectedBookIdx := slices.IndexFunc(exBookResults, f)
			addedBook = exBookResults[selectedBookIdx]
			bookAdded = true
		}
	}

	if !bookAdded {
		gbResults := searchGoogleBooks(bookTitle)
		tableGoogleResults(gbResults)

		var gbBookRes string
		fmt.Println("Enter # of book from Google Books:")
		fmt.Scanln(&gbBookRes)

		idx, _ := strconv.Atoi(gbBookRes)
		sb := gbResults[idx]

		addedBook, err := insertBook(sb)

		if err != nil {
			fmt.Println("Failed to add book", sb.title)
			return addedBook, err
		}

		bookAdded = true
	}

	logBook(addedBook)
	return addedBook, nil
}

func logBook(book bookRecord) {
	fmt.Println("How do you want to log %s?", book.Title)
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
