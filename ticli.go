package ticli

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Option struct {
	Addresses []string // TiDB addresses
	User      string   // Auth user
	Password  string   // Auth password
	DB        string   // Select database
	Timeout   int      // Timeout in seconds
}

type Client struct {
	ready         []string
	notReady      []string
	readyIndex    int
	notReadyIndex int
	closed        bool
	lock          *sync.Mutex
	option        *Option
}

func NewClient(opt *Option) *Client {
	opt.init()
	c := &Client{
		lock:   &sync.Mutex{},
		option: opt,
	}
	c.startChecker()
	return c
}

func (o *Option) init() {
	if o.Timeout == 0 {
		o.Timeout = 3
	}
}

func (c *Client) isServerReady(addr string) bool {
	connUri := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=2s", c.option.User, c.option.Password, addr, c.option.DB)
	db, err := sql.Open("mysql", connUri)
	if err != nil {
		return false
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return false
	}
	return true
}

func (c *Client) startChecker() {
	for _, addr := range c.option.Addresses {
		if c.isServerReady(addr) {
			c.ready = append(c.ready, addr)
		} else {
			c.notReady = append(c.notReady, addr)
		}
	}

	// check loop
	go func() {
		for !c.closed {
			readyCount := len(c.ready)
			for i := 0; i < readyCount; i++ {
				if !c.isServerReady(c.ready[i]) {
					c.lock.Lock()
					c.notReady = append(c.notReady, c.ready[i])
					c.ready = append(c.ready[0:i], c.ready[i+1:]...)
					readyCount -= 1
					c.lock.Unlock()
				}
			}
			notReadyCount := len(c.notReady)
			for i := 0; i < notReadyCount; i++ {
				if c.isServerReady(c.notReady[i]) {
					c.lock.Lock()
					c.ready = append(c.ready, c.notReady[i])
					c.notReady = append(c.notReady[0:i], c.notReady[i+1:]...)
					notReadyCount -= 1
					c.lock.Unlock()
				}
			}

			time.Sleep(1 * time.Second)
		}
	}()
}

func (c *Client) Close() {
	c.closed = true
}

func (c *Client) Open() (*sql.DB, error) {
	c.lock.Lock()
	if len(c.ready) == 0 {
		c.lock.Unlock()
		return nil, fmt.Errorf("no server available")
	}
	if c.readyIndex >= len(c.ready) {
		c.readyIndex = 0
	}
	addr := c.ready[c.readyIndex]
	c.readyIndex += 1
	c.lock.Unlock()

	connUri := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%ds", c.option.User, c.option.Password, addr, c.option.DB, c.option.Timeout)
	return sql.Open("mysql", connUri)
}
