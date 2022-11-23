package sqlite

import (
	"github.com/sombrerosheep/partsbundler/pkg/core"
)

var FakeLinks = [...]core.Link{
	{ID: 1, URL: "example.com"},
	{ID: 2, URL: "example.com/two"},
	{ID: 3, URL: "example.com/three"},
	{ID: 4, URL: "example.com/four"},
}

var FakeParts = [...]core.Part{
	{
		ID:    1,
		Kind:  "Resistor",
		Name:  "1k",
		Links: []core.Link(nil),
	},
	{
		ID:    2,
		Kind:  "Resistor",
		Name:  "2k",
		Links: []core.Link(nil),
	},
	{
		ID:    3,
		Kind:  "Resistor",
		Name:  "3k",
		Links: []core.Link(nil),
	},
}

var FakeKitParts = [...]core.KitPart{
	{Part: FakeParts[0], Quantity: 1},
	{Part: FakeParts[1], Quantity: 2},
	{Part: FakeParts[2], Quantity: 3},
}

var FakeKits = [...]core.Kit{
	{
		ID:        1,
		Parts:     []core.KitPart{},
		Name:      "The Alpha",
		Schematic: "example.com/the-alpha/schematic",
		Diagram:   "example.com/the-alpha/diagram",
		Links:     []core.Link{},
	},
	{
		ID:        2,
		Parts:     []core.KitPart{},
		Name:      "The Beta",
		Schematic: "example.com/the-beta/schematic",
		Diagram:   "example.com/the-beta/diagram",
		Links:     []core.Link{},
	},
}

type GreenSqliteMock struct {
	isqlitedb
}

func (db GreenSqliteMock) GetPart(partId int64) (core.Part, error) {
	if partId <= int64(len(FakeKitParts)) && partId > 0 {
		return FakeParts[partId-1], nil
	}

	return FakeParts[0], nil
}

func (db GreenSqliteMock) GetAllParts() ([]core.Part, error) {
	return FakeParts[:], nil
}

func (db GreenSqliteMock) GetPartLinks(partId int64) ([]core.Link, error) {
	return FakeLinks[:], nil
}

func (db GreenSqliteMock) AddLinkToPart(link string, partId int64) (int64, error) {
	return 1, nil
}

func (db GreenSqliteMock) RemoveLinkFromPart(linkId, partId int64) error {
	return nil
}

func (db GreenSqliteMock) CreatePart(name string, kind core.PartType) (int64, error) {
	return 1, nil
}

func (db GreenSqliteMock) RemovePart(partId int64) error {
	return nil
}

func (db GreenSqliteMock) GetKit(kitId int64) (core.Kit, error) {
	if kitId <= int64(len(FakeKits)) && kitId > 0 {
		return FakeKits[kitId-1], nil
	}
	return FakeKits[0], nil
}

func (db GreenSqliteMock) GetKitPartsForKit(kitId int64) ([]kitPartRef, error) {
	refs := []kitPartRef{
		{kitId: kitId, partId: 1, quantity: 1},
		{kitId: kitId, partId: 2, quantity: 2},
		{kitId: kitId, partId: 3, quantity: 3},
	}

	return refs, nil
}

func (db GreenSqliteMock) GetAllKits() ([]core.Kit, error) {
	return FakeKits[:], nil
}

func (db GreenSqliteMock) AddPartToKit(partId, kitId int64, quantity uint64) error {
	return nil
}

func (db GreenSqliteMock) GetKitPartUsage(partId int64) ([]int64, error) {
	ids := []int64{}

	for _, v := range FakeKits {
		ids = append(ids, v.ID)
	}

	return ids, nil
}

func (db GreenSqliteMock) UpdatePartQuantity(partId, kitId int64, quantity uint64) error {
	return nil
}

func (db GreenSqliteMock) RemovePartFromKit(partId, kitId int64) error {
	return nil
}

func (db GreenSqliteMock) GetKitLinks(kitId int64) ([]core.Link, error) {
	return FakeLinks[:], nil
}

func (db GreenSqliteMock) AddLinkToKit(link string, kitId int64) (int64, error) {
	return 1, nil
}

func (db GreenSqliteMock) RemoveLinkFromKit(linkId, kitId int64) error {
	return nil
}

func (db GreenSqliteMock) CreateKit(name, schematic, diagram string) (int64, error) {
	return 1, nil
}

func (db GreenSqliteMock) RemoveKit(kitId int64) error {
	return nil
}
