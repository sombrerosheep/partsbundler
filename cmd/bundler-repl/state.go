package main

import (
	"fmt"

	"github.com/sombrerosheep/partsbundler/internal/sqlite"

	"github.com/sombrerosheep/partsbundler/pkg/core"
	"github.com/sombrerosheep/partsbundler/pkg/service"
)

type KitNotFound struct {
	kitId int64
}

func (k KitNotFound) Error() string {
	return fmt.Sprintf("Kit %d not found", k.kitId)
}

type PartNotFound struct {
	partId int64
}

func (p PartNotFound) Error() string {
	return fmt.Sprintf("Part %d not found", p.partId)
}

type LinkNotFound struct {
	linkId, partId int64
}

func (l LinkNotFound) Error() string {
	return fmt.Sprintf("Link %d not found on Part %d", l.linkId, l.partId)
}

type PartInUse struct {
	partId int64
}

func (p PartInUse) Error() string {
	return fmt.Sprintf("Part %d is in use by one or more kits", p.partId)
}

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

func (s ReplState) GetKits() []core.Kit {
	return s.kits[:]
}

func (s ReplState) GetParts() []core.Part {
	return s.parts[:]
}

func (s ReplState) GetKit(kitId int64) (core.Kit, error) {
	for i := range s.kits {
		if s.kits[i].ID == kitId {
			return s.kits[i], nil
		}
	}

	return core.Kit{}, KitNotFound{kitId}
}

func (s ReplState) getPartRef(partId int64) (*core.Part, error) {
	for i := range s.parts {
		if s.parts[i].ID == partId {
			return &s.parts[i], nil
		}
	}

	return &core.Part{}, PartNotFound{partId}
}

func (s ReplState) GetPart(partId int64) (core.Part, error) {
	p, err := s.getPartRef(partId)
	if err != nil {
		return core.Part{}, err
	}

	if p == nil {
		return core.Part{}, PartNotFound{partId}
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
		return core.Link{}, PartNotFound{partId}
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
		return LinkNotFound{linkId, partId}
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
		return PartInUse{partId}
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
		return PartNotFound{partId}
	}

	s.parts = append(s.parts[:partIndex], s.parts[partIndex+1:]...)

	return nil
}
