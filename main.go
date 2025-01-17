// main.go
package main

import (
	"fmt"
	"go-pdf-analyzer/tools" // Replace with your module name
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Define the current working directory dynamically
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	fmt.Printf("Current working directory: %s\n", currentDir)

	// Define the PDF file and tools path
	pdfPath := filepath.Join(currentDir, "tools", "example.PDF")
	textOutputPath := filepath.Join(currentDir, "output.txt")

	// Verify if the file exists
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		log.Fatalf("PDF file not found: %s", pdfPath)
	}
	fmt.Printf("PDF path: %s\n", pdfPath)

	// Run pdfinfo to get metadata
	fmt.Println("Running pdfinfo...")
	pdfInfoOutput, err := tools.RunPDFInfo(pdfPath)
	if err != nil {
		log.Fatalf("Error running pdfinfo: %v", err)
	}
	fmt.Println("PDF Info Output:")
	fmt.Println(pdfInfoOutput)

	// Run pdftotext to extract text
	fmt.Println("Extracting text with pdftotext...")
	err = tools.RunPDFToText(pdfPath, textOutputPath)
	if err != nil {
		log.Fatalf("Error running pdftotext: %v", err)
	}
	fmt.Println("Text extraction completed.")

	// Read the extracted text
	extractedText, err := os.ReadFile(textOutputPath)
	if err != nil {
		log.Fatalf("Error reading extracted text: %v", err)
	}
	fmt.Println("Extracted Text:")
	fmt.Println(string(extractedText))

	// Get the port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if no environment variable is set
	}

	// Basic HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "PDF Analyzer Service is running!")
	})

	// Start the server
	fmt.Printf("Server is listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
