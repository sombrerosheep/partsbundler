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
func GetLinksForPart(int64) ([]Link, error) {
	return nil, nil
}

func GetLinksForKit(int64) ([]Link, error) {
	return nil, nil
}

func AddLinkToPart(string) (Link, error) {
	return Link{}, nil
}

func AddLinkToKit(string) (Link, error) {
	return Link{}, nil
}


// parts
func (d *SqlLiteDb) GetLinksForPart(partID int64) ([]string, error) {
	const stmt string = `
		select * from partlinks
		where partId = ?
	`
	var links = []string{}

	rows, err := d.db.Query(stmt, partID)
	if err != nil {
		return links, err
	}

	defer rows.Close()

	for rows.Next() {
		var link string
		err := rows.Scan(nil, &link)
		if err != nil {
			return links, err
		}

		links = append(links, link)
	}

	return links, nil
}

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

func (d *SqlLiteDb) PutPart(p Part) (Part, error) {
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
func GetKit(int64) (Kit, error) {
	return Kit{}, nil
}

func PutKit(Kit) (Kit, error) {
	return Kit{}, nil
}

func UpdateKit(Kit) (Kit, error) {
	return Kit{}, nil
}

func DeleteKit(int64) error {
	return nil
}