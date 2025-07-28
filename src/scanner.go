package main

import (
	"fmt"
	"log"
	"sync"
)

func scan(w chan string, urls map[string][]string, endpoints []string) {
	prepareM3u(w)

	var wg sync.WaitGroup

	var i int
	for _, urls := range urls {
		i++
		wg.Add(1)
		go scanUrl(&wg, w, urls, endpoints, fmt.Sprintf("N%d", i))
	}

	wg.Wait()
}

func check(wg *sync.WaitGroup, w chan string, ch, si int, url, name string) {
	defer wg.Done()
	if !checkStream || checkM3U8Status(url) {
		insertM3u(w, fmt.Sprintf("%sE%d", name, si), fmt.Sprintf("%sE%d CH%03d", name, si, ch), url)
		log.Println("+", url)
	} else {
		log.Println("-", url)
	}
}

func scanUrl(wg *sync.WaitGroup, w chan string, urls, endpoints []string, name string) {
	defer wg.Done()
	for ch_i, url := range urls {
		for st_i, ep := range endpoints {
			if len(ep) == 0 {
				continue
			}

			if checkUrl(url+ep, 400, 10) {
				wg.Add(1)
				go check(wg, w, ch_i, st_i, url+ep, name)
			} else {
				log.Println(" ", url+ep)
			}
		}
	}
}
