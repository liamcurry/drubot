package drubot

import (
	"github.com/liamcurry/drubot/plugins"
	"regexp"
	"strings"
)

// A chat bot
type Bot struct {
	Nick  string
	Rooms []string
	Conn  *Connection
	init  bool
}

// Connect to a server
func (b *Bot) Connect(uri string, pass string) (err error) {
	b.Conn, err = NewConnection(uri)
	if err != nil {
		return err
	}

	b.Conn.Write <- &Message{Command: "NICK", Params: b.Nick}
	b.Conn.Write <- &Message{Command: "MODE", Params: "+Bix " + b.Nick}
	b.Conn.Write <- &Message{
		Command:  "USER",
		Params:   b.Nick + " * *",
		Trailing: "a robot",
	}
	if pass != "" {
		b.Conn.Write <- &Message{Command: "PASS", Params: pass}
	}

	b.Listen()

	return
}

// Listen for incoming messages
func (b *Bot) Listen() {
	for {
		msg := <-b.Conn.Read
		switch msg.Command {
		case "PING":
			go b.handlePing(msg)
		case "MODE":
			go b.handleMode(msg)
		case "PRIVMSG":
			go b.handlePrivmsg(msg)
		}
	}
}

func (b *Bot) handlePing(msg *Message) {
	b.Conn.Write <- &Message{Command: "PONG", Trailing: msg.Trailing}
}

func (b *Bot) handleMode(msg *Message) {
	// Join rooms after receiving the first MODE response
	if !b.init && msg.Prefix != "" {
		b.init = true
		b.Conn.Write <- &Message{
			Command: "JOIN",
			Params:  strings.Join(b.Rooms, ","),
		}
	}
}

var cmdRe = regexp.MustCompile(
	`^(?P<nick>\w+)?\s*(?P<short>!)?(?P<cmd>\w+)\s+(?P<args>.+)$`,
)

func (b *Bot) handlePrivmsg(msg *Message) {
	match := cmdRe.FindStringSubmatch(msg.Trailing)

	if len(match) == 0 || match[1] != b.Nick && match[2] == "" {
		return
	}

	for _, p := range plugins.Plugins {
		for _, name := range p.Names {
			if match[3] == name {
				go func(p plugins.Plugin) {
					b.Conn.Write <- &Message{
						Command:  "PRIVMSG",
						Params:   msg.Params,
						Trailing: p.Run(&name, &match[4]),
					}
				}(p)
			}
		}
	}
}
