package main

import (
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"image"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	os.Exit(mainRun())
}

func mainRun() int {
	flag.Parse()

	inputPath := flag.Arg(0)
	if inputPath == "" {
		inputPath = "."
	}

	stat, err := os.Stat(inputPath)
	if err != nil {
		if os.IsNotExist(err) {
			_, _ = fmt.Fprintf(os.Stderr, "Input path %s is not exists.\n", inputPath)
			return 1
		}
		_, _ = fmt.Fprintf(os.Stderr, "Could not check stat of the input path: %s\n", inputPath)
		return 1
	}

	var filePathList []string

	if stat.IsDir() {
		fileEntries, err := os.ReadDir(inputPath)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Could not scan files in the directory: %v\n", err)
			return 1
		}

		filePathList = make([]string, 0, len(fileEntries))
		for _, entry := range fileEntries {
			if !entry.IsDir() {
				filePathList = append(filePathList, filepath.Join(inputPath, entry.Name()))
			}
		}
	} else {
		filePathList = []string{inputPath}
	}

	sort.Strings(filePathList)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"name", "width", "height"})

	for _, filePath := range filePathList {
		switch strings.ToLower(filepath.Ext(filePath)) {
		case ".png", ".jpg", ".jpeg", ".jpe", ".jfi", ".jfif", ".jif", ".gif":
			w, h, err := scanImageFile(filePath)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to decode image information: %v\n", err)
				continue
			}

			table.Append([]string{filepath.Base(filePath), strconv.Itoa(w), strconv.Itoa(h)})
		}
	}

	table.Render()

	return 0
}

func scanImageFile(filePath string) (w, h int, err error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		}
	}()

	img, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}
