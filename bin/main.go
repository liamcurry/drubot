package main

import (
	"flag"
	"github.com/liamcurry/drubot"
	"os"
)

var (
	uri  = flag.String("uri", "", "Hostname and port")
	nick = flag.String("nick", "drubot", "(optional) Name of the bot")
	pass = flag.String("pass", "", "(optional) Server password")
)

func main() {
	flag.Parse()

	rooms := flag.Args()

	if *uri == "" || len(rooms) == 0 {
		flag.PrintDefaults()
		os.Exit(2)
	}

	b := drubot.Bot{Nick: *nick, Rooms: rooms}
	b.Connect(*uri, *pass)
}
