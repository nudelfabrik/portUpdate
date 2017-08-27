package main

import (
	"fmt"
	"regexp"
)

var RegexDate *regexp.Regexp
var RegexAuthor *regexp.Regexp
var RegexAffectsLine *regexp.Regexp
var RegexAffects *regexp.Regexp
var RegexDescr *regexp.Regexp

func main() {
	RegexDate = regexp.MustCompile(`(([0-9]{8})):\n`)
	RegexAuthor = regexp.MustCompile(`AUTHORS?: (.*)\n`)
	RegexAffectsLine = regexp.MustCompile(`AFFECTS:.*\n`)
	RegexAffects = regexp.MustCompile(`[A-Za-z0-9-*]*/[A-Za-z0-9-*]*`)
	RegexDescr = regexp.MustCompile(`^.*\n.*\n.*\n.*\n((?s).*)`)

	entries := make(chan Entry)

	go GetEntries(entries)

	entrys := Consumer(entries)

	for i := 0; i < 10; i++ {
		fmt.Println(entrys[i].Date)
		fmt.Println(entrys[i].Ports)
	}

}
