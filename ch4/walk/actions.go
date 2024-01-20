package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func listFile(path string, out io.Writer) error {
	if _, err := fmt.Fprintln(out, path); err != nil {
		return err
	}
	return nil
}

func filterOut(path, ext string, minSize int64, info fs.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	if ext != "" && filepath.Ext(path) != ext {
		return true
	}

	return false
}

func delFile(path string, delLogger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	delLogger.Println(path)
	return nil
}

func archiveFile(destDir, root, path string) error {
	// First verify that the path exists and it is a dir
	info, err := os.Stat(destDir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}

	// Get the relative directory data, so that we replicate the structure into the backup
	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}

	// Create the destination file
	dest := fmt.Sprintf("%s.gz", filepath.Base(path))
	targetPath := filepath.Join(destDir, relDir, dest)

	// MkdirAll will create directories recursively as needed (noop if already exists)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	// Open the output file in write/create mode
	out, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	// Open the iput file
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	// Gzip writer writes compressed files directly
	// Also supports metadata line the filename for uncompressing
	zw := gzip.NewWriter(out)
	zw.Name = filepath.Base(path)

	// Copy from the file-in buffer directly into the gzip writer one
	if _, err = io.Copy(zw, in); err != nil {
		return err
	}

	// Not defering this one to bubble up the error so the caller handles it
	if err := zw.Close(); err != nil {
		return err
	}

	// Explicitely calling it to bubble up the error, although we also defer the
	// close in case anything else fails
	return out.Close()
}
