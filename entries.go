package portUpdate

import (
	"bufio"
	"os"
	"sync"
)

type Entry struct {
	Date        string
	Ports       []string
	Author      string
	Description string
}

func GetEntries(entries chan Entry) {
	/*
		resp, err := http.Get("https://svnweb.freebsd.org/ports/head/UPDATING?view=co")
		if err != nil {
			fmt.Println("Error Downloading UPDATING file", err)
		}
		defer resp.Body.Close()

		bodyScan := bufio.NewScanner(resp.Body)
	*/
	body, _ := os.Open("./UPDATING")
	bodyScan := bufio.NewScanner(body)

	// Looks for the next Entry by looking for the date
	// Token is one complete Entry. Including date until the beginning of the next date
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// Last Token. Otherwise an the next line causes a out of bounds
		if len(data) < 10 {
			return 0, nil, nil
		}
		// The first bytes of data are always a match, but we want to get the beginning
		// of the next entry.
		loc := regexDate.FindIndex(data[9:])
		if loc != nil {
			return loc[0] + 9, data[0 : loc[0]+8], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil

	}
	bodyScan.Split(split)

	// Ignore the header Text
	bodyScan.Scan()

	var wg sync.WaitGroup

	go Consumer(entries)
	bodyScan.Scan()

	for bodyScan.Scan() {
		wg.Add(1)
		go Parse(&wg, entries, bodyScan.Text())
	}
	wg.Wait()
	close(entries)
}

func Parse(wg *sync.WaitGroup, entries chan Entry, data string) {
	defer wg.Done()

	var e Entry
	date := regexDate.FindStringSubmatch(data)
	if date != nil {
		e.Date = date[1]
	}

	authors := regexAuthor.FindStringSubmatch(data)
	if authors != nil {
		e.Author = authors[1]
	}

	affects := regexAffectsLine.FindString(data)
	e.Ports = regexAffects.FindAllString(affects, -1)

	e.Description = regexDescr.FindStringSubmatch(data)[1]
	entries <- e
}

func Consumer(entries chan Entry) []Entry {
	entrys := make([]Entry, 0, 10)

	for e := range entries {
		entrys = append(entrys, e)
	}

	return entrys
}
