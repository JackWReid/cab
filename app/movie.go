package main

import (
	"errors"
	"fmt"
)

type movieRecord struct {
	Id    string
	Title string
	Year  string
	Slug  string
}

type movieEvent struct {
	Date  string `json:"date_updated"`
	Title string `json:"title"`
	Year  string `json:"year"`
}

type movieSearchResults struct {
	Count    int
	MovieIds []string
	Movies   []movieRecord
}

func getMoviesByStatus(movieStatus string) ([]movieEvent, error) {
	var qStr string
	var movieResults []movieEvent

	switch movieStatus {
	case "done":
		qStr = "SELECT title, year, date FROM view_movie_watched ORDER BY date DESC"
	case "todo":
		qStr = "SELECT title, year, date FROM view_movie_towatch ORDER BY date DESC"
	default:
		return movieResults, errors.New("Invalid movieStatus")
	}

	rows, err := DB.Query(qStr)

	for rows.Next() {
		var row movieEvent
		err = rows.Scan(&row.Title, &row.Year, &row.Date)
		checkErr(err)
		movieResults = append(movieResults, row)
	}

	rows.Close()

	if err != nil {
		return movieResults, err
	}

	return movieResults, nil
}

func searchMovies(query string) (movieSearchResults, error) {
	var movieResults []movieRecord
	movieQuery := `SELECT id, title, year, letterboxd_uri FROM movie WHERE title LIKE '%' || $1 || '%'`
	rows, err := DB.Query(movieQuery, query)
	checkErr(err)

	var movieIds []string
	for rows.Next() {
		var row movieRecord
		err = rows.Scan(&row.Id, &row.Title, &row.Year, &row.Slug)
		checkErr(err)
		movieIds = append(movieIds, row.Id)
		movieResults = append(movieResults, row)
	}

	rows.Close()

	return movieSearchResults{
		Count:    len(movieResults),
		MovieIds: movieIds,
		Movies:   movieResults,
	}, nil
}

func insertMovie(lb letterboxdMovie) (movieRecord, error) {
	var insertedMovie movieRecord
	insQ := `INSERT INTO movie(title, year, letterboxd_uri) VALUES(?, ?, ?)`
	result, err := DB.Exec(insQ, lb.Title, lb.Year, lb.Slug)

	if err != nil {
		return insertedMovie, err
	}

	id, _ := result.LastInsertId()

	reQ := `SELECT id, title, year FROM movie WHERE id = $1`
	row := DB.QueryRow(reQ, id)

	row.Scan(&insertedMovie.Id, &insertedMovie.Title, &insertedMovie.Slug)

	return insertedMovie, nil
}

func watchedMovie(movie movieRecord) error {
	cheQ := `SELECT COUNT(*) FROM movie_log WHERE movie_id = $1`
	row := DB.QueryRow(cheQ, movie.Id)
	var cheCount int
	row.Scan(&cheCount)

	if cheCount > 0 {
		fmt.Println(movie.Title, "already marked as watched")
		return nil
	} else {
		insQ := `INSERT INTO movie_log(movie_id) VALUES(?)`
		_, err := DB.Exec(insQ, movie.Id)

		if err != nil {
			return err
		}

		delQ := `DELETE FROM movie_to_movie_collection WHERE movie_id = $1 AND collecton_id = 1`
		_, err = DB.Exec(delQ, movie.Id)

		if err != nil {
			return err
		}

		return nil
	}
}

func towatchMovie(movie movieRecord) error {
	cheQ := `SELECT COUNT(*) FROM movie_to_movie_collection WHERE movie_id = $1 AND collection_id = 1`
	row := DB.QueryRow(cheQ, movie.Id)
	var cheCount int
	row.Scan(&cheCount)

	if cheCount > 0 {
		fmt.Println(movie.Title, "already marked as towatch")
		return nil
	} else {
		insQ := `INSERT INTO movie_to_movie_collection(movie_id, collection_id) VALUES(?, 1)`
		_, err := DB.Exec(insQ, movie.Id)

		if err != nil {
			return err
		}

		return nil
	}
}
