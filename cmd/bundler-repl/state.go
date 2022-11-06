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

func (s *ReplState) GetKits() []core.Kit {
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

func (s ReplState) GetPart(partId int64) (core.Part, error) {
	for i := range s.parts {
		if s.parts[i].ID == partId {
			return s.parts[i], nil
		}
	}

	return core.Part{}, PartNotFound{partId}
}
