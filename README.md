Simple Web Crawler
====================

* Base code is from [simple-webcrawler](https://schier.co/blog/2015/04/26/a-simple-web-scraper-in-go.html)

# Install
- Install [Redis](https://redis.io/)
- Install dependencies
```sh
$ go get golang.org/x/net   # html parser
$ go get gopkg.in/redis.v5  # redis client
```
- Install crawler
```sh
$ # note that $GOPATH/src/crawler must contain
$ # this source code files.
$ go install crawler
```
(Please read [go-setup](ihttps://golang.org/doc/code.html) for $GOPATH setup)

# Run
- Configuration --> TODO
- Run
```sh
$ $GOPATH/bin/crawler
```
