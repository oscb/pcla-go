package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
	<html>
		<head>
			<meta http-equiv="content-type" content="text/html;charset=utf-8">
			<title>Markdown Preview Tool</title>
		</head>
		<body>
	`
	footer = `
		</body>
	</html>
	`
)

func main() {
	filename := flag.String("file", "", "Markdown file to preview")
	skip := flag.Bool("s", false, "Skip Preview")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, os.Stdout, *skip); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string, out io.Writer, skipPreview bool) error {
	// Read all teh data from the input file and check for errors
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)

	// Create temporary file and check for errors
	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	outName := temp.Name()

	fmt.Fprintln(out, outName)

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outName)

	return preview(outName)
}

func parseContent(input []byte) []byte {
	// Parse the markdown file through blackfriday and bluemonday
	// to generate a valid and safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Create a buffer of bytes to write to file
	var buffer bytes.Buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHTML(outFname string, data []byte) error {
	// Write the bytes to the file
	return os.WriteFile(outFname, data, 0644)
}

func preview(fname string) error {
	cName := ""
	cParams := []string{}

	// define executable based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS Not supported")
	}

	// Append filename to parameters slice
	cParams = append(cParams, fname)

	// Locate executable in PATH
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	// Open the file with the default program
	err = exec.Command(cPath, cParams...).Run()

	// Give the browser time to open the file before deletion,
	// this is a very naive fix for race conditions
	time.Sleep(2 * time.Second)
	return err
}
