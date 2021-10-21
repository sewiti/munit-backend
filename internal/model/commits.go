package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/sewiti/munit-backend/internal/id"
)

var mockCommits = []Commit{
	{
		ID:       1,
		Title:    "Project setup",
		Message:  "",
		Created:  time.Now(),
		Modified: time.Now(),

		Project: "2a43ea0b",
	},
	{
		ID:       2,
		Title:    "Initial recording",
		Message:  "",
		Created:  time.Now(),
		Modified: time.Now(),

		Project: "2a43ea0b",
	},
}

type Commit struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Message  string    `json:"message"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`

	Project id.ID `json:"projectID"`
}

var ErrCommitNotFound = fmt.Errorf("commit %w", ErrNotFound)

func (c Commit) valid() error {
	const maxTitle = 72
	if c.Title == "" {
		return errors.New("commit title is empty")
	}
	if len(c.Title) > maxTitle {
		return fmt.Errorf("commit title is too long (max %d chars)", maxTitle)
	}
	_, err := GetProject(c.Project)
	return err
}

func GetCommit(project id.ID, commitID int) (Commit, error) {
	if _, err := GetProject(project); err != nil {
		return Commit{}, err
	}
	for _, c := range mockCommits {
		if c.Project == project && c.ID == commitID {
			return c, nil
		}
	}
	return Commit{}, ErrCommitNotFound
}

func getCommit(commitID int) (Commit, error) {
	for _, c := range mockCommits {
		if c.ID == commitID {
			return c, nil
		}
	}
	return Commit{}, ErrCommitNotFound
}

func GetAllCommits(project id.ID) ([]Commit, error) {
	if _, err := GetProject(project); err != nil {
		return nil, err
	}
	var commits []Commit
	for _, c := range mockCommits {
		if c.Project == project {
			commits = append(commits, c)
		}
	}
	return commits, nil
}

func AddCommit(c *Commit) error {
	if err := c.valid(); err != nil {
		return err
	}
	now := time.Now()
	c.ID = mockCommits[len(mockCommits)-1].ID + 1
	c.Created = now
	c.Modified = now
	mockCommits = append(mockCommits, *c)
	return nil
}

func EditCommit(c *Commit) error {
	if err := c.valid(); err != nil {
		return err
	}
	for i, mc := range mockCommits {
		if mc.ID == c.ID {
			c.Created = mc.Created
			c.Modified = time.Now()
			mockCommits[i] = *c
			return nil
		}
	}
	return ErrCommitNotFound
}

func DeleteCommit(project id.ID, commitID int) error {
	for i, c := range mockCommits {
		if c.Project == project && c.ID == commitID {
			copy(mockCommits[i:], mockCommits[i+1:])
			mockCommits = mockCommits[:len(mockCommits)-1]
			return deleteAllFiles(c.ID)
		}
	}
	return ErrCommitNotFound
}

func DeleteAllCommits(project id.ID) error {
	if _, err := GetProject(project); err != nil {
		return err
	}
	return deleteAllCommits(project)
}

func deleteAllCommits(project id.ID) error {
	for i := len(mockCommits) - 1; i >= 0; i-- {
		if mockCommits[i].Project == project {
			if err := deleteAllFiles(mockCommits[i].ID); err != nil {
				return err
			}
			copy(mockCommits[i:], mockCommits[i+1:])
			mockCommits = mockCommits[:len(mockCommits)-1]
		}
	}
	return nil
}
