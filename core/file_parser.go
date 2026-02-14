package core

import (
	"bufio"
	"os"
)

// LinesResult holds a channel of strings and the total number of lines.
type LinesResult struct {
	Lines      <-chan string
	TotalLines int
}

// ReadLines reads multiple files concurrently and sends the lines to a channel.
func ReadLines(paths []string, numWorkers int) <-chan string {
	lines := make(chan string)
	wp := NewWorkerPool(numWorkers)
	wp.Start()

	go func() {
		for _, path := range paths {
			path := path
			wp.Submit(Job{
				Fn: func() interface{} {
					file, err := os.Open(path)
					if err != nil {
						return nil
					}
					defer file.Close()

					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						lines <- scanner.Text()
					}
					return nil
				},
			})
		}
		wp.Stop()
	}()

	go func() {
		// This is not ideal, but for the sake of this test, we assume the results channel is drained
		// somewhere else. In a real application, we would need to drain the results channel.
		for range wp.Results() {
		}
		close(lines)
	}()

	return lines
}
