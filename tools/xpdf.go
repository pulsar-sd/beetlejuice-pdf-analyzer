package tools

import (
	"os/exec"
)

func RunPDFInfo(pdfPath string) (string, error) {
	cmd := exec.Command("pdfinfo", pdfPath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func RunPDFToText(pdfPath, outputPath string) error {
	cmd := exec.Command("pdftotext", pdfPath, outputPath)
	return cmd.Run()
}
