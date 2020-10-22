package portUpdate

import (
	"bufio"
	"os"
	"testing"
	"time"
)

var port1 = Entry{time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC), []string{"test/port1"}, "port1@freebsd.org", "portdescription1"}

func TestGetEntriesFromScanner(t *testing.T) {

	RegexCompile()

	tests := []struct {
		filename string
		time     time.Time
		entry    []Entry
	}{
		{"test1", time.Time{}, []Entry{port1}},
	}

	for _, test := range tests {

		t.Run(test.filename, func(t *testing.T) {
			body, err := os.Open("./tests/" + test.filename)
			if err != nil {
				t.Fatal("file does not exist: ", err)
			}
			bodyScan := bufio.NewScanner(body)

			entries := make(chan Entry)

			go GetEntriesFromScanner(test.time, entries, bodyScan)

			allEntries := Consumer(entries)

			if len(test.entry) != len(allEntries) {
				t.Fatalf("mismatch of entries: want %d, got %d", len(test.entry), len(allEntries))
			}

			for i, e := range test.entry {
				if e.Date != allEntries[i].Date {
					t.Errorf("Date mismatch: want %v, got %v", e.Date, allEntries[i].Date)
				}
				if e.Author != allEntries[i].Author {
					t.Errorf("Author mismatch: want %s, got %s", e.Author, allEntries[i].Author)
				}
				if e.Description != allEntries[i].Description {
					t.Errorf("Description mismatch: want %s, got %s", e.Description, allEntries[i].Description)
				}
				if len(e.Ports) != len(allEntries[i].Ports) {
					t.Fatalf("mismatch of Ports: want %d, got %d", len(test.entry), len(allEntries))
				}
				for j, p := range e.Ports {
					if p != allEntries[i].Ports[j] {
						t.Errorf("Port Mismatch : %dth Port is %s, want %s", j, p, allEntries[i].Ports[j])
					}
				}
			}

		})

	}
}
