package persister

import (
	"context"
	"fmt"
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
	createDatabase()
	createCollection()

	for _, meal := range meals {
		createDocument(meal)
	}

	ctx := driver.WithQueryCount(context.Background())
	query := "FOR d IN meals RETURN d"
	cursor, err := database.Query(ctx, query, nil)
	if err != nil {
		// handle error
	}
	defer cursor.Close()
	fmt.Printf("Query yields %d documents\n", cursor.Count())
}

func createDocument(meal crawler.Meal) {
	_, err := collection.CreateDocument(context.Background(), meal)

	if err != nil {
		log.Fatal(err)
	}
}

func createDatabase() {
	exists, _ := client.DatabaseExists(context.Background(), databaseName)
	if !exists {
		options := &driver.CreateDatabaseOptions{}
		db, err := client.CreateDatabase(context.Background(), databaseName, options)

		if err != nil {
			log.Fatal(err)
		}

		database = db
	} else {
		db, _ := client.Database(context.Background(), databaseName)
		database = db
	}
}

func createCollection() {
	exists, _ := database.CollectionExists(context.Background(), collectionName)
	if !exists {
		options := &driver.CreateCollectionOptions{}
		coll, err := database.CreateCollection(context.Background(), collectionName, options)

		if err != nil {
			log.Fatal(err)
		}

		collection = coll
	} else {
		coll, _ := database.Collection(context.Background(), collectionName)
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
