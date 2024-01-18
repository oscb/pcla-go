package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// Defining a boolean flag 0l to count lines instead of words
	lines := flag.Bool("l", false, "Count lines")
	flag.Parse()

	// Calling the count function to count the number of words
	// received from the Standard Input and printing it out
	fmt.Println(count(os.Stdin, *lines))
}

func count(r io.Reader, countLines bool) int {
	// A Scanner is used to read text from a Reader (such as files)
	scanner := bufio.NewScanner(r)

	// Define the scanner split type to words or lines (default is lines)
	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	wc := 0

	// Increment the counter for every word scanned
	for scanner.Scan() {
		wc++
	}

	return wc
}
