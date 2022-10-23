package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Datastore interface {
	Connect() error
}

type SqlLiteDb struct {
	DB *sql.DB
}

func (d *SqlLiteDb) Initialize() error {
	if d.DB == nil {
		return fmt.Errorf("Not connected to database")
	}

	{ // Create Parts Table
		const create string = `
			CREATE TABLE IF NOT EXISTS parts (
				id INTEGER PRIMARY KEY
				kind TEXT
				name TEXT
			);
		`

		if err := d.Exec(create); err != nil {
			return err
		}
	}

	{ // Create Kit Table
		const create string = `
			CREATE TABLE IF NOT EXISTS kits (
				id INTEGER PRIMARY KEY
				name TEXT
				schematic TEXT
				diagram TEXT
			);
		`

		if err := d.Exec(create); err != nil {
			return err
		}
	}

	{ // Create Links Table
		const create string = `
			CREATE TABLE IF NOT EXISTS links (
				id INTEGER PRIMARY KEY
				url TEXT
			);
		`

		if err := d.Exec(create); err != nil {
			return err
		}
	}

	{ // Create KitAssociations Table
		const create string = `
			CREATE TABLE IF NOT EXISTS kitassociations (
				id INTEGER PRIMARY KEY
				partId UNSIGNED BIG INT NOT NULL
				kitId UNSIGNED BIG INT NOT NULL
				quantity UNSIGNED BIG INT NOT NULL
			);
		`

		if err := d.Exec(create); err != nil {
			return err
		}
	}

	{ // Create KitLinks Table
		const create string = `
		CREATE TABLE IF NOT EXISTS kitlinks (
			id INTEGER PRIMARY KEY
			kitId UNSIGNED BIG INT NOT NULL
			linkId UNSIGNED BIG INT NOT NULL
		);
	`

		if err := d.Exec(create); err != nil {
			return err
		}
	}

	{ // Create PartLinks Table
		const create string = `
		CREATE TABLE IF NOT EXISTS partlinks (
			id INTEGER PRIMARY KEY
			partId UNSIGNED BIG INT NOT NULL
			linkId UNSIGNED BIG INT NOT NULL
		);
	`

		if err := d.Exec(create); err != nil {
			return err
		}
	}

	return nil
}

func (d *SqlLiteDb) Connect(filePath string) error {
	var err error

	d.DB, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return err
	}

	return nil
}

func (d *SqlLiteDb) Exec(statment string) error {
	if d.DB == nil {
		return fmt.Errorf("Not connected to database")
	}

	if _, err := d.DB.Exec(statment); err != nil {
		return err
	}

	return nil
}

func (d *SqlLiteDb) Query(query string) error {
	if d.DB == nil {
		return fmt.Errorf("Not connected to database")
	}

	return nil
}