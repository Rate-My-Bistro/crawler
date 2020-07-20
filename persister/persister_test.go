package persister

import (
	"github.com/loeffel-io/take"
	"testing"
)

func TestCreateDatabase(t *testing.T) {
	//setup
	establishConnection("localhost:28015")

	t.Run("create a database", func(t *testing.T) {
		createDatabase()

		exists := take.DatabaseExists(databaseName, session)
		if !exists {
			t.Errorf("database was could not created")
		}
	})
}

func TestCreateTable(t *testing.T) {
	//setup
	establishConnection("localhost:28015")
	createDatabase()

	t.Run("create a table", func(t *testing.T) {
		createTable()

		exists := take.TableExists(collectionName, session)
		if !exists {
			t.Errorf("table was could not created")
		}
	})
}

func TestInsert(t *testing.T) {
	//setup
	establishConnection("localhost:28015")
	createDatabase()
	createTable()

	t.Run("insert a record", func(t *testing.T) {
		meal := Meal{
			Name:       "TestName",
			Supplement: "TestSupplement",
			Price:      1.23,
			optionalSupplements: []Supplement{{
				Name:  "TestSupplement1",
				Price: 3.21,
			}, {
				Name:  "TestSupplement2",
				Price: 9.87,
			}},
		}

		insert(meal)
		exists := take.TableExists(collectionName, session)
		if !exists {
			t.Errorf("table was could not created")
		}
	})
}
