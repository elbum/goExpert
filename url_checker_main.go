package main

import (
	"errors"
	"fmt"
	"net/http"
)

var errRequestFailed = errors.New("Request Failed")

type result struct {
	url    string
	status string
}

func main() {
	// 1
	// var results = map[string]string{}

	// 2
	// results := make(map[string]string)
	c := make(chan result)
	urls := []string{
		"https://www.airbnb.com/",
		"https://www.google.com/",
		"https://www.amazon.com/",
		"https://www.instagram.com/",
		"https://soundflare.com",
		"https://heiuasdasdfg.com",
		"http://facebook.com/",
		"http://reddit.com",
	}

	for _, url := range urls {
		go hitURL(url, c)
	}

	for i := 0; i < len(urls); i++ {
		fmt.Println(<-c)
	}

}

func hitURL(url string, c chan<- result) { // send only

	fmt.Println("Checking :", url)

	if resp, err := http.Get(url); err != nil || resp.StatusCode >= 400 {
		fmt.Println(url, " error")
		c <- result{url: url, status: "failed"}
	} else {
		c <- result{url: url, status: "OK"}
	}

}
