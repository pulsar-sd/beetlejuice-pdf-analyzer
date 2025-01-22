package main

import (
	"encoding/json"
	"fmt"
	"go-pdf-analyzer/tools" // Replace with your module name
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Get the port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if no environment variable is set
	}

	// Analyze PDF handler
	http.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var requestData struct {
				ID  int    `json:"id"`
				URL string `json:"url"`
			}

			err := json.NewDecoder(r.Body).Decode(&requestData)
			if err != nil {
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				log.Printf("Error decoding request: %v", err)
				return
			}

			log.Printf("Received ID: %d, URL: %s", requestData.ID, requestData.URL)

			// Download the PDF from the provided URL
			pdfPath := filepath.Join(os.TempDir(), fmt.Sprintf("pdf_%d.pdf", requestData.ID))
			err = downloadFile(requestData.URL, pdfPath)
			if err != nil {
				http.Error(w, "Failed to download PDF", http.StatusInternalServerError)
				log.Printf("Error downloading PDF: %v", err)
				return
			}

			// Run pdfinfo to get metadata
			pdfInfoOutput, err := tools.RunPDFInfo(pdfPath)
			if err != nil {
				http.Error(w, "Failed to run pdfinfo", http.StatusInternalServerError)
				log.Printf("Error running pdfinfo: %v", err)
				return
			}
			log.Printf("PDF Info Output: %s", pdfInfoOutput)

			// Extract page count from pdfinfo
			pageCount := extractPageCount(pdfInfoOutput)

			// Run pdftotext to extract text
			textOutputPath := filepath.Join(os.TempDir(), fmt.Sprintf("text_%d.txt", requestData.ID))
			err = tools.RunPDFToText(pdfPath, textOutputPath)
			if err != nil {
				http.Error(w, "Failed to run pdftotext", http.StatusInternalServerError)
				log.Printf("Error running pdftotext: %v", err)
				return
			}

			// Read the extracted text
			extractedText, err := os.ReadFile(textOutputPath)
			if err != nil {
				http.Error(w, "Failed to read extracted text", http.StatusInternalServerError)
				log.Printf("Error reading extracted text: %v", err)
				return
			}

			// Send the response back to the client
			fmt.Fprintf(w, "PDF Analysis Complete: PageCount=%d, TextLength=%d", pageCount, len(extractedText))
		} else {
			fmt.Fprintln(w, "PDF Analyzer Service is running!")
		}
	})

	// Start the server
	fmt.Printf("Server is listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}

// downloadFile downloads a file from a URL and saves it to the specified path.
func downloadFile(url, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// extractPageCount extracts the page count from the pdfinfo output.
func extractPageCount(pdfInfo string) int {
	lines := strings.Split(pdfInfo, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Pages:") {
			var pageCount int
			fmt.Sscanf(line, "Pages: %d", &pageCount)
			return pageCount
		}
	}
	return 0 // Default if no page count is found
}
