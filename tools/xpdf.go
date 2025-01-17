package tools

import (
	"bytes"
	"fmt"
	"os/exec"
)

// RunPDFInfo executes the pdfinfo command to extract metadata from a PDF.
func RunPDFInfo(pdfPath string) (string, error) {
	// Use the absolute path to the `pdfinfo` binary
	pdfinfoPath := "./tools/pdfinfo"

	cmd := exec.Command(pdfinfoPath, pdfPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Command: %s\n", cmd.String())
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Stderr: %s\n", stderr.String())
	}
	return string(output), nil
}

// RunPDFToText executes the pdftotext command to extract text from a PDF.
func RunPDFToText(pdfPath, outputPath string) error {
	// Use the absolute path to the `pdftotext` binary
	pdftotextPath := "./tools/pdftotext"

	cmd := exec.Command(pdftotextPath, pdfPath, outputPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running pdftotext: %v\nstderr: %s", err, stderr.String())
	}
	return nil
}
