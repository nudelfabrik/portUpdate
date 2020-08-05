package main

import (
	pu "github.com/nudelfabrik/portUpdate"
	"github.com/nudelfabrik/portUpdate/server"
)

func main() {
	pu.RegexCompile()

	entries := make(chan pu.Entry)

	go pu.GetEntries(entries)

	entrys := pu.Consumer(entries)

	srv, _ := server.NewServer(entrys)
	srv.Start()
}
