package drubot

import (
	"bufio"
	"net"
	"regexp"
	"strings"
)

// Connection to an IRC server
type Connection struct {
	URI    string        // irc.freenode.net:6667
	Read   chan *Message // from server
	Write  chan *Message // to server
	reader *bufio.Reader
	writer *bufio.Writer
}

// NewConnection creates a new IRC server connection
func NewConnection(uri string) (c *Connection, err error) {
	println("connecting")
	c = new(Connection)
	c.URI = uri

	// Connect to the server
	addr, err := net.ResolveTCPAddr("tcp", c.URI)
	if err != nil {
		return nil, err
	}

	// Dial the socket
	socket, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	// Create reader and writer objects for the socket
	c.reader = bufio.NewReader(socket)
	c.writer = bufio.NewWriter(socket)

	// Create channels for the reader and writer to communicate on
	c.Read = make(chan *Message, 1000)
	c.Write = make(chan *Message, 1000)

	go c.readIncoming()
	go c.sendOutgoing()

	return c, nil
}

func (c *Connection) readIncoming() {
	for {
		// Checks for new messages to receive
		raw, _ := c.reader.ReadString(byte('\n'))

		// Bots are in charge of message reading
		m := NewMessage(raw)
		c.Read <- m
		println("<- " + m.Raw())
	}
}

func (c *Connection) sendOutgoing() {
	for {
		msg := <-c.Write
		println("-> " + msg.Raw())
		c.writer.WriteString(msg.Raw() + "\r\n")
		c.writer.Flush()
	}
}

// Message from the server
// :<prefix> <command> <params> :<trailing>
type Message struct{ Prefix, Command, Params, Trailing, raw, args string }

var msgRe = regexp.MustCompile(
	`^(?::(?P<prefix>[^\s]+) )?(?P<command>\w+)(?: (?P<params>[^:]+))?(?: ?:(?P<trailing>.+))?$`,
)

// NewMessage creates a new Message based on a string
func NewMessage(raw string) (m *Message) {
	m = new(Message)
	m.raw = strings.Trim(raw, "\r\n")

	match := msgRe.FindStringSubmatch(m.raw)
	m.Prefix = match[1]
	m.Command = match[2]
	m.Params = match[3]
	m.Trailing = match[4]

	return m
}

// Raw builds a raw string that can be sent to the server
func (m *Message) Raw() (raw string) {
	if m.raw == "" {
		if m.Prefix != "" {
			raw += ":%s "
		}
		raw += strings.TrimSpace(m.Command)
		if m.Params != "" {
			raw += " " + strings.TrimSpace(m.Params)
		}
		if m.Trailing != "" {
			raw += " :" + strings.TrimSpace(m.Trailing)
		}
		m.raw = raw
	}
	return m.raw
}
