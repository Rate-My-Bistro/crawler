package crawler

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

// Represents a meal
type Meal struct {
	// The Key is the identifier of each meal
	// it is composed like this: sha1( date + name )
	Key                 string       `json:"_key,omitempty"`
	Date                string       `json:"date"`
	Name                string       `json:"name"`
	Supplement          string       `json:"supplement"`
	Price               float64      `json:"price"`
	OptionalSupplements []Supplement `json:"optionalSupplements"`
}

//  Represents a supplement of an meal
type Supplement struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// The entry point of this crawler
//
// receives a reader that provides the content of a bistro website
// returns a map of meals
func Start(documentReader io.Reader) []Meal {
	doc := requestWebsiteDocument(documentReader)

	parsedDates := parseDates(doc)
	mealDates := parseMealsForAllDays(doc, parsedDates)

	return mealDates
}

// parses all meals found in the provided html document
// receives a document that holds the bistro website
// returns a map of meals
func parseMealsForAllDays(doc *goquery.Document, parsedDates []string) []Meal {
	meals := make([]Meal, 0)

	doc.Find("div#day").Each(func(i int, daySelection *goquery.Selection) {
		date := parsedDates[i]
		parsedMeals := parseMealsForDay(daySelection)

		for _, meal := range parsedMeals {
			meal.Date = date
			meal.Key = toSha1(meal.Date + meal.Name)
			meals = append(meals, meal)
		}
	})

	return meals
}

// parses all meals for a given day
// receives a selector that holds the meal data of a single day
// returns a set of meals
func parseMealsForDay(daySelection *goquery.Selection) []Meal {
	var meals []Meal

	daySelection.Find("div#meal").Each(func(i int, mealSelection *goquery.Selection) {
		meal := Meal{}
		meal.Name = mealSelection.Find("p.menuName").Text()
		meal.Supplement = mealSelection.Find("p.beschreibung").Text()
		meal.Price = convertToPrice(mealSelection.Find("p.preis > b").Text())
		meal.OptionalSupplements = parseOptionalSupplements(mealSelection)

		meals = append(meals, meal)
	})
	return meals
}

// parses all dates of the week
// receives a queryable html document
// returns a set of string dates in the format: yyyy-mm-dd
func parseDates(doc *goquery.Document) []string {
	parsedDates := make([]string, 0)
	// Parse date for this day
	doc.Find("div.table-col-header b").Each(func(i int, dateSelection *goquery.Selection) {
		dateString := dateSelection.Nodes[0].LastChild.Data
		parsedDate, _ := time.Parse("2.1.2006", dateString)
		dateString = parsedDate.Format("2006-01-02")
		parsedDates = append(parsedDates, dateString)
	})
	return parsedDates
}

// parses supplements of a given meal selection
// receives a queryable meal selection
// returns a set of supplements or an empty slice if no supplements found for a meal
func parseOptionalSupplements(mealSelection *goquery.Selection) (optionalSupplements []Supplement) {
	mealSelection.Find("div[style='padding-left:10px']").Each(func(i int, supplementSelection *goquery.Selection) {
		optionalSupplement := Supplement{}

		nameString := strings.TrimSpace(supplementSelection.Find("div").Nodes[0].FirstChild.Data)
		priceString := supplementSelection.Find("div").Nodes[1].FirstChild.Data

		optionalSupplement.Name = nameString
		optionalSupplement.Price = convertToPrice(priceString)

		optionalSupplements = append(optionalSupplements, optionalSupplement)
	})

	return optionalSupplements
}

// converts a price string that comes from the bistro page into a float
// receives a price as string
// returns a price as float
func convertToPrice(priceString string) float64 {
	priceString = strings.Replace(priceString, ",", ".", -1)
	priceString = strings.Replace(priceString, "€", "", -1)
	priceString = strings.TrimSpace(priceString)
	price, _ := strconv.ParseFloat(priceString, 64)
	return price
}

// requests a website content from a given reader
// receives a reader interface
// returns queryable html document
func requestWebsiteDocument(reader io.Reader) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		log.Fatal(err)
	}

	return doc
}

// generates a sha1 hash from the specified string 's'
// receives a string 's'
// returns a sha1 hash as string
func toSha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
