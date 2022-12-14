package main

import (
	"github.com/sombrerosheep/partsbundler/internal/sqlite"

	"github.com/sombrerosheep/partsbundler/pkg/core"
	"github.com/sombrerosheep/partsbundler/pkg/service"
)

type ReplState struct {
	kits    []core.Kit
	parts   []core.Part
	bundler *service.BundlerService
}

func (s *ReplState) Init() error {
	svc, err := sqlite.CreateSqliteService(dbPath)
	if err != nil {
		return err
	}

	s.bundler = svc

	if err = s.Refresh(); err != nil {
		return err
	}

	return nil
}

func (s *ReplState) Refresh() error {
	kits, err := s.bundler.Kits.GetAll()
	if err != nil {
		return err
	}

	parts, err := s.bundler.Parts.GetAll()
	if err != nil {
		return err
	}

	s.kits = kits
	s.parts = parts

	return nil
}

func (s ReplState) GetParts() []core.Part {
	return s.parts[:]
}

func (s ReplState) getPartRef(partId int64) (*core.Part, error) {
	for i := range s.parts {
		if s.parts[i].ID == partId {
			return &s.parts[i], nil
		}
	}

	return &core.Part{}, core.PartNotFound{PartID: partId}
}

func (s ReplState) GetPart(partId int64) (core.Part, error) {
	p, err := s.getPartRef(partId)
	if err != nil {
		return core.Part{}, err
	}

	if p == nil {
		return core.Part{}, core.PartNotFound{PartID: partId}
	}

	return *p, nil
}

func (s *ReplState) CreatePart(name string, kind core.PartType) (core.Part, error) {
	part, err := s.bundler.Parts.New(name, kind)
	if err != nil {
		return part, err
	}

	s.parts = append(s.parts, part)

	return part, nil
}

func (s *ReplState) AddLinkToPart(partId int64, link string) (core.Link, error) {
	part, err := s.getPartRef(partId)
	if err != nil {
		return core.Link{}, core.PartNotFound{PartID: partId}
	}

	newLink, err := s.bundler.Parts.AddLink(partId, link)
	if err != nil {
		return core.Link{}, err
	}

	part.Links = append(part.Links, newLink)

	return newLink, nil
}

func (s *ReplState) RemoveLinkFromPart(partId, linkId int64) error {
	part, err := s.getPartRef(partId)
	if err != nil {
		return err
	}

	err = s.bundler.Parts.RemoveLink(partId, linkId)
	if err != nil {
		return err
	}

	// update state
	linkIndex := int64(-1)
	for i := range part.Links {
		if part.Links[i].ID == linkId {
			linkIndex = int64(i)
			break
		}
	}

	if linkIndex < 0 {
		// this would mean the db had the entry (and didnt error) but
		// the state did not. The state is most likely out of date.
		return core.LinkNotFound{LinkID: linkId, OwnerID: partId}
	}

	part.Links = append(part.Links[:linkIndex], part.Links[linkIndex+1:]...)

	return nil
}

func (s *ReplState) DeletePart(partId int64) error {
	part, err := s.getPartRef(partId)
	if err != nil {
		return err
	}

	kitIds, err := s.bundler.Kits.GetPartUsage(part.ID)
	if err != nil {
		return err
	}

	if len(kitIds) > 0 {
		return core.PartInUse{PartID: partId}
	}

	err = s.bundler.Parts.Delete(partId)
	if err != nil {
		return err
	}

	partIndex := -1
	for i := range s.parts {
		if s.parts[i].ID == partId {
			partIndex = i
			break
		}
	}

	if partIndex < 0 {
		return core.PartNotFound{PartID: partId}
	}

	s.parts = append(s.parts[:partIndex], s.parts[partIndex+1:]...)

	return nil
}

func (s ReplState) GetKits() []core.Kit {
	return s.kits[:]
}

func (s ReplState) GetKit(kitId int64) (core.Kit, error) {
	for i := range s.kits {
		if s.kits[i].ID == kitId {
			return s.kits[i], nil
		}
	}

	return core.Kit{}, core.KitNotFound{KitID: kitId}
}

func (s *ReplState) CreateKit(name, schematic, diagram string) (core.Kit, error) {
	kit, err := s.bundler.Kits.New(name, schematic, diagram)
	if err != nil {
		return core.Kit{}, err
	}

	s.kits = append(s.kits, kit)

	return kit, nil
}

func (s ReplState) getKitRef(kitId int64) (*core.Kit, error) {
	for i := range s.kits {
		if s.kits[i].ID == kitId {
			return &s.kits[i], nil
		}
	}

	return nil, core.KitNotFound{KitID: kitId}
}

func (s *ReplState) AddLinkToKit(kitId int64, link string) (core.Link, error) {
	newLink, err := s.bundler.Kits.AddLink(kitId, link)
	if err != nil {
		return core.Link{}, err
	}

	ref, err := s.getKitRef(kitId)
	if err != nil {
		return core.Link{}, err
	}

	ref.Links = append(ref.Links, newLink)

	return newLink, nil
}

func (s *ReplState) RemoveLinkFromKit(kitId, linkId int64) error {
	kit, err := s.getKitRef(kitId)
	if err != nil {
		return err
	}

	err = s.bundler.Kits.RemoveLink(kitId, linkId)
	if err != nil {
		return err
	}

	linkIndex := int64(-1)
	for i, v := range kit.Links {
		if v.ID == linkId {
			linkIndex = int64(i)
		}
	}

	if linkIndex > 0 {
		return core.LinkNotFound{LinkID: linkId, OwnerID: kitId}
	}

	kit.Links = append(kit.Links[:linkIndex], kit.Links[linkIndex+1:]...)

	return nil
}

func (s *ReplState) AddPartToKit(partId, kitId int64, quantity uint64) error {
	kit, err := s.getKitRef(kitId)
	if err != nil {
		return err
	}

	part, err := s.GetPart(partId)
	if err != nil {
		return err
	}

	err = s.bundler.Kits.AddPart(kitId, partId, quantity)
	if err != nil {
		return err
	}

	kitPart := core.KitPart{
		Part:     part,
		Quantity: quantity,
	}

	kit.Parts = append(kit.Parts, kitPart)

	return nil
}

func (s *ReplState) UpdatePartQuantity(partId, kitId int64, quantity uint64) error {
	kit, err := s.getKitRef(kitId)
	if err != nil {
		return err
	}

	_, err = s.GetPart(partId)
	if err != nil {
		return err
	}

	err = s.bundler.Kits.SetPartQuantity(kitId, partId, quantity)
	if err != nil {
		return err
	}

	for i := range kit.Parts {
		if kit.Parts[i].ID == partId {
			kit.Parts[i].Quantity = quantity
			break
		}
	}

	return nil
}

func (s *ReplState) RemovePartFromKit(partId, kitId int64) error {
	kit, err := s.getKitRef(kitId)
	if err != nil {
		return err
	}

	_, err = s.GetPart(partId)
	if err != nil {
		return err
	}

	err = s.bundler.Kits.RemovePart(kitId, partId)
	if err != nil {
		return err
	}

	partIndex := int64(-1)
	for i, p := range kit.Parts {
		if p.ID == partId {
			partIndex = int64(i)
		}
	}

	if partIndex < 0 {
		return core.PartNotFound{PartID: partId}
	}

	kit.Parts = append(kit.Parts[:partIndex], kit.Parts[partIndex+1:]...)

	return nil
}

func (s *ReplState) DeleteKit(kitId int64) error {
	kitIndex := int64(-1)
	for i := range s.kits {
		if s.kits[i].ID == kitId {
			kitIndex = int64(i)
		}
	}

	if kitIndex < 0 {
		return core.KitNotFound{KitID: kitId}
	}

	err := s.bundler.Kits.Delete(kitId)
	if err != nil {
		return err
	}

	s.kits = append(s.kits[:kitIndex], s.kits[kitIndex+1:]...)

	return nil
}
