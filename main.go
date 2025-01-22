package main

import (
	"encoding/json"
	"fmt"
	"go-pdf-analyzer/tools" // Update this import path based on your module name
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	http.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			ID  int    `json:"id"`
			URL string `json:"url"`
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			log.Printf("Error decoding JSON: %v", err)
			return
		}

		// Download and process the PDF
		pdfPath := filepath.Join(os.TempDir(), fmt.Sprintf("pdf_%d.pdf", request.ID))
		err = downloadFile(request.URL, pdfPath)
		if err != nil {
			http.Error(w, "Failed to download PDF", http.StatusInternalServerError)
			log.Printf("Error downloading PDF: %v", err)
			return
		}

		// Get metadata using pdfinfo
		pageCount, err := tools.RunPDFInfo(pdfPath)
		if err != nil {
			http.Error(w, "Failed to analyze PDF", http.StatusInternalServerError)
			log.Printf("Error analyzing PDF: %v", err)
			return
		}

		// Extract text using pdftotext
		textPath := filepath.Join(os.TempDir(), fmt.Sprintf("text_%d.txt", request.ID))
		err = tools.RunPDFToText(pdfPath, textPath)
		if err != nil {
			http.Error(w, "Failed to extract text", http.StatusInternalServerError)
			log.Printf("Error extracting text: %v", err)
			return
		}

		extractedText, err := os.ReadFile(textPath)
		if err != nil {
			http.Error(w, "Failed to read extracted text", http.StatusInternalServerError)
			log.Printf("Error reading text file: %v", err)
			return
		}

		// Respond with JSON
		response := map[string]interface{}{
			"pageCount":     pageCount,
			"extractedText": string(extractedText),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if no environment variable is set
	}
	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
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
