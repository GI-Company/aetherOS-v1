package vfs

import (
	"fmt"
	"path"
	"sync"
)

type File struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Service struct {
	mu    sync.Mutex
	files map[string][]File
}

func NewService() *Service {
	files := make(map[string][]File)
	files["/"] = []File{}
	return &Service{
		files: files,
	}
}

func (s *Service) List(path string) ([]File, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if files, ok := s.files[path]; ok {
		return files, nil
	}
	return nil, fmt.Errorf("path not found: %s", path)
}

func (s *Service) CreateFile(filePath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir, name := path.Split(filePath)
	dir = path.Clean(dir)

	if _, ok := s.files[dir]; !ok {
		return fmt.Errorf("path not found: %s", dir)
	}

	s.files[dir] = append(s.files[dir], File{Name: name, Type: "file"})
	return nil
}

func (s *Service) CreateDirectory(dirPath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir, name := path.Split(dirPath)
	dir = path.Clean(dir)

	if _, ok := s.files[dir]; !ok {
		return fmt.Errorf("path not found: %s", dir)
	}

	s.files[dir] = append(s.files[dir], File{Name: name, Type: "folder"})
	s.files[dirPath] = []File{}
	return nil
}

func (s *Service) Delete(itemPath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir, name := path.Split(itemPath)
	dir = path.Clean(dir)

	if files, ok := s.files[dir]; ok {
		for i, file := range files {
			if file.Name == name {
				s.files[dir] = append(files[:i], files[i+1:]...)
				if file.Type == "folder" {
					delete(s.files, itemPath)
				}
				return nil
			}
		}
	}

	return fmt.Errorf("item not found: %s", itemPath)
}

func (s *Service) GetParent(itemPath string) string {
	if itemPath == "/" {
		return "/"
	}
	return path.Dir(itemPath)
}
