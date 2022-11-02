package main

import (
	"fmt"
)

var service BundlerService

func main() {
	fmt.Println("Hello")

	sqliteService, err := CreateSqliteService("./data/partsbundler.db")
	if err != nil {
		fmt.Printf("Error initializing sqlite service: %s", err)
		return
	}

	parts, err := sqliteService.Parts.GetAll()
	if err != nil {
		fmt.Printf("Error getting parts: $%s\n", err.Error())
		return
	}

	fmt.Println("Parts:")
	for _, v := range parts {
		fmt.Printf("| %3d | %20s | %25s |\n", v.ID, v.Kind, v.Name)
	}

	kits, err := sqliteService.Kits.GetAll()
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
