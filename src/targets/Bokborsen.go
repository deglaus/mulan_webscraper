package targets

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gosecondhand/src/database"
	"gosecondhand/src/utils"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GenerateBokborsenURL(searchString string, pageNum int) string {
	searchString = strings.ReplaceAll(searchString, " ", "+")
	concatenatedString := ("https://www.bokborsen.se/?_p=" + strconv.Itoa(pageNum) + "&c=0&q=" + searchString)
	return concatenatedString
}

func ScrapeBokborsenPageData(document *goquery.Document, searchString string) []utils.Item {
	// ALLOCATE A NEW DATA ITEM OBJECT
	data := utils.Item{}

	// Variable for storing all items as a list
	var dataList []utils.Item

	// DATABASE CONNECTION
	db := database.ConnectDB()

	// Suspend execution of this function until surround has ran
	defer db.Close()

	// ERROR HANDLING FOR DATABSE CONNECTION
	err := db.Ping()
	if err != nil {
		log.Fatal("Connection could not be verified with Ping(): ", err)
	}

	table_name := "Bokbörsen_" + strings.ReplaceAll(searchString, " ", "_")

	document.Find("div.single-product").Each(func(index int, item *goquery.Selection) {

		// URL //
		foundUrl := item.Find("div.header>h2>a")
		url, _ := foundUrl.Attr("href")
		processedUrl := "https://www.bokborsen.se" + url

		// PICTUREURL //
		foundPicUrl := item.Find("div.product_main_image>img")
		picUrl, _ := foundPicUrl.Attr("src")
		processedPicUrl := strings.ReplaceAll(picUrl, "_thumb", "")

		// TITLE //
		foundTitle := item.Find("div.header>h2>a>span")
		title := foundTitle.Text()
		processedTitle := strings.TrimSpace(title)

		// DESCRIPTION //
		// Get all Text within the tag <p> for the item
		descriptionRaw := item.Find("div.content-primary>div.header>p").Text()
		// Removes all long whitespaces
		description := strings.Join(strings.Fields(descriptionRaw), " ")

		// PRICE //
		price_div := item.Find("span.price")
		price := strings.TrimSpace(price_div.Text())
		// Convert price to float instead of string, the " SEK" substring has to be removed for this
		price = strings.ReplaceAll(price, " SEK", "")
		// THIS IS THE VARIABLE THAT WE WANT IN DATABASE (actual float of price)
		convertedPrice, _ := strconv.ParseFloat(price, 32)

		// INSERTS APPROPRIATE DATA TO STRUCT FIELDS
		data.SearchString = searchString
		data.Site = "Bokborsen"
		data.Description = description
		data.URL = processedUrl
		data.PictureURL = processedPicUrl
		data.Title = processedTitle
		data.Price = convertedPrice

		// Append item to list of all items
		dataList = append(dataList, data)
		// INSERT DATA TO DATABASE
		database.InsertData(db, table_name, data)
	})
	return dataList
}

func singlePageResponse(searchString string, pageNum int) {
	bokborsenURL := GenerateBokborsenURL(searchString, pageNum)

	response := utils.GetHTML(bokborsenURL)
	defer response.Body.Close()

	document, error := goquery.NewDocumentFromReader(response.Body)
	utils.CheckIfError(error)

	ScrapeBokborsenPageData(document, searchString)
}

func getTotalPageNum(searchString string, ch chan int) {
	bokborsenURL := GenerateBokborsenURL(searchString, 0)
	response := utils.GetHTML(bokborsenURL)
	defer response.Body.Close()

	document, error := goquery.NewDocumentFromReader(response.Body)
	utils.CheckIfError(error)

	foundMetaData := document.Find("div.search-meta>b")
	foundPageNumInfo := foundMetaData.Text()
	parts := strings.Split(foundPageNumInfo, "av ")

	totalListings, error := strconv.Atoi(parts[1])
	totalPages := math.Round(float64(totalListings) / 30)
	fmt.Println(totalPages)
	ch <- int(totalPages)
}

func BokborsenResponse(searchString string, wg *sync.WaitGroup) {
	defer wg.Done() // WG FOR CONCURRENCY BETWEEN SCRAPERS

	var wgBokborsen sync.WaitGroup // WG FOR INTRA-CONCURRENCY
	wgBokborsen.Add(1)             // Incr

	ch := make(chan int) // Create a channel
	go func() {
		getTotalPageNum(searchString, ch)
		wgBokborsen.Done()
	}()

	// Immediately start scrape of first page since there will always be at least one
	wgBokborsen.Add(1)
	go func() {
		singlePageResponse(searchString, 1)
		wgBokborsen.Done()
	}()
	// Will block until receives from channel
	totalPageResult := <-ch

	// Scrape the rest

	for i := 2; i <= totalPageResult; i++ {
		wgBokborsen.Add(1)
		go func() {
			singlePageResponse(searchString, i)
			wgBokborsen.Done()
		}()
	}

	// Be nice to the site
	time.Sleep(time.Millisecond * 100)

	// Makes sure goroutines actually stop (and start)
	// The function of this function is not relevant

	wgBokborsen.Wait() // Blocks until wait group counter is 0
	fmt.Println("BOKBÖRSEN RESPONSE COMPLETE!")
}
