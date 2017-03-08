Simple Web Crawler
====================

# Install
- Install [Redis](https://redis.io/)
- Install dependencies
```sh
$ go get golang.org/x/net                   # html parser
$ go get github.com/garyburd/redigo/redis   # redis client
$ go get github.com/BurntSushi/toml         # toml parser (for conf)
```
- Install crawler
```sh
$ # note that $GOPATH/src/crawler must contain
$ # this source code files.
$ go install crawler
```
(Please read [go-setup](ihttps://golang.org/doc/code.html) for $GOPATH setup)

# Run
- Configuration
    - conf.toml
```sh
nthread = 8
seed = "seed.txt"
redis = ["localhost:6379", ""]
```
    - seed.txt
```sh
http://cps.kaist.ac.kr
http://www.kaist.ac.kr
```
- Run
    - First, run redis server
    - Second, run crawler
```sh
$ $GOPATH/bin/crawler
```

# How crawler works
0. for each thread
1. root <-- pop(seed)
2. visit(root) and save links to ruid (i.e., root url id)
3. for link = pop(ruid) : visit(root + link) and save links to ruid
(except already visited or another root url)
