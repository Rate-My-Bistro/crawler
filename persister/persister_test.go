package persister

import (
	"context"
	"github.com/ansgarS/rate-my-bistro-crawler/config"
	"github.com/ansgarS/rate-my-bistro-crawler/crawler"
	"testing"
)

func TestInsertOrUpdate(t *testing.T) {
	//setup
	createClient()
	createDatabase()
	ensureCollection(config.Get().MealCollectionName)
	ensureCollection(config.Get().JobCollectionName)

	meal1Stub := crawler.Meal{
		Id: "abc",
	}

	meal1 := crawler.Meal{
		Id:    "abc",
		Date:  "2020-07-24",
		Name:  "Suppe",
		Price: 3.97,
		MandatorySupplements: []crawler.Supplement{
			{Name: "Reis", Price: 0}},
		OptionalSupplements: []crawler.Supplement{
			{Name: "Markklößchen", Price: 0.12},
			{Name: "Trokenes Brot", Price: 9.87}},
	}

	meal2 := crawler.Meal{
		Id:                   "b",
		Date:                 "2020-07-24",
		Name:                 "Reis",
		Price:                1.44,
		MandatorySupplements: []crawler.Supplement{{Name: "Salz", Price: 0}},
		OptionalSupplements:  []crawler.Supplement{{Name: "Chilli", Price: 1}},
	}

	t.Run("insert a record and update it", func(t *testing.T) {

		createOrUpdateDocument(config.Get().MealCollectionName, Identifiable(meal1Stub))

		if !DocumentExists(config.Get().MealCollectionName, meal1Stub.Id, nil) {
			t.Errorf("meal could not created")
		}

		createOrUpdateDocument(config.Get().MealCollectionName, meal1)

		var meal crawler.Meal
		ReadDocument(config.Get().MealCollectionName, meal1Stub.Id, context.Background(), &meal)

		if !(meal.Name == "Suppe") {
			t.Errorf("meal was not updated")
		}

		if !(meal.MandatorySupplements[0].Price == 0) {
			t.Errorf("mandadory supplement price shoud be 0")
		}

		if !(meal.OptionalSupplements[1].Price == 9.87) {
			t.Errorf("optional supplement price was not updated")
		}

		createOrUpdateDocument(config.Get().MealCollectionName, meal2)

		if !DocumentExists(config.Get().MealCollectionName, meal2.Id, nil) {
			t.Errorf("meal could not created")
		}
	})

	removeDocument(config.Get().MealCollectionName, meal1.Id)
	removeDocument(config.Get().MealCollectionName, meal2.Id)
}
