package main

import (
	"fmt"
	"go-pdf-analyzer/tools"
	"os"
	"path/filepath"
)

func main() {

	outputPath := "./output.txt"

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		return
	}
	fmt.Printf("Current working directory: %s\n", wd)

	// Construct absolute paths
	pdfPath := filepath.Join(wd, "example.pdf")
	fmt.Printf("PDF path: %s\n", pdfPath)

	// Run pdfinfo
	info, err := tools.RunPDFInfo(pdfPath)
	if err != nil {
		fmt.Printf("Error running pdfinfo: %v\n", err)
		return
	}
	fmt.Println("PDF Info Output:")
	fmt.Println(info)

	// Run pdftotext
	err = tools.RunPDFToText(pdfPath, outputPath)
	if err != nil {
		fmt.Printf("Error running pdftotext: %v\n", err)
		return
	}
	fmt.Println("Text extraction completed.")
}
