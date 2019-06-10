package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	Root    string
	Address string
}

type fileList struct {
	Files []string `json:"files"`
}

func main() {
	c := loadConfig()
	mux := http.NewServeMux()
	h := handlers{dir: c.Root}

	fs := http.FileServer(http.Dir(c.Root))
	mux.Handle("/files/", http.StripPrefix("/files/", fs))
	mux.HandleFunc("/files", h.handleFiles())

	err := http.ListenAndServe(c.Address, mux)
	if err != nil {
		panic(err)
	}
}

func loadConfig() config {
	executablePath := os.Args[0]
	executableName := filepath.Base(executablePath)
	executableName = strings.TrimSuffix(executableName, ".exe")
	b, err := ioutil.ReadFile(executableName + ".json")
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

type handlers struct {
	dir string
}

func (h *handlers) handleFiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			h.fileListJSON(w, r)
		case "POST":
			h.fileUpload(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (h *handlers) fileListJSON(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(h.dir)
	if err != nil {
		panic(err)
	}

	fileNames := fileList{make([]string, len(files))}

	for i, f := range files {
		fileNames.Files[i] = f.Name()
	}

	b, _ := json.Marshal(fileNames)
	w.Write(b)
}

func (h *handlers) fileUpload(w http.ResponseWriter, r *http.Request) {

}
