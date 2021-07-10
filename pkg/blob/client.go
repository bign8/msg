package blob

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	msg "github.com/bign8/msg/pkg"
)

// NewClient constructs a new blob client
func NewClient(conn *msg.Conn) *Client {
	return &Client{conn: conn}
}

// Client is a blob client
type Client struct {
	who  string // TODO: expire after long durations
	conn *msg.Conn
}

func (c *Client) findService(ctx msg.ContextOld) (addr string, err error) {
	if c.who == "" {
		var res *msg.Msg
		res, err = c.conn.Request(ctx, &msg.Msg{
			Title: "sd_lookup",    // service discovery
			Body:  []byte("blob"), // name of the service to lookup
		})
		if err == nil {
			c.who = string(res.Body)
		}
	}

	return c.who, err
}

// Load gets a large binary payload
func (c *Client) Load(ctx msg.ContextOld, id string) (bits []byte, err error) {
	addr, err := c.findService(ctx)
	if err != nil {
		return nil, err
	}
	res, err := http.Get(addr + "/load/" + id)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("blob: non-200: " + res.Status)
	}
	bits, err = ioutil.ReadAll(res.Body)
	if err == nil {
		err = res.Body.Close()
	}
	return bits, err
}

// Save stores a large binary payload
func (c *Client) Save(ctx msg.ContextOld, bits []byte) (id string, err error) {
	addr, err := c.findService(ctx)
	if err != nil {
		return "", err
	}
	res, err := http.Post(addr+"/save", "application/octet-stream", bytes.NewReader(bits))
	if err != nil {
		return "", err
	}
	bits, err = ioutil.ReadAll(res.Body)
	if err == nil {
		err = res.Body.Close()
	}
	return string(bits), err
}
