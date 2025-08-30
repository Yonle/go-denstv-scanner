package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var checkStream bool
var outputFile string
var hostsFile string
var endpointsFile string

var hc = HTTPClient{
	UserAgent: "Mozilla/5.0 (X11; Linux x86_64)",
	Referer:   "https://www.dens.tv",
	Client:    http.Client{Timeout: 5 * time.Second},
}
var b = sync.Pool{
	New: func() any {
		return make([]byte, 64)
	},
}

func main() {
	var helpCmd bool
	flag.BoolVar(&checkStream, "checkstream", false, "check the stream before saving. some channels might be skipped")
	flag.StringVar(&outputFile, "output", "denstv.m3u8", "m3u8 result output filename")
	flag.StringVar(&hostsFile, "hosts", "denstv_hosts.txt", "path to denstv_hosts.txt or similar")
	flag.StringVar(&endpointsFile, "endpoints", "denstv_endpoints.txt", "path to denstv_endpoints.txt or similar")
	flag.BoolVar(&helpCmd, "help", false, "show this.")

	flag.Parse()

	if helpCmd {
		flag.Usage()
		os.Exit(1)
	}

	hoststxt, err := os.ReadFile(hostsFile)
	if err != nil {
		fmt.Printf("I couldn't access %s on the current directory: %v\n", hostsFile, err)
		os.Exit(1)
		return
	}

	endpointstxt, err := os.ReadFile(endpointsFile)
	if err != nil {
		fmt.Printf("I couldn't access %s on the current directory: %v\n", endpointsFile, err)
		os.Exit(1)
		return
	}

	hosts := strings.Split(removeR(hoststxt), "\n")
	endpoints := strings.Split(removeR(endpointstxt), "\n")

	fmt.Println("Checking for host health...")

	activeHosts := checkHosts(hosts)

	types := []string{"h", "s"}

	urls := makeUrls(activeHosts, types)

	w := make(chan string)
	go makeFile(w, outputFile)

	fmt.Println("Now scanning...")
	scan(w, urls, endpoints)

	close(w)
	fmt.Println("Done. Result saved to", outputFile)

	if !checkStream {
		fmt.Println("Please note that some channels might be not able to load")
		fmt.Println("When that happens, All you need to do is just go to the next channel.")
		fmt.Println("\nOr rerun this program with -checkstream flag.")
	}
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
		if len(host) == 0 {
			continue
		}

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
		log.Printf("Failed to fetch %s: %v", u, err)
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode > maxStatusCode {
		return false
	}

	// get a buffer
	buf := b.Get().([]byte)
	defer b.Put(buf)

	n, _ := resp.Body.Read(buf)

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
}

func removeR(s []byte) string {
	return strings.ReplaceAll(string(s), "\r", "")
}
