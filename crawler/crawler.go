/*
Package crawler implements a simple parser that reads the cgm bistro website
and transforms them into manageable structs.
*/
package crawler

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Represents a meal
type Meal struct {
	// The Id is the identifier of each meal
	// it is composed like this: sha1( date + name )
	Id                   string       `json:"_key,omitempty"`
	Date                 string       `json:"date"`
	Name                 string       `json:"name"`
	Price                float64      `json:"price"`
	LowKcal              bool         `json:"lowKcal"`
	MandatorySupplements []Supplement `json:"mandatorySupplements"`
	OptionalSupplements  []Supplement `json:"optionalSupplements"`
}

//  Represents a supplement of an meal
type Supplement struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Crawls the content of the cgm bistro website for the current week
// returns a slice of meals
func CrawlCurrentWeek(bistroLocation string) []Meal {
	reader := createBistroReader(bistroLocation)

	doc := requestWebsiteDocument(reader)

	parsedDates := parseDates(doc)
	mealDates := parseMealsForAllDays(doc, parsedDates)

	return mealDates
}

// Receives a reader that provides the content of a bistro website for the specified date
// The date must have the format 'yyyy-mm-dd' example: '2020-12-31'
// returns a slice of meals for the week
func CrawlAtDate(bistroLocation string, date string) []Meal {
	if !strings.HasPrefix(bistroLocation, "http") {
		log.Fatal("Specific dates cannot parsed from an offline location, only urls are allowed.")
	}

	bistroLocation = buildDatedBistroLocation(bistroLocation, date)

	reader := createBistroReader(bistroLocation)

	doc := requestWebsiteDocument(reader)

	parsedDates := parseDates(doc)
	mealDates := parseMealsForAllDays(doc, parsedDates)

	return mealDates
}

func buildDatedBistroLocation(location string, date string) string {
	split := strings.Split(date, "-")
	year := split[0]
	month := split[1]
	day := split[2]

	if !strings.HasSuffix(location, "index.php") {
		location = strings.Replace(location+"/index.php", "//", "/", -1)
		location = strings.Replace(location, ":/", "://", 1)
	}

	location = location + fmt.Sprintf("?day=%s&month=%s&year=%s", day, month, year)

	return location
}

// creates an reader object based on the provided bistroUrl
func createBistroReader(bistroUrl string) (documentReader io.Reader) {
	if strings.HasPrefix(bistroUrl, "file://") {
		bistroUrl := strings.Replace(bistroUrl, "file://", "", -1)
		documentReader = readFile(bistroUrl)
	} else if strings.HasPrefix(bistroUrl, "/") {
		documentReader = readFile(bistroUrl)
	} else {
		documentReader = readUrl(bistroUrl).Body
	}

	return documentReader
}

// retrieves a http response from the specified url
func readUrl(url string) *http.Response {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return res
}

// retrieves a file handle from the specified file path
func readFile(filePath string) *os.File {
	bistroPageReader, err := os.Open(filePath)

	if err != nil {
		log.Fatal("Opening the following file failed: "+
			filePath, err)
	}

	return bistroPageReader
}

// Parses all meals found in the provided html document
// Receives a document that holds the bistro website
// Returns a map of meals
func parseMealsForAllDays(doc *goquery.Document, parsedDates []string) []Meal {
	meals := make([]Meal, 0)

	doc.Find("div#day").Each(func(i int, daySelection *goquery.Selection) {
		date := parsedDates[i]
		parsedMeals := parseMealsForDay(daySelection)

		for _, meal := range parsedMeals {
			meal.Date = date
			meal.Id = toSha1(meal.Date + meal.Name)
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
		meal.Price = convertToPrice(mealSelection.Find("p.preis > b").Text())
		meal.LowKcal = containsAttributeValue(mealSelection, "style", "background-color:greenyellow")
		meal.MandatorySupplements = parseMandatorySupplements(mealSelection)
		meal.OptionalSupplements = parseOptionalSupplements(mealSelection)

		// filters out days without meals (e.g. holidays)
		if meal.Price > 0 {
			meals = append(meals, meal)
		}
	})
	return meals
}

func parseMandatorySupplements(mealSelection *goquery.Selection) (mandatorySupplements []Supplement) {
	mandatorySupplementName := strings.TrimSpace(mealSelection.Find("p.beschreibung").Text())
	if mandatorySupplementName != "" {
		mandatorySupplements = append(mandatorySupplements, Supplement{
			Name:  mandatorySupplementName,
			Price: 0,
		})
	}
	return mandatorySupplements
}

// Checks if a html selection tag contains the specified attribute value
// Returns true if this is the case, otherwise false
func containsAttributeValue(selection *goquery.Selection, attributeName string, attributeValue string) bool {
	attr, exists := selection.Attr(attributeName)
	if exists {
		return strings.Contains(attr, attributeValue)
	}

	return false
}

// Parses all dates of the week
// Receives a queryable html document
// Returns a set of string dates in the format: yyyy-mm-dd
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

// Parses supplements of a given meal selection
// Receives a queryable meal selection
// Returns a set of supplements or an empty slice if no supplements found for a meal
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

// Converts a price string that comes from the bistro page into a float
// Receives a price as string
// Returns a price as float
func convertToPrice(priceString string) float64 {
	priceString = strings.Replace(priceString, ",", ".", -1)
	priceString = strings.Replace(priceString, "€", "", -1)
	priceString = strings.TrimSpace(priceString)
	price, _ := strconv.ParseFloat(priceString, 64)
	return price
}

// Requests a website content from a given reader
// Receives a reader interface
// Returns queryable html document
func requestWebsiteDocument(reader io.Reader) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		log.Fatal(err)
	}

	return doc
}

// Generates a sha1 hash from the specified string 's'
// Receives a string 's'
// Returns a sha1 hash as string
func toSha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func (meal Meal) GetId() string {
	return meal.Id
}
