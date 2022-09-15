/*
TESTS FOR UTILS MODULE USED FOR SCRAPING
MORE INFO IN MAKEFILE
*/

package tests

import (
	"net/http"
	"testing"
)

////////////////////////////////
// Adlibris ERROR CHECKING TESTS
////////////////////////////////

func TestAdlibrisGetHTML1(t *testing.T) {
	got, err := http.Get("https://www.adlibris.com/se/bok/skapa-med-mumin-9789177836971")

	if got.StatusCode != http.StatusOK || err != nil {
		t.Errorf("TEST FAILED: Got %d status error: %d!", got.StatusCode, err)
	}
}

func TestAdlibrisGetHTML2(t *testing.T) {
	got, err := http.Get("https://www.adlibris.com/se/produkt/arabia-mumin-mugg-hattifnattarna-30-cl-orange-32227929?article=P32227929")

	if got.StatusCode != http.StatusOK || err != nil {
		t.Errorf("TEST FAILED: Got %d status error: %d!", got.StatusCode, err)
	}
}

/////////////////////////////////
// Bokborsen ERROR CHECKING TESTS
/////////////////////////////////

func TestBokborsenGetHTML1(t *testing.T) {
	got, err := http.Get("https://www.bokborsen.se/view/Rey-Margret/Nicke-Nyfiken-I-Leksaksaff%C3%A4ren/11308018")

	if got.StatusCode != http.StatusOK || err != nil {
		t.Errorf("TEST FAILED: Got %d status error: %d!", got.StatusCode, err)
	}
}

func TestBokborsenGetHTML2(t *testing.T) {
	got, err := http.Get("https://www.bokborsen.se/view/H-A-Rey/Nicke-Nyfiken-P%C3%A5-")

	if got.StatusCode != http.StatusOK || err != nil {
		t.Errorf("TEST FAILED: Got %d status error: %d!", got.StatusCode, err)
	}
}

////////////////////////////
// Ebay ERROR CHECKING TESTS
////////////////////////////

func TestEbayGetHTML1(t *testing.T) {
	got, err := http.Get("https://www.ebay.com/itm/325145873628")

	if got.StatusCode != http.StatusOK || err != nil {
		t.Errorf("TEST FAILED: Got %d status error: %d!", got.StatusCode, err)
	}
}

func TestEbayGetHTML2(t *testing.T) {
	got, err := http.Get("https://www.ebay.com/itm/265440947686")

	if got.StatusCode != http.StatusOK || err != nil {
		t.Errorf("TEST FAILED: Got %d status error: %d!", got.StatusCode, err)
	}
}
