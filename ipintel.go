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

type IpIntel struct {
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

func NewIpIntel() *IpIntel {
	ipintel := &IpIntel{
		cache:   make(map[string]ipCacheEntry),
		timeout: time.Minute * 30,
	}

	return ipintel
}

func (intel *IpIntel) GetScore(addr string) (float64, bool, error) {
	intel.mutex.RLock()
	cache, ok := intel.cache[addr]
	intel.mutex.RUnlock()

	if !ok || cache.touched.Add(intel.timeout).Before(time.Now()) {
		score, err := getScore(addr)
		if err != nil {
			return 0.0, false, err
		}

		intel.mutex.Lock()
		defer intel.mutex.Unlock()
		intel.cache[addr] = ipCacheEntry{
			score:   score,
			touched: time.Now(),
		}

		return score, false, nil
	} else {
		intel.mutex.Lock()
		defer intel.mutex.Unlock()
		cache.touched = time.Now()

		return cache.score, true, nil
	}
}
