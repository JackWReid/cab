package main

import "errors"

type movieEvent struct {
	Date  string
	Title string
	Year  string
}

func getMoviesByStatus(movieStatus string) ([]movieEvent, error) {
	var qStr string
	var movieResults []movieEvent

	switch movieStatus {
	case "watched":
		qStr = "SELECT title, year, date FROM view_movie_watched ORDER BY date DESC"
	case "towatch":
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
