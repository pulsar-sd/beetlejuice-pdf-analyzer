// xpdf.go
package tools

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
)

// RunPDFInfo executes the pdfinfo command to extract metadata from a PDF.
func RunPDFInfo(pdfPath string) (string, error) {
	// Determine the absolute path of the `pdfinfo` binary
	pdfinfoPath := filepath.Join("tools", "pdfinfo")

	// Create the command
	cmd := exec.Command(pdfinfoPath, pdfPath)

	// Capture standard error
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Execute the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running pdfinfo: %v\nstderr: %s", err, stderr.String())
	}
	return string(output), nil
}

// RunPDFToText executes the pdftotext command to extract text from a PDF.
func RunPDFToText(pdfPath, outputPath string) error {
	// Determine the absolute path of the `pdftotext` binary
	pdftotextPath := filepath.Join("tools", "pdftotext")

	// Create the command
	cmd := exec.Command(pdftotextPath, pdfPath, outputPath)

	// Capture standard error
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running pdftotext: %v\nstderr: %s", err, stderr.String())
	}
	return nil
}
