package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type config struct {
	Root string
	Port int
}

func main() {
	c := loadConfig()
	mux := http.FileServer(http.Dir(c.Root))
	err := http.ListenAndServe(":"+strconv.Itoa(c.Port), mux)
	if err != nil {
		panic(err)
	}
}

func loadConfig() config {
	b, err := ioutil.ReadFile("static_file_server.json")
	if err != nil {
		panic(err)
	}
	c := config{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(err)
	}
	return c
}
