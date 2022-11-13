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

type CannotParseCommand struct {
	input string
}

func (cmd CannotParseCommand) Error() string {
	return fmt.Sprintf("Cannot parse into command: %s", cmd.input)
}

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

        return GetPartCmd{id}, nil
      }
    }

    case "new": {
      if words[1] == "part" && len(words) >= 4 {
        name := words[3]
        kind := words[2]

        return NewPartCmd{name, core.PartType(kind)}, nil
      } else if words[1] == "kit" && len(words) >= 5 {
        name := words[2]
        schematic := words[3]
        diagram := words[4]

        return NewKitCmd{name, schematic, diagram}, nil
      }
    }

    case "add": {
      if words[1] == "partlink" && len(words) >= 4 {
        id, err := strconv.ParseInt(words[2], 10, 64)
        if err != nil {
          return nil, err
        }
        link := words[3]

        return AddPartLinkCmd{id, link}, nil
      } else if words[1] == "kitlink" && len(words) >= 4 {
        id, err := strconv.ParseInt(words[2], 10, 64)
        if err != nil {
          return nil, err
        }
        link := words[3]

        return AddKitLinkCmd{id, link}, nil
      } else if words[1] == "kitpart" && len(words) >= 5 {
        kitId, err := strconv.ParseInt(words[2], 10, 64)
        if err != nil {
          return nil, err
        }

        partId, err := strconv.ParseInt(words[3], 10, 64)
        if err != nil {
          return nil, err
        }

        qty, err := strconv.ParseUint(words[4], 10, 64)
        if err != nil {
          return nil, err
        }

        return AddKitPartCmd{kitId, partId, qty}, nil
      }
    }

    case "remove": {
      if words[1] == "partlink" && len(words) >= 4 {
        partId, err := strconv.ParseInt(words[2], 10, 64)
        if err != nil {
          return nil, err
        }

        linkId, err := strconv.ParseInt(words[3], 10, 64)
        if err != nil {
          return nil, err
        }

        return RemovePartLinkCmd{partId: partId, linkId: linkId}, nil
      } else if words[1] == "kitlink" && len(words) >= 4 {
        kitId, err := strconv.ParseInt(words[2], 10, 64)
        if err != nil {
          return nil, err
        }

        linkId, err := strconv.ParseInt(words[3], 10, 64)
        if err != nil {
          return nil, err
        }

        return RemoveKitLinkCmd{kitId, linkId}, nil
      } else if words[1] == "kitpart" && len(words) >= 4 {
        kitId, err := strconv.ParseInt(words[2], 10, 64)
        if err != nil {
          return nil, err
        }

        partId, err := strconv.ParseInt(words[3], 10, 64)
        if err != nil {
          return nil, err
        }

        return RemoveKitPartCmd{kitId, partId}, nil
      }
    }

    case "set": {
      if words[1] == "kitpart" && len(words) >= 5 {
        kitId, err := strconv.ParseInt(words[2], 10, 64)
        if err != nil {
          return nil, err
        }

        partId, err := strconv.ParseInt(words[3], 10, 64)
        if err != nil {
          return nil, err
        }

        qty, err := strconv.ParseUint(words[4], 10, 64)
        if err != nil {
          return nil, err
        }

        return SetKitPartQuantityCmd{kitId, partId, qty}, nil
      }
    }

    case "delete": {
      if words[1] == "part" && len(words) >= 3 {
        id, err := strconv.ParseInt(words[2], 10, 64)
        if err != nil {
          return nil, err
        }

        return DeletePartCmd{id}, nil
      } else if words[1] == "kit" && len(words) >= 3 {
        id, err := strconv.ParseInt(words[2], 10, 64)
        if err != nil {
          return nil, err
        }

        return DeleteKitCmd{id}, nil
      }
    }
  }

	return nil, CannotParseCommand{input}
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
