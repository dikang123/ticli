package ticli

import (
	"testing"
)

func TestTiCli(t *testing.T) {
	opt := &Option{
		Addresses: []string{"127.0.0.1:4000"},
		User:      "tidb",
		Password:  "tidb",
		DB:        "demo",
		Timeout:   3,
	}
	cli := NewClient(opt)
	defer cli.Close()

	db, err := cli.Open()
	if err != nil {
		t.Fatalf("open db error: %s", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		t.Errorf("ping db error: %s", err)
	}
}
