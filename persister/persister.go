/*
Package persister implements a simple crud functionality for documents.

The crawler is able to analyze any data source as long as it complies with the 'io.Reader' contract.
*/
package persister

import (
	"context"
	"github.com/ansgarS/rate-my-bistro-crawler/config"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"log"
)

var client driver.Client
var database driver.Database
var collections = make(map[string]driver.Collection)

type Identifiable interface {
	GetId() string
}

func init() {
	createClient()
	createDatabase()
	ensureCollection(config.Get().MealCollectionName)
	ensureCollection(config.Get().JobCollectionName)
}

// persists the passed documents into the database
// the parameter databaseAddress defines the database target
func PersistDocuments(collectionName string, documents []Identifiable) {
	for _, document := range documents {
		createOrUpdateDocument(collectionName, document)
	}
}

// persists the passed document into the database
// the parameter databaseAddress defines the database target
func PersistDocument(collectionName string, document Identifiable) {
	createOrUpdateDocument(collectionName, document)
}

// Creates a new document document if it does not exists yet
// Otherwise it will updated, identified by the key
func createOrUpdateDocument(collectionName string, document Identifiable) {
	trxId, transactionContext := startTransaction(collectionName)

	if DocumentExists(collectionName, document.GetId(), transactionContext) {
		updateDocument(collectionName, document, transactionContext)
	} else {
		createDocument(collectionName, document, transactionContext)
	}

	if err := database.CommitTransaction(transactionContext, trxId, nil); err != nil {
		log.Printf("Failed to commit transaction for document %s: %s", document.GetId(), err)
	}
}

// initiate a new database transactions
// returns the transaction id and the transaction context
func startTransaction(collectionName string) (driver.TransactionID, context.Context) {
	bgContext := context.Background()
	trxId, err := database.BeginTransaction(bgContext, driver.TransactionCollections{Exclusive: []string{collectionName}}, nil)
	if err != nil {
		log.Printf("Failed to begin transaction: %s", err)
	}
	transactionContext := driver.WithTransactionID(bgContext, trxId)
	return trxId, transactionContext
}

// Removes a document by its identification key
// WARNING! Don't use this function from productive code.
func removeDocument(collectionName string, key string) {
	if DocumentExists(collectionName, key, nil) {
		_, err := collections[collectionName].RemoveDocument(context.Background(), key)

		if err != nil {
			log.Print(err)
		}
	}
}

// Checks if a document document exists by its key
func DocumentExists(collectionName string, key string, ctx context.Context) bool {
	if ctx == nil {
		ctx = context.Background()
	}

	exists, err := collections[collectionName].DocumentExists(ctx, key)
	if err != nil {
		log.Print(err)
	}

	return exists
}

// Updates an existing document document
// If it does not exists this function will fail
func updateDocument(collectionName string, document Identifiable, ctx context.Context) {

	if ctx == nil {
		ctx = context.Background()
	}

	_, err := collections[collectionName].UpdateDocument(ctx, document.GetId(), document)
	if err != nil {
		log.Print(err)
	}
}

// creates a new document document
// if a document with the same key already exists this function will fail
func createDocument(collectionName string, document Identifiable, ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	_, err := collections[collectionName].CreateDocument(ctx, document)

	if err != nil {
		log.Print(err)
	}
}

// Retrieve a document by its key
// If no document exists with given key, a NotFoundError is thrown.
func ReadDocument(collectionName string, key string, ctx context.Context, result interface{}) {
	if ctx == nil {
		ctx = context.Background()
	}

	_, err := collections[collectionName].ReadDocument(ctx, key, result)

	if err != nil {
		log.Print(err)
	}
}

// Retrieve a document by its key
// If no document exists with given key, an empty document is returned
func ReadDocumentIfExists(collectionName string, key string, result interface{}) {
	if DocumentExists(collectionName, key, nil) {
		ReadDocument(collectionName, key, nil, result)
	}
}

// Creates the specified database if it does not yet exist.
func createDatabase() {
	dbName := config.Get().DatabaseName
	exists, _ := client.DatabaseExists(context.Background(), dbName)
	if exists {
		db, _ := client.Database(context.Background(), dbName)
		database = db
	} else {
		options := &driver.CreateDatabaseOptions{}
		db, err := client.CreateDatabase(context.Background(), dbName, options)

		if err != nil {
			log.Print(err)
		}

		database = db
	}
}

// Creates the specified collection if it does not yet exist.
func ensureCollection(collectionName string) {
	exists, _ := database.CollectionExists(context.Background(), collectionName)
	var collection driver.Collection
	if exists {
		coll, _ := database.Collection(context.Background(), collectionName)
		collection = coll
	} else {
		options := &driver.CreateCollectionOptions{}
		coll, err := database.CreateCollection(context.Background(), collectionName, options)

		if err != nil {
			log.Print(err)
		}

		collection = coll
	}
	collections[collectionName] = collection
}

// Creates a new database connection client and keeps
// the instance as member variable alive
func createClient() {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{config.Get().DatabaseAddress},
	})
	if err != nil {
		log.Print(err)
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(config.Get().DatabaseUser, config.Get().DatabasePassword),
	})

	client = c

	if err != nil {
		log.Print(err)
	}
}
