package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type googleBooksIndustryIdentifier struct {
	IdentifierType string `json:"type"`
	Identifier     string `json:"identifier"`
}

type googleBooksVolumeInfo struct {
	Title               string                          `json:"title"`
	Subtitle            string                          `json:"subtitle"`
	Authors             []string                        `json:"authors"`
	Description         string                          `json:"description"`
	IndustryIdentifiers []googleBooksIndustryIdentifier `json:"industryIdentifiers"`
}

type googleBooksResult struct {
	Id         string                `json:"id"`
	VolumeInfo googleBooksVolumeInfo `json:"volumeInfo"`
}

type googleSearchResult struct {
	Items []googleBooksResult `json:"items"`
}

type condensedGoogleResult struct {
	isbn        string
	title       string
	subtitle    string
	author      string
	description string
}

func selectIsbn(identifiers []googleBooksIndustryIdentifier) (isbn string) {
	for _, i := range identifiers {
		if i.IdentifierType == "ISBN_13" {
			isbn = i.Identifier
		}
	}

	return isbn
}

func searchGoogleBooks(query string) []condensedGoogleResult {
	url := "https://www.googleapis.com/books/v1/volumes"

	googleClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("q", query)
	q.Add("max_results", "20")
	req.URL.RawQuery = q.Encode()

	res, getErr := googleClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	result := googleSearchResult{}

	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	var condensedResults []condensedGoogleResult
	for _, i := range result.Items {
		if len(i.VolumeInfo.Authors) > 0 && len(i.VolumeInfo.IndustryIdentifiers) > 0 {
			condensedResult := condensedGoogleResult{
				isbn:        selectIsbn(i.VolumeInfo.IndustryIdentifiers),
				title:       i.VolumeInfo.Title,
				subtitle:    i.VolumeInfo.Subtitle,
				author:      i.VolumeInfo.Authors[0],
				description: i.VolumeInfo.Description,
			}

			if len(condensedResult.author) > 0 && len(condensedResult.isbn) > 0 {
				condensedResults = append(condensedResults, condensedResult)
			}
		}
	}

	return condensedResults
}
