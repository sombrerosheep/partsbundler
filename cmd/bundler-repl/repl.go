package main

import (
	"fmt"

	"github.com/sombrerosheep/partsbundler/pkg/core"
)

type ReplCmd interface {
	Exec(state *ReplState) error
	String() string
}

// Kits Commands

// GetKitsCmd Repl Command to get kits
type GetKitsCmd struct{}

func (cmd GetKitsCmd) Exec(state *ReplState) error {
	for _, v := range state.GetKits() {
		fmt.Printf("| %3d | %15s | %10s | %10s | %3d |\n",
			v.ID, v.Name, v.Schematic, v.Diagram, len(v.Parts))
	}

	return nil
}

func (cmd GetKitsCmd) String() string {
	return "GetKits"
}

// GetKitCmd Repl Command to get a kit by ID
type GetKitCmd struct {
	kitId int64
}

func (cmd GetKitCmd) Exec(state *ReplState) error {
	kit, err := state.GetKit(cmd.kitId)
	if err != nil {
		return err
	}

	fmt.Printf("| %3d | %15s | %10s | %10s | %3d |\n",
		kit.ID, kit.Name, kit.Schematic, kit.Diagram, len(kit.Parts))

	return nil
}

func (cmd GetKitCmd) String() string {
	return fmt.Sprintf("GetKit(%d)", cmd.kitId)
}

// Part Commands

func printPart(part core.Part) {
	fmt.Printf("| %3d | %20s | %25s |\n", part.ID, part.Kind, part.Name)
}

// GetPartsCmd Repl Command to get all parts
type GetPartsCmd struct{}

func (cmd GetPartsCmd) Exec(state *ReplState) error {
	for _, v := range state.GetParts() {
		printPart(v)
	}

	return nil
}

func (cmd GetPartsCmd) String() string {
	return "GetParts"
}

// GetPartCmd Repl Command to get a part by its Id
type GetPartCmd struct {
	partId int64
}

func (cmd GetPartCmd) Exec(state *ReplState) error {
	part, err := state.GetPart(cmd.partId)
	if err != nil {
		return err
	}

	printPart(part)
	return nil
}

func (cmd GetPartCmd) String() string {
	return fmt.Sprintf("GetPart(%d)", cmd.partId)
}

// NewPartCmd Repl Command to create a part
type NewPartCmd struct {
	name string
	kind core.PartType
}

func (cmd NewPartCmd) Exec(state *ReplState) error {
	part, err := state.CreatePart(cmd.name, cmd.kind)
	if err != nil {
		return err
	}

	fmt.Println("Added Part:")
	printPart(part)

	return nil
}

func (cmd NewPartCmd) String() string {
	return fmt.Sprintf("NewPart: %s (%s)", cmd.name, cmd.kind)
}

// Misc Commands

type PrintUsageCmd struct{}

func (cmd PrintUsageCmd) Exec(_ *ReplState) error {
	fmt.Println("Usage:")
	fmt.Println("\tget kits")
	fmt.Println("\tget kit :kitId:")
	fmt.Println("\tget parts")
	fmt.Println("\tget part :partId:")
  fmt.Println("\tnew part :kind: :name:")

	return nil
}

func (cmd PrintUsageCmd) String() string {
	return "PrintUsage"
}

type ExitCmd struct{}

func (cmd ExitCmd) Exec(_ *ReplState) error {
	return nil
}

func (cmd ExitCmd) String() string {
	return "Exit"
}
