package main

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

		kitLinks, err := service.db.GetKitLinks(kits[i].ID)
		if err != nil {
			return nil, err
		}

		kits[i].Links = kitLinks
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

	kitLinks, err := service.db.GetKitLinks(kit.ID)
	if err != nil {
		return Kit{}, err
	}

	kit.Links = kitLinks


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

func (service SqliteKitService) Delete(kitId int64) error {
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
