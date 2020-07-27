package persister

import (
	"github.com/ansgarS/rate-my-bistro-crawler/crawler"
	"testing"
)

const testDatabaseName = "rate-by-bistro_test"
const testCollectionName = "meals_test"

func TestInsertOrUpdate(t *testing.T) {
	//setup
	createClient("http://localhost:8529")
	createDatabase(testDatabaseName)
	createCollection(testCollectionName)

	meal1Stub := crawler.Meal{
		Key: "abc",
	}

	meal1 := crawler.Meal{
		Key:   "abc",
		Date:  "2020-07-24",
		Name:  "Suppe",
		Price: 3.97,
		Supplements: []crawler.Supplement{
			{Name: "Brötchen", Price: 0},
			{Name: "Markklößchen", Price: 0.12},
			{Name: "Trokenes Brot", Price: 9.87}},
	}

	meal2 := crawler.Meal{
		Key:         "b",
		Date:        "2020-07-24",
		Name:        "Reis",
		Price:       1.44,
		Supplements: []crawler.Supplement{{Name: "Salz", Price: 0}},
	}

	t.Run("insert a record and update it", func(t *testing.T) {

		createOrUpdateMeal(meal1Stub)

		if !checkIfMealExists(meal1Stub.Key, nil) {
			t.Errorf("meal could not created")
		}

		createOrUpdateMeal(meal1)

		if !(getMeal(meal1Stub.Key).Name == "Suppe") {
			t.Errorf("meal was not updated")
		}

		if !(getMeal(meal1Stub.Key).Supplements[2].Price == 9.87) {
			t.Errorf("meal optional supplement price was not updated")
		}

		createOrUpdateMeal(meal2)

		if !checkIfMealExists(meal2.Key, nil) {
			t.Errorf("meal could not created")
		}
	})

	removeMeal(meal1.Key)
	removeMeal(meal2.Key)
}
