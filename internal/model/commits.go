package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var mockCommits = []Commit{
	{
		UUID:     uuid.MustParse("bd2267ec-901e-48db-b7fa-62e78198e73b"),
		Title:    "Project setup",
		Message:  "",
		Created:  time.Now(),
		Modified: time.Now(),

		Project: uuid.MustParse("2a43ea0b-6c12-4aeb-82e4-d5a4362d92fc"),
	},
	{
		UUID:     uuid.MustParse("824e8770-c7a9-44d8-8dea-90e0d1a35be4"),
		Title:    "Initial recording",
		Message:  "",
		Created:  time.Now(),
		Modified: time.Now(),

		Project: uuid.MustParse("2a43ea0b-6c12-4aeb-82e4-d5a4362d92fc"),
	},
}

type Commit struct {
	UUID     uuid.UUID `json:"uuid"`
	Title    string    `json:"title"`
	Message  string    `json:"message"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`

	Project uuid.UUID `json:"project"`
}

var ErrCommitNotFound = fmt.Errorf("commit %w", ErrNotFound)

func (c Commit) Valid() error {
	const maxTitle = 72
	if len(c.UUID) == 0 {
		return errors.New("project uuid is empty")
	}
	if c.Title == "" {
		return errors.New("commit title is empty")
	}
	if len(c.Title) > maxTitle {
		return fmt.Errorf("commit title is too long (max %d chars)", maxTitle)
	}
	_, err := GetProject(c.Project)
	return err
}

func GetCommit(project, commit uuid.UUID) (Commit, error) {
	if _, err := GetProject(project); err != nil {
		return Commit{}, err
	}
	for _, c := range mockCommits {
		if c.Project == project && c.UUID == commit {
			return c, nil
		}
	}
	return Commit{}, ErrCommitNotFound
}

func getCommit(commit uuid.UUID) (Commit, error) {
	for _, c := range mockCommits {
		if c.UUID == commit {
			return c, nil
		}
	}
	return Commit{}, ErrCommitNotFound
}

func GetAllCommits(project uuid.UUID) ([]Commit, error) {
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
	if err := c.Valid(); err != nil {
		return err
	}
	now := time.Now()
	c.Created = now
	c.Modified = now
	mockCommits = append(mockCommits, *c)
	return nil
}

func EditCommit(c *Commit) error {
	if err := c.Valid(); err != nil {
		return err
	}
	for i, mc := range mockCommits {
		if mc.UUID == c.UUID {
			c.Created = mc.Created
			c.Modified = time.Now()
			mockCommits[i] = *c
			return nil
		}
	}
	return ErrCommitNotFound
}

func DeleteCommit(project, commit uuid.UUID) error {
	for i, c := range mockCommits {
		if c.Project == project && c.UUID == commit {
			copy(mockCommits[i:], mockCommits[i+1:])
			mockCommits = mockCommits[:len(mockCommits)-1]
			return nil
		}
	}
	return ErrCommitNotFound
}

func DeleteAllCommits(project uuid.UUID) error {
	if _, err := GetProject(project); err != nil {
		return err
	}
	return deleteAllCommits(project)
}

func deleteAllCommits(project uuid.UUID) error {
	for i := len(mockCommits) - 1; i >= 0; i-- {
		if mockCommits[i].Project == project {
			copy(mockCommits[i:], mockCommits[i+1:])
			mockCommits = mockCommits[:len(mockCommits)-1]
		}
	}
	return nil
}
