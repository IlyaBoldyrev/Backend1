package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type getPost struct {
}

type Employee struct {
	Name   string  `json:"name" xml:"name"`
	Age    int     `json:"age" xml:"age"`
	Salary float32 `json:"salary" xml:"salary"`
}

type UploadHandler struct {
	HostAddr  string
	UploadDir string
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusBadRequest)
			return
		}
		filePath := h.UploadDir + "/" + header.Filename

		err = os.WriteFile(filePath, data, 0777)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "File %s has been successfully uploaded\n", header.Filename)
		fileLink := h.HostAddr + "/" + header.Filename

		req, err := http.NewRequest(http.MethodHead, fileLink, nil)

		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to check file", http.StatusInternalServerError)
			return
		}
		cli := &http.Client{}
		resp, err := cli.Do(req)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to check file", http.StatusInternalServerError)
			return
		}
		if resp.StatusCode != http.StatusOK {
			log.Println(err)
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, fileLink)

	case http.MethodGet:
		files, err := os.ReadDir(h.HostAddr)
		if err != nil {
			http.Error(w, "Unable to read files on server", http.StatusInternalServerError)
		}
		requestedExtension := r.Header.Get("Extension")

		for _, file := range files {
			info, err := file.Info()
			if err != nil {
				http.Error(w, "Unable to get info about file", http.StatusInternalServerError)
			}
			ext := filepath.Ext(h.HostAddr + info.Name())
			if ext == requestedExtension || requestedExtension == "" {
				size := fmt.Sprint(info.Size())
				fmt.Fprintf(w, info.Name()+" "+ext+" "+size+"\n")
			}
		}
	}
}

func (h *getPost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		name := r.FormValue("name")
		fmt.Fprintf(w, "Parsed query-param with key \"name\": %s", name)
	case http.MethodPost:
		var employee Employee
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			err := json.NewDecoder(r.Body).Decode(&employee)
			if err != nil {
				http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
				return
			}
		case "application/xml":
			err := xml.NewDecoder(r.Body).Decode(&employee)
			if err != nil {
				http.Error(w, "Unable to unmarshal XML", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "Unknown content type", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Got a new employee!\nName: %s\nAge: %dy.o.\nSalary %0.2f\n",
			employee.Name,
			employee.Age,
			employee.Salary,
		)
	}
}

func main() {
	uploadHandler := &UploadHandler{
		UploadDir: "/upload",
		HostAddr:  "/upload",
	}
	http.Handle("/upload", uploadHandler)

	h := &getPost{}
	http.Handle("/", h)

	dirToServe := http.Dir(uploadHandler.UploadDir)
	fs := &http.Server{
		Addr:         ":8080",
		Handler:      http.FileServer(dirToServe),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go fs.ListenAndServe()

	srv := &http.Server{
		Addr:         ":80",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	srv.ListenAndServe()

}
