package main

import (
	"fmt"
)

type PartType string

const (
	Resistor      PartType = "Resistor"
	Capacitor              = "Capacitor"
	IC                     = "IC"
	Transistor             = "Transistor"
	Diode                  = "Diode"
	Potentiometer          = "Potentiometer"
	Switch                 = "Switch"
)

type Part struct {
	ID    int64    `json:"id"`
	Kind  PartType `json:"kind"`
	Name  string   `json:"name"`
	Links []string `json:"links"`
}

func (d *SqlLiteDb) GetPart(partID uint64) (Part, error) {
	const stmt string = "SELECT id, kind, name from parts where id = ? limit 1;"
	var part = Part{}

	row := d.DB.QueryRow(stmt, partID)
	err := row.Scan(&part.ID, &part.Kind, &part.Name)
	if err != nil {
		return part, err
	}

	// todo: get part links

	return Part{}, nil
}

func (d *SqlLiteDb) PutPart(p Part) (Part, error) {
	const stmt string = `
		insert into parts(kind, name)
		values(?, ?)
	`
	newPart := p

	res, err := d.DB.Exec(stmt, p.Kind, p.Name)
	if err != nil {
		return newPart, err
	}

	newId, err := res.LastInsertId()
	if err != nil {
		return newPart, err
	}

	newPart.ID = newId

	// todo: links

	return newPart, nil
}

func (d *SqlLiteDb) UpdatePart(p Part) (Part, error) {
	const stmt string = `
		update parts
		set kind = ?,
		set name = ?
		where id = ?
	`
	res, err := d.DB.Exec(stmt, p.Kind, p.Name, p.ID)
	if err != nil {
		return p, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return p, err
	}

	if rows != 1 {
		return p, fmt.Errorf("Expected 1 row to be affected but %d were", rows)
	}

	// todo: links

	return p, nil
}

func (d *SqlLiteDb) DeletePart(partID uint64) error {
	const stmt string = `
		delete from parts where id = ?
	`

	res, err := d.DB.Exec(stmt, partID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("Expected 1 row to be affected but %d were", rows)
	}

	//todo: links
	return nil
}
