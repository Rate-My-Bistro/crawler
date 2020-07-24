package persister

import (
	"context"
	"github.com/ansgarS/rate-my-bistro-crawler/crawler"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"log"
)

const databaseName = "rate-my-bistro"
const collectionName = "meals"

var client driver.Client
var database driver.Database
var collection driver.Collection

func Start(databaseAddress string, meals []crawler.Meal) {
	createClient(databaseAddress)
	createDatabase(databaseName)
	createCollection(collectionName)

	for _, meal := range meals {
		createOrUpdateMeal(meal)
	}
}

func createOrUpdateMeal(meal crawler.Meal) {
	if checkIfMealExists(meal.Key) {
		updateMeal(meal)
	} else {
		createMeal(meal)
	}
}

func removeMeal(mealKey string) {
	if checkIfMealExists(mealKey) {
		_, err := collection.RemoveDocument(context.Background(), mealKey)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func checkIfMealExists(mealKey string) bool {
	exists, _ := collection.DocumentExists(context.Background(), mealKey)
	return exists
}

func updateMeal(meal crawler.Meal) {
	ctx := context.Background()
	_, err := collection.UpdateDocument(ctx, meal.Key, meal)
	if err != nil {
		log.Fatal(err)
	}
}

func createMeal(meal crawler.Meal) {
	_, err := collection.CreateDocument(context.Background(), meal)

	if err != nil {
		log.Fatal(err)
	}
}

func getMeal(mealKey string) (meal crawler.Meal) {
	_, err := collection.ReadDocument(context.Background(), mealKey, &meal)

	if err != nil {
		log.Fatal(err)
	}

	return meal
}

func createDatabase(dbName string) {
	exists, _ := client.DatabaseExists(context.Background(), dbName)
	if exists {
		db, _ := client.Database(context.Background(), dbName)
		database = db
	} else {
		options := &driver.CreateDatabaseOptions{}
		db, err := client.CreateDatabase(context.Background(), dbName, options)

		if err != nil {
			log.Fatal(err)
		}

		database = db
	}
}

func createCollection(colName string) {
	exists, _ := database.CollectionExists(context.Background(), colName)
	if exists {
		coll, _ := database.Collection(context.Background(), colName)
		collection = coll
	} else {
		options := &driver.CreateCollectionOptions{}
		coll, err := database.CreateCollection(context.Background(), colName, options)

		if err != nil {
			log.Fatal(err)
		}

		collection = coll
	}
}

func createClient(address string) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{address},
	})
	if err != nil {
		log.Fatal(err)
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})

	client = c

	if err != nil {
		log.Fatal(err)
	}
}
