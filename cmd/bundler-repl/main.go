package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
)

const (
	dbPath = "../../data/partsbundler.db"
)

type ReplCmd interface {
	Exec(state *ReplState) error
	String() string
}

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

type GetPartsCmd struct{}

func (cmd GetPartsCmd) Exec(state *ReplState) error {
	for _, v := range state.GetParts() {
		fmt.Printf("| %3d | %20s | %25s |\n", v.ID, v.Kind, v.Name)
	}

	return nil
}

func (cmd GetPartsCmd) String() string {
	return "GetParts"
}

type GetPartCmd struct {
	partId int64
}

func (cmd GetPartCmd) Exec(state *ReplState) error {
	part, err := state.GetPart(cmd.partId)
	if err != nil {
		return err
	}

	fmt.Printf("| %3d | %20s | %25s |\n", part.ID, part.Kind, part.Name)
	return nil
}

func (cmd GetPartCmd) String() string {
	return fmt.Sprintf("GetPart(%d)", cmd.partId)
}

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

type PrintUsageCmd struct {}

func (cmd PrintUsageCmd) Exec(_ *ReplState) error {
	fmt.Println("Usage:")
	fmt.Println("\tget kits")
	fmt.Println("\tget kit :kitId:")
	fmt.Println("\tget parts")
	fmt.Println("\tget part :partId:")
	
	return nil
}

func (cmd PrintUsageCmd) String() string {
	return "PrintUsage"
}

type ExitCmd struct {}

func (cmd ExitCmd) Exec(_ *ReplState) error {
	return nil
}

func (cmd ExitCmd) String() string {
	return "Exit"
}

// GetCommand parses the provided input and returns a
// ReplCmd to be Executed.
func GetCommand(input string) (ReplCmd, error) {
	nolines := strings.ReplaceAll(input, "\n", "")
	words := strings.Split(nolines, " ")

	// fmt.Printf("| ")
	// for _, word := range words {
	// 	fmt.Printf(" %s |", word)
	// }
	// fmt.Printf("\n")

	if len(words) > 0  && strings.ToLower(words[0]) == "exit"{
		return ExitCmd{}, nil
	}

	if len(words) < 2 {
		return PrintUsageCmd{}, nil
	}

	switch strings.ToLower(words[0]) {
		case "get": {
			if words[1] == "parts" {
				return GetPartsCmd{}, nil
			} else if words[1] == "kits" {
				return GetKitsCmd{}, nil
			} else if words[1] == "kit" && len(words) >= 3 {
				id, err := strconv.ParseInt(words[2], 10, 64)
				if err != nil {
					return nil, err
				}

				return GetKitCmd{kitId: id}, nil
			} else if words[1] == "part" && len(words) >= 3 {
				id, err := strconv.ParseInt(words[2], 10, 64)
				if err != nil {
					return nil, err
				}

				return GetPartCmd{partId: id}, nil
			}
		}
	}

	return nil, fmt.Errorf("Cannot parse cmd from input (%v)", words)
}

func main() {
	fmt.Println("Hello")

	state := &ReplState{}
	err := state.Init()
	if err != nil {
		fmt.Printf("Error initializing sqlite service: %s\n", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("bundler:> ")

		text, _ := reader.ReadString('\n')

		cmd, err := GetCommand(text)
		if err != nil {
			fmt.Println(err)
			continue
		}
		
		if _, ok := cmd.(ExitCmd); ok {
			break
		}

		err = cmd.Exec(state)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("byebye.")
}
