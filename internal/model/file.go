package model

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/sewiti/munit-backend/pkg/id"
)

const (
	fileSelect      = "SELECT id, path, data, created, modified, commit_id, project_id FROM file"
	fileSelectID    = fileSelect + " WHERE project_id=? AND commit_id=? AND id=?"
	fileSelectAllID = fileSelect + " WHERE project_id=? AND commit_id=?"

	fileInsert = "INSERT INTO file (id, path, data, created, modified, commit_id, project_id) VALUES (?,?,?,?,?,?,?)"
	fileUpdate = "UPDATE file SET path=?, data=?, modified=? WHERE project_id=? AND commit_id=? AND id=?"
)

type File struct {
	ID       id.ID     `json:"id"`
	Path     string    `json:"path"`
	Data     []byte    `json:"data"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`

	Commit  id.ID `json:"commitID"`
	Project id.ID `json:"projectID"`
}

func (f *File) scan(sc scanner) (*File, error) {
	return f, sc.Scan(
		&f.ID,
		&f.Path,
		&f.Data,
		&f.Created,
		&f.Modified,
		&f.Commit,
		&f.Project,
	)
}

func (f *File) validate() error {
	const (
		maxPath = 256
	)
	if !strings.HasPrefix(f.Path, "/") {
		return errors.New("file: path: must start with /")
	}
	if len(f.Path) > maxPath {
		return fmt.Errorf("file: path: too long, max %d", maxPath)
	}
	if _, file := path.Split(f.Path); file == "" {
		return errors.New("file: path: empty file name")
	}
	return nil
}

func GetFile(ctx context.Context, pid, cid, fid id.ID) (*File, error) {
	row := db.QueryRowContext(ctx, fileSelectID, pid, cid, fid)
	return new(File).scan(row)
}

func GetAllFiles(ctx context.Context, pid, cid id.ID) ([]File, error) {
	rows, err := db.QueryContext(ctx, fileSelectAllID, pid, cid)
	if err != nil {
		return nil, err
	}

	files := make([]File, 0)
	for rows.Next() {
		c, err := new(File).scan(rows)
		if err != nil {
			_ = rows.Close()
			return nil, err
		}
		files = append(files, *c)
	}

	if err = rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	if err = rows.Close(); err != nil {
		return nil, err
	}
	return files, nil
}

func InsertFile(ctx context.Context, f *File) error {
	if err := f.validate(); err != nil {
		return err
	}
	_, err := db.ExecContext(ctx, fileInsert,
		f.ID,
		f.Path,
		f.Data,
		f.Created,
		f.Modified,
		f.Commit,
		f.Project,
	)
	return err
}

func UpdateFile(ctx context.Context, pid, cid, fid id.ID, modifyFn func(*File) error) (*File, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, fileSelectID, pid, cid, fid)
	f, err := new(File).scan(row)
	if err != nil {
		return nil, err
	}

	if err = modifyFn(f); err != nil {
		return nil, err
	}
	if err = f.validate(); err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, fileUpdate,
		f.Path,
		f.Data,
		f.Modified,
		pid,
		cid,
		fid,
	)
	if err != nil {
		return nil, err
	}
	return f, tx.Commit()
}

func DeleteFile(ctx context.Context, pid, cid, fid id.ID) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, "DELETE FROM file WHERE project_id=? AND commit_id=? AND id=?", pid, cid, fid)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return tx.Commit()
}
