package persister

import (
	"github.com/ansgarS/rate-my-bistro-crawler/crawler"
	"testing"
)

const testDatabaseName = databaseName + "_test"
const testCollectionName = collectionName + "_test"

func TestInsertOrUpdate(t *testing.T) {
	//setup
	createClient("http://localhost:8529")
	createDatabase(testDatabaseName)
	createCollection(testCollectionName)

	meal1Stub := crawler.Meal{
		Key: "abc",
	}

	meal1 := crawler.Meal{
		Key:        "a",
		Date:       "2020-07-24",
		Name:       "Suppe",
		Supplement: "Brötchen",
		Price:      3.97,
		OptionalSupplements: []crawler.Supplement{
			{Name: "Markklößchen", Price: 0.12},
			{Name: "Trokenes Brot", Price: 9.87}},
	}

	meal2 := crawler.Meal{
		Key:        "b",
		Date:       "2020-07-24",
		Name:       "Reis",
		Supplement: "Salz",
		Price:      1.44,
	}

	t.Run("insert a record and update it", func(t *testing.T) {

		createOrUpdateMeal(meal1Stub)

		if !checkIfMealExists(meal1Stub) {
			t.Errorf("meal could not created")
		}

		createOrUpdateMeal(meal1)

		if !(getMeal(meal1Stub.Key).Name == "Suppe") {
			t.Errorf("meal was not updated")
		}

		if !(getMeal(meal1Stub.Key).OptionalSupplements[1].Price == 9.87) {
			t.Errorf("meal optional supplement price was not updated")
		}

		createOrUpdateMeal(meal2)

		if !checkIfMealExists(meal2) {
			t.Errorf("meal could not created")
		}
	})

	removeMeal(meal1.Key)
	removeMeal(meal2.Key)
}
