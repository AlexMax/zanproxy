package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/hpcloud/tail"
)

const banMessage = "You have been banned on suspicion of proxy use.  If you believe this is in error, please contact the administrators."

var ipIntel = NewIpIntel()

var reTimestamp = regexp.MustCompile(`^[0-9:;\[\]]+ `)
var reConnect = regexp.MustCompile(`Connect \(v[0-9.]+\): ([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)`)

func addBan(ip string, score float64) error {
	fmt.Printf("addBan(%s, %f)\n", ip, score)

	file, err := os.OpenFile("banlist.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			// Ban does not exist...
			break
		} else if err != nil {
			// We got an unexpected error, bail out with an error.
			return err
		}

		if strings.HasPrefix(line, ip) {
			// Ban exists, do nothing.
			return nil
		}
	}

	// Ban does not exist, append it.
	_, err = file.WriteString(fmt.Sprintf("%s:%s\n", ip, banMessage))
	if err != nil {
		return err
	}

	return nil
}

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
			continue
		}

		err = addBan(ip, score)
		if err != nil {
			log.Printf("addBan error: %#v", err)
			continue
		}
	}
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal("no arguments")
	}

	for _, arg := range flag.Args() {
		t, err := tail.TailFile(arg, tail.Config{
			Follow: true,
			Location: &tail.SeekInfo{
				Offset: 0,
				Whence: os.SEEK_END,
			},
			MustExist: true,
		})
		if err != nil {
			log.Fatal(err)
		}

		go parseTail(t)
	}

	select {}
}
