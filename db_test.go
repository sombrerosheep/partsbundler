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

func Test_Links(t *testing.T) {
	var stor Storage = &SqliteStorage{
		DBFilePath: test_db_path,
	}
	err := stor.Connect()
	if err != nil {
		t.Fatalf("Error connecting to test db (%s): %s", test_db_path, err)
	}
	defer func() {
		err := stor.Close()
		if err != nil {
			t.Errorf("Error closing test database: %s", err)
		}

		err = os.Remove(test_db_path)
		if err != nil {
			t.Errorf("Error removing test database: %s", err)
		}
	}()

	const (
		id       int64 = 42
		testLink       = "example.com"
	)

	// Links
	// // Part Links
	t.Run("Add/Get/RemoveLinksForPart", func(t *testing.T) {
		// add
		link, err := stor.AddLinkToPart(id, testLink)

		assert.Nil(t, err)
		assert.NotEqual(t, link.ID, 0)
		assert.Equal(t, link.URL, testLink)

		// get
		links, err := stor.GetLinksForPart(id)

		assert.Nil(t, err)
		assert.NotNil(t, links)
		assert.Len(t, links, 1)
		assert.Equal(t, links[0], link)

		// remove
		err = stor.RemoveLinkFromPart(id, link.ID)

		assert.Nil(t, err)

		// verify removal
		endLinks, err := stor.GetLinksForPart(id)

		assert.Nil(t, err)
		assert.NotNil(t, endLinks)
		assert.Len(t, endLinks, 0)
	})

	// // Kit Links
	t.Run("Add/Get/RemoveLinksForKit", func(t *testing.T) {
		// add
		kit, err := stor.AddLinkToKit(id, testLink)

		assert.Nil(t, err)
		assert.NotEqual(t, kit.ID, 0)
		assert.Equal(t, kit.URL, testLink)

		// get
		kits, err := stor.GetLinksForKit(id)

		assert.Nil(t, err)
		assert.NotNil(t, kits)
		assert.Len(t, kits, 1)
		assert.Equal(t, kits[0], kit)

		// remove
		err = stor.RemoveLinkFromKit(id, kit.ID)

		assert.Nil(t, err)

		// verify removal
		endKits, err := stor.GetLinksForKit(id)

		assert.Nil(t, err)
		assert.NotNil(t, endKits)
		assert.Len(t, endKits, 0)
	})

	// Parts
	t.Run("Parts", func(t *testing.T) {
		// Add Part
		inpart := Part{
			ID:    0,
			Kind:  "Resistor",
			Name:  "1k",
			Links: []Link(nil),
		}
		part, err := stor.AddPart(inpart)

		assert.Nil(t, err)
		assert.NotEqual(t, inpart.ID, part.ID)
		assert.Equal(t, inpart.Kind, part.Kind)
		assert.Equal(t, inpart.Name, part.Name)

		{ // Get Part
			gpart, err := stor.GetPart(part.ID)

			assert.Nil(t, err)
			assert.Equal(t, part, gpart)
		}

		{ // GetParts
			parts, err := stor.GetParts()

			assert.Nil(t, err)
			assert.Len(t, parts, 1)
			assert.Equal(t, part, parts[0])
		}

		{
			upart := Part{
				ID:    part.ID,
				Kind:  "Capacitor",
				Name:  "47nf",
				Links: []Link(nil),
			}

			updatedPart, err := stor.UpdatePart(upart)

			assert.Nil(t, err)
			assert.Equal(t, upart, updatedPart)

			getUpdatedPart, err := stor.GetPart(upart.ID)

			assert.Nil(t, err)
			assert.Equal(t, updatedPart, getUpdatedPart)
		}

		{
			err := stor.DeletePart(part.ID)

			assert.Nil(t, err)

			parts, err := stor.GetParts()

			assert.Nil(t, err)
			assert.Len(t, parts, 0)
		}

	})

	// Kits
	t.Run("Kits", func(t *testing.T) {
		// Add sample parts
		var testKitParts = []Part{
			{Name: "1k", Kind: "Resistor"},
			{Name: "4u7", Kind: "Capacitor"},
			{Name: "TL072", Kind: "IC"},
		}

		for i, v := range testKitParts {
			p, err := stor.AddPart(v)
			testKitParts[i] = p
			assert.Nilf(t, err, "Error inserting part into test db (%d:%#v): %s", i, v, err)
		}

		// Add Kit
		inkit := Kit{
			ID:        0,
			Parts:     []KitPart(nil),
			Name:      "TS808",
			Schematic: "example.com/test/schematic",
			Diagram:   "example.com/test-diagram",
			Links:     []Link(nil),
		}
		kit, err := stor.AddKit(inkit)

		assert.Nil(t, err)
		assert.NotEqual(t, kit.ID, inkit.ID)
		assert.Equal(t, inkit.Name, kit.Name)
		assert.Equal(t, inkit.Schematic, kit.Schematic)
		assert.Equal(t, inkit.Diagram, kit.Diagram)

		{ // Get Kit
			gotKit, err := stor.GetKit(kit.ID)

			assert.Nil(t, err)
			assert.Equal(t, kit, gotKit)
		}

		{ // Get Kits
			kits, err := stor.GetKits()

			assert.Nil(t, err)
			assert.Len(t, kits, 1)
			assert.Equal(t, kit, kits[0])
		}

		{ // UpdateKit
			ukit := Kit{
				ID:        kit.ID,
				Parts:     []KitPart(nil),
				Name:      kit.Name,
				Schematic: "example.com/moved/test/schematic",
				Diagram:   "example.com/moved/test/diagram",
				Links:     []Link{},
			}

			updatedKit, err := stor.UpdateKit(ukit)

			assert.Nil(t, err)
			assert.Equal(t, ukit, updatedKit)
		}

		{ // AddPartToKit
			// use parts added above
			quantity := uint64(1)
			expectedKitParts := make([]KitPart, len(testKitParts))

			for i, v := range testKitParts {
				err = stor.AddPartToKit(v.ID, kit.ID, quantity)

				assert.Nilf(t, err, "Error adding test part (%d:%#v) to kit: %s", i, v, err)

				expectedKitParts[i] = KitPart{
					Part:     v,
					Quantity: quantity,
				}
			}

			kitParts, err := stor.GetKitParts(kit.ID)

			assert.Nil(t, err)
			assert.Equal(t, len(testKitParts), len(kitParts))
			assert.Equal(t, expectedKitParts, kitParts)
		}

		{ // GetKitParts

		}

		{ // SetPartQuantityForKit

		}

		{ // RemovePartFromKit

		}

		{ // DeleteKit

		}

	})
}
