package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("path to file needs to be provided")
	}
	// File path to be read
	filePath := os.Args[1]
	isTest := false
	if len(os.Args) == 3 && os.Args[2] == "test" {
		isTest = true
		fmt.Println("Running in test mode")
	}
	// Get the file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatalf("failed to get file info: %v", err)
	}

	fileSize := fileInfo.Size()
	// Decide the number of workers and chunk size
	workers := runtime.NumCPU()
	//setting workers to 4 for integration test
	if isTest {
		workers = 2
	}
	chunkSize := fileSize / int64(workers)

	//timer
	start := time.Now()

	var readerWg sync.WaitGroup
	resultBitSet := NewShardedBitset()
	// Starting gouroutines, assigning start and end of each chunk of the file
	for i := 0; i < workers; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize
		if i == workers-1 {
			end = fileSize // last routine reads till the end of file
		}
		readerWg.Add(1)
		time.Sleep(500 * time.Microsecond)
		go worker(start, end, i, filePath, &readerWg, resultBitSet)
	}
	readerWg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("File processing complete.\nUnique IP count: %d\nTime elapsed: %s", resultBitSet.Count(), elapsed)
}

func worker(start, end int64, id int, filePath string, wg *sync.WaitGroup, bitset *ShardedBitset) {
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Worker %d failed to open file: %v\n", id, err)
	}
	defer file.Close()

	// Seek to the starting position
	_, err = file.Seek(start, 0)
	if err != nil {
		log.Fatalf("Worker %d failed to seek: %v\n", id, err)
		return
	}

	bufferSize := 4 * 1024 * 1024 // 4MB buffer
	reader := bufio.NewReaderSize(file, bufferSize)

	// this skipping to the next new line just in case we are in the middle of the line
	if start > 0 {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Worker %d failed to read: %v\n", id, err)
		}
		start += int64(len(line))
	}

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatalf("Worker %d failed to read: %v\n", id, err)
		}
		//egde case, last element in file
		if err == io.EOF && len(line) == 0 {
			break
		}

		ipIntArr := IpToIntArray(line)
		lastThreeInt, err := LastThreeBytesToInt(ipIntArr[1:])
		if err != nil {
			log.Fatalf("Worker %d error converting bytes to int %s\n", id, err)
		}

		bitset.Set(ipIntArr[0], lastThreeInt)

		//this exit condition is after the read since we want to process additional line
		if start >= end {
			break
		}
		start += int64(len(line))
	}
}
