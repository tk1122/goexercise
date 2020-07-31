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

func main() {
	start := time.Now()
	fileSizes := make(chan int64)
	var wg sync.WaitGroup

	// span worker
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go walkDirs("/usr", fileSizes, &wg)
	}

	//closer
	go func() {
		wg.Wait()
		close(fileSizes)
	}()

	numFile, totalFileSize := 0, 0
	tick := time.NewTicker(500 * time.Millisecond)

loop:
	for {
		select {
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
}

func walkDirs(dir string, fileSizes chan int64, wg *sync.WaitGroup) {
	defer wg.Done()
	walkDirsRecursive(dir, fileSizes)
}

func dirents(dir string) []os.FileInfo {
	gaurd <- struct{}{}
	entries, err := ioutil.ReadDir(dir)
	<- gaurd
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "reading directory %s failed: %s\n", dir, err)
		return nil
	}
	return entries
}

func walkDirsRecursive(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDirsRecursive(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}
