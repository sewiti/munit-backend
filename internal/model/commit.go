package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sewiti/munit-backend/pkg/id"
)

const (
	commitSelect       = "SELECT id, title, message, created, modified, project_id, user_id FROM commit"
	commitSelectID     = commitSelect + " WHERE project_id=? AND id=?"
	commitSelectAllPID = commitSelect + " WHERE project_id=?"

	commitInsert = "INSERT INTO commit (id, title, message, created, modified, project_id, user_id) VALUES (?,?,?,?,?,?,?)"
	commitUpdate = "UPDATE commit SET title=?, message=?, modified=? WHERE project_id=? AND id=?"
)

type Commit struct {
	ID       id.ID     `json:"id"`
	Title    string    `json:"title"`
	Message  string    `json:"message"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`

	Project id.ID `json:"projectID"`
	User    id.ID `json:"userID"`
}

func (c *Commit) scan(sc scanner) (*Commit, error) {
	return c, sc.Scan(
		&c.ID,
		&c.Title,
		&c.Message,
		&c.Created,
		&c.Modified,
		&c.Project,
		&c.User,
	)
}

func (c *Commit) validate() error {
	const (
		maxTitle   = 72
		maxMessage = 1024
	)
	if c.Title == "" {
		return errors.New("commit: title is empty")
	}
	if len(c.Title) > maxTitle {
		return fmt.Errorf("commit: title is too long, max %d", maxTitle)
	}
	if len(c.Message) > maxMessage {
		return fmt.Errorf("commit: message is too long, max %d", maxMessage)
	}

	if err := c.Project.Validate(); err != nil {
		return fmt.Errorf("commit: project: %w", err)
	}
	if err := c.User.Validate(); err != nil {
		return fmt.Errorf("commit: user: %w", err)
	}
	return nil
}

func GetCommit(ctx context.Context, pid, cid id.ID) (*Commit, error) {
	row := db.QueryRowContext(ctx, commitSelectID, pid, cid)
	return new(Commit).scan(row)
}

func GetAllCommits(ctx context.Context, pid id.ID) ([]Commit, error) {
	rows, err := db.QueryContext(ctx, commitSelectAllPID, pid)
	if err != nil {
		return nil, err
	}

	commits := make([]Commit, 0)
	for rows.Next() {
		c, err := new(Commit).scan(rows)
		if err != nil {
			_ = rows.Close()
			return nil, err
		}
		commits = append(commits, *c)
	}

	if err = rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	if err = rows.Close(); err != nil {
		return nil, err
	}
	return commits, nil
}

func InsertCommit(ctx context.Context, c *Commit) error {
	if err := c.validate(); err != nil {
		return err
	}
	_, err := db.ExecContext(ctx, commitInsert,
		c.ID,
		c.Title,
		c.Message,
		c.Created,
		c.Modified,
		c.Project,
		c.User,
	)
	return err
}

func EditCommit(ctx context.Context, pid, cid id.ID, modifyFn func(*Commit) error) (*Commit, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, commitSelectID, pid, cid)
	c, err := new(Commit).scan(row)
	if err != nil {
		return nil, err
	}

	if err = modifyFn(c); err != nil {
		return nil, err
	}
	if err = c.validate(); err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, commitUpdate,
		c.Title,
		c.Message,
		c.Modified,
		pid,
		cid,
	)
	if err != nil {
		return nil, err
	}
	return c, tx.Commit()
}

func DeleteCommit(ctx context.Context, pid, cid id.ID) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM file WHERE project_id=? AND id=?", pid, cid)
	if err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, "DELETE FROM commit WHERE project_id=? AND id=?", pid, cid)
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
