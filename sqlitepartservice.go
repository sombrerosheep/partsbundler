package main

type SqlitePartService struct {
	db isqlitedb
}

func (service SqlitePartService) GetAll() ([]Part, error) {
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

func (service SqlitePartService) GetParts(ids []int64) ([]Part, error) {
	parts := make([]Part, len(ids))

	for i, id := range ids {
		part, err := service.Get(id)
		if err != nil {
			return nil, err
		}

		parts[i] = part
	}

	return parts, nil
}

func (service SqlitePartService) Get(partId int64) (Part, error) {
	part, err := service.db.GetPart(partId)
	if err != nil {
		return Part{}, err
	}

	links, err := service.db.GetPartLinks(part.ID)
	if err != nil {
		return Part{}, err
	}

	part.Links = links

	return part, nil
}

func (service SqlitePartService) AddLink(partId int64, link string) (Link, error) {
	linkId, err := service.db.AddLinkToPart(link, partId)
	if err != nil {
		return Link{}, err
	}

	l := Link{
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

func (service SqlitePartService) New(name string, kind PartType) (Part, error) {
	part := Part{
		Name:  name,
		Kind:  kind,
		Links: []Link{},
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
