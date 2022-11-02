package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type isqlitedb interface {
	Connect() error
	Close() error

	GetPart(partId int64) (Part, error)
	GetAllParts() ([]Part, error)
	GetPartLinks(partId int64) ([]Link, error)
	AddLinkToPart(link string, partId int64) (int64, error)
	RemoveLinkFromPart(linkId, partId int64) error
	CreatePart(name string, kind PartType) (int64, error)
	DeletePart(partId int64) error

	GetKit(kitId int64) (Kit, error)
	GetKitPartsForKit(kitId int64) ([]kitPartRef, error)
	GetAllKits() ([]Kit, error)
	AddPartToKit(partId, kitId int64, quantity uint64) error
	UpdatePartQuantity(partId, kitId int64, quantity uint64) error
	RemovePartFromKit(partId, kitId int64) error
	GetKitLinks(kitId int64) ([]Link, error)
	AddLinkToKit(link string, kitId int64) (int64, error)
	RemoveLinkFromKit(linkId, kitId int64) error
	CreateKit(name, schematic, diagram string) (int64, error)
	RemoveKit(kitId int64) error
}

func CreateSqliteDB(dbPath string) (isqlitedb, error) {
	sq := &sqlitedb{
		DBFilePath: dbPath,
	}

	err := sq.Connect()
	if err != nil {
		return nil, err
	}

	return sq, nil
}

type sqlitedb struct {
	db         *sql.DB
	DBFilePath string
}

func (db *sqlitedb) Connect() error {
	var err error

	db.db, err = sql.Open("sqlite3", db.DBFilePath)
	if err != nil {
		return err
	}

	return nil
}

func (db sqlitedb) Close() error {
	return db.db.Close()
}

func (db sqlitedb) GetPart(partId int64) (Part, error) {
	const query string = `
		select id, name, kind from parts
			where id = ?
	`
	part := Part{}

	row := db.db.QueryRow(query, partId)
	err := row.Scan(&part.ID, &part.Name, &part.Kind)

	if err == sql.ErrNoRows {
		err = nil
	}

	return part, err
}

func (db sqlitedb) GetAllParts() ([]Part, error) {
	const query string = `
		select id, name, kind from parts
	`
	parts := []Part{}

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		part := Part{}

		// handle no rows err
		rows.Scan(&part.ID, &part.Name, &part.Kind)

		parts = append(parts, part)
	}

	return parts, nil
}

func (db sqlitedb) GetPartLinks(partId int64) ([]Link, error) {
	const query string = `
		select id, link from partlinks
			where partId = ?
	`
	links := []Link{}

	rows, err := db.db.Query(query, partId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		link := Link{}

		err := rows.Scan(&link.ID, &link.URL)
		if err != nil {
			return nil, nil
		}

		links = append(links, link)
	}

	return links, nil
}

func (db sqlitedb) AddLinkToPart(link string, partId int64) (int64, error) {
	const stmt string = `
		insert into partlinks(partId, link)
			values(?, ?)
	`

	res, err := db.db.Exec(stmt, partId, link)
	if err != nil {
		return -1, nil
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, nil
	}

	return id, nil
}

func (db sqlitedb) RemoveLinkFromPart(linkId, partId int64) error {
	const stmt string = `
		delete from partlinks
			where id = ? and partId = ?
	`

	_, err := db.db.Exec(stmt, linkId, partId)

	return err
}

func (db sqlitedb) CreatePart(name string, kind PartType) (int64, error) {
	const stmt string = `
		insert into parts(name, kind)
			values(?, ?)
	`

	res, err := db.db.Exec(stmt, name, kind)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (db sqlitedb) DeletePart(partId int64) error {
	const stmt string = `
		delete from parts where id = ?
	`

	_, err := db.db.Exec(stmt, partId)

	return err
}

func (db sqlitedb) GetKit(kitId int64) (Kit, error) {
	const query string = `
		select id, name, schematic, diagram from kits
			where id = ?
	`
	kit := Kit{}

	row := db.db.QueryRow(query, kitId)

	err := row.Scan(&kit.ID, &kit.Name, &kit.Schematic, &kit.Diagram)

	if err == sql.ErrNoRows {
		err = nil
	}

	return kit, nil
}

type kitPartRef struct {
	kitId    int64
	partId   int64
	quantity uint64
}

func (db sqlitedb) GetKitPartsForKit(kitId int64) ([]kitPartRef, error) {
	const query string = `
		select kitId, partId, quantity from kitparts
			where kitId = ?
	`

	parts := []kitPartRef{}

	rows, err := db.db.Query(query, kitId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		part := kitPartRef{}
		rows.Scan(&part.kitId, &part.partId, &part.quantity)

		parts = append(parts, part)
	}

	return parts, nil
}

func (db sqlitedb) GetAllKits() ([]Kit, error) {
	const query string = `
		select id, name, schematic, diagram from kits
	`

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}

	kits := []Kit{}
	for rows.Next() {
		kit := Kit{}

		err := rows.Scan(&kit.ID, &kit.Name, &kit.Schematic, &kit.Diagram)
		if err != nil {
			return nil, err
		}

		kits = append(kits, kit)
	}

	return kits, nil
}

func (db sqlitedb) AddPartToKit(partId, kitId int64, quantity uint64) error {
	const stmt string = `
		insert into kitparts(partId, kitId, quantity)
			values(?, ?, ?)
	`

	_, err := db.db.Exec(stmt, partId, kitId, quantity)

	return err
}

func (db sqlitedb) UpdatePartQuantity(partId, kitId int64, quantity uint64) error {
	const stmt string = `
		update kitparts
			set quantity = ?
			where partId = ? and kitId = ?
	`
	_, err := db.db.Exec(stmt, quantity, partId, kitId)

	return err
}

func (db sqlitedb) RemovePartFromKit(partId, kitId int64) error {
	const stmt string = `
		delete from kitparts
			where partId = ? and kitId = ?
	`

	_, err := db.db.Exec(stmt, partId, kitId)

	return err
}

func (db sqlitedb) GetKitLinks(kitId int64) ([]Link, error) {
	const query string = `
		select id, link from kitlinks
			where kitId = ?
	`

	links := []Link{}

	rows, err := db.db.Query(query, kitId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		link := Link{}

		rows.Scan(&link.ID, &link.URL)

		links = append(links, link)
	}

	return links, nil
}

func (db sqlitedb) AddLinkToKit(link string, kitId int64) (int64, error) {
	const stmt string = `
		insert into kitlinks(kitId, link)
			values(?, ?)
	`

	res, err := db.db.Exec(stmt, kitId, link)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (db sqlitedb) RemoveLinkFromKit(linkId, kitId int64) error {
	const stmt string = `
		delete from kitlinks
			where kitId = ? and id = ?
	`

	_, err := db.db.Exec(stmt, kitId, linkId)

	return err
}

func (db sqlitedb) CreateKit(name, schematic, diagram string) (int64, error) {
	const stmt string = `
		insert into kits(name, schematic, diagram)
			values(?, ?, ?)
	`

	res, err := db.db.Exec(stmt, name, schematic, diagram)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (db sqlitedb) RemoveKit(kitId int64) error {
	const stmt string = `
		delete from kits
			where id = ?
	`

	_, err := db.db.Exec(stmt, kitId)

	return err
}
