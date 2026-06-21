package database

import (
	"testing"

	"argument/internal/config"
)

func setupTest() {
	Disconnect()
	Connect("mysql", config.Conf.Database.URI)
}

func TestBrokenDB(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("If database is unreachable we should crash")
		}
	}()

	Disconnect()
	Connect("root:asdf@tcp(127.0.0.1:12345)/argument")
}

func TestConnectingTwice(t *testing.T) {
	Connect("mysql", config.Conf.Database.URI)
}
