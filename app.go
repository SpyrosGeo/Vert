package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/chai2010/webp"
)

func downloadFile(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	file, err := os.CreateTemp("", "image-*.webp")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func convertWebpToJpg(webpPath, jpgPath string) error {
	reader, err := os.Open(webpPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	img, err := webp.Decode(reader)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(jpgPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = webp.Encode(outputFile, img, &webp.Options{Quality: 100})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var url, outputDir string
	var deleteAll bool
	outputDir = getOSPath()
	flag.StringVar(&url, "url", "", "URL of the webp image")
	flag.BoolVar(&deleteAll, "delete", false, "Delete all files in the output directory before saving the converted image")
	flag.Parse()

	if url == "" {
		fmt.Println("Please provide the URL of the webp image using the -url flag")
		return
	}

	if deleteAll {
		err := deleteFilesInDirectory(outputDir)
		if err != nil {
			fmt.Println("Error deleting files in the output directory:", err)
			return
		}
	}

	err := os.MkdirAll(outputDir, 0755) // Create the directory if it doesn't exist
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	webpPath, err := downloadFile(url)
	if err != nil {
		fmt.Println("Error downloading webp file:", err)
		return
	}
	defer os.Remove(webpPath)

	filename := filepath.Base(url)
	jpgPath := filepath.Join(outputDir, filename[:len(filename)-len(filepath.Ext(filename))]+".jpg")

	err = convertWebpToJpg(webpPath, jpgPath)
	if err != nil {
		fmt.Println("Error converting webp to jpg:", err)
		return
	}

	fmt.Println("Image successfully converted and saved to", jpgPath)
}

func deleteFilesInDirectory(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := os.RemoveAll(filepath.Join(dir, file.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}
func getOSPath() string {

	currentOs := runtime.GOOS
	if currentOs == "windows" {
		return "C:\\Users\\ThatGuy\\Desktop\\sterea"
	}
	return "/home/thatguy/Downloads/sterea/"
}
