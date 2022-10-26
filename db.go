package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDb struct {
	db         *sql.DB
	DBFilePath string
}

func (d *SqliteDb) Connect() error {
	var err error

	d.db, err = sql.Open("sqlite3", d.DBFilePath)
	if err != nil {
		return err
	}

	return nil
}

func (d *SqliteDb) Close() error {
	return d.db.Close()
}

// Links
func (d *SqliteDb) GetLinksForPart(partId int64) ([]Link, error) {
	const query string = `
		select id, link from partlinks
			where partId = ?;
	`
	var links = []Link{}

	rows, err := d.db.Query(query, partId)
	if err != nil {
		return links, err
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

func (d *SqliteDb) GetLinksForKit(kitId int64) ([]Link, error) {
	const query string = `
		select id, link from kitlinks
			where kitId = ?;
	`
	var links = []Link{}

	rows, err := d.db.Query(query, kitId)
	if err != nil {
		return links, err
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

func (d *SqliteDb) AddLinkToPart(partId int64, url string) (Link, error) {
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

func (d *SqliteDb) AddLinkToKit(kitId int64, url string) (Link, error) {
	const stmt string = `
	insert into kitlinks(kitId, link)
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

func (d *SqliteDb) RemoveLinkFromPart(partId int64, linkId int64) error {
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

func (d *SqliteDb) RemoveLinkFromKit(kitId int64, linkId int64) error {
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
func (d *SqliteDb) GetParts() ([]Part, error) {
	const query string = `
		select id, kind, name from parts;
	`
	parts := []Part{}

	rows, err := d.db.Query(query)
	if err != nil {
		return parts, nil
	}

	for rows.Next() {
		part := Part{}

		err = rows.Scan(&part.ID, &part.Kind, &part.Name)
		if err != nil {
			break
		}

		parts = append(parts, part)
	}

	return parts, nil
}

func (d *SqliteDb) GetPart(partId int64) (Part, error) {
	const stmt string = "SELECT id, kind, name from parts where id = ? limit 1;"
	var part = Part{}

	row := d.db.QueryRow(stmt, partId)
	err := row.Scan(&part.ID, &part.Kind, &part.Name)
	if err != nil {
		return part, err
	}

	return part, nil
}

func (d *SqliteDb) AddPart(p Part) (Part, error) {
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

func (d *SqliteDb) UpdatePart(p Part) (Part, error) {
	const stmt string = `
		update parts
			set kind = ?,
					name = ?
		where id = ?;
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

func (d *SqliteDb) DeletePart(partId int64) error {
	const stmt string = `
		delete from parts where id = ?
	`

	res, err := d.db.Exec(stmt, partId)
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
func (d *SqliteDb) GetKits() ([]Kit, error) {
	const query string = `
		select id, name, schematic, diagram from kits;
	`
	var kits = []Kit{}

	rows, err := d.db.Query(query)
	if err != nil {
		return kits, err
	}

	for rows.Next() {
		kit := Kit{}

		err = rows.Scan(&kit.ID, &kit.Name, &kit.Schematic, &kit.Diagram)
		if err != nil {
			/* when we have an error, should we continue processing records?
			   if that happens, we should return the error as well.
			*/
			return kits, err
		}

		kits = append(kits, kit)
	}

	return kits, nil
}

func (d *SqliteDb) GetKit(kitId int64) (Kit, error) {
	const query string = `
		select id, name, schematic, diagram from kits where id = ?
	`
	var kit = Kit{}

	row := d.db.QueryRow(query, kitId)

	err := row.Scan(&kit.ID, &kit.Name, &kit.Schematic, &kit.Diagram)
	if err != nil {
		return kit, err
	}

	return kit, nil
}

func (d *SqliteDb) AddKit(kit Kit) (Kit, error) {
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

func (d *SqliteDb) UpdateKit(kit Kit) (Kit, error) {
	const stmt string = `
		update kits
			set name = ?,
				  schematic = ?,
				  diagram = ?
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

func (d *SqliteDb) DeleteKit(kitId int64) error {
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

func (d *SqliteDb) GetKitParts(kitId int64) ([]KitPart, error) {
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
			return parts, err
		}

		p, err := d.GetPart(pid)
		if err != nil {
			return parts, err
		}

		kp.Part = p
		kp.Quantity = quantity

		parts = append(parts, kp)
	}

	return parts, nil
}

func (d *SqliteDb) AddPartToKit(partId, kitId int64, quantity uint64) error {
	// does it already exist?
	const qexist string = "select id from kitparts where partId = ? and kitId = ?"

	exists, err := d.db.Query(qexist, partId, kitId)
	if err != nil {
		return err
	}

	if exists.Next() {
		return fmt.Errorf("Part already exists.")
	}

	// add it
	const stmt string = `
	insert into
		kitparts(partId, kitId, quantity)
		values(?, ?, ?)
	`
	res, err := d.db.Exec(stmt, partId, kitId, quantity)
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

func (d *SqliteDb) SetPartQuantityForKit(partId int64, kitId uint64, quantity int64) error {
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

func (d *SqliteDb) RemovePartFromKit(partId, kitId int64) error {
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
