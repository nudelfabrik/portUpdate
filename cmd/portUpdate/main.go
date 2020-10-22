package main

import (
	"fmt"
	"time"

	pu "github.com/nudelfabrik/portUpdate"
	"github.com/nudelfabrik/portUpdate/postgres"
)

func main() {
	pu.RegexCompile()

	srv, err := postgres.NewBackendService()
	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		fmt.Println(err)
	}

	entries := make(chan pu.Entry)

	go pu.GetEntries(time.Time{}, entries)

	entrys := pu.Consumer(entries)

	err = srv.AddEntries(entrys[:10])
	if err != nil {
		fmt.Println(err)
	}
}
