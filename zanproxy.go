/*
 *  zanproxy: a proxy detector for Zandronum
 *  Copyright (C) 2016  Alex Mayfield <alexmax2742@gmail.com>
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU Affero General Public License for more details.
 *
 *  You should have received a copy of the GNU Affero General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hpcloud/tail"
)

var config *Config

var ipIntel = NewIPIntel()

var reTimestamp = regexp.MustCompile(`^[0-9:;\[\]]+ `)
var reConnect = regexp.MustCompile(`Connect \(v[0-9.]+\): ([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)`)

func addBan(ip string, score float64) error {
	file, err := os.OpenFile(config.Banlist, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
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
			log.Printf("%s is greater than or equal to MinScore, already exists in banlist. (%f >= %f)", ip, score, config.MinScore)
			return nil
		}
	}

	// Ban does not exist, append it.
	_, err = file.WriteString(fmt.Sprintf("\n%s:%s", ip, config.BanMessage))
	if err != nil {
		return err
	}

	log.Printf("%s is greater than or equal to MinScore, added to banlist. (%f >= %f)", ip, score, config.MinScore)
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
		score, _, err := ipIntel.GetScore(ip, config.Email)
		if err != nil {
			log.Printf("ipIntel error: %#v", err)
			continue
		}

		// Don't add to banlist unless we meet the minimum score.
		if score < config.MinScore {
			log.Printf("%s is less than MinScore. (%f < %f)", ip, score, config.MinScore)
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
	if len(os.Args) != 2 {
		log.Print("Missing parameter - config file")
		os.Exit(1)
	}

	var err error
	config, err = NewConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for _, globarg := range config.Logfiles {
		matches, err := filepath.Glob(globarg)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		if matches == nil {
			log.Fatalf("No matches for path %s", globarg)
			os.Exit(1)
		}

		for _, arg := range matches {
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
	}

	select {}
}
