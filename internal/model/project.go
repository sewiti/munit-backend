package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sewiti/munit-backend/pkg/id"
)

type Project struct {
	ID          id.ID     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`

	Owner        id.ID   `json:"ownerID"`
	Contributors []id.ID `json:"contributors"`
}

func (p *Project) scan(sc scanner) error {
	var uid *id.ID
	err := sc.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Created,
		&p.Modified,
		&p.Owner,
		&uid,
	)
	if uid != nil {
		p.Contributors = append(p.Contributors, *uid)
	}
	return err
}

func (p *Project) validate() error {
	const (
		maxName        = 72
		maxDescription = 1024
	)

	if err := p.ID.Validate(); err != nil {
		return fmt.Errorf("project: %w", err)
	}
	if err := p.Owner.Validate(); err != nil {
		return fmt.Errorf("project: owner: %w", err)
	}

	// Name
	if p.Name == "" {
		return errors.New("project: name is empty")
	}
	if len(p.Name) > maxName {
		return fmt.Errorf("project: name is too long, max %d", maxName)
	}

	// Description
	if len(p.Description) > maxDescription {
		return fmt.Errorf("project: description is too long, max %d", maxDescription)
	}

	// Contributors
	for _, c := range p.Contributors {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("project: contributor: %w", err)
		}
		if c == p.Owner {
			return fmt.Errorf("project: owner cannot be a contributor")
		}
	}
	return nil
}

func GetProject(ctx context.Context, pid id.ID) (*Project, error) {
	rows, err := db.QueryContext(ctx,
		"SELECT p.id, p.name, p.description, p.created, p.modified, p.owner_id, c.user_id "+
			"FROM project p LEFT JOIN contributor c ON p.id=c.project_id "+
			"WHERE p.id=?",
		pid,
	)
	if err != nil {
		return nil, err
	}
	return getProject(rows)
}

func getProject(rows *sql.Rows) (*Project, error) {
	defer rows.Close()

	p := new(Project)
	var err error
	for rows.Next() {
		if err = p.scan(rows); err != nil {
			_ = rows.Close()
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	if err = p.ID.Validate(); err != nil {
		return nil, ErrNotFound
	}
	if p.Contributors == nil {
		p.Contributors = make([]id.ID, 0)
	}
	return p, nil
}

func GetAllProjects(ctx context.Context, uid id.ID) ([]Project, error) {
	// TODO: Refactor
	rows, err := db.QueryContext(ctx,
		"SELECT p.id, p.name, p.description, p.created, p.modified, p.owner_id, c.user_id "+
			"FROM project p LEFT JOIN contributor c ON p.id=c.project_id "+
			"WHERE p.owner_id=? OR c.user_id=?",
		uid, uid,
	)
	if err != nil {
		return nil, err
	}

	projects := make([]Project, 0)
	var p Project
	scan := Project{
		Contributors: make([]id.ID, 0),
	}

	for rows.Next() {
		if err = scan.scan(rows); err != nil {
			_ = rows.Close()
			return nil, err
		}

		if p.ID != "" && p.ID != scan.ID {
			scan.Contributors = scan.Contributors[len(p.Contributors):]
			projects = append(projects, p)
		}
		p = scan
	}
	if p.ID != "" {
		projects = append(projects, p)
	}

	if err = rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	if err = rows.Close(); err != nil {
		return nil, err
	}
	return projects, nil
}

func InsertProject(ctx context.Context, p *Project) error {
	if err := p.validate(); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Project
	_, err = tx.ExecContext(ctx,
		"INSERT INTO project (id, name, description, created, modified, owner_id) VALUES (?,?,?,?,?,?)",
		p.ID,
		p.Name,
		p.Description,
		p.Created,
		p.Modified,
		p.Owner,
	)
	if err != nil {
		return err
	}

	// Contributors
	if len(p.Contributors) > 0 {
		var query strings.Builder
		args := make([]interface{}, 0, 2*len(p.Contributors))
		for i, cid := range p.Contributors {
			if i == 0 {
				query.WriteString("INSERT INTO contributor (project_id, user_id) VALUES")
			} else {
				query.WriteString(",")
			}
			query.WriteString(" (?,?)")
			args = append(args, p.ID, cid)
		}

		_, err = tx.ExecContext(ctx, query.String(), args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func UpdateProject(ctx context.Context, pid id.ID, modifyFn func(*Project) error) (*Project, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx,
		"SELECT p.id, p.name, p.description, p.created, p.modified, p.owner_id, c.user_id "+
			"FROM project p LEFT JOIN contributor c ON p.id=c.project_id WHERE p.id=?",
		pid,
	)
	if err != nil {
		return nil, err
	}
	p, err := getProject(rows)
	if err != nil {
		return nil, err
	}
	if err = modifyFn(p); err != nil {
		return nil, err
	}
	if err = p.validate(); err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE project SET name=?, description=?, modified=?, owner_id=? WHERE id=?",
		p.Name,
		p.Description,
		p.Modified,
		p.Owner,
		pid,
	)
	if err != nil {
		return nil, err
	}

	// Contributors
	_, err = tx.ExecContext(ctx, "DELETE FROM contributor WHERE project_id=?", pid)
	if err != nil {
		return nil, err
	}
	if len(p.Contributors) > 0 {
		var query strings.Builder
		args := make([]interface{}, 0, 2*len(p.Contributors))
		for i, cid := range p.Contributors {
			if i == 0 {
				query.WriteString("INSERT INTO contributor (project_id, user_id) VALUES")
			} else {
				query.WriteString(",")
			}
			query.WriteString(" (?,?)")
			args = append(args, p.ID, cid)
		}

		_, err = tx.ExecContext(ctx, query.String(), args...)
		if err != nil {
			return nil, err
		}
	}

	return p, tx.Commit()
}

func DeleteProject(ctx context.Context, pid id.ID) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM file WHERE project_id=?", pid)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM commit WHERE project_id=?", pid)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM contributor WHERE project_id=?", pid)
	if err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, "DELETE FROM project WHERE id=?", pid)
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
