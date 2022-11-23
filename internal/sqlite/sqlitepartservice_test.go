package sqlite

import (
	"testing"

	"github.com/sombrerosheep/partsbundler/pkg/core"

	"github.com/stretchr/testify/assert"
)

func Test_sqlitepartservice_GetAll(t *testing.T) {
	t.Run("When no errors are returned", func(t *testing.T) {
		sut := SqlitePartService{
			db: GreenSqliteMock{},
		}

		expectedParts := make([]core.Part, len(FakeParts))
		for i := range FakeParts {
			p := FakeParts[i]
			p.Links = FakeLinks[:]

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
			db: GreenSqliteMock{},
		}

		expectedPart := FakeParts[0]

		part, err := sut.Get(FakeParts[0].ID)

		assert.Nil(t, err)
		assert.Equal(t, expectedPart, part)
	})
}

func Test_sqlitepartservice_AddLink(t *testing.T) {
	t.Run("AddLink", func(t *testing.T) {
		sut := SqlitePartService{
			db: GreenSqliteMock{},
		}

		expectedLinkId := FakeLinks[0]

		linkId, err := sut.AddLink(FakeParts[0].ID, FakeLinks[0].URL)

		assert.Nil(t, err)
		assert.Equal(t, expectedLinkId, linkId)
	})
}

func Test_sqlitepartservice_RemoveLink(t *testing.T) {
	t.Run("RemoveLink", func(t *testing.T) {
		sut := SqlitePartService{
			db: GreenSqliteMock{},
		}

		err := sut.RemoveLink(1, 1)

		assert.Nil(t, err)
	})
}

func Test_sqlitepartservice_New(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		sut := SqlitePartService{
			db: GreenSqliteMock{},
		}

		expectedPart := FakeParts[0]
		expectedPart.Links = []core.Link{}

		part, err := sut.New(expectedPart.Name, expectedPart.Kind)

		assert.Nil(t, err)
		assert.Equal(t, part, expectedPart)
	})
}

func Test_sqlitepartservice_Delete(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {
		sut := SqlitePartService{
			db: GreenSqliteMock{},
		}

		err := sut.Delete(1)

		assert.Nil(t, err)
	})
}
