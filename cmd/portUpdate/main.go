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

	for i := 0; i < 10; i++ {
		err = srv.AddEntries(entrys)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(entrys[i].Date)
		fmt.Println(entrys[i].Ports)
		fmt.Println(entrys[i].Author)
		fmt.Println(entrys[i].Description)
	}
}
