package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	OkuReadingURL   string `json:"oku_reading_url"`
	OkuReadURL      string `json:"oku_read_url"`
	OkuToreadURL    string `json:"oku_toread_url"`
	DbFile          string `json:"db_file"`
	BookPublishDir  string `json:"book_publish_dir"`
	MoviePublishDir string `json:"movie_publish_dir"`
}

func loadConfig() (c Config, error error) {
	var config Config

	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return config, errors.New("Failed to find user home directory")
	}

	content, err := ioutil.ReadFile(fmt.Sprintf("%s/.cab", dirname))
	if err != nil {
		fmt.Println(err)
		return config, errors.New("Failed to open ~/.cab")
	}

	err = json.Unmarshal([]byte(content), &config)
	if err != nil {
		fmt.Println(err)
		return config, errors.New("Failed to parse JSON in ~/.cab")
	}

	return config, nil
}
