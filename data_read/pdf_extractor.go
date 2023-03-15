package data_read

import (
	"code.sajari.com/docconv"
	"io"
	"log"
	"net/http"
	"os"
)

type PdfExtractor interface {
	Extract(url, name string) (string, error)
}

type PdfExtractorImpl struct {
}

func (p PdfExtractorImpl) Extract(url, name string) (string, error) {
	fileName := "./" + name + ".pdf"
	err := downloadFile(url, fileName)
	if err != nil {
		return "", err
	}
	output, err := readPdf(fileName)
	if err != nil {
		return "", err
	}
	// Now delete the file
	err = os.Remove(fileName)
	if err != nil {
		return output, err
	}
	return output, nil
}

func NewPdfExtractor() PdfExtractor {
	return &PdfExtractorImpl{}
}

// Download a file from a URL and save it to a local file
func downloadFile(url string, filepath string) error {
	// Get the HTTP response from the URL
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Create the local file to save the downloaded content
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the contents of the HTTP response to the local file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func readPdf(path string) (string, error) {
	res, err := docconv.ConvertPath(path)
	if err != nil {
		log.Fatal(err)
	}
	return res.Body, nil
}
