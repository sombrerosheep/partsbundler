package sqlite

import (
	"github.com/sombrerosheep/partsbundler/pkg/core"
)

type SqlitePartService struct {
	db isqlitedb
}

func (service SqlitePartService) GetAll() ([]core.Part, error) {
	parts, err := service.db.GetAllParts()
	if err != nil {
		return nil, err
	}

	for i := range parts {
		part := &parts[i]

		links, err := service.db.GetPartLinks(part.ID)
		if err != nil {
			return nil, err
		}

		part.Links = links
	}

	return parts, nil
}

func (service SqlitePartService) GetParts(ids []int64) ([]core.Part, error) {
	parts := make([]core.Part, len(ids))

	for i, id := range ids {
		part, err := service.Get(id)
		if err != nil {
			return nil, err
		}

		parts[i] = part
	}

	return parts, nil
}

func (service SqlitePartService) Get(partId int64) (core.Part, error) {
	part, err := service.db.GetPart(partId)
	if err != nil {
		return core.Part{}, err
	}

	links, err := service.db.GetPartLinks(part.ID)
	if err != nil {
		return core.Part{}, err
	}

	part.Links = links

	return part, nil
}

func (service SqlitePartService) AddLink(partId int64, link string) (core.Link, error) {
	linkId, err := service.db.AddLinkToPart(link, partId)
	if err != nil {
		return core.Link{}, err
	}

	l := core.Link{
		ID:  linkId,
		URL: link,
	}

	return l, nil
}

func (service SqlitePartService) RemoveLink(partId int64, linkId int64) error {
	err := service.db.RemoveLinkFromPart(linkId, partId)
	if err != nil {
		return err
	}

	return nil
}

func (service SqlitePartService) New(name string, kind core.PartType) (core.Part, error) {
	part := core.Part{
		Name:  name,
		Kind:  kind,
		Links: []core.Link{},
	}

	partId, err := service.db.CreatePart(name, kind)
	if err != nil {
		return part, err
	}

	part.ID = partId

	return part, nil
}

func (service SqlitePartService) Delete(partId int64) error {
	err := service.db.RemovePart(partId)

	return err
}
