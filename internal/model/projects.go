package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/sewiti/munit-backend/internal/id"
)

var mockProjects = []Project{
	{
		ID:          "2a43ea0b",
		Name:        "H3",
		Description: "Ethan & Hila",
		Created:     time.Now(),
		Modified:    time.Now(),
	},
	{
		ID:          "1cee6f88",
		Name:        "Tom Scott's project",
		Description: "Brittisssshhh",
		Created:     time.Now(),
		Modified:    time.Now(),
	},
	{
		ID:          "8666e20b",
		Name:        "Wirtual",
		Description: "How to pull a Wirtual",
		Created:     time.Now(),
		Modified:    time.Now(),
	},
}

type Project struct {
	ID          id.ID     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

var ErrProjectNotFound = fmt.Errorf("project %w", ErrNotFound)

func (p Project) valid() error {
	const maxName = 72
	if len(p.ID) == 0 {
		return errors.New("project uuid is empty")
	}
	if p.Name == "" {
		return errors.New("project name is empty")
	}
	if len(p.Name) > maxName {
		return fmt.Errorf("project name is too long (max %d chars)", maxName)
	}
	return nil
}

func GetProject(project id.ID) (Project, error) {
	for _, p := range mockProjects {
		if p.ID == project {
			return p, nil
		}
	}
	return Project{}, ErrProjectNotFound
}

func GetAllProjects() ([]Project, error) {
	return mockProjects, nil
}

func AddProject(p *Project) error {
	if err := p.valid(); err != nil {
		return err
	}
	now := time.Now()
	p.Created = now
	p.Modified = now
	mockProjects = append(mockProjects, *p)
	return nil
}

func EditProject(p *Project) error {
	if err := p.valid(); err != nil {
		return err
	}
	for i, mp := range mockProjects {
		if mp.ID == p.ID {
			p.Created = mp.Created
			p.Modified = time.Now()
			mockProjects[i] = *p
			return nil
		}
	}
	return ErrProjectNotFound
}

func DeleteProject(project id.ID) error {
	for i := range mockProjects {
		if mockProjects[i].ID == project {
			copy(mockProjects[i:], mockProjects[i+1:])
			mockProjects = mockProjects[:len(mockProjects)-1]
			return deleteAllCommits(project)
		}
	}
	return ErrProjectNotFound
}
