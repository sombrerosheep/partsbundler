package mock

import (
	"github.com/sombrerosheep/partsbundler/pkg/core"
	"github.com/sombrerosheep/partsbundler/pkg/service"
)

var linkIdCounter = int64(99)
var FakeLinks = [...]core.Link{
	{ID: 1, URL: "example.com/one"},
	{ID: 2, URL: "example.com/two"},
	{ID: 3, URL: "example.com/three"},
}

var partIdCounter = int64(99)
var FakeParts = [...]core.Part{
	{
		ID:    1,
		Kind:  "Resistor",
		Name:  "1k",
		Links: FakeLinks[:],
	},
	{
		ID:    2,
		Kind:  "Capacitor",
		Name:  "47pf",
		Links: FakeLinks[:],
	},
}

var kitIdCounter = int64(99)
var FakeKits = [...]core.Kit{
	{
		ID: 1,
		Parts: []core.KitPart{
			{
				Part:     FakeParts[0],
				Quantity: 1,
			},
		},
		Name:      "MyKit",
		Schematic: "example.com/mykit-schematic",
		Diagram:   "example.com/mykit-diagram",
		Links:     FakeLinks[:],
	},
}

type stubPartService struct {
	service.IPartService
}

type stubKitService struct {
	service.IKitService
}

var stubParts = stubPartService{}
var stubKits = stubKitService{}

var StubBundlerService = &service.BundlerService{
	Parts: &stubParts,
	Kits:  &stubKits,
}

func (s *stubPartService) GetAll() ([]core.Part, error) {
	return FakeParts[:], nil
}

func (s *stubPartService) New(name string, kind core.PartType) (core.Part, error) {
	id := partIdCounter
	partIdCounter += 1

	part := core.Part{
		ID:    id,
		Kind:  kind,
		Name:  name,
		Links: []core.Link{},
	}

	return part, nil
}

func (s *stubPartService) AddLink(partId int64, link string) (core.Link, error) {
	linkId := linkIdCounter
	linkIdCounter += 1

	newLink := core.Link{
		ID:  linkId,
		URL: link,
	}

	return newLink, nil
}

func (s *stubPartService) RemoveLink(linkId, partId int64) error {
	return nil
}

func (s *stubPartService) Delete(partId int64) error {
	return nil
}

func (s *stubKitService) GetAll() ([]core.Kit, error) {
	return FakeKits[:], nil
}

func (s *stubKitService) GetPartUsage(partId int64) ([]int64, error) {
	ids := []int64{}

	for _, k := range FakeKits {
		for _, p := range k.Parts {
			if p.ID == partId {
				ids = append(ids, k.ID)
				continue
			}
		}
	}

	return ids, nil
}

func (s *stubKitService) New(name, schematic, diagram string) (core.Kit, error) {
	kitId := kitIdCounter
	kitIdCounter += 1

	kit := core.Kit{
		ID:        kitId,
		Parts:     []core.KitPart{},
		Name:      name,
		Schematic: schematic,
		Diagram:   diagram,
		Links:     []core.Link{},
	}

	return kit, nil
}

func (s *stubKitService) AddLink(kitId int64, link string) (core.Link, error) {
	linkId := linkIdCounter
	linkIdCounter += 1

	newLink := core.Link{
		ID:  linkId,
		URL: link,
	}

	return newLink, nil
}

func (s *stubKitService) RemoveLink(kitId int64, linkId int64) error {
	return nil
}

func (s *stubKitService) AddPart(kitId, partId int64, quantity uint64) error {
	return nil
}

func (s *stubKitService) SetPartQuantity(kitId, partId int64, quantity uint64) error {
	return nil
}

func (s *stubKitService) RemovePart(kitId, partId int64) error {
	return nil
}

func (s *stubKitService) DeleteKit(kitId int64) error {
	return nil
}
