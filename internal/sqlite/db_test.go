package sqlite

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sombrerosheep/partsbundler/pkg/core"

	"github.com/stretchr/testify/assert"
)

const (
  test_db_setup = "./import/setup.sql"
  testLink = "example.com"
)

func getTestDbConnection(t *testing.T, dbPath string) (*sqlitedb, error) {
	// setup test database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to test db (%s): %s", dbPath, err)
	}

	// prepare test.db
	b, err := ioutil.ReadFile(test_db_setup)
	if err != nil {
		return nil, fmt.Errorf("Error reading setup script: %s", err)
	}

	// execute setup script
	_, err = db.Exec(string(b))
	if err != nil {
		return nil, fmt.Errorf("Error executing setup script: %s", err)
	}

	var testdb = sqlitedb{
		DBFilePath: dbPath,
	}

	err = testdb.Connect()
	if err != nil {
		return &testdb, nil
	}

	return &testdb, nil
}

func testDbDeferredCleanup(t *testing.T, db *sqlitedb, dbPath string) {
	err := db.Close()
	if err != nil {
		t.Errorf("Error closing test database: %s", err)
	}

	err = os.Remove(dbPath)
	if err != nil {
		t.Errorf("Error removing test database: %s", err)
	}
}

func TestSqliteParts(t *testing.T) {
	const dbPath = "./import/dbparttest.db"
	testdb, err := getTestDbConnection(t, dbPath)
	if err != nil {
		t.Fatalf("Error connecting to test db (%s): %s", dbPath, err)
	}
	defer testDbDeferredCleanup(t, testdb, dbPath)

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
		expected := core.Part{
			ID:    partId,
			Kind:  partKind,
			Name:  partName,
			Links: []core.Link(nil),
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
		expected := []core.Link{
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

	t.Run("RemovePart", func(t *testing.T) {
		emptyPart := core.Part{
			ID:    0,
			Kind:  "",
			Name:  "",
			Links: []core.Link(nil),
		}
		err := testdb.RemovePart(partId)

		assert.Nil(t, err)

		part, err := testdb.GetPart(partId)

		assert.Nil(t, err)
		assert.Equal(t, emptyPart, part)
	})

	t.Run("GetAllParts", func(t *testing.T) {
		expectedParts := []core.Part{
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
}

func Test_SQLiteKits(t *testing.T) {
	const dbPath = "./import/dbkittest.db"
	testdb, err := getTestDbConnection(t, dbPath)
	if err != nil {
		t.Fatalf("Error connecting to test db (%s): %s", dbPath, err)
	}
	defer testDbDeferredCleanup(t, testdb, dbPath)

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
			expectedKit := core.Kit{
				ID:        kitId,
				Parts:     []core.KitPart(nil),
				Name:      kitName,
				Schematic: testLink,
				Diagram:   testLink,
				Links:     []core.Link(nil),
			}
			kit, err := testdb.GetKit(kitId)

			assert.Nil(t, err)
			assert.Equal(t, kit, expectedKit)
		})

		t.Run("AddPartToKit", func(t *testing.T) {
			err := testdb.AddPartToKit(partId, kitId, quantity)

			assert.Nil(t, err)
		})

		t.Run("GetKitPartUsage", func(t *testing.T) {
			kitIds, err := testdb.GetKitPartUsage(partId)

			assert.Nil(t, err)
			assert.Len(t, kitIds, 1)
			assert.Equal(t, []int64{kitId}, kitIds)
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
			expectedLink := core.Link{ID: linkId, URL: testLink}

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
			emptyKit := core.Kit{
				ID:        0,
				Parts:     []core.KitPart(nil),
				Name:      "",
				Schematic: "",
				Diagram:   "",
				Links:     []core.Link(nil),
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
