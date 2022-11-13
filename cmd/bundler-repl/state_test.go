package main

import (
	"testing"

	"github.com/sombrerosheep/partsbundler/pkg/core"
	"github.com/sombrerosheep/partsbundler/pkg/service/mock"
	"github.com/stretchr/testify/assert"
)

func Test_GetParts(t *testing.T) {
	t.Run("should return parts", func(t *testing.T) {

		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		expected := mock.FakeParts[:]

		parts := sut.GetParts()

		assert.Equal(t, expected, parts)
	})
}

func Test_GetPart(t *testing.T) {
	t.Run("should return each part by id", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		parts := sut.GetParts()

		assert.Len(t, parts, len(mock.FakeParts))
		for _, p := range parts {
			expected := p

			actual, err := sut.GetPart(expected.ID)

			assert.Nil(t, err)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("should return PartNotFound when requested part is not in state", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		partId := int64(99)

		_, err := sut.GetPart(partId)

		assert.IsType(t, core.PartNotFound{}, err)
		assert.Equal(t, err.(core.PartNotFound).PartID, partId)
	})
}

func Test_CreatePart(t *testing.T) {
	t.Run("should create part and add it to the state", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
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
		sut := &ReplState{bundler: mock.StubBundlerService}
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
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		partId := int64(7777)

		_, err := sut.AddLinkToPart(partId, "example.com")

		assert.NotNil(t, err)
		assert.IsType(t, core.PartNotFound{}, err)
		assert.Equal(t, partId, err.(core.PartNotFound).PartID)
	})
}

func Test_RemoveLinkFromPart(t *testing.T) {
	t.Run("should remove link from part", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
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
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		partId := int64(999)
		linkId := int64(1)

		err := sut.RemoveLinkFromPart(partId, linkId)

		assert.NotNil(t, err)
		assert.IsType(t, core.PartNotFound{}, err)
		assert.Equal(t, partId, err.(core.PartNotFound).PartID)
	})

	t.Run("should return LinkNotFound when partId doesn't include linkId", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		partId := int64(1)
		linkId := int64(999)

		err := sut.RemoveLinkFromPart(partId, linkId)

		assert.NotNil(t, err)
		assert.IsType(t, core.LinkNotFound{}, err)
		assert.Equal(t, partId, err.(core.LinkNotFound).OwnerID)
		assert.Equal(t, linkId, err.(core.LinkNotFound).LinkID)
	})
}

func Test_DeletePart(t *testing.T) {
	t.Run("should remove part", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		part, err := sut.CreatePart("test", "Resistor")

		assert.Nil(t, err)

		err = sut.DeletePart(part.ID)

		assert.Nil(t, err)

		_, err = sut.GetPart(part.ID)

		assert.NotNil(t, err)
		assert.IsType(t, core.PartNotFound{}, err)
		assert.Equal(t, part.ID, err.(core.PartNotFound).PartID)
	})

	t.Run("should return PartInUse when part is in use by a kit", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		partId := mock.FakeParts[0].ID

		err := sut.DeletePart(partId)

		assert.NotNil(t, err)
		assert.IsType(t, core.PartInUse{}, err)
		assert.Equal(t, partId, err.(core.PartInUse).PartID)
	})
}

func Test_GetKits(t *testing.T) {
	t.Run("should return kits", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		expected := mock.FakeKits[:]

		kits := sut.GetKits()

		assert.Equal(t, expected, kits)
	})
}

func Test_GetKit(t *testing.T) {
	t.Run("should return each kit by id", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		kits := sut.GetKits()

		assert.Len(t, kits, len(mock.FakeKits))
		for _, k := range kits {
			expected := k

			actual, err := sut.GetKit(expected.ID)

			assert.Nil(t, err)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("should return KitNotFound when requested kit is not in state", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		kitId := int64(99)

		_, err := sut.GetKit(kitId)

		assert.IsType(t, core.KitNotFound{}, err)
		assert.Equal(t, err.(core.KitNotFound).KitID, kitId)
	})
}

func Test_CreateKit(t *testing.T) {
	t.Run("should create kit", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		name := "my kit"
		schematic := "example.com/mykit/schematic"
		diagram := "example.com/mykit/diagram"

		kit, err := sut.CreateKit(name, schematic, diagram)

		assert.Nil(t, err)
		assert.Greater(t, kit.ID, int64(0))
		assert.Equal(t, name, kit.Name)
		assert.Equal(t, kit.Schematic, schematic)
		assert.Equal(t, kit.Diagram, diagram)

		gotKit, err := sut.GetKit(kit.ID)

		assert.Nil(t, err)
		assert.Equal(t, kit, gotKit)
	})
}

func Test_AddLinkToKit(t *testing.T) {
	t.Run("should add link and add it to the kit", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		url := "example.com/test"
		kitId := mock.FakeKits[0].ID

		link, err := sut.AddLinkToKit(kitId, url)

		assert.Nil(t, err)
		assert.Equal(t, url, link.URL)
		assert.Greater(t, link.ID, int64(0))

		kit, err := sut.GetKit(kitId)

		assert.Nil(t, err)
		assert.Condition(t, linksContainsLink(kit.Links, link))
	})

	t.Run("should return KitNotFound when kitId doesn't exist", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		url := "example.com/test"
		kitId := int64(99)

		_, err := sut.AddLinkToKit(kitId, url)

		assert.NotNil(t, err)
		assert.IsType(t, core.KitNotFound{}, err)
		assert.Equal(t, err.(core.KitNotFound).KitID, kitId)
	})
}

func Test_RemoveLinkFromKit(t *testing.T) {
	t.Run("should remove link from Kit", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		kit := mock.FakeKits[0]
		link := kit.Links[0]

		err := sut.RemoveLinkFromKit(kit.ID, link.ID)

		assert.Nil(t, err)

		gotKit, err := sut.GetKit(kit.ID)

		assert.Nil(t, err)
		assert.False(t, linksContainsLink(gotKit.Links, link)())
	})
}

func Test_AddPartToKit(t *testing.T) {
	t.Run("should add kit to part", func(t *testing.T) {

		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		kitId := mock.FakeKits[0].ID
		partName := "my part"
		quantity := uint64(42)

		newPart, err := sut.CreatePart(partName, "Resistor")

		assert.Nil(t, err)

		err = sut.AddPartToKit(newPart.ID, kitId, quantity)

		assert.Nil(t, err)

		refs, err := sut.bundler.Kits.GetPartUsage(newPart.ID)

		assert.Nil(t, err)

		assert.Condition(t, func() bool {
			for _, ref := range refs {
				if ref == kitId {
					return true
				}
			}

			return false
		})
	})

	t.Run("should return KitNotFound when kit doesn't exist", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		partId := mock.FakeParts[0].ID
		kitId := int64(9999)

		err := sut.AddPartToKit(partId, kitId, 1)

		assert.NotNil(t, err)
		assert.IsType(t, core.KitNotFound{}, err)
		assert.Equal(t, kitId, err.(core.KitNotFound).KitID)
	})

	t.Run("should return PartNotFound when part doesn't exist", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		partId := int64(9999)
		kitId := mock.FakeKits[0].ID

		err := sut.AddPartToKit(partId, kitId, 1)

		assert.NotNil(t, err)
		assert.IsType(t, core.PartNotFound{}, err)
		assert.Equal(t, partId, err.(core.PartNotFound).PartID)
	})
}

func Test_UpdatePartQuantity(t *testing.T) {
	t.Run("should updated part quantity", func(t *testing.T) {

		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		kit := mock.FakeKits[0]
		part := kit.Parts[0]
		newQty := part.Quantity * 2

		err := sut.UpdatePartQuantity(part.ID, kit.ID, newQty)

		assert.Nil(t, err)

		gotKit, err := sut.GetKit(kit.ID)

		assert.Nil(t, err)

		var thatPart *core.KitPart = nil
		for i := range gotKit.Parts {
			if gotKit.Parts[i].ID == part.ID {
				thatPart = &gotKit.Parts[i]
			}
		}

		assert.NotNil(t, thatPart)
		assert.Equal(t, newQty, thatPart.Quantity)
	})

	t.Run("should return KitNotFound when kit does not exist", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		kitId := int64(9999)
		part := mock.FakeKits[0].Parts[0]
		newQty := part.Quantity * 2

		err := sut.UpdatePartQuantity(part.ID, kitId, newQty)

		assert.NotNil(t, err)
		assert.IsType(t, core.KitNotFound{}, err)
		assert.Equal(t, kitId, err.(core.KitNotFound).KitID)
	})

	t.Run("should return PartNotFound when part does not exist", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		kitId := mock.FakeKits[0].ID
		partId := int64(9999)
		newQty := uint64(9876)

		err := sut.UpdatePartQuantity(partId, kitId, newQty)

		assert.NotNil(t, err)
		assert.IsType(t, core.PartNotFound{}, err)
		assert.Equal(t, partId, err.(core.PartNotFound).PartID)
	})
}

func Test_RemovePartFromKit(t *testing.T) {
	t.Run("should remove part from kit", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		kit := mock.FakeKits[0]
		part := kit.Parts[0]

		err := sut.RemovePartFromKit(part.ID, kit.ID)

		assert.Nil(t, err)

		gotKit, err := sut.GetKit(kit.ID)

		assert.Nil(t, err)
		assert.Condition(t, func() bool {
			for _, p := range gotKit.Parts {
				if p.ID == part.ID {
					return false
				}
			}

			return true
		})
	})

	t.Run("should return KitNotFound when kit does not exist", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		kitId := int64(9999)
		part := mock.FakeKits[0].Parts[0]

		err := sut.RemovePartFromKit(part.ID, kitId)

		assert.NotNil(t, err)
		assert.IsType(t, core.KitNotFound{}, err)
		assert.Equal(t, kitId, err.(core.KitNotFound).KitID)
	})

	t.Run("should return PartNotFound when part does not exist", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		kitId := mock.FakeKits[0].ID
		partId := int64(9999)

		err := sut.RemovePartFromKit(partId, kitId)

		assert.NotNil(t, err)
		assert.IsType(t, core.PartNotFound{}, err)
		assert.Equal(t, partId, err.(core.PartNotFound).PartID)
	})
}

func Test_DeleteKit(t *testing.T) {
	t.Run("should delete kit", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		kitId := mock.FakeKits[0].ID

		err := sut.DeleteKit(kitId)

		assert.Nil(t, err)

		_, err = sut.GetKit(kitId)

		assert.NotNil(t, err)
		assert.IsType(t, core.KitNotFound{}, err)
		assert.Equal(t, kitId, err.(core.KitNotFound).KitID)
	})

	t.Run("should return KitNotFound when kit doesn't exist", func(t *testing.T) {
		sut := &ReplState{bundler: mock.StubBundlerService}
		sut.Refresh()

		kitId := int64(9999)

		err := sut.DeleteKit(kitId)

		assert.NotNil(t, err)
		assert.IsType(t, core.KitNotFound{}, err)
		assert.Equal(t, kitId, err.(core.KitNotFound).KitID)
	})
}
