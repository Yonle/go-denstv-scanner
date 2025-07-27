package main

func prepareM3u(w chan string) {
	w <- "#EXTM3U\n"
	w <- "#PLAYLIST: Dens.TV\n"
}

func insertM3u(w chan string, name, url string) {
	w <- "#EXTINF:-1, " + name + "\n"
	w <- url + "\n"
}
