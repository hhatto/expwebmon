package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (c *Collector) indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if strings.HasSuffix(r.RequestURI, ".css") {
		w.Header().Set("Content-Type", "text/css")
	}

	filename := "./static/index.html"
	if r.RequestURI == "/" || r.RequestURI == "/index.html" {
		t := template.Must(template.ParseFiles(filename))
		err := t.Execute(w, c.Keys)
		if err != nil {
			log.Printf("template.ExecuteTemplate() error. err: %v", err)
		}
	} else {
		filename = fmt.Sprintf(".%s", r.RequestURI)
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(w, "index.html is not found. err: %v", err)
			return
		}
		fmt.Fprintf(w, "%s", b)
	}
	return
}

const DefaultVars = "memstats.Alloc,memstats.Sys,memstats.HeapAlloc,memstats.HeapInuse,memstats.PauseTotalNs,memstats.NumGC"

func parseVars(varsStr string) (ret []string) {
	s := strings.Split(varsStr, ",")
	for _, str := range s {
		ret = append(ret, str)
	}
	return
}

func parsePorts(portsStr string) (ret []string, ok bool) {
	ok = true
	s := strings.Split(portsStr, ",")
	for _, str := range s {
		if str == "" {
			return ret, false
		}
		ret = append(ret, str)
	}
	return
}

func Usage() {
	flag.PrintDefaults()
}

func main() {
	var vars = flag.String("vars", DefaultVars, "expvar keys")
	var fetchPorts = flag.String("fetchports", "", "Ports/URLs")
	var bindAddr = flag.String("bind", "localhost:9999", "host:port")
	flag.Parse()

	interval := 1000 * time.Millisecond

	c := Collector{
		Interval: interval,
	}

	portsStr, ok := parsePorts(*fetchPorts)
	if !ok {
		Usage()
		return
	}

	c.Data = map[string]*ExpvarData{}
	c.Keys = parseVars(*vars)

	hostAndPort := fmt.Sprintf("localhost:%s", portsStr[0])
	url := fmt.Sprintf("http://%s/debug/vars", hostAndPort)
	c.URLs = append(c.URLs, url)

	router := httprouter.New()
	router.GET("/", c.indexHandler)
	router.GET("/static/*filename", c.indexHandler)
	router.GET("/expvars", c.ServeHTTP)

	http.Handle("/", router)
	go http.ListenAndServe(*bindAddr, nil)
	log.Printf("listen on %v", *bindAddr)

	c.Run()
}
