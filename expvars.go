package main

import (
	"strings"

	"github.com/antonholmquist/jason"
)

type MemStats struct {
	Alloc int64 `json:"Alloc"`
}

type ExpvarData struct {
	*jason.Object
}

type ExpvarResponseData struct {
	Keys  []string          `json:"keys"`
	Datas map[string]string `json:"data"`
}

func (e *ExpvarData) getFlattenData(keys []string, resp *ExpvarResponseData) {
	resp.Datas = map[string]string{}
	resp.Keys = keys
	resp.Datas["cmdline"], _ = e.GetString("cmdline")
	for _, key := range keys {
		k := strings.Split(key, ".")
		valueStr, err := e.GetString(k...)
		if err == nil {
			resp.Datas[key] = valueStr
			continue
		}
		valueNumber, err := e.GetNumber(k...)
		if err == nil {
			resp.Datas[key] = valueNumber.String()
			continue
		}
	}
}
