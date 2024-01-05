package internal

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func Zip(dstFile *os.File, srcPath string) error {
	z := zip.NewWriter(dstFile)
	defer z.Close()

	err := filepath.Walk(srcPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name, err = filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		dst, err := z.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		src, err := os.Open(path)

		if err != nil {
			return err
		}
		defer src.Close()

		_, err = io.Copy(dst, src)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}