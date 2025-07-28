package main

import (
	"bufio"
	"net/url"
	"strings"
	"unicode/utf8"
)

func checkM3U8Status(indexURL string) (ok bool) {
	// Therr are 3 paths need to be checked:
	// index.m3u8 -> resolution's m3u8 -> actual ts itself

	// Step 1: Get the index.m3u8 content
	resp, err := hc.Get(indexURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return
	}

	// Step 2: Parse variant playlists (resolutions)
	scanner := bufio.NewScanner(resp.Body)
	var resolutionM3U8 string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		// Assume first non-comment is the variant m3u8
		resolutionM3U8 = line
		break
	}

	if resolutionM3U8 == "" {
		return
	}

	// Handle relative URL
	resolutionM3U8URL := joinURL(indexURL, resolutionM3U8)

	// Step 3: Get resolution m3u8 content
	resp, err = hc.Get(resolutionM3U8URL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return
	}

	// Step 4: Get first TS segment
	scanner = bufio.NewScanner(resp.Body)
	var firstTS string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		if !utf8.ValidString(line) { // That's it. You win. Fucker.
			return true
		}

		firstTS = line
		break
	}

	if firstTS == "" {
		return
	}

	firstTSURL := joinURL(resolutionM3U8URL, firstTS)

	// Step 5: Check TS segment status
	resp, err = hc.Head(firstTSURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return
	}

	return true
}

// Join relative paths to base URLs
func joinURL(baseURL, ref string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return ref
	}
	refURL, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	return u.ResolveReference(refURL).String()
}
