package main

import (
	"os"

	// toml parser
	"github.com/BurntSushi/toml"
)

type Config struct {
	Nthread int
	Seed string
	Redis []string	// addr (e.g., "localhost:6379")
                    // password (e.g., "")
}
var conf Config

func main() {
	// get config file
	conf_toml := ""
	switch len(os.Args) {
		case 1: conf_toml = "conf.toml"
		case 2: conf_toml = os.Args[1]
		default: {
			Dbg.E("Usage: %v [conf.toml]\n", os.Args[0])
		}
	}

	// parse config toml
	if _, err := toml.DecodeFile(conf_toml, &conf); err != nil {
		Dbg.E("TOML error: %v\n", err)
		return
	}
	Dbg.I("Conf: %v\n", conf)

	// open redis and insert seed urls
	if err := DBInit(); err != nil {
		Dbg.E("DB error: %v\n", err)
		return
	}

	// spawn work threads
	done := make(chan int, conf.Nthread)
	for i := 0;i < conf.Nthread; i++ {
		go run(i, done)
	}

	// wait until crawl done
	for i := 0;i < conf.Nthread; i++ {
		<-done
	}
}
