package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var exampleLinks = [...]Link{
	{ID: 1, URL: "example.com"},
	{ID: 2, URL: "example.com/two"},
	{ID: 3, URL: "example.com/three"},
	{ID: 4, URL: "example.com/four"},
}

var exampleParts = [...]Part{
	{
		ID:    1,
		Kind:  "Resistor",
		Name:  "1k",
		Links: []Link(nil),
	},
	{
		ID:    2,
		Kind:  "Resistor",
		Name:  "2k",
		Links: []Link(nil),
	},
	{
		ID:    3,
		Kind:  "Resistor",
		Name:  "3k",
		Links: []Link(nil),
	},
}

type greenSqliteMock struct {
	isqlitedb
}

func (db greenSqliteMock) GetPart(partId int64) (Part, error) {
	if partId <= int64(len(exampleKitParts)) && partId > 0 {
		return exampleParts[partId - 1], nil
	}

	return exampleParts[0], nil
}

func (db greenSqliteMock) GetAllParts() ([]Part, error) {
	return exampleParts[:], nil
}

func (db greenSqliteMock) GetPartLinks(partId int64) ([]Link, error) {
	return exampleLinks[:], nil
}

func (db greenSqliteMock) AddLinkToPart(link string, partId int64) (int64, error) {
	return 1, nil
}

func (db greenSqliteMock) RemoveLinkFromPart(linkId, partId int64) error {
	return nil
}

func (db greenSqliteMock) CreatePart(name string, kind PartType) (int64, error) {
	return 1, nil
}

func (db greenSqliteMock) RemovePart(partId int64) error {
	return nil
}

func Test_sqlitepartservice_GetAll(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqlitePartService{
			db: greenSqliteMock{},
		}

		expectedParts := make([]Part, len(exampleParts))
		for i := range exampleParts {
			p := exampleParts[i]
			p.Links = exampleLinks[:]

			expectedParts[i] = p
		}

		parts, err := sut.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, parts, expectedParts)
	})
}

func Test_sqlitepartservice_Get(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqlitePartService{
			db: greenSqliteMock{},
		}

		expectedPart := exampleParts[0]

		part, err := sut.Get(exampleParts[0].ID)

		assert.Nil(t, err)
		assert.Equal(t, expectedPart, part)
	})
}

func Test_sqlitepartservice_AddLink(t *testing.T) {
	t.Run("AddLink", func(t *testing.T) {
		sut := SqlitePartService{
			db: greenSqliteMock{},
		}

		expectedLinkId := exampleLinks[0]

		linkId, err := sut.AddLink(exampleParts[0].ID, exampleLinks[0].URL)

		assert.Nil(t, err)
		assert.Equal(t, expectedLinkId, linkId)
	})
}

func Test_sqlitepartservice_RemoveLink(t *testing.T) {
	t.Run("RemoveLink", func(t *testing.T) {
		sut := SqlitePartService{
			db: greenSqliteMock{},
		}

		err := sut.RemoveLink(1, 1)

		assert.Nil(t, err)
	})
}

func Test_sqlitepartservice_New(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		sut := SqlitePartService{
			db: greenSqliteMock{},
		}

		expectedPart := exampleParts[0]
		expectedPart.Links = []Link{}

		part, err := sut.New(expectedPart.Name, expectedPart.Kind)

		assert.Nil(t, err)
		assert.Equal(t, part, expectedPart)
	})
}

func Test_sqlitepartservice_Delete(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {
		sut := SqlitePartService{
			db: greenSqliteMock{},
		}

		err := sut.Delete(1)

		assert.Nil(t, err)
	})
}
