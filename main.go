package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const version = "v.1.0.1"

func main() {

	fmt.Printf("m3u8-dl %s\n", version)

	var (
		inputUrl   string
		outputFile string
	)
	flag.StringVar(&inputUrl, "url", "", "url path to m3u8")
	flag.StringVar(&outputFile, "file", "", "path to output file")

	flag.Parse()
	if inputUrl == "" {
		log.Fatal("url is required")
	}

	if outputFile == "" {
		log.Fatal("file is required")
	}

	// received file from server
	resp, err := http.Get(inputUrl)
	if err != nil {
		log.Fatal("Download error: ", err)
	}
	defer resp.Body.Close()

	// create output file
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatal("Download error: ", err)
	}
	defer f.Close()

	// read server response line by line
	scanner := bufio.NewScanner(resp.Body)
	i := 0
	for scanner.Scan() {
		l := scanner.Text()

		// if line contains url address
		if strings.HasPrefix(l, "http") {
			// download file part
			part, err := downloadFilePart(l)
			if err != nil {
				log.Fatal("Download part error: ", err)
			}

			// write part to output file
			if _, err = f.Write(part); err != nil {
				log.Fatal("Write part to output file: ", err)
			}
			log.Printf("Download part %d\n", i)
			i++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// downloadFilePart download file part from server
func downloadFilePart(url string) ([]byte, error) {
	result := make([]byte, 0)
	var err error
	var resp *http.Response

	for i := 1; i < 5; i++ {
		if i > 1 {
			fmt.Println("===================RETRY===================")
		}
		resp, err = http.Get(url)
		if err == nil {
			result, err = io.ReadAll(resp.Body)
			if err == nil {
				return result, nil
			}
		}
	}

	return result, err
}
