package model

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/sewiti/munit-backend/pkg/id"
)

var mockFiles = []File{
	{
		ID:       1,
		Path:     "/project.als",
		Data:     []byte("project\n"),
		Created:  time.Now(),
		Modified: time.Now(),

		Commit:  1,
		Project: "2a43ea0b",
	},
	{
		ID:       2,
		Path:     "/project info/project.cfg",
		Data:     []byte("ProjectInfo ProjectInfo\n"),
		Created:  time.Now(),
		Modified: time.Now(),

		Commit:  1,
		Project: "2a43ea0b",
	},
}

type File struct {
	ID       int       `json:"id"`
	Path     string    `json:"path"`
	Data     []byte    `json:"data"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`

	Commit  int   `json:"commitID"`
	Project id.ID `json:"projectID"`
}

var ErrFileNotFound = fmt.Errorf("file %w", ErrNotFound)

func (f File) valid() error {
	const maxPath = 1000
	if !strings.HasPrefix(f.Path, "/") {
		return errors.New("file path: must start with /")
	}
	if len(f.Path) > maxPath {
		return fmt.Errorf("file path: too long: max %d", maxPath)
	}
	if _, file := path.Split(f.Path); file == "" {
		return errors.New("file path: empty file name")
	}
	_, err := getCommit(f.Commit)
	return err
}

func GetFile(commitID, fileID int) (File, error) {
	if _, err := getCommit(commitID); err != nil {
		return File{}, err
	}
	for _, f := range mockFiles {
		if f.Commit == commitID && f.ID == fileID {
			return f, nil
		}
	}
	return File{}, ErrFileNotFound
}

func GetAllFiles(commitID int) ([]File, error) {
	if _, err := getCommit(commitID); err != nil {
		return nil, err
	}
	var files []File
	for _, f := range mockFiles {
		if f.Commit == commitID {
			files = append(files, f)
		}
	}
	return files, nil
}

func AddFile(f *File) error {
	if err := f.valid(); err != nil {
		return err
	}
	now := time.Now()
	f.ID = mockFiles[len(mockFiles)-1].ID + 1
	f.Created = now
	f.Modified = now
	mockFiles = append(mockFiles, *f)
	return nil
}

func EditFile(f *File) error {
	if err := f.valid(); err != nil {
		return err
	}
	for i, mf := range mockFiles {
		if mf.ID == f.ID {
			f.Created = mf.Created
			f.Modified = time.Now()
			mockFiles[i] = *f
			return nil
		}
	}
	return ErrFileNotFound
}

func DeleteFile(commitID, fileID int) error {
	for i, f := range mockFiles {
		if f.Commit == commitID && f.ID == fileID {
			copy(mockFiles[i:], mockFiles[i+1:])
			mockFiles = mockFiles[:len(mockFiles)-1]
			return nil
		}
	}
	return ErrFileNotFound
}

func DeleteAllFiles(commitID int) error {
	if _, err := getCommit(commitID); err != nil {
		return err
	}
	return deleteAllFiles(commitID)
}

func deleteAllFiles(commitID int) error {
	for i := len(mockFiles) - 1; i >= 0; i-- {
		if mockFiles[i].Commit == commitID {
			copy(mockFiles[i:], mockFiles[i+1:])
			mockFiles = mockFiles[:len(mockFiles)-1]
		}
	}
	return nil
}
