package main

import (
	"fmt"
	"gosecondhand/src/database"
	"gosecondhand/src/targets"
	"gosecondhand/src/utils"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// Wait group variable
var wg sync.WaitGroup

// Search String variable for command-line input
var searchString string

func main() {
	// The first argument os.Args[0] is always the program name
	// This variable stores the search string argument when calling the scraper from the Makefile
	searchString := os.Args[1]
	// Scan for user search string

	// Measures how long it takes to scrape the enabled scrapers
	// Can be useful for comparing concurrency performance
	start := time.Now()
	searchString = strings.ReplaceAll(searchString, "\n", "")
	searchString = strings.ReplaceAll(searchString, "\r", "")

	db := database.ConnectDB()
	defer db.Close()
	err := db.Ping()
	if err != nil {
		log.Fatal("Connection could not be verified with Ping(): ", err)
	}

	isAlreadyInDB := database.IsInDatabase(db, searchString)
	if isAlreadyInDB {
		fmt.Println("Keyword was already in database; scraping not performed.")
		return
	}

	searchString = strings.ReplaceAll(searchString, "_", " ")

	// Create all tables in database if they don't already exist
	stores := []string{"Adlibris_", "Biblio_", "Blocket_", "Bokbörsen_", "Citiboard_", "Etsy_", "FacebookMarket_", "Tradera_"}
	amountOfSites := 8
	i := 0
	for i < amountOfSites {
		table_name := stores[i] + strings.ReplaceAll(searchString, " ", "_")
		err = database.CreateTables(db, table_name)
		utils.CheckIfError(err)
		i++
	}

	// Collecting performance data for executing go scrapers concurrently
	wg.Add(7)
	
	go targets.AdlibrisResponse(searchString, &wg)
	go targets.BiblioResponse(searchString, &wg)
	go targets.BokborsenResponse(searchString, &wg)
	go targets.CitiboardResponse(searchString, &wg)
	go targets.EtsyResponse(searchString, &wg)
	go targets.ExecuteBlocketandTradera(searchString, &wg)
	go targets.FacebookMarketResponse(searchString, "uppsala", &wg)

	wg.Wait()

	// Gets the time taken for scrapers to be done executing
	duration := time.Since(start)
	fmt.Println(duration)
}
