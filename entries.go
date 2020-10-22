package portUpdate

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type Entry struct {
	Date        time.Time
	Ports       []string
	Author      string
	Description string
}

func GetEntries(lastUpdate time.Time, entries chan Entry) {
	/*
		resp, err := http.Get("https://svnweb.freebsd.org/ports/head/UPDATING?view=co")
		if err != nil {
			fmt.Println("Error Downloading UPDATING file", err)
		}
		defer resp.Body.Close()

		bodyScan := bufio.NewScanner(resp.Body)
	*/
	body, _ := os.Open("./UPDATING2")
	bodyScan := bufio.NewScanner(body)

	GetEntriesFromScanner(lastUpdate, entries, bodyScan)

}

// Retrieves all Entries from the Scanner on of after the date specified in lastUpdate.
// All Entries before that date are assumed to already being stored in the BackendService.
// Entries at exactly the day of lastUpdate have to be checked seprarately before adding them to the database.
func GetEntriesFromScanner(lastUpdate time.Time, entries chan Entry, bodyScan *bufio.Scanner) {

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

	// If no last Time is given, parse all entries
	if lastUpdate.IsZero() {
		for bodyScan.Scan() {
			wg.Add(1)
			go parse(&wg, &lastUpdate, entries, bodyScan.Text())
		}
	}
	wg.Wait()
	close(entries)
}

func parseDate(data string) time.Time {
	date := regexDate.FindStringSubmatch(data)
	if date != nil {
		datum, err := time.Parse("20060102", date[1])
		if err == nil {
			return datum
		}
	}
	return time.Time{}
}

func parse(wg *sync.WaitGroup, lastUpdate *time.Time, entries chan Entry, data string) {
	defer wg.Done()

	var e Entry

	e.Date = parseDate(data)

	// If the Date is before lastUpdate,
	// it was already present the last time portUpdate updated and can be skipped
	if e.Date.Before(*lastUpdate) {
		return
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
