package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const test_db_path = "./data/test.db"
const test_db_setup = "./data/import/setup.sql"

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(code)
}

func run(m *testing.M) (int, error) {
	// setup test database
	db, err := sql.Open("sqlite3", test_db_path)
	if err != nil {
		return -1, fmt.Errorf("Error connecting to test db (%s): %s", test_db_path, err)
	}

	// prepare test.db
	b, err := ioutil.ReadFile(test_db_setup)
	if err != nil {
		return -1, fmt.Errorf("Error reading setup script: %s", err)
	}

	// execute setup script
	_, err = db.Exec(string(b))
	if err != nil {
		return -1, fmt.Errorf("Error executing setup script: %s", err)
	}

	// teardown
	defer func() {
		db.Close()
	}()

	return m.Run(), nil
}

func getTestDbConnection(t *testing.T) (sqlitedb, error) {
	var testdb = sqlitedb{
		DBFilePath: test_db_path,
	}
	err := testdb.Connect()
	if err != nil {
		return testdb, nil
	}

	return testdb, nil
}

func testDbDeferredCleanup(t *testing.T, db sqlitedb) {
	err := db.Close()
	if err != nil {
		t.Errorf("Error closing test database: %s", err)
	}

	err = os.Remove(test_db_path)
	if err != nil {
		t.Errorf("Error removing test database: %s", err)
	}
}

func Test_sqlite(t *testing.T) {
	testdb, err := getTestDbConnection(t)
	if err != nil {
		t.Fatalf("Error connecting to test db (%s): %s", test_db_path, err)
	}
	defer testDbDeferredCleanup(t, testdb)

	const (
		testLink = "example.com"
	)

	t.Run("Parts", func(t *testing.T) {
		var partId int64 = 0
		var linkId int64 = 0
		const partName = "4.7k"
		const partKind = "Resistor"

		t.Run("CreatePart", func(t *testing.T) {
			pid, err := testdb.CreatePart(partName, partKind)

			partId = pid

			assert.Nil(t, err)
			assert.Greater(t, pid, int64(0))
		})

		t.Run("GetPart", func(t *testing.T) {
			expected := Part{
				ID:    partId,
				Kind:  partKind,
				Name:  partName,
				Links: []Link(nil),
			}

			part, err := testdb.GetPart(partId)

			assert.Nil(t, err)
			assert.Equal(t, expected, part)
		})

		t.Run("AddLinkToPart", func(t *testing.T) {
			lid, err := testdb.AddLinkToPart(testLink, partId)

			linkId = lid

			assert.Nil(t, err)
			assert.Greater(t, linkId, int64(0))
		})

		t.Run("GetPartLinks", func(t *testing.T) {
			expected := []Link{
				{ID: linkId, URL: testLink},
			}
			links, err := testdb.GetPartLinks(partId)

			assert.Nil(t, err)
			assert.Equal(t, expected, links)
		})

		t.Run("RemoveLinkFromPart", func(t *testing.T) {
			err := testdb.RemoveLinkFromPart(linkId, partId)

			assert.Nil(t, err)

			links, err := testdb.GetPartLinks(partId)

			assert.Nil(t, err)
			assert.Len(t, links, 0)
		})

		t.Run("DeletePart", func(t *testing.T) {
			emptyPart := Part{
				ID:    0,
				Kind:  "",
				Name:  "",
				Links: []Link(nil),
			}
			err := testdb.DeletePart(partId)

			assert.Nil(t, err)

			part, err := testdb.GetPart(partId)

			assert.Nil(t, err)
			assert.Equal(t, emptyPart, part)
		})

		t.Run("GetAllParts", func(t *testing.T) {
			expectedParts := []Part{
				{Name: "4.7k", Kind: "Resistor"},
				{Name: "47uf", Kind: "Capacitor"},
				{Name: "TL072", Kind: "IC"},
			}

			for i := range expectedParts {
				part := &expectedParts[i]

				id, err := testdb.CreatePart(part.Name, part.Kind)
				if err != nil {
					t.Fatalf("Error inserting test part (%d:%#v): %s",
						i, part, err)
				}

				part.ID = id
			}

			parts, err := testdb.GetAllParts()

			assert.Nil(t, err)
			assert.Len(t, parts, len(expectedParts))
			assert.Equal(t, expectedParts, parts)
		})
	})

	t.Run("Kits", func(t *testing.T) {
		var kitId int64 = 0
		var linkId int64 = 0
		const quantity = 7
		const partId int64 = 42
		const kitName = "The Burninator"
		const partName = "4.7k"
		const partKind = "Resistor"

		t.Run("CreateKit", func(t *testing.T) {
			kid, err := testdb.CreateKit(kitName, testLink, testLink)

			kitId = kid

			assert.Nil(t, err)
			assert.Greater(t, kid, int64(0))
		})

		t.Run("GetKit", func(t *testing.T) {
			expectedKit := Kit{
				ID:        kitId,
				Parts:     []KitPart(nil),
				Name:      kitName,
				Schematic: testLink,
				Diagram:   testLink,
				Links:     []Link(nil),
			}
			kit, err := testdb.GetKit(kitId)

			assert.Nil(t, err)
			assert.Equal(t, kit, expectedKit)
		})

		t.Run("AddPartToKit", func(t *testing.T) {
			err := testdb.AddPartToKit(partId, kitId, quantity)

			assert.Nil(t, err)
		})

		t.Run("GetPartKitsForKit", func(t *testing.T) {
			expectedParts := []kitPartRef{
				{
					kitId:    kitId,
					partId:   partId,
					quantity: quantity,
				},
			}

			partRefs, err := testdb.GetKitPartsForKit(kitId)

			assert.Nil(t, err)
			assert.Len(t, partRefs, len(expectedParts))
			assert.Equal(t, partRefs, expectedParts)
		})

		t.Run("UpdatePartQuantity", func(t *testing.T) {
			newquantity := quantity * uint64(2)
			expectedParts := []kitPartRef{
				{
					kitId:    kitId,
					partId:   partId,
					quantity: newquantity,
				},
			}

			err := testdb.UpdatePartQuantity(partId, kitId, newquantity)

			assert.Nil(t, err)

			partsRefs, err := testdb.GetKitPartsForKit(kitId)

			assert.Nil(t, err)
			assert.Equal(t, partsRefs, expectedParts)
		})

		t.Run("RemovePartsFromKit", func(t *testing.T) {
			err := testdb.RemovePartFromKit(partId, kitId)

			assert.Nil(t, err)

			partRefs, err := testdb.GetKitPartsForKit(kitId)

			assert.Nil(t, err)
			assert.Len(t, partRefs, 0)
		})

		t.Run("AddLinkToKit", func(t *testing.T) {
			id, err := testdb.AddLinkToKit(testLink, kitId)

			linkId = id

			assert.Nil(t, err)
			assert.Greater(t, id, int64(0))
		})

		t.Run("GetKitLinks", func(t *testing.T) {
			expectedLink := Link{ID: linkId, URL: testLink}

			links, err := testdb.GetKitLinks(kitId)

			assert.Nil(t, err)
			assert.Equal(t, links[0], expectedLink)
		})

		t.Run("RemoveLinkFromKit", func(t *testing.T) {
			err := testdb.RemoveLinkFromKit(linkId, kitId)

			assert.Nil(t, err)

			links, err := testdb.GetKitLinks(kitId)

			assert.Nil(t, err)
			assert.Len(t, links, 0)
		})

		t.Run("RemoveKit", func(t *testing.T) {
			emptyKit := Kit{
				ID:        0,
				Parts:     []KitPart(nil),
				Name:      "",
				Schematic: "",
				Diagram:   "",
				Links:     []Link(nil),
			}

			err := testdb.RemoveKit(kitId)

			assert.Nil(t, err)

			kit, err := testdb.GetKit(kitId)

			assert.Nil(t, err)
			assert.Equal(t, kit, emptyKit)
		})

		t.Run("GetAllKits", func(t *testing.T) {

		})
	})
}
