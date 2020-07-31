package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var gaurd = make(chan struct{}, 30)
var done = make(chan struct{})

func main() {
	start := time.Now()
	fileSizes := make(chan int64)
	var wg sync.WaitGroup

	// span worker
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go walkDirs("/usr", fileSizes, &wg)
	}

	//closer
	go func() {
		wg.Wait()
		close(fileSizes)
	}()

	//interupter
	go func() {
		_, _ = os.Stdin.Read(make([]byte, 1))
		close(done)
	}()

	numFile, totalFileSize := 0, 0
	tick := time.NewTicker(500 * time.Millisecond)

loop:
	for {
		select {
		case <-done:
			// drain channel to make goroutines exist
			for range fileSizes {
			}
			return
		case fs, ok := <-fileSizes:
			if !ok {
				tick.Stop()
				break loop
			}
			numFile++
			totalFileSize += int(fs)
		case <-tick.C:
			fmt.Println(numFile)
			fmt.Println(totalFileSize)
		}
	}

	fmt.Println(numFile)
	fmt.Println(totalFileSize)
	fmt.Println(time.Since(start))
	panic("")
}

func walkDirs(dir string, fileSizes chan int64, wg *sync.WaitGroup) {
	defer wg.Done()
	if cancelled() {
		return
	}
	walkDirsRecursive(dir, fileSizes)
}

func dirents(dir string) []os.FileInfo {
	select {
	case gaurd <- struct{}{}:
	case <-done:
		return nil
	}

	defer func() {
		<-gaurd
	}()
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "reading directory %s failed: %s\n", dir, err)
		return nil
	}
	return entries
}

func walkDirsRecursive(dir string, fileSizes chan<- int64) {
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDirsRecursive(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}
