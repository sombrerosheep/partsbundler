package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var exampleKitParts = [...]KitPart{
	{Part: exampleParts[0], Quantity: 1},
	{Part: exampleParts[1], Quantity: 2},
	{Part: exampleParts[2], Quantity: 3},
}

var exampleKits = [...]Kit{
	{
		ID:        1,
		Parts:     []KitPart{},
		Name:      "The Alpha",
		Schematic: "example.com/the-alpha/schematic",
		Diagram:   "example.com/the-alpha/diagram",
		Links:     []Link{},
	},
	{
		ID:        2,
		Parts:     []KitPart{},
		Name:      "The Beta",
		Schematic: "example.com/the-beta/schematic",
		Diagram:   "example.com/the-beta/diagram",
		Links:     []Link{},
	},
}

func (db greenSqliteMock) GetKit(kitId int64) (Kit, error) {
	if kitId <= int64(len(exampleKits)) && kitId > 0 {
		return exampleKits[kitId-1], nil
	}
	return exampleKits[0], nil
}

func (db greenSqliteMock) GetKitPartsForKit(kitId int64) ([]kitPartRef, error) {
	refs := []kitPartRef{
		{kitId: kitId, partId: 1, quantity: 1},
		{kitId: kitId, partId: 2, quantity: 2},
		{kitId: kitId, partId: 3, quantity: 3},
	}

	return refs, nil
}

func (db greenSqliteMock) GetAllKits() ([]Kit, error) {
	return exampleKits[:], nil
}

func (db greenSqliteMock) AddPartToKit(partId, kitId int64, quantity uint64) error {
	return nil
}

func (db greenSqliteMock) UpdatePartQuantity(partId, kitId int64, quantity uint64) error {
	return nil
}

func (db greenSqliteMock) RemovePartFromKit(partId, kitId int64) error {
	return nil
}

func (db greenSqliteMock) GetKitLinks(kitId int64) ([]Link, error) {
	return exampleLinks[:], nil
}

func (db greenSqliteMock) AddLinkToKit(link string, kitId int64) (int64, error) {
	return 1, nil
}

func (db greenSqliteMock) RemoveLinkFromKit(linkId, kitId int64) error {
	return nil
}

func (db greenSqliteMock) CreateKit(name, schematic, diagram string) (int64, error) {
	return 1, nil
}

func (db greenSqliteMock) RemoveKit(kitId int64) error {
	return nil
}

func Test_sqlitekitservice_GetAll(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: greenSqliteMock{},
			partservice: SqlitePartService{
				db: greenSqliteMock{},
			},
		}

		expectedKitParts := make([]KitPart, len(exampleKitParts))
		for i := range exampleKitParts {
			kp := exampleKitParts[i]
			kp.Part.Links = exampleLinks[:]

			expectedKitParts[i] = kp
		}

		expectedKits := make([]Kit, len(exampleKits))
		for i := range exampleKits {
			k := exampleKits[i]
			k.Parts = expectedKitParts[:]
			k.Links = exampleLinks[:]

			expectedKits[i] = k
		}

		kits, err := sut.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, expectedKits, kits)
	})
}

func Test_sqlitekitservice_Get(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: greenSqliteMock{},
			partservice: SqlitePartService{
				db: greenSqliteMock{},
			},
		}

		expectedKitParts := make([]KitPart, len(exampleKitParts))
		for i := range exampleKitParts {
			kp := exampleKitParts[i]
			kp.Part.Links = exampleLinks[:]

			expectedKitParts[i] = kp
		}

		expectedKit := exampleKits[0]
		expectedKit.Parts = expectedKitParts[:]
		expectedKit.Links = exampleLinks[:]

		kit, err := sut.Get(expectedKit.ID)

		assert.Nil(t, err)
		assert.Equal(t, expectedKit, kit)
	})
}

func Test_sqlitekitservice_AddLink(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: greenSqliteMock{},
			partservice: SqlitePartService{
				db: greenSqliteMock{},
			},
		}

		expectedLink := exampleLinks[0]

		link, err := sut.AddLink(exampleKits[0].ID, exampleLinks[0].URL)

		assert.Nil(t, err)
		assert.Equal(t, expectedLink, link)
	})
}

func Test_sqlitekitservice_RemoveLink(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: greenSqliteMock{},
			partservice: SqlitePartService{
				db: greenSqliteMock{},
			},
		}

		err := sut.RemoveLink(exampleKits[0].ID, exampleLinks[0].ID)

		assert.Nil(t, err)
	})
}

func Test_sqlitekitservice_AddPart(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: greenSqliteMock{},
			partservice: SqlitePartService{
				db: greenSqliteMock{},
			},
		}

		err := sut.RemovePart(exampleKits[0].ID, exampleParts[0].ID)

		assert.Nil(t, err)
	})
}

func Test_sqlitekitservice_SetPartQuantity(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: greenSqliteMock{},
			partservice: SqlitePartService{
				db: greenSqliteMock{},
			},
		}

		err := sut.SetPartQuantity(exampleKits[0].ID, exampleParts[0].ID, 42)

		assert.Nil(t, err)
	})
}

func Test_sqlitekitservice_RemovePart(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: greenSqliteMock{},
			partservice: SqlitePartService{
				db: greenSqliteMock{},
			},
		}

		err := sut.RemovePart(exampleKits[0].ID, exampleParts[0].ID)

		assert.Nil(t, err)
	})
}

func Test_sqlitekitservice_New(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: greenSqliteMock{},
			partservice: SqlitePartService{
				db: greenSqliteMock{},
			},
		}

		expectedKit := Kit{
			ID:        exampleKits[0].ID,
			Parts:     []KitPart{},
			Name:      exampleKits[0].Name,
			Schematic: exampleKits[0].Schematic,
			Diagram:   exampleKits[0].Diagram,
			Links:     []Link{},
		}

		kit, err := sut.New(expectedKit.Name, expectedKit.Schematic, expectedKit.Diagram)

		assert.Nil(t, err)
		assert.Equal(t, expectedKit, kit)
	})
}

func Test_sqlitekitservice_Delete(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: greenSqliteMock{},
			partservice: SqlitePartService{
				db: greenSqliteMock{},
			},
		}

		err := sut.Delete(exampleKits[0].ID)

		assert.Nil(t, err)
	})
}
