package main

import (
	"fmt"
	"go-pdf-analyzer/tools"
	"os"
	"path/filepath"
)

func main() {

	outputPath := "./output.txt"

	wd, _ := os.Getwd()                         // Get the working directory
	pdfPath := filepath.Join(wd, "example.pdf") // Construct the absolute path
	fmt.Printf("Current working directory: %s\n", wd)
	fmt.Printf("PDF path: %s\n", pdfPath)

	// Simulate the tool execution
	output, err := tools.RunPDFInfo(pdfPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("PDF Info Output:")
	fmt.Println(output)

	// Run pdftotext
	err = tools.RunPDFToText(pdfPath, outputPath)
	if err != nil {
		fmt.Printf("Error running pdftotext: %v\n", err)
		return
	}
	fmt.Println("Text extraction completed.")
}
