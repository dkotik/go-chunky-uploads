package sqlite

import (
	"database/sql"
	"fmt"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
	_ "github.com/mattn/go-sqlite3" // database driver
)

var _ chunkyUploads.FileRepository = (*SQLiteDriver)(nil)

type SQLiteDriver struct {
	tableFiles            string
	tableChunks           string
	tableChunkAttachments string
	db                    *sql.DB

	sqlFileCreate            *sql.Stmt
	sqlFileRetrieve          *sql.Stmt
	sqlFileUpdate            *sql.Stmt
	sqlFileDelete            *sql.Stmt
	sqlChunkCreate           *sql.Stmt
	sqlChunkCreateAttachment *sql.Stmt
	sqlChunkAttachmentList   *sql.Stmt
	sqlChunkRetrieve         *sql.Stmt
	sqlChunkDelete           *sql.Stmt
	sqlChunkStorageUsage     *sql.Stmt
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
			hash BLOB NOT NULL,
			path TEXT NOT NULL,
		    title TEXT NOT NULL,
			description TEXT NOT NULL,
			content_type TEXT NOT NULL,
			status INTEGER NOT NULL,
			size INTEGER NOT NULL,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			deleted_at INTEGER NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS ` + driver.tableChunks + ` (
            uuid BLOB PRIMARY KEY,
			hash BLOB NOT NULL,
			content BLOB NOT NULL,
			referenced INTEGER NOT NULL,
			size INTEGER NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS ` + driver.tableChunkAttachments + ` (
            file BLOB KEY,
			chunk BLOB KEY,
			start INTEGER NOT NULL,
			end INTEGER NOT NULL
        )`,
	} {
		_, err = db.Exec(statement)
		if err != nil {
			db.Close()
			return nil, err
		}
	}

	if err = driver.setupStatements(); err != nil {
		db.Close()
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
	s.sqlChunkCreateAttachment.Close()
	s.sqlChunkAttachmentList.Close()
	s.sqlChunkRetrieve.Close()
	s.sqlChunkDelete.Close()
	s.sqlChunkStorageUsage.Close()
	return s.db.Close()
}
