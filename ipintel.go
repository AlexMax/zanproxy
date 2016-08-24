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
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type ipCacheEntry struct {
	score   float64
	touched time.Time
}

// IPIntel is a cache of results from getIPIntel.
type IPIntel struct {
	cache   map[string]ipCacheEntry
	mutex   sync.RWMutex
	timeout time.Duration
}

func getScore(addr string) (float64, error) {
	v := url.Values{}

	v.Set("ip", addr)
	v.Set("contact", "alexmax2742@gmail.com")
	v.Set("flags", "f")

	res, err := http.Get("http://check.getipintel.net/check.php?" + v.Encode())
	if err != nil {
		return 0.0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0.0, err
	}

	score, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		return 0.0, err
	}

	return score, nil
}

// NewIPIntel creates a new instance of IPIntel.
func NewIPIntel() *IPIntel {
	ipintel := &IPIntel{
		cache:   make(map[string]ipCacheEntry),
		timeout: time.Minute * 30,
	}

	return ipintel
}

// GetScore retrieves a score from getIPintel.
//
// The score is on a scale from 0.0 to 1.0.  If the score was retrieved
// from the cache, the second return value will return true.
func (intel *IPIntel) GetScore(addr string) (float64, bool, error) {
	intel.mutex.RLock()
	cache, ok := intel.cache[addr]
	intel.mutex.RUnlock()

	if !ok || cache.touched.Add(intel.timeout).Before(time.Now()) {
		// Grab the score from getIPintel.
		score, err := getScore(addr)
		if err != nil {
			return 0.0, false, err
		}

		// Save score to the cache.
		intel.mutex.Lock()
		defer intel.mutex.Unlock()
		intel.cache[addr] = ipCacheEntry{
			score:   score,
			touched: time.Now(),
		}

		// Return fresh score.
		return score, false, nil
	}

	// Return cached value
	intel.mutex.Lock()
	defer intel.mutex.Unlock()
	cache.touched = time.Now()

	return cache.score, true, nil
}
