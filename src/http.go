package main

import "net/http"

type HTTPClient struct {
	UserAgent string
	Referer   string
	Client    http.Client
}

func (h *HTTPClient) DoReq(method, url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", h.UserAgent)
	req.Header.Set("Referer", h.Referer)
	return h.Client.Do(req)
}

func (h *HTTPClient) Get(url string) (resp *http.Response, err error) {
	return h.DoReq("GET", url)
}

func (h *HTTPClient) Head(url string) (resp *http.Response, err error) {
	return h.DoReq("HEAD", url)
}
