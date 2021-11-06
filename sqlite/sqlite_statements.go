package sqlite

func (s *SQLiteDriver) setupStatements() (err error) {
	if s.sqlChunkCreate, err = s.db.Prepare(
		`INSERT INTO ` + s.tableChunks + ` (uuid, hash, content) VALUES(?,?,?)`); err != nil {
		return err
	}
	// ChunkRetrieve(context.Context, UUID) (*Chunk, error)

	if s.sqlChunkDelete, err = s.db.Prepare(
		`DELETE FROM ` + s.tableChunks + ` WHERE uuid=?`); err != nil {
		return err
	}

	// ChunkAttachmentList(context.Context, UUID) ([]*ChunkAttachment, error)

	if s.sqlFileDelete, err = s.db.Prepare(
		`DELETE FROM ` + s.tableFiles + ` WHERE uuid=?`); err != nil {
		return err
	}

	return nil
}
