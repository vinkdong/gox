package gzip

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// GetFileFromGzip 从给定的 .tgz 文件中提取指定文件的内容。
// zipfile 是要读取的 .tgz 文件路径。
// filename 是需要从归档中提取的文件名。
func GetFileFromGzip(zipfile string, filename string) ([]byte, error) {
	file, err := os.Open(zipfile)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", zipfile, err)
	}
	defer file.Close()

	// 创建 gzip reader
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("error creating gzip reader for %s: %w", zipfile, err)
	}
	defer gzReader.Close()

	// 创建 tar reader
	tarReader := tar.NewReader(gzReader)

	// 遍历 tar 归档中的所有文件
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // 已经读到文件末尾
		}
		if err != nil {
			return nil, fmt.Errorf("error reading tar entry: %w", err)
		}

		// 检查是否为想要提取的文件
		if filepath.Base(header.Name) == filename {
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, tarReader); err != nil {
				return nil, fmt.Errorf("error extracting file %s: %w", filename, err)
			}

			return buf.Bytes(), nil
		}
	}

	return nil, fmt.Errorf("file %s not found in archive %s", filename, zipfile)
}
