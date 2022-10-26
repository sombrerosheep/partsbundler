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
	var stor Storage = &SqliteDb{
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

	// AddLinkToPart
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

	// AddLinkToKit
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
}