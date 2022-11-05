package main

import (
	"fmt"
)

const (
	dbPath = "./data/partsbundler.db"
)

// type IState interface {

// }

type ReplState struct {
	kits    []Kit
	parts   []Part
	service *BundlerService
}

func (s *ReplState) Init() error {
	svc, err := CreateSqliteService(dbPath)
	if err != nil {
		return err
	}

	s.service = svc

	if err = s.Refresh(); err != nil {
		return err
	}

	return nil
}

func (s *ReplState) Refresh() error {
	kits, err := s.service.Kits.GetAll()
	if err != nil {
		return err
	}

	parts, err := s.service.Parts.GetAll()
	if err != nil {
		return err
	}

	s.kits = kits
	s.parts = parts

	return nil
}

func (s *ReplState) GetKits() []Kit {
	return s.kits[:]
}

func (s ReplState) GetParts() []Part {
	return s.parts[:]
}

func (s ReplState) GetKit(kitId int64) (Kit, error) {
	return Kit{}, nil
}

func (s ReplState) GetPart(partId int64) (Part, error) {
	return Part{}, nil
}

var service BundlerService

func main() {
	fmt.Println("Hello")

	state := ReplState{}
	err := state.Init()
	if err != nil {
		fmt.Printf("Error initializing sqlite service: %s", err)
		return
	}

	parts := state.GetParts()

	fmt.Println("Parts:")
	for _, v := range parts {
		fmt.Printf("| %3d | %20s | %25s |\n", v.ID, v.Kind, v.Name)
	}

	kits := state.GetKits()
	if err != nil {
		fmt.Printf("Error getting kits: $%s\n", err.Error())
		return
	}

	fmt.Println("Kits:")
	for _, v := range kits {
		fmt.Printf("| %3d | %15s | %10s | %10s | %3d |\n",
			v.ID, v.Name, v.Schematic, v.Diagram, len(v.Parts))
	}

	fmt.Println("byebye.")
}
