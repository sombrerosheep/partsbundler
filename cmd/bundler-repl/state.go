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

func (s *ReplState) GetKits() []core.Kit {
	return s.kits[:]
}

func (s ReplState) GetParts() []core.Part {
	return s.parts[:]
}

func (s ReplState) GetKit(kitId int64) (core.Kit, error) {
	return core.Kit{}, nil
}

func (s ReplState) GetPart(partId int64) (core.Part, error) {
	return core.Part{}, nil
}
