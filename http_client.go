package main

import "net/http"

type HttpClientHelper interface {
	Get(url string) (*http.Response, error)
}

type HttpClient struct {
	c HttpClientHelper
}

func CreateHttpClient() *HttpClient {
	return &HttpClient{
		c: &http.Client{},
	}
}

func (client *HttpClient) Get(url string) (*http.Response, error) {
	resp, err := client.c.Get(url)
	return resp, err
}
