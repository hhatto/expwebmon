package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Collector struct {
	sync.RWMutex
	Keys     []string               // expvar keys
	Data     map[string]*ExpvarData // key: url
	Interval time.Duration
	URLs     []string
}

func (c *Collector) Run() {
	for {
		wg := &sync.WaitGroup{}
		for _, url := range c.URLs {
			wg.Add(1)
			c.CollectExpvar(url, wg)
		}

		wg.Wait()
		time.Sleep(c.Interval)
	}
}

func (c *Collector) getResponseJsonData() string {
	var ret []byte
	var err error
	exp := new(ExpvarResponseData)

	c.RLock()
	defer c.RUnlock()
	for _, data := range c.Data {
		data.getFlattenData(c.Keys, exp)
		ret, err = json.Marshal(exp)
		if err != nil {
			log.Printf("json encode error!! err: %v", err)
			return ""
		}
	}

	return string(ret)
}

func (c *Collector) CollectExpvar(url string, wg *sync.WaitGroup) {
	x, err := getExpvarData(url)
	if err != nil {
		return
	}

	c.Lock()
	c.Data[url] = x
	c.Unlock()

	wg.Done()
}

func (c *Collector) ServeHTTP(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ret := c.getResponseJsonData()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", ret)
}
