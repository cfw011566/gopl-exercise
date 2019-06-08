// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 250.

// The du3 command computes the disk usage of the files in a directory.
package main

// The du3 variant traverses all directories in parallel.
// It uses a concurrency-limiting counting semaphore
// to avoid opening too many files at once.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var vFlag = flag.Bool("v", false, "show verbose progress messages")

//!+
func main() {
	// ...determine roots...

	//!-
	flag.Parse()

	// Determine the initial directories.
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	for _, root := range roots {
		duonepath(root)
	}
}

func duonepath(filepath string) {
	//!+
	// Traverse each root of the file tree in parallel.
	fileSizes := make(chan int64)
	var n sync.WaitGroup
	n.Add(1)
	go walkDir(filepath, &n, fileSizes)
	go func() {
		n.Wait()
		close(fileSizes)
	}()
	//!-

	// Print the results periodically.
	var ticker *time.Ticker
	var tick <-chan time.Time
	if *vFlag {
		ticker = time.NewTicker(500 * time.Millisecond)
		tick = ticker.C
	}
	var nfiles, nbytes int64
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(filepath, nfiles, nbytes)
		}
	}

	if ticker != nil {
		ticker.Stop()
	}

	printDiskUsage(filepath, nfiles, nbytes) // final totals
	//!+
	// ...select loop...
}

//!-

func printDiskUsage(filepath string, nfiles, nbytes int64) {
	if nbytes > 1e9 {
		fmt.Printf("%s: %d files  %.2f GB\n", filepath, nfiles, float64(nbytes)/1e9)
	} else if nbytes > 1e6 {
		fmt.Printf("%s: %d files  %.2f MB\n", filepath, nfiles, float64(nbytes)/1e6)
	} else {
		fmt.Printf("%s: %d files  %.2f KB\n", filepath, nfiles, float64(nbytes)/1e3)
	}
}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
//!+walkDir
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

//!-walkDir

//!+sema
// sema is a counting semaphore for limiting concurrency in dirents.
var sema = make(chan struct{}, 20)

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token
	// ...
	//!-sema

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}
