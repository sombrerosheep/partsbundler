package main

import (
	"testing"

	"github.com/sombrerosheep/partsbundler/pkg/core"
	"github.com/sombrerosheep/partsbundler/pkg/service"
	"github.com/stretchr/testify/assert"
)

var linkIdCounter = int64(99)
var fakeLinks = [...]core.Link{
	{ID: 1, URL: "example.com/one"},
	{ID: 2, URL: "example.com/two"},
	{ID: 3, URL: "example.com/three"},
}

var partIdCounter = int64(99)
var fakeParts = [...]core.Part{
	{
		ID:    1,
		Kind:  "Resistor",
		Name:  "1k",
		Links: fakeLinks[:],
	},
	{
		ID:    2,
		Kind:  "Capacitor",
		Name:  "47pf",
		Links: fakeLinks[:],
	},
}

var kitIdCounter = int64(99)
var fakeKits = [...]core.Kit{
	{
		ID: 1,
		Parts: []core.KitPart{
			{
				Part:     fakeParts[0],
				Quantity: 1,
			},
		},
		Name:      "MyKit",
		Schematic: "example.com/mykit-schematic",
		Diagram:   "example.com/mykit-diagram",
		Links:     fakeLinks[:],
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

var stubBundlerService = &service.BundlerService{
	Parts: &stubParts,
	Kits:  &stubKits,
}

func (s *stubPartService) GetAll() ([]core.Part, error) {
	return fakeParts[:], nil
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

func (s *stubKitService) GetAll() ([]core.Kit, error) {
	return fakeKits[:], nil
}

func Test_GetKits(t *testing.T) {
	t.Run("should return kits", func(t *testing.T) {
		sut := &ReplState{bundler: stubBundlerService}
		sut.Refresh()

		expected := fakeKits[:]

		kits := sut.GetKits()

		assert.Equal(t, expected, kits)
	})
}

func Test_GetParts(t *testing.T) {
	t.Run("should return parts", func(t *testing.T) {

		sut := &ReplState{bundler: stubBundlerService}
		sut.Refresh()

		expected := fakeParts[:]

		parts := sut.GetParts()

		assert.Equal(t, expected, parts)
	})
}

func Test_GetKit(t *testing.T) {
	t.Run("should return each kit by id", func(t *testing.T) {
		sut := &ReplState{bundler: stubBundlerService}
		sut.Refresh()

		kits := sut.GetKits()

		assert.Len(t, kits, len(fakeKits))
		for _, k := range kits {
			expected := k

			actual, err := sut.GetKit(expected.ID)

			assert.Nil(t, err)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("should return KitNotFound when requested kit is not in state", func(t *testing.T) {
		sut := &ReplState{bundler: stubBundlerService}
		sut.Refresh()

		kitId := int64(99)

		_, err := sut.GetKit(kitId)

		assert.IsType(t, KitNotFound{}, err)
		assert.Equal(t, err.(KitNotFound).kitId, kitId)
	})
}

func Test_GetPart(t *testing.T) {
	t.Run("should return each part by id", func(t *testing.T) {
		sut := &ReplState{bundler: stubBundlerService}
		sut.Refresh()

		parts := sut.GetParts()

		assert.Len(t, parts, len(fakeParts))
		for _, p := range parts {
			expected := p

			actual, err := sut.GetPart(expected.ID)

			assert.Nil(t, err)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("should return PartNotFound when requested part is not in state", func(t *testing.T) {
		sut := &ReplState{bundler: stubBundlerService}
		sut.Refresh()

		partId := int64(99)

		_, err := sut.GetPart(partId)

		assert.IsType(t, PartNotFound{}, err)
		assert.Equal(t, err.(PartNotFound).partId, partId)
	})
}

func Test_CreatePart(t *testing.T) {
	t.Run("should create part and add it to the state", func(t *testing.T) {
		sut := &ReplState{bundler: stubBundlerService}
		sut.Refresh()

		const name string = "part name"
		const kind core.PartType = "Resistor"

		newPart, err := sut.CreatePart(name, kind)

		assert.Nil(t, err)
		assert.Greater(t, newPart.ID, int64(0))
		assert.Equal(t, name, newPart.Name)
		assert.Equal(t, kind, newPart.Kind)
		assert.Len(t, newPart.Links, 0)

		verify, err := sut.GetPart(newPart.ID)

		assert.Nil(t, err)
		assert.Equal(t, newPart, verify)
	})
}

func linksContainsLink(links []core.Link, find core.Link) assert.Comparison {
	return func() bool {
		for _, v := range links {
			if v.ID == find.ID && v.URL == find.URL {
				return true
			}
		}

		return false
	}
}

func Test_AddLinkToPart(t *testing.T) {
	t.Run("should add link and add it to the part", func(t *testing.T) {
		sut := &ReplState{bundler: stubBundlerService}
		sut.Refresh()

		expectedPartId := sut.parts[0].ID
		testLink := "example.com/a-new-link"

		newLink, err := sut.AddLinkToPart(expectedPartId, testLink)

		assert.Nil(t, err)
		assert.Greater(t, newLink.ID, int64(0))
		assert.Equal(t, testLink, newLink.URL)

		part, err := sut.GetPart(expectedPartId)

		assert.Nil(t, err)
		assert.Condition(t, linksContainsLink(part.Links, newLink))
	})

	t.Run("should return PartNotFound when partId doesn't exist", func(t *testing.T) {
		sut := &ReplState{bundler: stubBundlerService}
		sut.Refresh()

		partId := int64(7777)

		_, err := sut.AddLinkToPart(partId, "example.com")

		assert.NotNil(t, err)
		assert.IsType(t, PartNotFound{}, err)
		assert.Equal(t, partId, err.(PartNotFound).partId)
	})
}

func Test_RemoveLinkFromPart(t *testing.T) {
	t.Run("should remove link from part", func (t *testing.T) {
		sut := &ReplState{bundler: stubBundlerService}
		sut.Refresh()

		part := sut.parts[0]
		link := part.Links[0]

		err := sut.RemoveLinkFromPart(part.ID, link.ID)

		assert.Nil(t, err)

		actual, err := sut.GetPart(part.ID)

		assert.Nil(t, err)
		assert.False(t, linksContainsLink(actual.Links, link)())
	})

	t.Run("should return PartNotFound when partId doesn't exist", func(t *testing.T) {
		sut := &ReplState{bundler: stubBundlerService}
		sut.Refresh()
		
		partId := int64(999)
		linkId := int64(1)

		err := sut.RemoveLinkFromPart(partId, linkId)
		
		assert.NotNil(t, err)
		assert.IsType(t, PartNotFound{}, err)
		assert.Equal(t, partId, err.(PartNotFound).partId)
	})

	t.Run("should return LinkNotFound when partId doesn't include linkId", func(t *testing.T) {
		sut := &ReplState{bundler: stubBundlerService}
		sut.Refresh()
		
		partId := int64(1)
		linkId := int64(999)

		err := sut.RemoveLinkFromPart(partId, linkId)
		
		assert.NotNil(t, err)
		assert.IsType(t, LinkNotFound{}, err)
		assert.Equal(t, partId, err.(LinkNotFound).partId)
		assert.Equal(t, linkId, err.(LinkNotFound).linkId)
	})
}