/*
Package persister implements a simple crud functionality for meals.

The crawler is able to analyze any data source as long as it complies with the 'io.Reader' contract.
*/
package persister

import (
	"context"
	"github.com/ansgarS/rate-my-bistro-crawler/crawler"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"log"
)

var client driver.Client
var database driver.Database
var collection driver.Collection

// persists the passed meals into the database
// the parameter databaseAddress defines the database target
func PersistMeals(databaseAddress string,
	databaseName string,
	collectionName string,
	meals []crawler.Meal) {
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
	trxId, transactionContext := startTransaction()

	if checkIfMealExists(meal.Key, transactionContext) {
		updateMeal(meal, transactionContext)
	} else {
		createMeal(meal, transactionContext)
	}

	if err := database.CommitTransaction(transactionContext, trxId, nil); err != nil {
		log.Fatalf("Failed to commit transaction: %s", err)
	}
}

// initiate a new database transactions
// returns the transaction id and the transaction context
func startTransaction() (driver.TransactionID, context.Context) {
	bgContext := context.Background()
	trxId, err := database.BeginTransaction(bgContext, driver.TransactionCollections{Exclusive: []string{collection.Name()}}, nil)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %s", err)
	}
	transactionContext := driver.WithTransactionID(bgContext, trxId)
	return trxId, transactionContext
}

// Removes a meal by its identification key
// WARNING! Don't use this function from productive code.
func removeMeal(mealKey string) {
	if checkIfMealExists(mealKey, nil) {
		_, err := collection.RemoveDocument(context.Background(), mealKey)

		if err != nil {
			log.Fatal(err)
		}
	}
}

// Checks if a meal document exists by its key
func checkIfMealExists(mealKey string, ctx context.Context) bool {

	if ctx == nil {
		ctx = context.Background()
	}

	exists, _ := collection.DocumentExists(ctx, mealKey)
	return exists
}

// Updates an existing meal document
// If it does not exists this function will fail
func updateMeal(meal crawler.Meal, ctx context.Context) {

	if ctx == nil {
		ctx = context.Background()
	}

	_, err := collection.UpdateDocument(ctx, meal.Key, meal)
	if err != nil {
		log.Fatal(err)
	}
}

// creates a new meal document
// if a document with the same key already exists this function will fail
func createMeal(meal crawler.Meal, ctx context.Context) {

	if ctx == nil {
		ctx = context.Background()
	}

	_, err := collection.CreateDocument(ctx, meal)

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
