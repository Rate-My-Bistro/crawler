/*
Package persister implements a simple crud functionality for documents.

The crawler is able to analyze any data source as long as it complies with the 'io.Reader' contract.
*/
package persister

import (
	"context"
	"fmt"
	"github.com/ansgarS/rate-my-bistro-crawler/config"
	"github.com/ansgarS/rate-my-bistro-crawler/crawler"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"log"
)

var client driver.Client
var database driver.Database
var collection driver.Collection

type Identifiable interface {
	GetId() string
}

func init() {
	createClient(config.Cfg.DatabaseAddress)
	createDatabase(config.Cfg.DatabaseName)
}

// persists the passed documents into the database
// the parameter databaseAddress defines the database target
func PersistDocuments(collectionName string, documents []Identifiable) {
	ensureCollection(collectionName)

	for _, document := range documents {
		createOrUpdateDocument(document)
	}
}

// persists the passed document into the database
// the parameter databaseAddress defines the database target
func PersistDocument(collectionName string, document Identifiable) {
	ensureCollection(collectionName)

	createOrUpdateDocument(document)
}

// Creates a new document document if it does not exists yet
// Otherwise it will updated, identified by the key
func createOrUpdateDocument(document Identifiable) {
	trxId, transactionContext := startTransaction()

	if checkIfDocumentExists(document.GetId(), transactionContext) {
		updateDocument(document, transactionContext)
	} else {
		createDocument(document, transactionContext)
	}

	if err := database.CommitTransaction(transactionContext, trxId, nil); err != nil {
		log.Fatalf("Failed to commit transaction for document %s: %s", document.GetId(), err)
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

// Removes a document by its identification key
// WARNING! Don't use this function from productive code.
func removeDocument(key string) {
	if checkIfDocumentExists(key, nil) {
		_, err := collection.RemoveDocument(context.Background(), key)

		if err != nil {
			log.Fatal(err)
		}
	}
}

// Checks if a document document exists by its key
func checkIfDocumentExists(mealKey string, ctx context.Context) bool {
	if ctx == nil {
		ctx = context.Background()
	}

	exists, _ := collection.DocumentExists(ctx, mealKey)
	return exists
}

// Updates an existing document document
// If it does not exists this function will fail
func updateDocument(document Identifiable, ctx context.Context) {

	if ctx == nil {
		ctx = context.Background()
	}

	_, err := collection.UpdateDocument(ctx, document.GetId(), document)
	if err != nil {
		log.Fatal(err)
	}
}

// creates a new document document
// if a document with the same key already exists this function will fail
func createDocument(document Identifiable, ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	_, err := collection.CreateDocument(ctx, document)

	if err != nil {
		log.Fatal(err)
	}
}

// Retrieve a document by its key
// If no document exists with given key, a NotFoundError is thrown.
func ReadDocument(mealKey string, result interface{}) {
	_, err := collection.ReadDocument(context.Background(), mealKey, result)

	if err != nil {
		log.Fatal(err)
	}
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
func ensureCollection(colName string) {
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

// Prints all meals that were found in the database for the specified date
func PrintMealsForDate(date string) {
	ensureCollection(config.Cfg.MealCollectionName)

	ctx := context.Background()
	query := "FOR d IN meals FILTER d.date == @date RETURN d"
	bindVars := map[string]interface{}{
		"date": date,
	}
	cursor, err := database.Query(ctx, query, bindVars)
	if err != nil {
		// handle error
	}
	_ = cursor.Close()
	for {
		var doc crawler.Meal
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			// handle other errors
		}
		fmt.Printf("Got doc with key '%s' from query\n", doc.Name)
	}
}
