package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var thread = 4

var hc = http.Client{Timeout: 5 * time.Second}
var q = make(chan int, 4)
var b = sync.Pool{
	New: func() any {
		return make([]byte, 64)
	},
}

func main() {
	hoststxt, err := os.ReadFile("denstv_hosts.txt")
	if err != nil {
		fmt.Println("I couldn't access denstv_hosts.txt on the current directory.")
		os.Exit(1)
		return
	}

	hosts := strings.Split(string(hoststxt), "\n")

	fmt.Println("Checking for host health...")

	activeHosts := checkHosts(hosts)

	types := []string{"h", "s"}

	urls := makeUrls(activeHosts, types)

	w := make(chan string)
	go makeFile(w, "denstv.m3u8")

	fmt.Println("Now scanning...")
	scan(w, urls)

	close(w)
	fmt.Println("Done. Result saved to denstv.m3u8")
	fmt.Println("Please note that some channels might be not able to load")
	fmt.Println("When that happens, All you need to do is just go to the next channel.")
}

func makeUrls(hosts, types []string) (urls map[string][]string) {
	urls = make(map[string][]string)

	for _, host := range hosts {
		for _, typ := range types {
			for i := range make([]struct{}, 500) {
				if _, ok := urls[host]; !ok {
					urls[host] = []string{}
				}

				urls[host] = append(urls[host], fmt.Sprintf("%s/%s/%s%02d", host, typ, typ, i))
			}
		}
	}

	return
}

func checkHosts(hosts []string) (activeHosts []string) {
	for _, host := range hosts {
		if checkUrl(host, 500, 0) {
			activeHosts = append(activeHosts, host)
			fmt.Println("--  OK:", host)
		} else {
			fmt.Println("-- NOK:", host)
		}
	}

	return
}

func checkUrl(u string, maxStatusCode, maxLen int) bool {
	resp, err := hc.Get(u)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode > maxStatusCode {
		return false
	}

	// get a buffer
	buf := b.Get().([]byte)
	defer b.Put(buf)

	n, err := resp.Body.Read(buf)

	return n >= maxLen
}

func makeFile(w chan string, fn string) {
	f, err := os.Create(fn)
	if err != nil {
		panic(err)
	}

	for d := range w {
		f.WriteString(d)
	}

	f.Close()

	return
}

