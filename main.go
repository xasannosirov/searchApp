package main

import (
	"code.sajari.com/docconv"
	"fmt"
	"github.com/otiai10/gosseract/v2"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func searchInPDF(filePath, searchTerm string) (bool, error) {
	res, err := docconv.ConvertPath(filePath)
	if err != nil {
		return false, err
	}

	return strings.Contains(strings.ToLower(res.Body), searchTerm), nil
}

func searchInDOCX(filepath, searchTerm string) (bool, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return false, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			_ = err
		}
	}(f)

	var r io.Reader
	r = f

	tmpl, _, err := docconv.ConvertDocx(r)
	if err != nil {
		return false, err
	}

	return strings.Contains(strings.ToLower(tmpl), searchTerm), nil
}

func searchInImage(filepath, searchTerm string) (bool, error) {
	client := gosseract.NewClient()
	defer func(client *gosseract.Client) {
		err := client.Close()
		if err != nil {
			_ = err
		}
	}(client)

	err := client.SetImage(filepath)
	if err != nil {
		return false, err
	}

	text, err := client.Text()
	if err != nil {
		return false, err
	}

	return strings.Contains(strings.ToLower(text), searchTerm), nil
}

func main() {
	searchWord := "lorem"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	folderPath := filepath.Join(homeDir, "projects/go/src/searchApp/files")

	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filePath := folderPath + "/" + file.Name()
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Name()), "."))

		switch ext {
		case "pdf":
			status, err := searchInPDF(filePath, searchWord)
			if err != nil {
				log.Fatal(err)
			}
			if status {
				println(filePath)
			}
		case "docx":
			status, err := searchInDOCX(filePath, searchWord)
			if err != nil {
				log.Fatal(err)
			}
			if status {
				println(filePath)
			}
		case "jpeg", "jpg", "png":
			status, err := searchInImage(filePath, searchWord)
			if err != nil {
				log.Fatal(err)
			}
			if status {
				println(filePath)
			}
		default:
			fmt.Printf("Unsupported file type: %s\n", file.Name())
		}
	}
}
