package sqlite

func (s *SQLiteDriver) setupStatements() (err error) {
	if s.sqlChunkCreate, err = s.db.Prepare(
		`INSERT INTO ` + s.tableChunks + ` (uuid, hash, content, referenced, size) VALUES(?,?,?,1,?)`); err != nil {
		return err
	}

	if s.sqlChunkCreateAttachment, err = s.db.Prepare(`INSERT INTO ` + s.tableChunkAttachments + `(file, chunk, start, end) VALUES(?,?,?,?)`); err != nil {
		return err
	}

	if s.sqlChunkRetrieve, err = s.db.Prepare(`SELECT hash, content FROM ` + s.tableChunks); err != nil {
		return err
	}

	if s.sqlChunkAttachmentList, err = s.db.Prepare(`SELECT chunk, start, end FROM ` + s.tableChunkAttachments + ` WHERE file=?`); err != nil {
		return err
	}

	if s.sqlChunkDelete, err = s.db.Prepare(
		`DELETE FROM ` + s.tableChunks + ` WHERE uuid=?`); err != nil {
		return err
	}

	if s.sqlFileCreate, err = s.db.Prepare(`INSERT INTO ` + s.tableFiles + ` (uuid, hash, path, title, description, content_type, status, size, created_at, updated_at, deleted_at) VALUES(?,?,?,?,?,?,?,0,?,?,0)`); err != nil {
		return err
	}

	if s.sqlFileRetrieve, err = s.db.Prepare(`SELECT hash, path, title, description, content_type, status, size, created_at, updated_at, deleted_at FROM ` + s.tableFiles + ` WHERE uuid=?`); err != nil {
		return err
	}

	if s.sqlFileUpdate, err = s.db.Prepare(`UPDATE ` + s.tableFiles + ` SET
		path=?,
	    title=?,
		description=?,
		content_type=?,
		status=?,
		size=?,
		updated_at=?
		WHERE uuid=?
	 `); err != nil {
		return err
	}

	if s.sqlFileDelete, err = s.db.Prepare(
		`DELETE FROM ` + s.tableFiles + ` WHERE uuid=?`); err != nil {
		return err
	}

	if s.sqlChunkStorageUsage, err = s.db.Prepare(`SELECT SUM("size") FROM ` + s.tableChunks); err != nil {
		return err
	}

	return nil
}
