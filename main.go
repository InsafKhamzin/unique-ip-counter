package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/bits-and-blooms/bitset"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("path to file needs to be provided")
	}

	// File path to be read
	filePath := os.Args[1]
	// Get the file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatalf("failed to get file info: %v", err)
	}

	fileSize := fileInfo.Size()

	// Decide the number of workers and chunk size
	workers := 10 // TODO should be based on cpu and disk type
	chunkSize := fileSize / int64(workers)

	var readerWg sync.WaitGroup
	var workerCount = make(chan *bitset.BitSet)

	start := time.Now()

	//global bitset
	size := uint(1) << 32
	resultBitSet := bitset.New(size)
	go func() {
		// TODO we can have just one bitset and sharing it accross routines applying locks
		for val := range workerCount {
			resultBitSet = resultBitSet.Union(val)
		}
	}()

	// Start reader gouroutines
	for i := 0; i < workers; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize
		if i == workers-1 {
			end = fileSize // last routine reads till the end
		}
		readerWg.Add(1)
		go worker(start, end, i, filePath, &readerWg, workerCount)
	}
	readerWg.Wait()
	close(workerCount)

	elapsed := time.Since(start)
	fmt.Println("File processing complete. Total count", resultBitSet.Count(), "Time elapsed", elapsed)
}

func worker(start, end int64, id int, filePath string, wg *sync.WaitGroup, workerBitset chan<- *bitset.BitSet) {
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

	size := uint(1) << 32
	bitset := bitset.New(size)

	for {
		line, err := reader.ReadBytes('\n')
		intIP := ipBytesToInt(line)
		bitset.Set(uint(intIP))
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Worker %d failed to read: %v\n", id, err)
		}

		start += int64(len(line))
		//this exit condition is after the read since we want to process additional line
		if start >= end {
			break
		}
	}
	//communicating local bitset
	workerBitset <- bitset
}

func ipBytesToInt(ip []byte) uint32 {
	var result uint32
	pointer := 0
	octet := 3
	for i := 0; i < len(ip); i++ {
		if ip[i] == '.' || ip[i] == '\n' {
			octetSegment := ip[pointer:i]
			val := 0
			for i := 0; i < len(octetSegment); i++ {
				// Subtracting the ASCII value, getting digit
				digit := int(octetSegment[i] - '0')
				val = val*10 + digit
			}
			result |= uint32(val) << (uint32(octet) * 8)
			pointer = i + 1
			octet--
		}
	}
	return result
}
