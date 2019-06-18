package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type R struct {
	Result string   `json:"result,omitempty"`
	Files  []string `json:"files,omitempty"`
}

type config struct {
	Root    string
	Address string
}

func main() {
	c := loadConfig()
	mux := http.NewServeMux()
	h := handlers{dir: c.Root}

	fs := http.FileServer(http.Dir(c.Root))
	mux.Handle("/files/", http.StripPrefix("/files/", fs))
	mux.HandleFunc("/files", h.HandleFileListJSON())
	mux.HandleFunc("/upload", h.HandleUploadFiles())

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

func jsonResponse(w http.ResponseWriter, j R) {
	b, _ := json.Marshal(j)
	w.Write(b)
}

type handlers struct {
	dir string
}

func (h *handlers) HandleUploadFiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			b, err := Asset("templates/upload.html")
			if err != nil {
				panic(err)
			}
			w.Write(b)
		case "POST":
			h.fileUpload(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (h *handlers) HandleFileListJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files, err := ioutil.ReadDir(h.dir)
		if err != nil {
			panic(err)
		}

		fileNames := make([]string, len(files))

		for i, f := range files {
			fileNames[i] = f.Name()
		}

		jsonResponse(w, R{Files: fileNames})
	}
}

func (h *handlers) fileUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(16 * 1024 * 1024)
	uploadFile, handler, err := r.FormFile("upload_file")
	if err != nil {
		jsonResponse(w, R{Result: err.Error()})
		return
	}
	defer uploadFile.Close()

	outputFile, err := os.OpenFile(filepath.Join(h.dir, handler.Filename),
		os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		jsonResponse(w, R{Result: err.Error()})
		return
	}
	defer outputFile.Close()
	written, err := io.Copy(outputFile, uploadFile)
	if err != nil {
		jsonResponse(w, R{Result: err.Error()})
		return
	}
	jsonResponse(w, R{Result: strconv.FormatInt(written, 10)})
}
