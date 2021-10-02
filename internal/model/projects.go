package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var mockProjects = []Project{
	{
		UUID:        uuid.MustParse("2a43ea0b-6c12-4aeb-82e4-d5a4362d92fc"),
		Name:        "H3",
		Description: "Ethan & Hila",
		Created:     time.Now(),
		Modified:    time.Now(),
	},
	{
		UUID:        uuid.MustParse("1cee6f88-ff6d-4a46-8f4f-fb752b36cafa"),
		Name:        "Tom Scott's project",
		Description: "Brittisssshhh",
		Created:     time.Now(),
		Modified:    time.Now(),
	},
	{
		UUID:        uuid.MustParse("8666e20b-db30-456e-95ce-8e87b3a58a40"),
		Name:        "Wirtual",
		Description: "How to pull a Wirtual",
		Created:     time.Now(),
		Modified:    time.Now(),
	},
}

type Project struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

var ErrProjectNotFound = fmt.Errorf("project %w", ErrNotFound)

func (p Project) Valid() error {
	const maxName = 72
	if len(p.UUID) == 0 {
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

func GetProject(project uuid.UUID) (Project, error) {
	for _, p := range mockProjects {
		if p.UUID == project {
			return p, nil
		}
	}
	return Project{}, ErrProjectNotFound
}

func GetAllProjects() ([]Project, error) {
	return mockProjects, nil
}

func AddProject(p *Project) error {
	if err := p.Valid(); err != nil {
		return err
	}
	now := time.Now()
	p.Created = now
	p.Modified = now
	mockProjects = append(mockProjects, *p)
	return nil
}

func EditProject(p *Project) error {
	if err := p.Valid(); err != nil {
		return err
	}
	for i, mp := range mockProjects {
		if mp.UUID == p.UUID {
			p.Created = mp.Created
			p.Modified = time.Now()
			mockProjects[i] = *p
			return nil
		}
	}
	return ErrProjectNotFound
}

func DeleteProject(project uuid.UUID) error {
	for i := range mockProjects {
		if mockProjects[i].UUID == project {
			copy(mockProjects[i:], mockProjects[i+1:])
			mockProjects = mockProjects[:len(mockProjects)-1]
			return deleteAllCommits(project)
		}
	}
	return ErrProjectNotFound
}
