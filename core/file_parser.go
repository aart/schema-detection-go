package core

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)



// ReadLines reads multiple files concurrently and sends the lines to a channel.
func ReadLines(paths []string, numWorkers int) <-chan string {
	lines := make(chan string)
	var wg sync.WaitGroup

	for _, path := range paths {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			file, err := os.Open(filePath)
			if err != nil {
				// In a real application, you might want to log this error
				fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filePath, err)
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines <- scanner.Text()
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filePath, err)
			}
		}(path)
	}

	go func() {
		wg.Wait()
		close(lines)
	}()

	return lines
}
