package utils

import (
	"fmt"
	"net/http"
)

type Item struct {
	// Metadata
	Id           int
	SearchString string
	Site         string // Look at this first to determine relevant data
	URL          string
	PictureURL   string

	Title       string
	Description string
	Price       float64
	Category    string
}

func CheckIfError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func GetHTML(url string) *http.Response {

	response, error := http.Get(url)
	CheckIfError(error)

	if response.StatusCode != 200 {
		fmt.Println("Status Code: ", response.StatusCode)
	}

	return response
}
