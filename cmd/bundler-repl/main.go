package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sombrerosheep/partsbundler/pkg/core"
)

const (
	dbPath = "../../data/partsbundler.db"
)

func printInput(words []string) {
	fmt.Printf("| ")
	for _, word := range words {
		fmt.Printf(" %s |", word)
	}
	fmt.Printf("\n")
}

// GetCommand parses the provided input and returns a
// ReplCmd to be Executed.
func GetCommand(input string) (ReplCmd, error) {
	nolines := strings.ReplaceAll(input, "\n", "")
	words := strings.Split(nolines, " ")

	// printInput(words)

	if len(words) > 0 && strings.ToLower(words[0]) == "exit" {
		return ExitCmd{}, nil
	}

	if len(words) < 2 {
		return PrintUsageCmd{}, nil
	}

	switch strings.ToLower(words[0]) {
	case "get":
		{
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

	case "new":
		{
			if words[1] == "part" && len(words) >= 4 {
				name := words[3]
				kind := words[2]

				return NewPartCmd{name, core.PartType(kind)}, nil
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
