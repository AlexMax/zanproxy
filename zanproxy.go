package main

import (
	"flag"
	"log"
	"regexp"

	"github.com/hpcloud/tail"
)

var ipIntel = NewIpIntel()

var reTimestamp = regexp.MustCompile(`^[0-9:;\[\]]+ `)
var reConnect = regexp.MustCompile(`Connect \(v[0-9.]+\): ([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)`)

func parseTail(t *tail.Tail) {
	// Goroutine-specific regexes
	recTimestamp := reTimestamp.Copy()
	recConnect := reConnect.Copy()

	// Try and find a connecting IP in every single line.
	for line := range t.Lines {
		indexes := recTimestamp.FindStringIndex(line.Text)
		var testString string
		if indexes != nil {
			testString = line.Text[indexes[1]:]
		} else {
			testString = line.Text[:]
		}

		cGroups := recConnect.FindStringSubmatch(testString)
		if cGroups == nil {
			continue
		}
		ip := cGroups[1]

		// Get IP Intel on given IP.
		score, _, err := ipIntel.GetScore(ip)
		if err != nil {
			log.Printf("ipIntel error: %#v", err)
			return
		}

		log.Printf("GetScore(%s): %f", ip, score)
	}
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal("no arguments")
	}

	for _, arg := range flag.Args() {
		t, err := tail.TailFile(arg, tail.Config{
			MustExist: true,
			Follow:    true,
		})
		if err != nil {
			log.Fatal(err)
		}

		go parseTail(t)
	}

	select {}
}
