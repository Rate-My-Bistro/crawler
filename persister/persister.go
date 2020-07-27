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

// persists the passed meals into the database
// the parameter databaseAddress defines the database target
func Start(databaseAddress string, meals []crawler.Meal) {
	createClient(databaseAddress)
	createDatabase(databaseName)
	createCollection(collectionName)

	for _, meal := range meals {
		createOrUpdateMeal(meal)
	}
}

// Creates a new meal document if it does not exists yet
// Otherwise it will updated, identified by the key
func createOrUpdateMeal(meal crawler.Meal) {
	if checkIfMealExists(meal.Key) {
		updateMeal(meal)
	} else {
		createMeal(meal)
	}
}

// Removes a meal by its identification key
func removeMeal(mealKey string) {
	if checkIfMealExists(mealKey) {
		_, err := collection.RemoveDocument(context.Background(), mealKey)

		if err != nil {
			log.Fatal(err)
		}
	}
}

// Checks if a meal document exists by its key
func checkIfMealExists(mealKey string) bool {
	exists, _ := collection.DocumentExists(context.Background(), mealKey)
	return exists
}

// Updates an existing meal document
// If it does not exists this function will fail
func updateMeal(meal crawler.Meal) {
	ctx := context.Background()
	_, err := collection.UpdateDocument(ctx, meal.Key, meal)
	if err != nil {
		log.Fatal(err)
	}
}

// creates a new meal document
// if a document with the same key already exists this function will fail
func createMeal(meal crawler.Meal) {
	_, err := collection.CreateDocument(context.Background(), meal)

	if err != nil {
		log.Fatal(err)
	}
}

// Retrieve a meal by its key
// If no meal exists with given key, a NotFoundError is thrown.
func getMeal(mealKey string) (meal crawler.Meal) {
	_, err := collection.ReadDocument(context.Background(), mealKey, &meal)

	if err != nil {
		log.Fatal(err)
	}

	return meal
}

// Creates the specified database if it does not yet exist.
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

// Creates the specified collection if it does not yet exist.
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

// Creates a new database connection client and keeps
// the instance as member variable alive
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
