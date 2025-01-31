package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

type SimplifiedData struct {
    ID       int       `json:"id"`
    Name     string    `json:"name"`
    Value    float64   `json:"value"`
    Tags     []string  `json:"tags"`
    Metadata Metadata `json:"metadata"`
}

type Metadata struct {
    Created  string `json:"created"`
    Priority int    `json:"priority"`
    Active   bool   `json:"active"`
}

// ProcessJSONStream processes a large JSON file with minimal memory overhead
// Implements sequential processing
func ProcessJSONStream(filePath string, processFn func(SimplifiedData) error) (int, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a decoder
	decoder := json.NewDecoder(file)

	// Ensure we're starting with an array
	token, err := decoder.Token()
	if err != nil {
		return 0, fmt.Errorf("failed to read opening token: %w", err)
	}
	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		return 0, fmt.Errorf("expected start of JSON array, got %v", token)
	}

	// Counter for processed items
	count := 0

	// Process each item
	for decoder.More() {
		var item SimplifiedData
		if err := decoder.Decode(&item); err != nil {
			return count, fmt.Errorf("error decoding item at position %d: %w", count, err)
		}

		// Optional processing of each item
		if err := processFn(item); err != nil {
			return count, err
		}

		count++
	}

	return count, nil
}

// ParallelProcessJSONStream processes a large JSON file using parallel processing
// Implements parallel processing using goroutines
func ParallelProcessJSONStream(filePath string, workerCount int, processFn func(SimplifiedData) error) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// Ensure we're starting with an array
	token, err := decoder.Token()
	if err != nil {
		return 0, fmt.Errorf("failed to read opening token: %w", err)
	}
	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		return 0, fmt.Errorf("expected start of JSON array, got %v", token)
	}

	// Buffered channel for items
	itemChan := make(chan SimplifiedData, 100)
	var wg sync.WaitGroup
	var mu sync.Mutex
	count := 0
	errChan := make(chan error, workerCount)

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range itemChan {
				mu.Lock()
				count++
				mu.Unlock()

				if err := processFn(item); err != nil {
					errChan <- err
					return
				}
			}
		}()
	}

	// Start a goroutine to read JSON items
	go func() {
		defer close(itemChan)
		defer close(errChan)

		for decoder.More() {
			var item SimplifiedData
			if err := decoder.Decode(&item); err != nil {
				errChan <- err
				return
			}
			itemChan <- item
		}
	}()

	// Wait for all workers to finish
	wg.Wait()

	// Check for any errors
	select {
	case err := <-errChan:
		return count, err
	default:
		return count, nil
	}
}

func main() {
	// Command-line flags
	parallelFlag := flag.Bool("parallel", false, "Enable parallel processing")
	flag.Parse()

	// Optional: CPU profiling
	cpuProfile, err := os.Create("cpu_profile.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()

	// Memory profiling
	memProfile, err := os.Create("mem_profile.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer memProfile.Close()

	// Start timing
	start := time.Now()

	var count int
	var processingErr error

	// Choose processing method based on flag
	if *parallelFlag {
		count, processingErr = ParallelProcessJSONStream("large-file.json", runtime.NumCPU(), func(item SimplifiedData) error {
			// Optional per-item processing
			return nil
		})
	} else {
		count, processingErr = ProcessJSONStream("large-file.json", func(item SimplifiedData) error {
			// Optional per-item processing
			return nil
		})
	}

	if processingErr != nil {
		log.Fatalf("Error processing file: %v", processingErr)
	}

	// Log performance metrics
	duration := time.Since(start)
	log.Printf("Processed %d items in %v", count, duration)

	// Write memory profile
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(memProfile); err != nil {
		log.Fatal(err)
	}
}
