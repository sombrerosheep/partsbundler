package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sombrerosheep/partsbundler/pkg/core"
)

type isqlitedb interface {
	Connect() error
	Close() error

	GetPart(partId int64) (core.Part, error)
	GetAllParts() ([]core.Part, error)
	GetPartLinks(partId int64) ([]core.Link, error)
	AddLinkToPart(link string, partId int64) (int64, error)
	RemoveLinkFromPart(linkId, partId int64) error
	CreatePart(name string, kind core.PartType) (int64, error)
	RemovePart(partId int64) error

	GetKit(kitId int64) (core.Kit, error)
	GetKitPartUsage(partId int64) ([]int64, error)
	GetKitPartsForKit(kitId int64) ([]kitPartRef, error)
	GetAllKits() ([]core.Kit, error)
	AddPartToKit(partId, kitId int64, quantity uint64) error
	UpdatePartQuantity(partId, kitId int64, quantity uint64) error
	RemovePartFromKit(partId, kitId int64) error
	GetKitLinks(kitId int64) ([]core.Link, error)
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

func (db sqlitedb) GetPart(partId int64) (core.Part, error) {
	const query string = `
		select id, name, kind from parts
			where id = ?
	`
	part := core.Part{}

	row := db.db.QueryRow(query, partId)
	err := row.Scan(&part.ID, &part.Name, &part.Kind)

	if err != nil && err == sql.ErrNoRows {
		return part, core.PartNotFound{PartID: partId}
	}

	return part, err
}

func (db sqlitedb) GetAllParts() ([]core.Part, error) {
	const query string = `
		select id, name, kind from parts
	`
	parts := []core.Part{}

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		part := core.Part{}

		// handle no rows err
		err := rows.Scan(&part.ID, &part.Name, &part.Kind)
		if err != nil {
			if err == sql.ErrNoRows {
				return parts, nil
			}

			return nil, err
		}

		parts = append(parts, part)
	}

	return parts, nil
}

func (db sqlitedb) GetPartLinks(partId int64) ([]core.Link, error) {
	const query string = `
		select id, link from partlinks
			where partId = ?
	`
	links := []core.Link{}

	_, err := db.GetPart(partId)
	if err != nil {
		return nil, err
	}

	rows, err := db.db.Query(query, partId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		link := core.Link{}

		err := rows.Scan(&link.ID, &link.URL)
		if err != nil {
			if err == sql.ErrNoRows {
				return []core.Link{}, nil
			}

			return nil, err
		}

		links = append(links, link)
	}

	return links, nil
}

// todo: check other funcs for "does <part/kit> exist"
func (db sqlitedb) AddLinkToPart(link string, partId int64) (int64, error) {
	const stmt string = `
		insert into partlinks(partId, link)
			values(?, ?)
	`

	_, err := db.GetPart(partId)
	if err != nil {
		return -1, err
	}

	res, err := db.db.Exec(stmt, partId, link)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (db sqlitedb) RemoveLinkFromPart(linkId, partId int64) error {
	const stmt string = `
		delete from partlinks
			where id = ? and partId = ?
	`

	_, err := db.GetPart(partId)
	if err != nil {
		return err
	}

	_, err = db.db.Exec(stmt, linkId, partId)

	return err
}

func (db sqlitedb) CreatePart(name string, kind core.PartType) (int64, error) {
	const stmt string = `
		insert into parts(name, kind)
			values(?, ?)
	`

	if err := kind.IsValid(); err != nil {
		return -1, err
	}

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

func (db sqlitedb) RemovePart(partId int64) error {
	const stmt string = `
		delete from parts where id = ?
	`

	_, err := db.GetPart(partId)
	if err != nil {
		return err
	}

	_, err = db.db.Exec(stmt, partId)

	return err
}

func (db sqlitedb) GetKit(kitId int64) (core.Kit, error) {
	const query string = `
		select id, name, schematic, diagram from kits
			where id = ?
	`
	kit := core.Kit{}

	row := db.db.QueryRow(query, kitId)

	err := row.Scan(&kit.ID, &kit.Name, &kit.Schematic, &kit.Diagram)

	if err != nil && err == sql.ErrNoRows {
		return kit, core.KitNotFound{KitID: kitId}
	}

	return kit, err
}

func (db sqlitedb) GetKitPartUsage(partId int64) ([]int64, error) {
	const query string = `
		select kitId from kitparts
			where partId = ?
	`

	_, err := db.GetPart(partId)
	if err != nil {
		return nil, err
	}

	rows, err := db.db.Query(query, partId)
	if err != nil {
		return nil, err
	}

	// todo: will .Next be true if there are no rows?
	//       add a test to find out!
	ids := []int64{}
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)

		if err != nil {
			if err == sql.ErrNoRows {
				return []int64{}, nil
			}

			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
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

	_, err := db.GetKit(kitId)
	if err != nil {
		return nil, err
	}

	parts := []kitPartRef{}

	rows, err := db.db.Query(query, kitId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		part := kitPartRef{}
		err = rows.Scan(&part.kitId, &part.partId, &part.quantity)
		if err != nil {
			if err == sql.ErrNoRows {
				return []kitPartRef{}, nil
			}

			return nil, err
		}

		parts = append(parts, part)
	}

	return parts, nil
}

func (db sqlitedb) GetAllKits() ([]core.Kit, error) {
	const query string = `
		select id, name, schematic, diagram from kits
	`

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}

	kits := []core.Kit{}
	for rows.Next() {
		kit := core.Kit{}

		err := rows.Scan(&kit.ID, &kit.Name, &kit.Schematic, &kit.Diagram)
		if err != nil {
			if err == sql.ErrNoRows {
				return []core.Kit{}, nil
			}

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

	_, err := db.GetKit(kitId)
	if err != nil {
		return err
	}

	_, err = db.GetPart(partId)
	if err != nil {
		return err
	}

	_, err = db.db.Exec(stmt, partId, kitId, quantity)

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

func (db sqlitedb) GetKitLinks(kitId int64) ([]core.Link, error) {
	const query string = `
		select id, link from kitlinks
			where kitId = ?
	`

	links := []core.Link{}

	rows, err := db.db.Query(query, kitId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		link := core.Link{}

		err = rows.Scan(&link.ID, &link.URL)
		if err != nil {
			if err == sql.ErrNoRows {
				return []core.Link{}, nil
			}

			return nil, err
		}

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
