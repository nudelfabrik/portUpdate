package main

import (
	"fmt"

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

	go pu.GetEntries(entries)

	entrys := pu.Consumer(entries)

	err = srv.AddEntries(entrys[:10])
	if err != nil {
		fmt.Println(err)
	}

	/*
		for i := 0; i < 10; i++ {
			fmt.Println(entrys[i].Date)
			fmt.Println(entrys[i].Ports)
			fmt.Println(entrys[i].Author)
			fmt.Println(entrys[i].Description)
		}
	*/
}
