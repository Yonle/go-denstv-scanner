package main

import (
	"fmt"
	"log"
	"sync"
)

func scan(w chan string, urls map[string][]string) {
	prepareM3u(w)

	var wg sync.WaitGroup

	var i int
	for _, urls := range urls {
		i++
		wg.Add(1)
		go scanUrl(&wg, w, urls, fmt.Sprintf("N%d", i))
	}

	wg.Wait()
}

func scanUrl(wg *sync.WaitGroup, w chan string, urls []string, name string) {
	defer wg.Done()
	for i, url := range urls {
		if checkUrl(url+"/index.m3u8", 400, 10) {
			wg.Add(1)
			go check(wg, w, i, url, name)
		} else {
			log.Println(" ", url)
		}
	}
}
