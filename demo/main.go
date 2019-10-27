package main

import (
	"log"
	"sync"
	"ticli"
	"time"
)

func main() {
	opt := &ticli.Option{
		Addresses: []string{"10.9.44.100:4000", "10.9.172.178:4000", "10.9.68.213:4000"},
		User:      "tidb",
		Password:  "tidb",
		DB:        "demo",
		Timeout:   3,
	}
	cli := ticli.NewClient(opt)
	defer cli.Close()

	wg := &sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go startThread(wg, i, cli)
	}
	wg.Wait()
}

func startThread(wg *sync.WaitGroup, id int, cli *ticli.Client) {
	defer wg.Done()
	for {
		err := doSomethingWithDB(id, cli)
		if err != nil {
			log.Printf("[Thread %d] reconnect db after 3s", id)
			time.Sleep(3 * time.Second)
		} else {
			log.Printf("[Thread %d] exit", id)
			return
		}
	}
}

func doSomethingWithDB(id int, cli *ticli.Client) error {
	db, err := cli.Open()
	if err != nil {
		log.Printf("[Thread %d] open db error: %s", id, err)
		return err
	}
	defer db.Close()

	for {
		err := db.Ping()
		if err != nil {
			log.Printf("[Thread %d] ping db error: %s", id, err)
			return err
		} else {
			log.Printf("[Thread %d] ping db ok", id)
		}
		time.Sleep(3 * time.Second)
	}
}
