package main

import (
	"fmt"
	// "io/ioutil"
	"log"
	"net/http"
 )
 
 func main() {
	test_urls := []string{
		"https://exampleio",
	}

	failed400Urls := []string{}
	lower400Urls := []string{}
	for _, test_url := range test_urls {
		resp, err := http.Get(test_url)
		if err != nil {
			log.Fatalln(fmt.Errorf("hhtp.get failed: %s", err))
		}

		if resp.StatusCode >= 500 && resp.StatusCode <= 600 {
			fmt.Printf("failed 5xx for %s\n  resp: %+v\n\n", test_url, resp)
		}

		if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
			failed400Urls = append(failed400Urls, test_url);
		}

		if resp.StatusCode <= 400 {
			lower400Urls = append(lower400Urls, fmt.Sprintf("%d %s", resp.StatusCode, test_url));
		}
	}
	
	for _, failed400Url := range failed400Urls {
		fmt.Printf("failed 4xx for %s\n", failed400Url)
	}

	fmt.Println("")

	fmt.Println("lower 400 urls:")
	for _, lower400Url := range lower400Urls {
		fmt.Printf("  %s\n", lower400Url)
	}
}
