package zipfile

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func unzipFile(file *zip.File, dstDir string) error {
	filePath := filepath.Join(dstDir, file.Name)
	if file.FileInfo().IsDir() {
		return os.MkdirAll(filePath, os.ModePerm)
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, rc)
	return err
}
func Unzip(zipPath, dstDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if err := unzipFile(file, dstDir); err != nil {
			return err
		}
	}
	return nil
}
