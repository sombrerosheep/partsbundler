package main

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// var exampleLinks = [...]Link{
// 	{ID: 1, URL: "example.com"},
// 	{ID: 2, URL: "example.com/two"},
// 	{ID: 3, URL: "example.com/three"},
// 	{ID: 4, URL: "example.com/four"},
// }

// var exampleParts = [...]Part{
// 	{
// 		ID:    1,
// 		Kind:  "Resistor",
// 		Name:  "1k",
// 		Links: []Link(nil),
// 	},
// 	{
// 		ID:    2,
// 		Kind:  "Resistor",
// 		Name:  "2k",
// 		Links: []Link(nil),
// 	},
// 	{
// 		ID:    3,
// 		Kind:  "Resistor",
// 		Name:  "3k",
// 		Links: []Link(nil),
// 	},
// }

// var exampleKitParts = [...]KitPart{
// 	{Part: exampleParts[0], Quantity: 1},
// 	{Part: exampleParts[1], Quantity: 2},
// 	{Part: exampleParts[2], Quantity: 3},
// }

// var exampleKits = [...]Kit{
// 	{
// 		ID:        1,
// 		Parts:     []KitPart(nil),
// 		Name:      "Kit Alpha",
// 		Schematic: "example.com/schematic-alpha",
// 		Diagram:   "example.com/diagram-alpha",
// 		Links:     []Link(nil),
// 	},
// 	{
// 		ID:        2,
// 		Parts:     []KitPart(nil),
// 		Name:      "Kit Beta",
// 		Schematic: "example.com/schematic-beta",
// 		Diagram:   "example.com/diagram-beta",
// 		Links:     []Link(nil),
// 	},
// }

// type greenSqliteMock struct {
// 	isqlitedb
// }

// func (db greenSqliteMock) GetPart(partId int64) (Part, error) {
// 	return exampleParts[0], nil
// }

// func (db greenSqliteMock) GetAllParts() ([]Part, error) {
// 	return exampleParts[:], nil
// }

// func (db greenSqliteMock) GetPartLinks(partId int64) ([]Link, error) {
// 	return exampleLinks[:], nil
// }

// func (db greenSqliteMock) AddLinkToPart(link string, partId int64) (int64, error) {
// 	return 1, nil
// }

// func (db greenSqliteMock) RemoveLinkFromPart(linkId, partId int64) error {
// 	return nil
// }

// func (db greenSqliteMock) CreatePart(name string, kind PartType) (int64, error) {
// 	return 1, nil
// }

// func (db greenSqliteMock) DeletePart(partId int64) error {
// 	return nil
// }

// func (db greenSqliteMock) GetKit(kitId int64) (Kit, error) {
// 	return exampleKits[0], nil
// }

// func (db greenSqliteMock) GetKitPartsForKit(kitId int64) ([]kitPartRef, error) {
// 	refs := []kitPartRef{
// 		{
// 			kitId:    kitId,
// 			partId:   1,
// 			quantity: 1,
// 		},
// 		{
// 			kitId:    kitId,
// 			partId:   2,
// 			quantity: 2,
// 		},
// 		{
// 			kitId:    kitId,
// 			partId:   3,
// 			quantity: 3,
// 		},
// 	}

// 	return refs, nil
// }

// func (db greenSqliteMock) GetAllKits() ([]Kit, error) {
// 	return exampleKits[:], nil
// }

// func (db greenSqliteMock) AddPartToKit(partId, kitId int64, quantity uint64) error {
// 	return nil
// }

// func (db greenSqliteMock) UpdatePartQuantity(partId, kitId int64, quantity uint64) error {
// 	return nil
// }

// func (db greenSqliteMock) RemovePartFromKit(partId, kitId int64) error {
// 	return nil
// }

// func (db greenSqliteMock) GetKitLinks(kitId int64) ([]Link, error) {
// 	return exampleLinks[:], nil
// }

// func (db greenSqliteMock) AddLinkToKit(link string, kitId int64) (int64, error) {
// 	return 1, nil
// }

// func (db greenSqliteMock) RemoveLinkFromKit(linkId, kitId int64) error {
// 	return nil
// }

// func (db greenSqliteMock) CreateKit(name, schematic, diagram string) (int64, error) {
// 	return 1, nil
// }

// func (db greenSqliteMock) RemoveKit(kitId int64) error {
// 	return nil
// }

// func Test_SqlitePartsService(t *testing.T) {

// 	sut := SqlitePartService{
// 		db: greenSqliteMock{},
// 	}

// 	t.Run("GetAll", func(t *testing.T) {
// 		expectedKitParts := make([]KitPart, len(exampleParts))
// 		for i := range exampleParts {
// 			kp := KitPart{
// 				Part:     exampleParts[i],
// 				Quantity: uint64(i) + 1,
// 			}

// 			kp.Links = exampleLinks[:]

// 			expectedKitParts[i] = kp
// 		}
// 		expectedKits := make([]Kit, len(exampleKits))
// 		for i := range exampleKits {
// 			expectedKits[i] = exampleKits[i]
// 			expectedKits[i].Parts = expectedKitParts
// 			expectedKits[i].Links = exampleLinks[:]
// 		}

// 		kits, err := sut.GetAll()

// 		assert.Nil(t, err)
// 		// for i := range kits {
// 		// 	assert.Equalf(t, kits[i], expectedKits[i])
// 		// }
// 		assert.Equal(t, kits, expectedKits)
// 	})
// }