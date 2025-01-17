// xpdf.go
package tools

import (
	"bytes"
	"fmt"
	"os/exec"
)

// RunPDFInfo executes the pdfinfo command to extract metadata from a PDF.
func RunPDFInfo(pdfPath string) (string, error) {
	cmd := exec.Command("./tools/pdfinfo", pdfPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running pdfinfo: %v\nstderr: %s", err, stderr.String())
	}
	return string(output), nil
}

// RunPDFToText executes the pdftotext command to extract text from a PDF.
func RunPDFToText(pdfPath, outputPath string) error {
	cmd := exec.Command("./tools/pdftotext", pdfPath, outputPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running pdftotext: %v\nstderr: %s", err, stderr.String())
	}
	return nil
}
