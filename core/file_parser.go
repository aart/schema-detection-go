package core

import (
	"bufio"
	"os"
	"sync"
)

// ReadLines reads a file and returns its content as a slice of strings, with each string representing a line.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// ReadLinesConcurrently reads multiple files concurrently and sends the lines to a channel.
func ReadLinesConcurrently(paths []string) <-chan string {
	lines := make(chan string)
	var wg sync.WaitGroup

	for _, path := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			file, err := os.Open(path)
			if err != nil {
				// In a real application, you'd want to handle this error better.
				// For now, we'll just print it and move on.
				// Note: This is not ideal for production code.
				// A better approach would be to send errors on a separate channel.
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines <- scanner.Text()
			}
		}(path)
	}

	go func() {
		wg.Wait()
		close(lines)
	}()

	return lines
}
