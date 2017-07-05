package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// `getWriter` takes an output path and returns
// pointer to a bufio writer.
func getWriter(path string) *bufio.Writer {
	if path == "" {
		fmt.Println("Empty output. Using os.Stdout")
		return bufio.NewWriter(os.Stdout)
	}

	// Check if the output file exists.
	if _, err := os.Stat(path); err == nil {
		fmt.Println("The file already exists")
		os.Exit(1)
	}

	// Create file for writing.
	outputFile, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Valid output file. Using file %v\n", outputFile.Name())
	return bufio.NewWriter(outputFile)
}

// `writeContents` is a niladic function which
// writes some content to a bufio.Writer object.
func writeContents(content []byte, w *bufio.Writer) {
	defer w.Flush()
	w.Write(content)
	w.Write([]byte("\n"))
}

func main() {
	outputPath := flag.String(
		"o",
		"",
		"Output file (leave blank to output to stdout)",
	)
	flag.Parse()
	// fmt.Printf("Output file is %v\n", *outputPath)

	// Check that positional arg `url` is given.
	if flag.NArg() != 1 {
		fmt.Println("Url must be specified")
		os.Exit(1)
	}

	// Validate the positional arg `url`.
	dataURL, err := url.ParseRequestURI(flag.Args()[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// fmt.Printf("dataURL: %v %T\n", dataURL, dataURL)

	// Get the writer object to write content to.
	writer := getWriter(*outputPath)

	// Fetch page content.
	response, err := http.Get(dataURL.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read the contents of the response.
	content, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Write content to file.
	writeContents(content, writer)
}
