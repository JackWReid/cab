package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"
)

func addBook(db *sql.DB, bookTitle string) (bookRecord, error) {
	var exBookResults []bookRecord
	var addedBook bookRecord
	var bookAdded bool

	if len(bookTitle) == 0 {
		return addedBook, errors.New("No title given to add book")
	}

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
	} else {
		fmt.Println("No books found, searching Google Books.")
	}

	if !bookAdded {
		gbResults := searchGoogleBooks(bookTitle)
		tableGoogleResults(gbResults)

		var gbBookRes string
		fmt.Println("Enter # of book from Google Books:")
		fmt.Scanln(&gbBookRes)

		idx, err := strconv.Atoi(gbBookRes)

		if err != nil || idx > len(gbResults) {
			err = errors.New("Invalid Google Books ID")
			return addedBook, err
		}

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
	var logRes string
	var logErr error
	fmt.Println("How do you want to log", book.Title, "?")
	fmt.Println("[1] Read\n[2] Reading\n[3] To read")
	fmt.Print("> ")
	fmt.Scan(&logRes)

	switch logRes {
	case "1":
		logErr = readBook(book)
	case "2":
		logErr = readingBook(book)
	case "3":
		logErr = toreadBook(book)
	default:
		panic("Invalid log choice")
	}

	if logErr != nil {
		fmt.Println("Error logging book")
		fmt.Println(logErr)
	}
}

func addMovie(db *sql.DB, movieTitle string) (movieRecord, error) {
	var exMovieResults []movieRecord
	var addedMovie movieRecord
	var movieAdded bool

	if len(movieTitle) == 0 {
		return addedMovie, errors.New("No title given to add book")
	}

	searchRes, err := searchMovies(movieTitle)

	if err != nil {
		fmt.Println("Failed to search movies for", movieTitle)
		return addedMovie, err
	}

	if searchRes.Count > 0 {
		exMovieResults = searchRes.Movies
		tableMovieRecord(exMovieResults)
		var exMovieRes string
		fmt.Println("Enter existing movie ID, or [x] to search Letterboxd:")
		fmt.Print("> ")
		fmt.Scan(&exMovieRes)

		if slices.Contains(searchRes.MovieIds, exMovieRes) {
			f := func(b movieRecord) bool {
				return exMovieRes == b.Id
			}
			selectedMovieIdx := slices.IndexFunc(exMovieResults, f)
			addedMovie = exMovieResults[selectedMovieIdx]
			movieAdded = true
		}
	} else {
		fmt.Println("No movies found, searching Letterboxd.")
	}

	if !movieAdded {
		lbResults := searchLb(movieTitle)
		tableLbResults(lbResults)

		var lbMovieRes string
		fmt.Println("Enter # of movie from Letterboxd:")
		fmt.Print("> ")
		fmt.Scan(&lbMovieRes)

		idx, err := strconv.Atoi(lbMovieRes)

		if err != nil || idx > len(lbResults) {
			err = errors.New("Invalid Letterboxd ID")
			return addedMovie, err
		}

		sm := lbResults[idx]

		addedMovie, err = insertMovie(sm)

		if err != nil {
			fmt.Println("Failed to add movie", sm.Title)
			return addedMovie, err
		}

		movieAdded = true
	}

	logMovie(addedMovie)
	return addedMovie, nil
}

func logMovie(movie movieRecord) {
	var logRes string
	var logErr error
	fmt.Println("How do you want to log", movie.Title, "?")
	fmt.Println("[1] Watched\n[2] To watch")
	fmt.Print("> ")
	fmt.Scan(&logRes)

	switch logRes {
	case "1":
		logErr = watchedMovie(movie)
	case "2":
		logErr = towatchMovie(movie)
	default:
		panic("Invalid log choice")
	}

	if logErr != nil {
		fmt.Println("Error logging movie")
		fmt.Println(logErr)
	}
}
