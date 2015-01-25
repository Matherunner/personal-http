package main

import (
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func logConnection(r *http.Request) {
	log.Printf("[%s] %s %s\n", r.RemoteAddr, r.Method, r.URL)
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	logConnection(r)
	w.WriteHeader(http.StatusNotFound)
}

func serveList(w http.ResponseWriter, r *http.Request) {
	logConnection(r)
}

func serveUpload(w http.ResponseWriter, r *http.Request) {
	logConnection(r)
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("upload.gtpl")
		if err != nil {
			log.Fatalf("Failed to parse upload.gtpl (%s)\n", err.Error())
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Printf("[%s] Failed to write upload form to client (%s)\n:",
				r.RemoteAddr, err.Error())
			return;
		}
	} else if r.Method == "POST" {
		w.WriteHeader(http.StatusNoContent)
		mpReader, err := r.MultipartReader()
		if err != nil {
			log.Printf("[%s] Failed to read multipart form (%s)\n",
				r.RemoteAddr, err.Error())
			return
		}

		log.Printf("[%s] Got multipart form of size %d\n",
			r.RemoteAddr, r.ContentLength)

		var part *multipart.Part
		for {
			part, err = mpReader.NextPart()
			if err == io.EOF {
				return
			} else if err != nil {
				continue
			}

			if part.FileName() != "" {
				log.Printf("[%s] Downloading \"%s\"\n", r.RemoteAddr, part.FileName())
				break
			}

			part.Close()
		}
		defer part.Close()

		outfile, err := os.OpenFile(filepath.Base(part.FileName()),
			os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0600)
		if err != nil {
			log.Printf("[%s] Failed to create file (%s)\n", r.RemoteAddr, err.Error())
			return
		}
		defer outfile.Close()

		_, err = io.Copy(outfile, part)
		if err != nil {
			log.Printf("[%s] Error occured when transferring \"%s\" (%s)\n",
				r.RemoteAddr, part.FileName(), err.Error())
			return
		}

		log.Printf("[%s] Transfer of \"%s\" completed\n", r.RemoteAddr, part.FileName())
	}
}

func main() {
	http.HandleFunc("/", serveRoot)
	http.HandleFunc("/list", serveList)
	http.HandleFunc("/upload", serveUpload)

	log.Println("Server started.")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatalf("Failed to listen to port (%s)\n", err.Error())
	}
}
