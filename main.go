package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: drag_png_to_exe <input_png_file>")
		os.Exit(1)
	}
	
	inputFile := os.Args[1]
	if !strings.HasSuffix(strings.ToLower(inputFile), ".png") {
		fmt.Println("Error: Input file must be a PNG.")
		os.Exit(1)
	}
	
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found.\n", inputFile)
		os.Exit(1)
	}
	
	outputFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ".ico"
	batFilePath := filepath.Join(os.TempDir(), "conversion.bat")
	err := createBatchFile(batFilePath, inputFile, outputFile)
	if err != nil {
		fmt.Printf("Error creating batch file: %v\n", err)
		os.Exit(1)
	}
	
	cmd := exec.Command("cmd", "/C", batFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error executing batch file: %v\n", err)
		os.Exit(1)
	}
	err = os.Remove(batFilePath)
	if err != nil {
		fmt.Printf("Error deleting batch file: %v\n", err)
	}
	fmt.Printf("PNG file converted to ICO: %s\n", outputFile)
}

func createBatchFile(batFilePath, inputFile, outputFile string) error {
	batchContent := fmt.Sprintf("ffmpeg -i \"%s\" \"%s\"", inputFile, outputFile)
	return ioutil.WriteFile(batFilePath, []byte(batchContent), 0644)
}
