package main

import (
	"bytes"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

type File struct {
	Path         string
	Name         string
	CheckSum     []byte
	IsDuplicated bool
}

func main() {
	srcDir := flag.String("s", "", "source directory to compare")
	dstDir := flag.String("d", "", "destination directory to compare")
	outDir := flag.String("o", "", "output directory for duplicate files from dst that exist in src")
	flag.Parse()
	if *srcDir == "" || *dstDir == "" || *outDir == "" {
		flag.Usage()
		return
	}
	var dstFiles []File
	sha1V := sha1.New()
	err := filepath.Walk(*dstDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		sha1V.Reset()
		f, _ := os.Open(path)
		_, err = io.Copy(sha1V, f)
		f.Close()
		if err != nil {
			return err
		}
		file := File{
			Path:         path,
			Name:         info.Name(),
			CheckSum:     sha1V.Sum(nil),
			IsDuplicated: false,
		}
		dstFiles = append(dstFiles, file)
		return nil
	})
	check(err)
	err = filepath.Walk(*srcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		sha1V.Reset()
		f, _ := os.Open(path)
		_, err = io.Copy(sha1V, f)
		f.Close()
		if err != nil {
			return err
		}
		checkSum := sha1V.Sum(nil)
		for index, dstFile := range dstFiles {
			if bytes.Equal(checkSum, dstFile.CheckSum) {
				dstFiles[index].IsDuplicated = true
			}
		}
		return nil
	})
	check(err)
	err = os.MkdirAll(*outDir, os.ModePerm)
	check(err)
	count := 0
	for _, dstFile := range dstFiles {
		if dstFile.IsDuplicated {
			count++
			fmt.Println("Duplicate file: " + dstFile.Path)
			err = os.Rename(dstFile.Path, path.Join(*outDir, dstFile.Name))
			check(err)
		}
	}
	fmt.Printf("%d duplicate files moved.\n", count)
}
func check(err error) {
	if err != nil {
		panic(err)
	}
}
