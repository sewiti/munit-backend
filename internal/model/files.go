package model

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
)

var mockFiles = []File{
	{
		UUID:     uuid.MustParse("84ef2c38-24c3-47dc-9991-4ad1961c47bb"),
		Path:     "/project.als",
		Data:     []byte("project\n"),
		Created:  time.Now(),
		Modified: time.Now(),

		Commit: uuid.MustParse("bd2267ec-901e-48db-b7fa-62e78198e73b"),
	},
	{
		UUID:     uuid.MustParse("a9e08373-a4c7-49bb-b9fd-ac672aeb6cc1"),
		Path:     "/project info/project.cfg",
		Data:     []byte("ProjectInfo ProjectInfo\n"),
		Created:  time.Now(),
		Modified: time.Now(),

		Commit: uuid.MustParse("bd2267ec-901e-48db-b7fa-62e78198e73b"),
	},
}

type File struct {
	UUID     uuid.UUID `json:"uuid"`
	Path     string    `json:"path"`
	Data     []byte    `json:"data"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`

	Commit uuid.UUID `json:"commit"`
}

var ErrFileNotFound = fmt.Errorf("file %w", ErrNotFound)

func (f File) Valid() error {
	const maxPath = 1000
	if len(f.UUID) == 0 {
		return errors.New("project uuid is empty")
	}
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

func GetFile(commit, file uuid.UUID) (File, error) {
	if _, err := getCommit(commit); err != nil {
		return File{}, err
	}
	for _, f := range mockFiles {
		if f.Commit == commit && f.UUID == file {
			return f, nil
		}
	}
	return File{}, ErrFileNotFound
}

func GetAllFiles(commit uuid.UUID) ([]File, error) {
	if _, err := getCommit(commit); err != nil {
		return nil, err
	}
	var files []File
	for _, f := range mockFiles {
		if f.Commit == commit {
			files = append(files, f)
		}
	}
	return files, nil
}

func AddFile(f *File) error {
	if err := f.Valid(); err != nil {
		return err
	}
	now := time.Now()
	f.Created = now
	f.Modified = now
	mockFiles = append(mockFiles, *f)
	return nil
}

func EditFile(f *File) error {
	if err := f.Valid(); err != nil {
		return err
	}
	for i, mf := range mockFiles {
		if mf.UUID == f.UUID {
			f.Created = mf.Created
			f.Modified = time.Now()
			mockFiles[i] = *f
			return nil
		}
	}
	return ErrFileNotFound
}

func DeleteFile(commit, file uuid.UUID) error {
	for i, f := range mockFiles {
		if f.Commit == commit && f.UUID == file {
			copy(mockFiles[i:], mockFiles[i+1:])
			mockFiles = mockFiles[:len(mockFiles)-1]
			return nil
		}
	}
	return ErrFileNotFound
}

func DeleteAllFiles(commit uuid.UUID) error {
	if _, err := getCommit(commit); err != nil {
		return err
	}
	return deleteAllFiles(commit)
}

func deleteAllFiles(commit uuid.UUID) error {
	for i := len(mockFiles) - 1; i >= 0; i-- {
		if mockFiles[i].Commit == commit {
			copy(mockFiles[i:], mockFiles[i+1:])
			mockFiles = mockFiles[:len(mockFiles)-1]
		}
	}
	return nil
}
