package crawler

import (
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Represents a meal
type Meal struct {
	date                string
	name                string
	supplement          string
	price               float64
	optionalSupplements []Supplement
}

//  Represents a supplement of an meal
type Supplement struct {
	name  string
	price float64
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
			meal.date = date
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
		meal.name = mealSelection.Find("p.menuName").Text()
		meal.supplement = mealSelection.Find("p.beschreibung").Text()
		meal.price = convertToPrice(mealSelection.Find("p.preis > b").Text())
		meal.optionalSupplements = parseOptionalSupplements(mealSelection)

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

		optionalSupplement.name = nameString
		optionalSupplement.price = convertToPrice(priceString)

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
