package zip

import (
	"archive/zip"
	"github.com/vinkdong/gox/log"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/*
compress file or folder to zip compress file
*/
func Compress(filename string, target string, addRootDir bool) error {
	targetFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer func() {
		if err := targetFile.Close(); err != nil {
			log.Error(err)
		}
	}()
	archive := zip.NewWriter(targetFile)
	defer func() {
		if err := archive.Close(); err != nil {
			log.Error(err)
		}
	}()
	stat, err := os.Stat(filename)
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	walkFunc := func(path string, info os.FileInfo, err error) error {
		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		if !addRootDir && stat.IsDir() && abs == absPath {
			return nil
		}
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(abs, absPath+"/")
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		w, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(w, f)
		return err
	}
	return filepath.Walk(filename, walkFunc)
}
