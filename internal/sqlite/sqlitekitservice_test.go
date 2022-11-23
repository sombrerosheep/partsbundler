package sqlite

import (
	"testing"

	"github.com/sombrerosheep/partsbundler/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_sqlitekitservice_GetAll(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: GreenSqliteMock{},
			partservice: SqlitePartService{
				db: GreenSqliteMock{},
			},
		}

		expectedKitParts := make([]core.KitPart, len(FakeKitParts))
		for i := range FakeKitParts {
			kp := FakeKitParts[i]
			kp.Part.Links = FakeLinks[:]

			expectedKitParts[i] = kp
		}

		expectedKits := make([]core.Kit, len(FakeKits))
		for i := range FakeKits {
			k := FakeKits[i]
			k.Parts = expectedKitParts[:]
			k.Links = FakeLinks[:]

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
			db: GreenSqliteMock{},
			partservice: SqlitePartService{
				db: GreenSqliteMock{},
			},
		}

		expectedKitParts := make([]core.KitPart, len(FakeKitParts))
		for i := range FakeKitParts {
			kp := FakeKitParts[i]
			kp.Part.Links = FakeLinks[:]

			expectedKitParts[i] = kp
		}

		expectedKit := FakeKits[0]
		expectedKit.Parts = expectedKitParts[:]
		expectedKit.Links = FakeLinks[:]

		kit, err := sut.Get(expectedKit.ID)

		assert.Nil(t, err)
		assert.Equal(t, expectedKit, kit)
	})
}

func Test_sqlitekitservice_AddLink(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: GreenSqliteMock{},
			partservice: SqlitePartService{
				db: GreenSqliteMock{},
			},
		}

		expectedLink := FakeLinks[0]

		link, err := sut.AddLink(FakeKits[0].ID, FakeLinks[0].URL)

		assert.Nil(t, err)
		assert.Equal(t, expectedLink, link)
	})
}

func Test_sqlitekitservice_RemoveLink(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: GreenSqliteMock{},
			partservice: SqlitePartService{
				db: GreenSqliteMock{},
			},
		}

		err := sut.RemoveLink(FakeKits[0].ID, FakeLinks[0].ID)

		assert.Nil(t, err)
	})
}

func Test_sqlitekitservice_AddPart(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: GreenSqliteMock{},
			partservice: SqlitePartService{
				db: GreenSqliteMock{},
			},
		}

		newPartId := int64(999)

		err := sut.AddPart(FakeKits[0].ID, newPartId, 5)

		assert.Nil(t, err)
	})
}

func Test_GetPartUsage(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: GreenSqliteMock{},
			partservice: SqlitePartService{
				db: GreenSqliteMock{},
			},
		}

		expectedIds := []int64{
			FakeKits[0].ID,
			FakeKits[1].ID,
		}

		ids, err := sut.GetPartUsage(1)

		assert.Nil(t, err)
		assert.Equal(t, expectedIds, ids)
	})
}

func Test_sqlitekitservice_SetPartQuantity(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: GreenSqliteMock{},
			partservice: SqlitePartService{
				db: GreenSqliteMock{},
			},
		}

		err := sut.SetPartQuantity(FakeKits[0].ID, FakeParts[0].ID, 42)

		assert.Nil(t, err)
	})
}

func Test_sqlitekitservice_RemovePart(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: GreenSqliteMock{},
			partservice: SqlitePartService{
				db: GreenSqliteMock{},
			},
		}

		err := sut.RemovePart(FakeKits[0].ID, FakeParts[0].ID)

		assert.Nil(t, err)
	})
}

func Test_sqlitekitservice_New(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: GreenSqliteMock{},
			partservice: SqlitePartService{
				db: GreenSqliteMock{},
			},
		}

		expectedKit := core.Kit{
			ID:        FakeKits[0].ID,
			Parts:     []core.KitPart{},
			Name:      FakeKits[0].Name,
			Schematic: FakeKits[0].Schematic,
			Diagram:   FakeKits[0].Diagram,
			Links:     []core.Link{},
		}

		kit, err := sut.New(expectedKit.Name, expectedKit.Schematic, expectedKit.Diagram)

		assert.Nil(t, err)
		assert.Equal(t, expectedKit, kit)
	})
}

func Test_sqlitekitservice_Delete(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqliteKitService{
			db: GreenSqliteMock{},
			partservice: SqlitePartService{
				db: GreenSqliteMock{},
			},
		}

		err := sut.Delete(FakeKits[0].ID)

		assert.Nil(t, err)
	})
}
