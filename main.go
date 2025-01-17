package main

import (
	"fmt"
	"go-pdf-analyzer/tools"
)

func main() {
	pdfPath := "./example.pdf"
	outputPath := "./output.txt"

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
