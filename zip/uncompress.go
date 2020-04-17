package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Uncompress(src, dest string, perm os.FileMode) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	if err := os.MkdirAll(dest, perm); err != nil {
		return err
	}

	walk := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			return os.MkdirAll(path, perm)
		}

		if _, err := os.Stat(filepath.Dir(path)); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(filepath.Dir(path), perm); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		nf, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
		if err != nil {
			return err
		}
		defer nf.Close()

		_, err = io.Copy(nf, rc)
		return err
	}

	for _, f := range r.File {
		if err := walk(f); err != nil {
			return err
		}
	}

	return nil
}
