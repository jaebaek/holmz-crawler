package main

import (
	"net/http"
	"bytes"
	"strings"
	"sync/atomic"

	// html parser
	"golang.org/x/net/html"

	// redis client
	"github.com/garyburd/redigo/redis"
)

// root url id
var RUID uint64 = 0
func getRUID() uint64 {
	return atomic.AddUint64(&RUID, 1)
}

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

// Extract all http** links from a given webpage
func visit(url string, ruid uint64, c redis.Conn) {
	// check if already visited
	if _, err := redis.String(c.Do("HGET", HTML, url)); err == nil {
		Dbg.V("HGET html %v: %v\n", url, err)
		return
	}
	Dbg.V("visit: %v\n", url)

	// get html page
	resp, err := http.Get(url)
	if err != nil {
		Dbg.E("ERROR: Failed to crawl \"%v\"\n", url)
		Dbg.E("ERROR: %v\n", err)
		return
	}

	b := resp.Body
	defer b.Close()

	// get links
	var htmlBuf bytes.Buffer
	z := html.NewTokenizer(b)
	for {
		tt := z.Next()
		htmlBuf.Write(z.Raw())

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			page := htmlBuf.String()
			c.Do("HSET", HTML, url, page)
			Dbg.V(page)
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			if t.Data != "a" {
				continue
			}

			// Extract the href value, if there is one
			if ok, link := getHref(t); ok {
				if strings.Index(link, "http") != 0 && len(link) != 0 && link[0] != '#' && link != "./" && link != "/" {

					// ruid -- LIST [children ...]
					if _, err := c.Do("LPUSH", ruid, link); err != nil {
						Dbg.E("LPUSH error: %v\n", err)
					}
				}
			}
		}
	}
}

func run(tid int, done chan int) {

	// open new connection
	c := redisPool.Get()
	defer c.Close()

	// start from root url
	for {
		url, _ := redis.String(c.Do("LPOP", SEED))
		if url == "" {
			break
		}

		ruid := getRUID()
		visit (url, ruid, c)
		for {
			child, _ := redis.String(c.Do("LPOP", ruid))
			if child == "" {
				break
			}
			visit (url + child, ruid, c)
		}
	}
	Dbg.I("Thread #%v done\n", tid)
	done <- tid
}
