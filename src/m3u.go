package main

import (
	"fmt"
)

func prepareM3u(w chan string) {
	w <- "#EXTM3U\n"
}

func insertM3u(w chan string, group, name, url string) {
	w <- fmt.Sprintf("#EXTINF:-1 group-title=\"%s\",%s\n", group, name)
	w <- "#EXTVLCOPT:http-referrer=https://dens.tv\n"
	w <- "#EXTVLCOPT:http-user-agent=Mozilla/5.0 (X11; Linux x86_64)\n"
	w <- "#KODIPROP:inputstream=inputstream.adaptive\n"
	w <- "#KODIPROP:inputstreamaddon=inputstream.adaptive\n"
	w <- "#KODIPROP:inputstream.adaptive.manifest_type=hls\n"
	w <- url + "\n"
}
