package meals

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const bistroUrl = "https://bistro.cgm.ag/index.php"

type Meal struct {
	name                string
	supplement          string
	price               float64
	optionalSupplements []OptionalSupplement
}

type OptionalSupplement struct {
	name  string
	price float64
}

func Start() map[string][]Meal {
	doc := requestWebsiteDocument(bistroUrl)

	parsedDates := parseDates(doc)
	mealDates := parseMealsForAllDays(doc, parsedDates)

	return mealDates
}

func parseMealsForAllDays(doc *goquery.Document, parsedDates []string) map[string][]Meal {
	mealDates := make(map[string][]Meal)

	doc.Find("div#day").Each(func(i int, daySelection *goquery.Selection) {
		mealDates[parsedDates[i]] = parseMealsForDay(daySelection)
	})

	return mealDates
}

func parseMealsForDay(daySelection *goquery.Selection) []Meal {
	var meals []Meal

	daySelection.Find("div#crawler").Each(func(i int, mealSelection *goquery.Selection) {
		meal := Meal{}
		meal.name = mealSelection.Find("p.menuName").Text()
		meal.supplement = mealSelection.Find("p.beschreibung").Text()
		meal.price = convertToPrice(mealSelection.Find("p.preis > b").Text())
		meal.optionalSupplements = parseOptionalSupplements(mealSelection)

		meals = append(meals, meal)
	})
	return meals
}

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

func parseOptionalSupplements(mealSelection *goquery.Selection) (optionalSupplements []OptionalSupplement) {
	mealSelection.Find("div[style='padding-left:10px']").Each(func(i int, supplementSelection *goquery.Selection) {
		optionalSupplement := OptionalSupplement{}

		nameString := strings.TrimSpace(supplementSelection.Find("div").Nodes[0].FirstChild.Data)
		priceString := supplementSelection.Find("div").Nodes[1].FirstChild.Data

		optionalSupplement.name = nameString
		optionalSupplement.price = convertToPrice(priceString)

		optionalSupplements = append(optionalSupplements, optionalSupplement)
	})

	return optionalSupplements
}

func convertToPrice(priceString string) float64 {
	priceString = strings.Replace(priceString, ",", ".", -1)
	priceString = strings.Replace(priceString, "€", "", -1)
	priceString = strings.TrimSpace(priceString)
	price, _ := strconv.ParseFloat(priceString, 64)
	return price
}

func requestWebsiteDocument(url string) *goquery.Document {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	res.Body.Close()

	return doc
}
