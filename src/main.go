package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	// Download file from url
	resp, err := http.Get("https://www.github.com")
	if err != nil {
		panic(err)
	}

	// Save response to file
	file, err := os.Create("github.html")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}

	defer file.Close()
}
