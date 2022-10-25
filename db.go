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
func (d *SqlLiteDb) GetLinksForPart(partID int64) ([]Link, error) {
	const query string = `
		select id, link from partlinks
			where partId = ?;
	`
	var links = []Link{}

	rows, err := d.db.Query(query, partID)
	if err != nil {
		return links, nil
	}
	defer rows.Close()

	for rows.Next() {
		var l = Link{}

		err = rows.Scan(&l.ID, &l.URL)
		if err != nil {
			break
		}

		links = append(links, l)
	}

	return links, nil
}

func (d *SqlLiteDb) GetLinksForKit(kitID int64) ([]Link, error) {
	const query string = `
		select id, link from kitlinks
			where partId = ?;
	`
	var links = []Link{}

	rows, err := d.db.Query(query, kitID)
	if err != nil {
		return links, nil
	}
	defer rows.Close()

	for rows.Next() {
		var l = Link{}

		err = rows.Scan(&l.ID, &l.URL)
		if err != nil {
			break
		}

		links = append(links, l)
	}

	return links, nil
}

func (d *SqlLiteDb) AddLinkToPart(partId int64, url string) (Link, error) {
	const stmt string = `
		insert into partlinks(partId, link)
			values(?, ?);
	`
	var link = Link{}

	res, err := d.db.Exec(stmt, partId, url)
	if err != nil {
		return link, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return link, err
	}

	link.ID = id
	link.URL = url

	return link, nil
}

func (d *SqlLiteDb) AddLinkToKit(kitId int64, url string) (Link, error) {
	const stmt string = `
	insert into kitlinks(partId, link)
		values(?, ?);
`
	var link = Link{}

	res, err := d.db.Exec(stmt, kitId, url)
	if err != nil {
		return link, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return link, err
	}

	link.ID = id
	link.URL = url

	return link, nil
}

func (d *SqlLiteDb) RemoveLinkFromPart(partId int64, linkId int64) error {
	const stmt string = `
		delete from partlinks where id = ? and partId = ?
	`

	res, err := d.db.Exec(stmt, linkId, partId)
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

	return nil
}

func (d *SqlLiteDb) RemoveLinkFromKit(kitId int64, linkId int64) error {
	const stmt string = `
		delete from kitlinks where id = ? and kitId = ?
	`

	res, err := d.db.Exec(stmt, linkId, kitId)
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

	return part, nil
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

	return newPart, nil
}

func (d *SqlLiteDb) UpdatePart(p Part) (Part, error) {
	const stmt string = `
		update parts
			set kind = ?,
			set name = ?,
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

	return nil
}

// kit
func (d *SqlLiteDb) GetKit(kitId int64) (Kit, error) {
	const query string = `
		select id, name, schematic, diagram from kits where kitId = ?
	`
	var kit = Kit{}

	rows, err := d.db.Query(query, kitId)
	if err != nil {
		return kit, nil
	}
	defer rows.Close()

	err = rows.Scan(&kit.ID, &kit.Name, &kit.Schematic, &kit.Diagram)
	if err != nil {
		return kit, err
	}

	return kit, nil
}

func (d *SqlLiteDb) AddKit(kit Kit) (Kit, error) {
	const stmt string = `
		insert into kits(name, schematic, diagram)
			values(?, ?, ?)
	`
	newKit := kit

	res, err := d.db.Exec(stmt, kit.Name, kit.Schematic, kit.Diagram)
	if err != nil {
		return kit, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return kit, err
	}

	newKit.ID = id

	return newKit, nil
}

func (d *SqlLiteDb) UpdateKit(kit Kit) (Kit, error) {
	const stmt string = `
		update kits
			set name = ?,
			set schematic = ?,
			set diagram = ?
		where id = ?
	`

	res, err := d.db.Exec(stmt, kit.Name, kit.Schematic, kit.Diagram, kit.ID)
	if err != nil {
		return kit, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return kit, err
	}

	if rows != 1 {
		return kit, fmt.Errorf("Expected 1 row to be affected but %d were", rows)
	}

	return kit, nil
}

func (d *SqlLiteDb) DeleteKit(kitId int64) error {
	const stmt string = `
		delete from parts where id = ?
	`

	res, err := d.db.Exec(stmt, kitId)
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

	return nil
}

func (d *SqlLiteDb) GetKitParts(kitId uint64) ([]KitPart, error) {
	const qkitparts string = `
		select partId, quantity from kitparts
			where kitId = ?
	`
	var parts = []KitPart{}

	rows, err := d.db.Query(qkitparts, kitId)
	if err != nil {
		return parts, err
	}
	defer rows.Close()

	for rows.Next() {
		var kp = KitPart{}
		var pid int64
		var quantity uint64

		err := rows.Scan(&pid, &quantity)
		if err != nil {
			return parts, nil
		}

		p, err := d.GetPart(pid)
		if err != nil {
			return parts, nil
		}

		kp.Part = p
		kp.Quantity = quantity

		parts = append(parts, kp)
	}

	return parts, nil
}

func (d *SqlLiteDb) AddPartToKit(partId, kitId int64, quanity uint64) error {
	// does it already exist?
	const qexist string = "select id from kitparts where partId = ? and kitId = ?"
	res, err := d.db.Exec(qexist)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return nil
	}

	if count > 0 {
		return fmt.Errorf("Part already exists.")
	}

	// add it
	const stmt string = `insert into kitparts(partId, kitId, quantity) values(?, ? ?)`
	res, err = d.db.Exec(stmt, partId, kitId, quanity)
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

	return nil
}

func (d *SqlLiteDb) SetPartQuantityForKit(partId int64, quantity uint64, kitId int64) error {
	const stmt string = `
		update kitparts
			set(quantity = ?)
		where partId = ? and kitId = ?
	`

	res, err := d.db.Exec(stmt, quantity, partId, kitId)
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

	return nil
}

func (d *SqlLiteDb) RemovePartFromKit(partId, kitId int64) error {
	const stmt string = `
		delete from kitparts
		where partId = ? and kitId = ?
	`

	res, err := d.db.Exec(stmt, partId, kitId)
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

	return nil
}
