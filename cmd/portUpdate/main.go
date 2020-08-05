package main

import (
	"fmt"

	pu "github.com/nudelfabrik/portUpdate"
)

func main() {

	pu.RegexCompile()

	entries := make(chan pu.Entry)

	go pu.GetEntries(entries)

	entrys := pu.Consumer(entries)

	for i := 0; i < 10; i++ {
		fmt.Println(entrys[i].Date)
		fmt.Println(entrys[i].Ports)
		fmt.Println(entrys[i].Author)
		fmt.Println(entrys[i].Description)
	}

}
