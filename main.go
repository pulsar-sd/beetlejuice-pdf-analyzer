package main

import (
	"encoding/json"
	"fmt"
	"go-pdf-analyzer/tools" // Replace with your module name
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	supabaseURL    = "https://vlubhhhbjtarkedepxgz.supabase.co"
	serviceRoleKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InZsdWJoaGhianRhcmtlZGVweGd6Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3MjM1NDU5MTIsImV4cCI6MjAzOTEyMTkxMn0.BXNsyBM-_O_nBoM4I71aYZlUN5Dy73IaYMLf0TBMIlM"
)

func main() {
	// Get the port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if no environment variable is set
	}

	// Basic HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "PDF Analyzer Service is running!")
	})

	// Analyze PDF handler
	http.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Parse request body
		var reqBody struct {
			ID int `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		// Fetch PDF URL from Supabase
		pdfURL, err := fetchPDFURL(reqBody.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching PDF URL: %v", err), http.StatusInternalServerError)
			return
		}

		// Download PDF locally
		pdfPath, err := downloadPDF(pdfURL, reqBody.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error downloading PDF: %v", err), http.StatusInternalServerError)
			return
		}
		defer os.Remove(pdfPath) // Clean up downloaded file

		// Run pdfinfo to get metadata
		pdfInfoOutput, err := tools.RunPDFInfo(pdfPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error running pdfinfo: %v", err), http.StatusInternalServerError)
			return
		}

		// Extract text using pdftotext
		textOutputPath := filepath.Join(os.TempDir(), fmt.Sprintf("output_%d.txt", reqBody.ID))
		defer os.Remove(textOutputPath) // Clean up text output file
		if err := tools.RunPDFToText(pdfPath, textOutputPath); err != nil {
			http.Error(w, fmt.Sprintf("Error running pdftotext: %v", err), http.StatusInternalServerError)
			return
		}

		// Read extracted text
		extractedText, err := os.ReadFile(textOutputPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading extracted text: %v", err), http.StatusInternalServerError)
			return
		}

		// Return the results
		response := map[string]interface{}{
			"pdfInfo":       pdfInfoOutput,
			"extractedText": string(extractedText),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Start the server
	fmt.Printf("Server is listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}

// fetchPDFURL retrieves the PDF URL for the given ID from Supabase.
func fetchPDFURL(id int) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/v1/pdf_uploaded_files?id=eq.%d", supabaseURL, id), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("apikey", serviceRoleKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", serviceRoleKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var results []struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no PDF URL found for ID %d", id)
	}
	return results[0].URL, nil
}

// downloadPDF downloads the PDF from the given URL to a temporary file.
func downloadPDF(url string, id int) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download PDF: status code %d", resp.StatusCode)
	}

	pdfPath := filepath.Join(os.TempDir(), fmt.Sprintf("pdf_%d.pdf", id))
	file, err := os.Create(pdfPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}
	return pdfPath, nil
}
