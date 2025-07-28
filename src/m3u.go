package main

import (
	"fmt"
)

func prepareM3u(w chan string) {
	w <- "#EXTM3U\n"
}

func insertM3u(w chan string, group, name, url string) {
	m3uinfo := fmt.Sprintf("#EXTINF:-1 group-title=\"%s\",%s\n", group, name)
	m3uinfo += fmt.Sprintf("#EXTVLCOPT:http-referer=%s\n", hc.Referer)
	m3uinfo += fmt.Sprintf("#EXTVLCOPT:http-user-agent=%s\n", hc.UserAgent)
	m3uinfo += "#KODIPROP:inputstream=inputstream.adaptive\n"
	m3uinfo += "#KODIPROP:inputstreamaddon=inputstream.adaptive\n"
	m3uinfo += "#KODIPROP:inputstream.adaptive.manifest_type=hls\n"
	m3uinfo += url + "\n"
	w <- m3uinfo
}
