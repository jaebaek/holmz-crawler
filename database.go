package main

import (
	"fmt"
	"strings"
	"io/ioutil"
	"time"

	// redis client
	"github.com/garyburd/redigo/redis"
)

type DBError string

func (e DBError) Error() string {
	return string(e)
}

const SEED = "seed"
const HTML = "html"

var redisPool *redis.Pool

func DBInit() error {

	// init redisPool
    redisPool = &redis.Pool{
        MaxIdle:     conf.Nthread,
        IdleTimeout: 240 * time.Second,
        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", conf.Redis[0])
            if err != nil {
                return nil, err
            }
			/*
            if _, err := c.Do("AUTH", password); err != nil {
                c.Close()
                return nil, err
            }
			*/
            return c, err
        },
    }

	// get a connection
	conn := redisPool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return err
	}

	// get seed urls
	if exists, err := redis.Int(conn.Do("EXISTS", SEED));
	exists == 1 || err != nil {
		if err != nil {
			return err
		}
		Dbg.V("Try Del %v\n", SEED)
		if del, err2 := redis.Int(conn.Do("DEL", SEED));
		del != 1 || err2 != nil {
			if err2 != nil {
				return err2
			}
			return DBError(fmt.Sprintf("Del %v fails\n", SEED))
		}
	}
	Dbg.V("insert: %v from %v\n", SEED, conf.Seed)
	tmp, err := ioutil.ReadFile(conf.Seed)
	if err != nil {
		return err
	}

	// insert seed urls
	seed := strings.Split(string(tmp), "\n")
	for _, s := range seed {
		if strings.Index(s, "http") == 0 {
			Dbg.V("insert seed url \"%v\"\n", s)
			if _, err := conn.Do("LPUSH", SEED, s); err != nil {
				return err
			}
		}
	}

	return nil
}
