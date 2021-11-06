package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // database driver
)

type SQLiteDriver struct {
	tableFiles            string
	tableChunks           string
	tableChunkAttachments string
	db                    *sql.DB

	sqlFileCreate    *sql.Stmt
	sqlFileRetrieve  *sql.Stmt
	sqlFileUpdate    *sql.Stmt
	sqlFileDelete    *sql.Stmt
	sqlChunkCreate   *sql.Stmt
	sqlChunkRetrieve *sql.Stmt
	sqlChunkUpdate   *sql.Stmt
	sqlChunkDelete   *sql.Stmt
}

func NewSQLiteDriver(tablePrefix, path string) (*SQLiteDriver, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared", path))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	driver := &SQLiteDriver{
		tableFiles:            tablePrefix + "_files",
		tableChunks:           tablePrefix + "_chunks",
		tableChunkAttachments: tablePrefix + "_chunk_attachments",
		db:                    db,
	}

	for _, statement := range []string{
		`CREATE TABLE IF NOT EXISTS ` + driver.tableFiles + ` (
            uuid BLOB PRIMARY KEY,
			hash BLOB,
			path TEXT,
		    title TEXT,
			description TEXT,
			content_type TEXT,
			status INTEGER,
			size INTEGER,
			created_at INTEGER,
			updated_at INTEGER,
			deleted_at INTEGER
        )`,
		`CREATE TABLE IF NOT EXISTS ` + driver.tableChunks + ` (
            uuid BLOB PRIMARY KEY,
			hash BLOB,
			content BLOB
        )`,
		`CREATE TABLE IF NOT EXISTS ` + driver.tableChunkAttachments + ` (
            file BLOB,
			chunk BLOB,
			start INTEGER,
			end INTEGER
        )`,
	} {
		_, err = db.Exec(statement)
		if err != nil {
			return nil, err
		}
	}

	if err = driver.setupStatements(); err != nil {
		return nil, err
	}

	return driver, nil
}

func (s *SQLiteDriver) Close() error {
	s.sqlFileCreate.Close()
	s.sqlFileRetrieve.Close()
	s.sqlFileUpdate.Close()
	s.sqlFileDelete.Close()
	s.sqlChunkCreate.Close()
	s.sqlChunkRetrieve.Close()
	s.sqlChunkUpdate.Close()
	s.sqlChunkDelete.Close()
	return s.db.Close()
}
