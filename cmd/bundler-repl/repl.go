package main

import (
	"fmt"

	"github.com/sombrerosheep/partsbundler/pkg/core"
)

type ReplCmd interface {
	Exec(state *ReplState) error
	String() string
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

type DeletePartCmd struct {
	partId int64
}

func (cmd DeletePartCmd) Exec(state *ReplState) error {
	err := state.DeletePart(cmd.partId)

	return err
}

func (cmd DeletePartCmd) String() string {
	return fmt.Sprintf("DeletePart: %d", cmd.partId)
}

type AddPartLinkCmd struct {
	partId int64
	link   string
}

func (cmd AddPartLinkCmd) Exec(state *ReplState) error {
	_, err := state.AddLinkToPart(cmd.partId, cmd.link)

	return err
}

func (cmd AddPartLinkCmd) String() string {
	return fmt.Sprintf("AddPartLink: %d (%s)", cmd.partId, cmd.link)
}

type RemovePartLinkCmd struct {
	partId int64
	linkId int64
}

func (cmd RemovePartLinkCmd) Exec(state *ReplState) error {
	err := state.RemoveLinkFromPart(cmd.partId, cmd.linkId)

	return err
}

func (cmd RemovePartLinkCmd) String() string {
	return fmt.Sprintf("RemovePartLink: %d (%d)", cmd.partId, cmd.linkId)
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

// NewKitCmd
type NewKitCmd struct {
	name      string
	schematic string
	diagram   string
}

func (cmd NewKitCmd) Exec(state *ReplState) error {
	kit, err := state.CreateKit(cmd.name, cmd.schematic, cmd.diagram)
	if err != nil {
		return err
	}

	state.kits = append(state.kits, kit)

	return nil
}

func (cmd NewKitCmd) String() string {
	return fmt.Sprintf("NewKit: %s | %s | %s", cmd.name, cmd.schematic, cmd.diagram)
}

// AddKitLinkCmd
type AddKitLinkCmd struct {
	kitId int64
	link  string
}

func (cmd AddKitLinkCmd) Exec(state *ReplState) error {
	_, err := state.AddLinkToKit(cmd.kitId, cmd.link)

	return err
}

func (cmd AddKitLinkCmd) String() string {
	return fmt.Sprintf("AddKitLink: %d, (%s)", cmd.kitId, cmd.link)
}

// RemoveKitLinkCmd
type RemoveKitLinkCmd struct {
	kitId  int64
	linkId int64
}

func (cmd RemoveKitLinkCmd) Exec(state *ReplState) error {
	err := state.RemoveLinkFromKit(cmd.kitId, cmd.linkId)

	return err
}

func (cmd RemoveKitLinkCmd) String() string {
	return fmt.Sprintf("RemoveKitLink: %d (%d)", cmd.kitId, cmd.linkId)
}

// AddKitPartCmd
type AddKitPartCmd struct {
	kitId    int64
	partId   int64
	quantity uint64
}

func (cmd AddKitPartCmd) Exec(state *ReplState) error {
	err := state.AddPartToKit(cmd.partId, cmd.kitId, cmd.quantity)

	return err
}

func (cmd AddKitPartCmd) String() string {
	return fmt.Sprintf("AddKitPart: %d:%d (%d)", cmd.kitId, cmd.partId, cmd.quantity)
}

// SetKitPartQuantityCmd
type SetKitPartQuantityCmd struct {
	kitId    int64
	partId   int64
	quantity uint64
}

func (cmd SetKitPartQuantityCmd) Exec(state *ReplState) error {
	err := state.UpdatePartQuantity(cmd.partId, cmd.kitId, cmd.quantity)

	return err
}

func (cmd SetKitPartQuantityCmd) String() string {
	return fmt.Sprintf("SetKitPartQuantity: %d:%d (%d)", cmd.kitId, cmd.partId, cmd.quantity)
}

// RemoveKitPartCmd
type RemoveKitPartCmd struct {
	kitId  int64
	partId int64
}

func (cmd RemoveKitPartCmd) Exec(state *ReplState) error {
	err := state.RemovePartFromKit(cmd.partId, cmd.kitId)

	return err
}

func (cmd RemoveKitPartCmd) String() string {
	return fmt.Sprintf("RemoveKitPart: %d:%d", cmd.kitId, cmd.partId)
}

// DeleteKitCmd
type DeleteKitCmd struct {
	kitId int64
}

func (cmd DeleteKitCmd) Exec(state *ReplState) error {
	err := state.DeleteKit(cmd.kitId)

	return err
}

func (cmd DeleteKitCmd) String() string {
	return fmt.Sprintf("DeleteKid: %d", cmd.kitId)
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
	fmt.Println("\tdelete part :partId:")

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
