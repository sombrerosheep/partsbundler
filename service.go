package main

// may pull up logic from the storage part
// i'd like to thin out the storage part heavily
// but i'm not sure how to maintain the interface
// without making it too db-specific
// does the data help me out?

type IPartService interface {
	GetAll() ([]Part, error)
	Get(partId int64) (Part, error)

	AddLink(partId int64, link string) (Link, error)
	RemoveLink(partId int64, linkId int64) error

	New(name string, kind PartType) (Part, error)
	DeletePart(partId int64) error
}

type IKitService interface {
	GetAll() ([]Kit, error)
	Get(kitId int64) (Kit, error)

	AddLink(kitId int64, link string) (Link, error)
	RemoveLink(kitId int64, linkId int64) error

	AddPart(kitId int64, partId int64, quantity uint64) error
	SetPartQuantity(kitId int64, partId int64, quantity uint64) error
	RemovePart(kitId int64, partId int64) error

	New(name string, schematic string, diagram string) (Kit, error)
	RemoveKit(kitId int64) error
}

type BundlerService struct {
	Parts IPartService
	Kits  IKitService
}

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
		return Part{}, nil
	}

	part.Links = links

	return part, nil
}

func (service SqlitePartService) AddLink(partId int64, link string) (Link, error) {
	linkId, err := service.db.AddLinkToPart(link, partId)
	if err != nil {
		return Link{}, nil
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

func (service SqlitePartService) DeletePart(partId int64) error {
	err := service.db.DeletePart(partId)

	return err
}

type SqliteKitService struct {
	db          isqlitedb
	partservice SqlitePartService
}

func (service SqliteKitService) getKitParts(kitId int64) ([]KitPart, error) {
	partRefs, err := service.db.GetKitPartsForKit(kitId)
	if err != nil {
		return nil, err
	}

	kitParts := make([]KitPart, len(partRefs))

	for i, partRef := range partRefs {
		kitPart := KitPart{
			Quantity: partRef.quantity,
		}

		part, err := service.partservice.Get(partRef.partId)
		if err != nil {
			return nil, err
		}

		kitPart.Part = part
		kitParts[i] = kitPart
	}

	return kitParts, nil
}

func (service SqliteKitService) GetAll() ([]Kit, error) {
	kits, err := service.db.GetAllKits()
	if err != nil {
		return nil, err
	}

	for i := range kits {
		kitParts, err := service.getKitParts(kits[i].ID)
		if err != nil {
			return nil, err
		}

		kits[i].Parts = kitParts
	}

	return kits, nil
}

func (service SqliteKitService) Get(kitId int64) (Kit, error) {
	kit, err := service.db.GetKit(kitId)
	if err != nil {
		return Kit{}, err
	}

	kit.Parts, err = service.getKitParts(kit.ID)
	if err != nil {
		return Kit{}, err
	}

	return kit, nil
}

func (service SqliteKitService) AddLink(kitId int64, link string) (Link, error) {
	l := Link{
		URL: link,
	}

	id, err := service.db.AddLinkToKit(link, kitId)
	if err != nil {
		return l, err
	}

	l.ID = id

	return l, nil
}

func (service SqliteKitService) RemoveLink(kitId int64, linkId int64) error {
	return service.db.RemoveLinkFromKit(linkId, kitId)
}

func (service SqliteKitService) AddPart(kitId, partId int64, quantity uint64) error {
	return service.db.AddPartToKit(partId, kitId, quantity)
}

func (service SqliteKitService) SetPartQuantity(kitId int64, partId int64, quantity uint64) error {
	return service.db.UpdatePartQuantity(partId, kitId, quantity)
}

func (service SqliteKitService) RemovePart(kitId, partId int64) error {
	return service.db.RemovePartFromKit(partId, kitId)
}

func (service SqliteKitService) New(name string, schematic string, diagram string) (Kit, error) {
	kit := Kit{
		ID:        0,
		Parts:     []KitPart{},
		Name:      name,
		Schematic: schematic,
		Diagram:   diagram,
		Links:     []Link{},
	}

	kitId, err := service.db.CreateKit(name, schematic, diagram)
	if err != nil {
		return kit, err
	}

	kit.ID = kitId

	return kit, nil
}

func (service SqliteKitService) RemoveKit(kitId int64) error {
	return service.db.RemoveKit(kitId)
}

func CreateSqliteService(dbPath string) (*BundlerService, error) {
	stor, err := CreateSqliteDB(dbPath)
	if err != nil {
		return nil, err
	}

	parts := SqlitePartService{
		db: stor,
	}
	kits := &SqliteKitService{
		db:          stor,
		partservice: parts,
	}

	svc := &BundlerService{
		Parts: parts,
		Kits:  kits,
	}

	return svc, nil
}
