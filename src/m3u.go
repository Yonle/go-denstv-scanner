package main

import (
	"fmt"
	"net/url"
)

var m3ugprops string

func init() {
	kodiheaders := fmt.Sprintf("Referer=%s&User-Agent=%s",
		url.QueryEscape(hc.Referer),
		url.QueryEscape(hc.UserAgent),
	)

	// Fun Fact:
	// "Referer" is actually a typo from the HTTP RFC 7231 Section 5.5.2.
	// It was supposedly "Referrer", but for historical reasons, It is and remains "Referer".
	//
	// Your brain is now damaged. English to Engrisu
	m3ugprops = fmt.Sprintf("#EXTVLCOPT:http-referer=%s\n", hc.Referer)
	m3ugprops += fmt.Sprintf("#EXTVLCOPT:http-user-agent=%s\n", hc.UserAgent)
	m3ugprops += "#KODIPROP:inputstream=inputstream.adaptive\n"
	m3ugprops += "#KODIPROP:inputstreamaddon=inputstream.adaptive\n"
	m3ugprops += "#KODIPROP:inputstream.adaptive.manifest_type=hls\n"
	m3ugprops += fmt.Sprintf("#KODIPROP:inputstream.adaptive.manifest_headers=%s\n", kodiheaders)
	m3ugprops += fmt.Sprintf("#KODIPROP:inputstream.adaptive.stream_headers=%s\n", kodiheaders)
}

func prepareM3u(w chan string) {
	w <- "#EXTM3U\n"
}

func insertM3u(w chan string, group, name, stream_url string) {
	m3uinfo := fmt.Sprintf("#EXTINF:-1 group-title=\"%s\",%s\n", group, name)

	m3uinfo += m3ugprops
	m3uinfo += stream_url + "\n"
	w <- m3uinfo
}
