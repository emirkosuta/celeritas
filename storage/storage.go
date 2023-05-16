package storage

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Storage struct {
	BaseDir string
}

func (s Storage) createFile(filename string, data []byte, mode os.FileMode) error {
	filePath := filepath.Join(s.BaseDir, filename)
	return os.WriteFile(filePath, data, mode)
}

// CreateFile creates a new file with the given filename and data.
func (s Storage) CreateFile(filename string, data []byte) error {
	return s.createFile(filename, data, 0644)
}

// ReadFile reads the content of the given filename.
func (s Storage) ReadFile(filename string) ([]byte, error) {
	filePath := filepath.Join(s.BaseDir, filename)
	return os.ReadFile(filePath)
}

// UpdateFile updates the content of the given filename with the provided data.
func (s Storage) UpdateFile(filename string, data []byte) error {
	return s.createFile(filename, data, 0644)
}

// DeleteFile deletes the file with the given filename.
func (s Storage) DeleteFile(filename string) error {
	filePath := filepath.Join(s.BaseDir, filename)
	return os.Remove(filePath)
}

// CreateDir creates a new directory with the given dirname.
func (s Storage) CreateDir(dirname string) error {
	dirPath := filepath.Join(s.BaseDir, dirname)
	return os.MkdirAll(dirPath, 0755)
}

// ReadDir reads the content of the given directory.
func (s Storage) ReadDir(dirname string) ([]os.DirEntry, error) {
	dirPath := filepath.Join(s.BaseDir, dirname)
	return os.ReadDir(dirPath)
}

// DeleteDir deletes the directory with the given dirname.
func (s Storage) DeleteDir(dirname string) error {
	dirPath := filepath.Join(s.BaseDir, dirname)
	return os.RemoveAll(dirPath)
}

// EmptyDir deletes all files and subdirectories within the given directory.
func (s *Storage) EmptyDir(dirname string) error {
	dirPath := filepath.Join(s.BaseDir, dirname)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())
		err := os.RemoveAll(entryPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// Move moves the src directory or file to the dst directory.
func (s *Storage) Move(src string, dst string) error {
	srcPath := filepath.Join(s.BaseDir, src)
	dstPath := filepath.Join(s.BaseDir, dst, filepath.Base(src))

	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(dstPath), 0755)
		if err != nil {
			return err
		}
	}

	return os.Rename(srcPath, dstPath)
}

// Copy a directory or file from one path to another
func (s *Storage) Copy(src string, dst string) error {
	srcPath := filepath.Join(s.BaseDir, src)
	dstPath := filepath.Join(s.BaseDir, dst, filepath.Base(src))

	info, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return s.copyDir(srcPath, dstPath)
	}

	return s.copyFile(srcPath, dstPath)
}

func (s *Storage) copyDir(src string, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = s.copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = s.copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Storage) copyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Sync()
}

// Download writes the file from the specified source path to an http.ResponseWriter.
func (s *Storage) Download(srcPath string, w http.ResponseWriter) error {
	srcPath = filepath.Join(s.BaseDir, srcPath)

	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(w, file); err != nil {
		return err
	}
	return nil
}

// FullPath returns the full path of a file or directory given its relative path.
func (s *Storage) FullPath(relPath string) string {
	return filepath.Join(s.BaseDir, relPath)
}
