package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SqlLiteDb struct {
	db *sql.DB
}

func (d *SqlLiteDb) Connect(filePath string) error {
	var err error

	d.db, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return err
	}

	return nil
}

// Links
func (d *SqlLiteDb) GetLinksForPart(int64) ([]Link, error) {
	return nil, nil
}

func (d *SqlLiteDb) GetLinksForKit(int64) ([]Link, error) {
	return nil, nil
}

func (d *SqlLiteDb) AddLinkToPart(string) (Link, error) {
	return Link{}, nil
}

func (d *SqlLiteDb) AddLinkToKit(string) (Link, error) {
	return Link{}, nil
}

func (d *SqlLiteDb) RemoveLinkFromPart(int64, int64) error {
	return nil
}

func (d *SqlLiteDb) RemoveLinkFromKit(int64, int64) error {
	return nil
}

// parts
func (d *SqlLiteDb) GetPart(partID int64) (Part, error) {
	const stmt string = "SELECT id, kind, name from parts where id = ? limit 1;"
	var part = Part{}

	row := d.db.QueryRow(stmt, partID)
	err := row.Scan(&part.ID, &part.Kind, &part.Name)
	if err != nil {
		return part, err
	}

	// todo: get part links

	return Part{}, nil
}

func (d *SqlLiteDb) AddPart(p Part) (Part, error) {
	const stmt string = `
		insert into parts(kind, name)
		values(?, ?)
	`
	newPart := p

	res, err := d.db.Exec(stmt, p.Kind, p.Name)
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
	res, err := d.db.Exec(stmt, p.Kind, p.Name, p.ID)
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

func (d *SqlLiteDb) DeletePart(partID int64) error {
	const stmt string = `
		delete from parts where id = ?
	`

	res, err := d.db.Exec(stmt, partID)
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

// kit
func (d *SqlLiteDb) GetKit(int64) (Kit, error) {
	return Kit{}, nil
}

func (d *SqlLiteDb) PutKit(Kit) (Kit, error) {
	return Kit{}, nil
}

func (d *SqlLiteDb) UpdateKit(Kit) (Kit, error) {
	return Kit{}, nil
}

func (d *SqlLiteDb) DeleteKit(int64) error {
	return nil
}

func (d *SqlLiteDb) GetKitParts(kitID uint64) ([]Part, error) {

	return []Part{}, nil
}

func (d *SqlLiteDb) AddPartToKit(p Part) (Part, error) {
	return p, nil
}

func (d *SqlLiteDb) RemovePartFromKit(p Part) (Part, error) {
	return p, nil
}
