package tests

import (
	"github.com/PuerkitoBio/goquery"
	"gosecondhand/src/database"
	"gosecondhand/src/targets"
	"gosecondhand/src/utils"
	"log"
	"strings"
	"testing"
)

//////////////////////////////////////////////////////////////
// TEST HELPER FUNCTION FOR RETURNING LIST OF ITEMS SCRAPED //
//////////////////////////////////////////////////////////////

func HelperBokbörsen(searchString string) []utils.Item {
	createAllTables(searchString)
	bokbörsenURL := targets.GenerateBokborsenURL(searchString, 1)

	response := utils.GetHTML(bokbörsenURL)
	defer response.Body.Close()

	document, error := goquery.NewDocumentFromReader(response.Body)
	utils.CheckIfError(error)

	result := targets.ScrapeBokborsenPageData(document, searchString)
	return result
}

func HelperAdlibris(searchString string) []utils.Item {
	createAllTables(searchString)
	adlibrisURL := targets.GenerateAdlibrisURL(searchString)

	response := utils.GetHTML(adlibrisURL)
	defer response.Body.Close()

	document, error := goquery.NewDocumentFromReader(response.Body)
	utils.CheckIfError(error)

	result := targets.ScrapeAdlibrisPageData(document, searchString)
	return result
}

func HelperEbay(searchString string) []utils.Item {
	createAllTables(searchString)
	ebayURL := targets.GenerateEbayURL(searchString)

	response := utils.GetHTML(ebayURL)
	defer response.Body.Close()

	document, error := goquery.NewDocumentFromReader(response.Body)
	utils.CheckIfError(error)

	result := targets.ScrapeEbayPageData(document, searchString)
	return result
}

//////////////////////////////////////////////////////
// BOKBORSEN CHECK IF RESULT LIST OF ITEMS IF EMPTY //
//////////////////////////////////////////////////////

func TestBokbörsenEmpty1(t *testing.T) {
	result := HelperBokbörsen("strindberg")
	if len(result) == 0 {
		t.Errorf("BOKBORSEN EMPTY TEST FAILED!")
	}
}

// THIS TEST ENSURES THAT THE SEARCH STRING IS ODD
// AND DOESN'T RETURN ANY RESULTS
func TestBokbörsenEmpty2(t *testing.T) {
	result := HelperBokbörsen("jdowaijdoiajwdoiajwodijaowidj")
	if len(result) != 0 {
		t.Errorf("BOKBORSEN EMPTY TEST FAILED!")
	}
}

/////////////////////////////////////////////////////
// ADLIBRIS CHECK IF RESULT LIST OF ITEMS IF EMPTY //
/////////////////////////////////////////////////////

func TestAdlibrisEmpty1(t *testing.T) {
	result := HelperAdlibris("strindberg")
	if len(result) == 0 {
		t.Errorf("ADLIBRIS EMPTY TEST FAILED!")
	}
}

// THIS TEST ENSURES THAT THE SEARCH STRING IS ODD
// AND DOESN'T RETURN ANY RESULTS
func TestAdlibrisEmpty2(t *testing.T) {
	result := HelperAdlibris("jdowaijdoiajwdoiajwodijaowidj")
	if len(result) != 0 {
		t.Errorf("ADLIBRIS EMPTY TEST FAILED!")
	}
}

///////////////////////////////////////////////////
// Bokbörsen CHECK IF ITEM DATA IS VALID/CORRECT //
///////////////////////////////////////////////////

func TestBokbörsenTitle(t *testing.T) {
	result := HelperBokbörsen("Harry Potter")
	if !strings.Contains(result[0].Title, "Harry Potter") {
		t.Errorf("BOKBORSEN TITLE TEST FAILED!")
	}
}

func TestBokbörsenDescription(t *testing.T) {
	result := HelperBokbörsen("mumin")
	if !strings.Contains(result[0].Description, "skick") {
		t.Errorf("BOKBORSEN DESCRIPTION TEST FAILED!")
	}
}

func TestBokbörsenPrice(t *testing.T) {
	result := HelperBokbörsen("strindberg")
	if result[0].Price <= 0 {
		t.Errorf("BOKBORSEN PRICE TEST FAILED!")
	}
}

//////////////////////////////////////////////////
// Adlibris CHECK IF ITEM DATA IS VALID/CORRECT //
//////////////////////////////////////////////////

func TestAdlibrisTitle(t *testing.T) {
	result := HelperAdlibris("Harry Potter")
	if !strings.Contains(result[0].Title, "Harry Potter") {
		t.Errorf("ADLIBRIS TITLE TEST FAILED!")
	}
}

func TestAdlibrisDescription(t *testing.T) {
	result := HelperAdlibris("mumin")
	if !strings.Contains(result[0].Description, "bok") {
		t.Errorf("ADLIBRIS DESCRIPTION TEST FAILED!, Got: %q", result[0].Description)
	}
}

func TestAdlibrisPrice(t *testing.T) {
	result := HelperAdlibris("mumin")
	if result[0].Price <= 0 {
		t.Errorf("ADLIBRIS PRICE TEST FAILED!")
	}
}

func createAllTables(searchString string) {
	db := database.ConnectDB()
	defer db.Close()
	err := db.Ping()
	if err != nil {
		log.Fatal("Connection could not be verified with Ping(): ", err)
	}

	// Create all tables in database if they don't already exist
	stores := []string{"Adlibris_", "Biblio_", "Blocket_", "Bokbörsen_", "Citiboard_", "Etsy_", "FacebookMarket_", "Tradera_"}
	i := 0
	for i < 8 {
		table_name := stores[i] + strings.ReplaceAll(searchString, " ", "_")
		err = database.CreateTables(db, table_name)
		utils.CheckIfError(err)
		i++
	}
}
