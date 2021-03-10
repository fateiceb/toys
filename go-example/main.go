package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	findflag string
)

func init() {
	flag.StringVar(&findflag, "find", "1111", "find file")
}
func main() {
	flag.Parse()
	files := directoryList()
	findfileBygotask(files)
}

func findfileBygo(files []string) {
	ch := make(chan struct{}, 4)
	for i := 0; i < len(files)-1; i++ {
		go func(path string) {
			ch <- struct{}{}
			// fmt.Println(path)
			filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				// fmt.Println(info.Name())
				return nil
			})
		}(files[i])
		<-ch
	}
}
func findfileBynomal(files []string) {
	for i := 0; i < len(files)-1; i++ {
		// fmt.Println(files[i])
		filepath.Walk(files[i], func(path string, info os.FileInfo, err error) error {
			// fmt.Println(info.Name())
			return nil
		})

	}
}
func findfileBygotask(files []string) {
	group := sync.WaitGroup{}
	ch := make(chan string, 4)
	for i := 0; i < len(files)-1; i++ {
		group.Add(1)
		ch <- files[i]
	}

	for i := 0; i < 6; i++ {
		go func(chan string) {
			path := <-ch
			// fmt.Println(path)
			filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				// fmt.Println(info.Name())
				return nil
			})
			group.Done()
		}(ch)
	}
	group.Wait()
}
func directoryList() []string {
	var files []string
	homepath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(errors.New("get user home fail"))
	}
	if os.Chdir(homepath) != nil {
		log.Fatal(errors.New("homepath change fail"))
	}
	err = filepath.Walk(homepath+"/toys", func(path string, info os.FileInfo, err error) error {
		if info.Name() == ".git" {
			return filepath.SkipDir
		}
		if info.IsDir() && info.Name() != "toys" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
