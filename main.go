package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func main() {

	var domains []string

	var full bool
	flag.BoolVar(&full, "full", false, "show full info")

	flag.Parse()

	if flag.NArg() > 0 {
		domains = []string{flag.Arg(0)}
	} else {

		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			domains = append(domains, sc.Text())
		}

		if err := sc.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
		}
	}

	fetchFns := []fetchFn{
		getAlienvaultURLs,
	}

	for _, domain := range domains {

		var wg sync.WaitGroup
		wurls := make(chan wurl)

		for _, fn := range fetchFns {
			wg.Add(1)
			fetch := fn
			go func() {
				defer wg.Done()
				resp, err := fetch(domain)
				if err != nil {
					return
				}

				for _, r := range resp {
					wurls <- r
				}
			}()
		}

		go func() {
			wg.Wait()
			close(wurls)
		}()

		seen := make(map[string]bool)
		for w := range wurls {
			if _, ok := seen[w.url]; ok {
				continue
			}
			seen[w.url] = true

			if full {
				fmt.Printf("%s %s %d\n", w.date, w.url, w.httpcode)
			} else {
				fmt.Println(w.url)
			}
		}
	}
}

type wurl struct {
	date     string
	url      string
	httpcode int
}

type alienvault struct {
	HasNext    bool       `json:"has_next"`
	ActualSize int        `json:"actual_size"`
	URLs       []url_list `json:"url_list"`
	PageNum    int        `json:"page_num"`
	Limit      int        `json:"limit"`
	FullSize   int        `json:"full_size"`
	Paged      bool       `json:"paged"`
}

type url_list struct {
	Date       string `json:"date"`
	URL        string `json:"url"`
	StatusCode int    `json:"httpcode"`
}

type fetchFn func(string) ([]wurl, error)

func getAlienvaultURLs(domain string) ([]wurl, error) {
	out := make([]wurl, 0)

	wrapper := alienvault{}
	wrapper.HasNext = true
	page := 1

	for wrapper.HasNext == true {
		client := &http.Client{}

		// 10/27/2020 - 'limit' and 'page' params are switched in the API
		url := fmt.Sprintf("https://otx.alienvault.com/otxapi/indicator/domain/url_list/%s?limit=%d&page=50", domain, page)
		req, err := http.NewRequest("GET", url, nil)

		req.Header.Add("X-OTX-API-KEY", os.Getenv("OTX"))

		res, err := client.Do(req)

		if err != nil {
			return []wurl{}, err
		}

		raw, err := ioutil.ReadAll(res.Body)

		res.Body.Close()
		if (err != nil) || strings.Contains(string(raw), "malformed") || strings.Contains(string(raw), "endpoint not found") {
			return []wurl{}, err
		}

		if strings.Contains(string(raw), "Over throttling limit") {
			fmt.Println("[-] API Limit")
			os.Exit(1)
		}

		err = json.Unmarshal(raw, &wrapper)
		if err != nil {
			return []wurl{}, err
		}

		for _, url := range wrapper.URLs {
			out = append(out, wurl{date: url.Date, url: url.URL, httpcode: url.StatusCode})
		}

		page++
	}

	return out, nil

}

func isSubdomain(rawUrl, domain string) bool {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return false
	}

	return strings.ToLower(u.Hostname()) != strings.ToLower(domain)
}
